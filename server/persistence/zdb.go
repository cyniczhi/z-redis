/**
This file defines the structure of a zdb file
*********** Note ***********
1. Only string type value implemented now
2. No compress for key_value_pairs

3. Multiple databases supported
****************************

**************************** ZDB file ******************************
| ZREDIS | db_version | databases | EOF |
ZREDIS: 	'Z', 'R', 'E', 'D', 'I', 'S'  5 byte
db_version: 0  							  1 byte

**************************** databases ******************************
| database0 | database3 | ... |

**************************** database ******************************
| SELECTDB | db_number | key_value_pairs |
SELECTDB:    							  8 byte
db_number: number  						  1 byte

**************************** key value pair ******************************
| TYPE | KEY | VALUE |
TYPE: const [0, 1, ..., n]: 			  1 byte
KEY: | LEN | string | 					  int32, []byte
VALUE: | LEN | string | 				  int32, []byte

 */
package persistence

import (
	"io/ioutil"
	"log"
	"github.com/cyniczhi/z-redis/server/core"
	"os"
	"bufio"
	"fmt"
	"io"
)

var dbFlag = [8]byte{'S', 'E', 'L', 'E', 'C', 'T', 'D', 'B'}
var zFlag = [6]byte{'Z', 'R', 'E', 'D', 'I', 'S'}

type zDbFile struct {
	databases []*zDatabase
}

type zDatabase struct {
	id      byte
	content []*zKvPair
}

type zKvPair struct {
	valType byte   // type: 0 string
	key     []byte // key: string
	val     []byte // val: string by default
}

// return the buffer of one database
func (db *zDatabase) buff() []byte {
	ret := make([]byte, 0)
	ret = append(ret, dbFlag[:]...)
	ret = append(ret, db.id)
	for _, p := range db.content {
		ret = append(ret, p.valType)
		ret = append(ret, p.key...)
		ret = append(ret, p.val...)
	}
	return ret
}

// Add a database to zdb file from a hash map dict
func (z *zDbFile) AddDatabase(dbNum int, hMap map[string]*core.ZObject) {
	db := new(zDatabase)
	db.id = byte(dbNum)
	for k, v := range hMap {
		db.add(k, v.Ptr.(string))
	}
	z.databases = append(z.databases, db)
}

// persistent Databases
func Persistence(databases []*core.Database) {
	zdb := new(zDbFile)
	zdb.databases = make([]*zDatabase, 0)

	// buffer to be persisted
	result := make([]byte, 0)
	result = append(result, zFlag[:]...)
	result = append(result, core.DbVersion)
	for _, db := range databases {
		zdb.AddDatabase(int(db.ID), db.Dict)
	}
	for _, db := range zdb.databases {
		result = append(result, db.buff()...)
	}

	err := ioutil.WriteFile(core.DefaultZdbFilePath, append(result, core.DbVersion), 0644)
	check(err)
}

func LoadDatabases() (ret []*core.Database, ok bool) {
	if fileObj, err := os.Open(core.DefaultZdbFilePath); err == nil {
		defer fileObj.Close()

		reader := bufio.NewReader(fileObj)

		// TODO: validation the db file
		buf := make([]byte, 6)
		if _, err := reader.Read(buf); err == nil {
			for i := 0; i < 6; i++ {
				if buf[i] != zFlag[i] {
					log.Printf("Database file <%s> illegal", err)
					return nil, false
				}
			}
		} else {
			log.Printf("Database file <%s> illegal", err)
			return nil, false
		}

		// read databases into content
		content := make([]byte, 0)
		buf = make([]byte, core.BuffSize)
		for idx := 0; ; idx++ {
			if n, err := reader.Read(buf); err == nil {
				content = append(content, buf[0:n]...)
				if n < core.BuffSize {
					content = content[:n+idx*core.BuffSize]
					break
				}
			} else if err == io.EOF {
				content = append(content, buf[0:n]...)
				break
			} else {
				log.Printf("Load database file <%s> error ", err)
				panic(err)
			}
		}

		// parse content into databases
		ret := make([]*core.Database, 0)

		// read zdb version
		log.Printf("Version: %d", content[0])
		content = content[1:]

		// TODO: more strict boundary condition check, not enough now
		for {
			// validation
			if len(content) < 8 {
				log.Printf("Database file <%s> illegal, load err: %s", core.DefaultZdbFilePath, err)
				return nil, false
			}
			for i, c := range content[0:8] {
				if c != dbFlag[i] {
					log.Printf("Database file <%s> illegal, load err: %s", core.DefaultZdbFilePath, err)
					return nil, false
				}
			}

			// init a blank database
			dbTmp := new(core.Database)
			dbTmp.Dict = make(map[string]*core.ZObject, 100)
			lru := new(core.LRUDict)
			lru.Head = nil
			lru.Tail = nil
			lru.Max = core.MaxCachedSize
			lru.Dict = make(map[string]*core.Node, 100)
			dbTmp.ExpireDict = lru

			dbTmp.ID = int32(content[8])
			content = content[9:]
			newDbFlag := false
			// parse key_value_pairs into a database
			for {
				if len(content) < 2 {
					// no key_val_pairs
					return ret, true
				}
				for i, c := range content[0:8] {
					if c != dbFlag[i] {
						break
					} else if i == 7 {
						// dbFlag read
						// start a new database
						newDbFlag = true
					}
				}
				if newDbFlag {
					break
				}

				vType := content[0]
				content = content[1:]

				if vType != core.ObjectTypeString {
					panic(fmt.Sprintf("Value type <%s> not supported yet", vType))
				}

				lenK := core.Byte2Int(content[0:4])
				key := content[4:lenK+4]
				lenV := core.Byte2Int(content[lenK+4:lenK+8])
				val := content[lenK+8:lenK+8+lenV]
				dbTmp.Set(string(key), core.CreateObject(core.ObjectTypeString, string(val)))

				content = content[lenK+lenV+8:]
			}

			// append database
			ret = append(ret, dbTmp)
		}
		return ret, true
	} else {
		log.Printf("Error occured when loading zdb file: %s", err)
	}
	return nil, false

}

// add a key_value pair to a zDatabase
func (db *zDatabase) add(key string, val string) {
	pair := new(zKvPair)
	pair.valType = core.ObjectTypeString
	pair.key = append(core.Int2Byte(uint32(len(key))), key...)
	pair.val = append(core.Int2Byte(uint32(len(val))), val...)
	db.content = append(db.content, pair)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
