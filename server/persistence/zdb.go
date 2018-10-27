/**
This file defines the structure of a zdb file
*********** Note ***********
1. Only string type value implemented now
2. No compress for key_value_pairs
****************************

**************************** ZDB file ******************************
| ZREDIS | db_version | databases | EOF |
ZREDIS: 	'Z', 'R', 'E', 'D', 'I', 'S'  5byte
db_version: 0  							  1byte

**************************** databases ******************************
| database0 | database3 | ... |

**************************** database ******************************
| SELECTDB | db_number | key_value_pairs |
SELECTDB: 1  							  1 byte
db_number: number  						  1 byte

**************************** key value pair ******************************
| TYPE | KEY | VALUE |
TYPE: const [0, 1, ..., n]: 			  1byte
KEY: string								  []byte
VALUE: string							  []byte

 */
package persistence

import (
	"io/ioutil"
)

const (
	dbVersion byte = 0
	defaultFilePath = "test.zdb"
)

type zDbFile struct {
	startFlag [6]byte // "ZREDIS"
	dbVersion byte    // 0
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
func (db *zDatabase)buff() []byte {
	ret := make([]byte, 0)
	for _, p := range db.content {
		ret = append(ret, p.valType)
		ret = append(ret, p.key...)
		ret = append(ret, p.val...)
	}
	return ret
}

// persistent a zdb file
func (z *zDbFile) Persistence() {
	result := make([]byte, 0)
	result = append(result, z.startFlag[:]...)
	result = append(result, z.dbVersion)
	for _, db := range z.databases {
		result = append(result, 1)
		result = append(result, db.id)
		result = append(result, db.buff()...)
	}

	err := ioutil.WriteFile(defaultFilePath, append(result, dbVersion), 0644)
	check(err)
}

// Add a key_value pair to a zDatabase
func (db *zDatabase)Add(key string, val string) {
	pair := new(zKvPair)
	pair.valType = 0
	pair.key = []byte(key)
	pair.val = []byte(val)
	db.content = append(db.content, pair)
}

func Test() {
	db := new(zDatabase)
	db.id = 0
	db.Add("aaa", "aaa")
	db.Add("bbb", "aaa")
	db.Add("ccc", "aaa")
	db.Add("ddd", "aaa")
	db.Add("eee", "aaa")

	zdb := new(zDbFile)
	zdb.startFlag = [6]byte{'Z', 'R', 'E', 'D', 'I', 'S'}
	zdb.dbVersion = 2
	zdb.databases = make([]*zDatabase, 0)
	zdb.databases = append(zdb.databases, db)
	zdb.Persistence()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
