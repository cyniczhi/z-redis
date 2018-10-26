package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// handle os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go sigHandler(c)

	// network handler
	socket, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println("server init err ")
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

// worker
func handle(conn net.Conn) {
	for {
		buff, err := parseQuery(conn)
		if err != nil {
			log.Println("obtain and parse query err")
			return
		}
		result := handleQuery(buff)
		response(conn, result)
	}
}

// TODO: parse query from client
func parseQuery(conn net.Conn) (query string, err error) {
	buff := make([]byte, 512)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("parse query: conn.Read err!=nil", err, "---len---", n, conn)
		conn.Close()
		return "", err
	}
	//log.Println(string(buff))
	return string(buff), nil
}

// TODO: handle request and response
func handleQuery(buff string) string {
	resp := buff + " from Client"
	return resp
}

// response
func response(conn net.Conn, buff string) {
	conn.Write([]byte(buff))
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

func exitHandler() {
	log.Println("exiting z-redis...")
	log.Println("bye ")
	os.Exit(0)
}
