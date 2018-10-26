package server

type Database struct {
	Dict dict
	ID   int32

	//Expires dict
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
