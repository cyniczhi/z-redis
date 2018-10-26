package main

import "github.com/cyniczhi/z-redis/server"

func main()  {
	server := new(server.Server)
	server.Start()
}
