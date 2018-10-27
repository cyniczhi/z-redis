package core

// TODO: diff type error handling

type Database struct {
	Dict       BaseDict
	ExpireDict *LRUDict
	ID         int32
}


func CreateObject(t int, ptr interface{}) (o *ZObject) {
	o = new(ZObject)
	o.ObjectType = t
	o.Ptr = ptr
	return
}

func (d *Database) Get(key string) (*ZObject, *error) {

	if o, ok := d.Dict[key]; ok && (o != nil) {
		// update core dict
		d.ExpireDict.Renew(key)
		return o, nil
	} else {
		return nil, new(error)
	}
}

func (d *Database) Set(key string, val *ZObject) (*ZObject, bool) {
	if val, ok := val.Ptr.(string); ok {
		if d.ExpireDict.Has(key) {
			// if exist key, renew core dict
			d.ExpireDict.Renew(key)
		} else {
			// if new key, insert into core dict
			d.ExpireDict.Insert(key, d)
		}
		valObj := CreateObject(ObjectTypeString, val)
		d.Dict[key] = valObj
		return valObj, true
	} else {
		return nil, false
	}
}

func (d *Database) Del(key string) {
	d.ExpireDict.Del(key, d)
}
