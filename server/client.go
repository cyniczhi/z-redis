package server

import (
	"net"
	"regexp"
	"log"
	"strings"
)

// one client is mapped to one TCP connection
type Client struct {
	Cmd      *Command
	Argv     []*ZObject
	Argc     int
	Db       *Database
	QueryBuf string
	Buf      string
	Conn     net.Conn
}

func (c *Client) parseQuery() (err error) {
	// TODO: buffer here may not enough
	buff := make([]byte, 512)
	n, err := c.Conn.Read(buff)

	if err != nil {
		log.Println("parse query: conn.Read err!=nil", err, "---len---", n, c.Conn)
		c.Conn.Close()
		return err
	}
	tmp := string(buff[0:n])
	log.Println(tmp)
	parts := strings.Split(tmp, "\n")
	c.QueryBuf = parts[0]
	return nil
}

func (c *Client) handleQuery(){
	r := regexp.MustCompile("[^\\s]+")
	parts := r.FindAllString(strings.Trim(c.QueryBuf, " "), -1)
	argc, argv := len(parts), parts
	c.Argc = argc
	//c.Argv = make([]*object.GodisObject, 5)
	j := 0
	for _, v := range argv {
		c.Argv[j] = CreateObject(ObjectTypeString, v)
		j++
	}
}

func (c *Client) addReply(o *ZObject) {
	c.Buf = o.Ptr.(string)
}

func (c *Client)Run(s *Server) {
	// TODO: use chanel to communicate here
	for {
		err := c.parseQuery()
		if err != nil {
			log.Println("Obtain and parse query err")
			log.Println("Connection closed.")
			return
		}
		c.handleQuery()
		s.ProcessCommand(c)
		c.Conn.Write([]byte(c.Buf))
	}
}
