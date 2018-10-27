package server

import (
	"net"
	"fmt"
	"log"
	"github.com/cyniczhi/z-redis/server/core"
)

type Server struct {
	Addr         string
	Db           []*Database
	DbNum        int
	Start        int64
	Port         int32
	NextClientID int32
	Clients      int32
	Pid          int
	Commands     map[string]*Command
	Dirty        int64
}

// record and maintain a connection
func (s *Server) CreateClient(conn net.Conn) (c *Client) {
	c = new(Client)
	c.Db = s.Db[0]
	c.Argv = make([]*core.ZObject, 5)
	c.QueryBuf = ""
	c.Conn = conn
	return c
}

func (s *Server) ProcessCommand(c *Client) {
	v := c.Argv[0].Ptr
	name, ok := v.(string)
	if !ok {
		log.Println("(error) ERR unknown command ", name)
		c.addReply(CreateObject(ObjectTypeString, fmt.Sprintf("(error) ERR unknown command '%s'", name)))
	}
	if cmd, ok := s.Commands[name]; !ok {
		c.addReply(CreateObject(ObjectTypeString, fmt.Sprintf("(error) ERR unknown command '%s'", name)))
	} else {
		c.Cmd = cmd
		call(c)
	}
}

func call(c *Client) {
	c.Cmd.Proc(c)
}
