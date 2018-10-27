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
	"github.com/cyniczhi/z-redis/server/core"
	"github.com/cyniczhi/z-redis/server/persistence"
)

func CreateServer() (server *Server) {
	server = new(Server)

	server.Addr = "0.0.0.0:9999"
	server.Pid = os.Getpid()

	// allocate mem for databases
	if dbs, ok := persistence.LoadDatabases(); ok {
		server.Db = dbs
		server.DbNum = len(dbs)
	} else {
		server.DbNum = core.DefaultDbNumber
		server.Db = make([]*core.Database, core.DefaultDbNumber)
		for i := 0; i < server.DbNum; i++ {
			server.Db[i] = new(core.Database)
			server.Db[i].Dict = make(map[string]*core.ZObject, 100)
			server.Db[i].ID = int32(i)
		}
	}

	// init LRU
	for i := 0; i < server.DbNum; i++ {
		lru := new(core.LRUDict)
		lru.Head = nil
		lru.Tail = nil
		lru.Max = core.MaxCachedSize
		lru.Dict = make(map[string]*core.Node, 100)
		server.Db[i].ExpireDict = lru
	}

	server.Start = time.Now().UnixNano() / 1000000

	// add commands
	getCommand := &Command{Name: "get", Proc: GetCommand}
	setCommand := &Command{Name: "set", Proc: SetCommand}
	delCommand := &Command{Name: "Del", Proc: DelCommand}

	server.Commands = map[string]*Command{
		"get": getCommand,
		"set": setCommand,
		"Del": delCommand,
	}
	return server
}

func (s *Server) Run() {

	// handle os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go sigHandler(c, s)

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


func sigHandler(c chan os.Signal, server *Server) {
	for s := range c {
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			exitHandler(server)
		default:
			log.Println("unexpected signal ", s)
		}
	}
}

// TODO: persistence
func exitHandler(s *Server) {
	log.Println("exiting z-redis...")
	log.Println("Persisting databases...")
	persistence.Persistence(s.Db)
	log.Println("Databases saved")
	log.Println("bye ")
	os.Exit(0)
}
