package core

// TODO: really support different type obj
const (
	ObjectTypeString = 1
	ObjectTypeList   = 2
	ObjectTypeSet    = 3
	ObjectTypeZSet   = 4
	ObjectTypeHash   = 5
)

type ZObject struct {
	ObjectType int
	Ptr        interface{}
}

type BaseDict map[string]*ZObject
