package main

import (
	//"fmt"
	"github.com/cyniczhi/z-redis/server"
	"fmt"
)

//import "github.com/cyniczhi/z-redis/server"

func main()  {
	fmt.Println("Welcome to z-redis server!")
	server := server.CreateServer()
	server.Run()
}
