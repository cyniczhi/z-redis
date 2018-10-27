package server

import (
	"github.com/cyniczhi/z-redis/server/expire"
)

// TODO: diff type error handling


type Database struct {
	Dict       dict
	ExpireDict expire.LRUDict
	ID         int32
}

type dict map[string]*ZObject

type ZObject struct {
	ObjectType int
	Ptr        interface{}
}

// TODO: really support different type obj
const ObjectTypeString = 1

func CreateObject(t int, ptr interface{}) (o *ZObject) {
	o = new(ZObject)
	o.ObjectType = t
	o.Ptr = ptr
	return
}

func (d *Database)get(key string) (*ZObject, *error) {

	if o, ok := d.Dict[key]; ok && (o != nil) {
		// update expire dict
		return o, nil
	} else {
		return nil, new(error)
	}
}

func (d *Database)set(key string, val *ZObject) (*ZObject, *error) {
	if val, ok := val.Ptr.(string); ok {
		valObj := CreateObject(ObjectTypeString, val)
		d.Dict[key] = valObj
		return valObj, nil
	} else {
		return nil, new(error)
	}
}

func (d* Database)del(key string) {
	delete(d.Dict, key)
}
