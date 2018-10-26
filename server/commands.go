// FIXME: db.Dict(map) is unsafe

package server

import "log"

type Command struct {
	Name string
	Proc cmdFunc
}

type cmdFunc func(c *Client, s *Server)

func SetCommand(c *Client, s *Server) {
	objKey := c.Argv[1]
	objVal := c.Argv[2]
	if c.Argc != 3 {
		c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong number of arguments for 'set' command"))
		return
	}
	if stringKey, ok1 := objKey.Ptr.(string); ok1 {
		if stringValue, ok2 := objVal.Ptr.(string); ok2 {
			c.Db.Dict[stringKey] = CreateObject(ObjectTypeString, stringValue)
			c.addReply(CreateObject(ObjectTypeString, "OK"))
			return
		} else {
			c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong <value> of arguments for 'set' command"))
			return
		}
	} else {
		c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong <key> of arguments for 'set' command"))
		return
	}
}

func GetCommand(c *Client, s *Server) {
	db := c.Db
	objKey := c.Argv[1]
	log.Println(objKey.Ptr.(string))
	if o, ok := db.Dict[objKey.Ptr.(string)]; ok && (o != nil) {
		c.addReply(o)
	} else {
		c.addReply(CreateObject(ObjectTypeString, "nil"))
	}
}

func DelCommand(c *Client, s *Server) {
	db := c.Db
	objKey := c.Argv[1]
	delete(db.Dict, objKey.Ptr.(string))
	//c.addReply(CreateObject(ObjectTypeString, "nil"))
}
