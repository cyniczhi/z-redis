// FIXME: db.Dict(map) is unsafe

package server

import (
	"fmt"
)

type Command struct {
	Name string
	Proc cmdFunc
}

type cmdFunc func(c *Client)

func SetCommand(c *Client) {
	objKey := c.Argv[1]
	objVal := c.Argv[2]
	if c.Argc != 3 {
		c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong number of arguments for 'set' command"))
		return
	}
	if stringKey, ok1 := objKey.Ptr.(string); ok1 {
		if o, ok2 := c.Db.set(stringKey, objVal); ok2 && o != nil {
			c.addReply(o)
		} else {
			c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong <value> of arguments for 'set' command"))
		}
	} else {
		c.addReply(CreateObject(ObjectTypeString, "(error) ERR wrong <key> of arguments for 'set' command"))
	}
}

func GetCommand(c *Client) {
	db := c.Db
	objKey := c.Argv[1]
	if o, ok := db.get(objKey.Ptr.(string)); ok != nil {
		c.addReply(CreateObject(ObjectTypeString, "nil"))
	} else {
		c.addReply(o)
	}
}

func DelCommand(c *Client) {
	if key, ok1 := c.Argv[1].Ptr.(string); ok1 {
		c.Db.del(key)
		c.addReply(CreateObject(ObjectTypeString, fmt.Sprintf("Key %s deleted", key)))
	} else {
		c.addReply(CreateObject(ObjectTypeString, fmt.Sprintf("(error) ERR Del %s error", key)))
	}
}
