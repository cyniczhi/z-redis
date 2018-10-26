package main

import (
	"fmt"
	"bufio"
	"os"
	"net"
	"strings"
	"log"
)

const (
)

func main() {
	// init
	log.Println("Welcome to z-redis!")
	localAddr, err:= os.Hostname()
	checkError(err)
	serverAddr := "127.0.0.1:9999"

	reader := bufio.NewReader(os.Stdin)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", serverAddr)
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	for {
		// send query
		fmt.Print(localAddr + "[client]> ")
		query, _ := reader.ReadString('\n')

		// TODO: parse query to satisfy redis protocol
		query = strings.Replace(query, "\n", "", -1)
		sendQuery(query, conn)

		// print response
		resp := make([]byte, 1024)
		n, err := conn.Read(resp)
		checkError(err)
		if n == 0 {
			fmt.Println(serverAddr + "[server]> ", "nil")
		} else {
			fmt.Println(serverAddr + "[server]> ", string(resp))
		}
	}

}
func sendQuery(query string, conn net.Conn) (n int, err error) {
	data := []byte(query)
	n, err = conn.Write(data)
	return n, err
}
func checkError(err error) {
	if err != nil {
		log.Println("err occured.", err.Error())
		os.Exit(1)
	}
}
