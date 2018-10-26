/**

 */
package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func CreateServer() (server *Server) {
	server = new(Server)

	server.Addr = "0.0.0.0:9999"
	server.Pid = os.Getpid()
	server.DbNum = 8

	// allocate mem for databases
	server.Db = make([]*Database, server.DbNum)
	for i := 0; i < server.DbNum; i++ {
		server.Db[i] = new(Database)
		server.Db[i].Dict = make(map[string]*ZObject, 100)
	}
	//log.Println("init db begin-->", server.Db)

	server.Start = time.Now().UnixNano() / 1000000

	// add commands
	getCommand := &Command{Name: "get", Proc: GetCommand}
	setCommand := &Command{Name: "set", Proc: SetCommand}
	delCommand := &Command{Name: "del", Proc: DelCommand}

	server.Commands = map[string]*Command{
		"get": getCommand,
		"set": setCommand,
		"del": delCommand,
	}
	return server
}

func (s *Server) Run() {

	// handle os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go sigHandler(c)

	// network handler
	socket, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Println("(ERR) Error occured when initializing network")
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			continue
		}
		c := s.CreateClient(conn)
		go c.Run(s)
	}
}


func sigHandler(c chan os.Signal) {
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitHandler()
		default:
			log.Println("unexpected signal ", s)
		}
	}
}

// TODO: serialize mem to disk
func exitHandler() {
	log.Println("exiting z-redis...")
	log.Println("bye ")
	os.Exit(0)
}
