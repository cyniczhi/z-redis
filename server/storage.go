package server

import (
)

// TODO: diff type error handling


type Database struct {
	Dict       dict
	ExpireDict *LRUDict
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
		// update baseds dict
		d.ExpireDict.Renew(key)
		return o, nil
	} else {
		return nil, new(error)
	}
}

func (d *Database)set(key string, val *ZObject) (*ZObject, bool) {
	if val, ok := val.Ptr.(string); ok {
		if d.ExpireDict.Has(key) {
			// if exist key, renew baseds dict
			d.ExpireDict.Renew(key)
		} else {
			// if new key, insert into baseds dict
			d.ExpireDict.Insert(key, d)
		}
		valObj := CreateObject(ObjectTypeString, val)
		d.Dict[key] = valObj
		return valObj, true
	} else {
		return nil, false
	}
}

func (d* Database) del(key string) {
	d.ExpireDict.Del(key, d)
}
