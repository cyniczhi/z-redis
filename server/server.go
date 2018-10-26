package server

import (
	"net"
	"fmt"
	"log"
)

type Server struct {
	Addr  string
	Db    []*Database
	DbNum int
	Start            int64
	Port int32
	//RdbFilename      string
	//AofFilename      string
	NextClientID int32
	//SystemMemorySize int32
	Clients  int32
	Pid      int
	Commands map[string]*Command
	Dirty    int64
	//AofBuf           []string
}

// record and maintain a connection
func (s *Server) CreateClient(conn net.Conn) (c *Client) {
	c = new(Client)
	c.Db = s.Db[0]
	c.Argv = make([]*ZObject, 5)
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
		call(c, s)
	}
}

func call(c *Client, s *Server) {
	c.Cmd.Proc(c, s)
}
