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

func (z *zDbFile) Persistence() {

}

func (db *zDatabase)Add(key string, val string) {

}

func Test() {
	t := []byte{'1', '2', 3, 4, 5}
	err := ioutil.WriteFile("test.zdb", append(t, dbVersion), 0644)
	check(err)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
