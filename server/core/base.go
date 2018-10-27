package core

// TODO: really support different type obj
const (
	ObjectTypeString = 1
	ObjectTypeList   = 2
	ObjectTypeSet    = 3
	ObjectTypeZSet   = 4
	ObjectTypeHash   = 5
)

var MaxCachedSize int32 = 3
var DefaultZdbFilePath = "default.zdb"
var DbVersion byte = 0

type ZObject struct {
	ObjectType int
	Ptr        interface{}
}

type BaseDict map[string]*ZObject

// Little Endian: uint32 to []byte
func Int2Byte(data uint32) (ret []byte) {
	var len uint32 = 4
	ret = make([]byte, len)
	var tmp uint32 = 0xff
	var index uint32 = 0
	for index = 0; index < len; index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

// Little Endian: []byte to uint32
func Byte2Int(data []byte) uint32 {
	var ret uint32 = 0
	var len uint32 = 4
	var i uint32 = 0
	for i = 0; i < len; i++ {
		ret = ret | (uint32(data[i]) << (i * 8))
	}
	return ret
}
