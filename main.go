package main

import (
	"fmt"
	//"github.com/cyniczhi/z-redis/server"
	"github.com/cyniczhi/z-redis/server/core"
	"github.com/cyniczhi/z-redis/server"
)

//import "github.com/cyniczhi/z-redis/server"

func main() {
	fmt.Println("Welcome to z-redis server!")
	server := server.CreateServer()
	server.Run()
	//persistence.TestWrite()
	//persistence.TestRead()
	//test()
}

// FIXME: I don't understand the output when run the code below
type tserver struct {
	dbs []*database
}
type database struct {
	dict map[string]*core.ZObject
}

func test() {
	dbs := make([]*database, 0)
	db := new(database)
	db.dict = make(map[string]*core.ZObject, 0)
	db.dict["a"] = core.CreateObject(core.ObjectTypeString, "1")
	dbs = append(dbs, db)

	s := new(tserver)
	s.dbs = dbs
	fmt.Println(dbs[0].dict)
	fmt.Println(dbs[0].dict["a"].Ptr.(string))
	fmt.Println(s.dbs[0].dict)
	fmt.Println(s.dbs[0].dict["a"].Ptr.(string))

	db.dict["c"] = core.CreateObject(core.ObjectTypeString, "1")
	s.dbs[0].dict["b"] = core.CreateObject(core.ObjectTypeString, "1")

	for k, v := range db.dict {
		tmp := s.dbs[0]
		tmp.dict[k] = v
		//s.dbs[0].dict[k] = v
		s.dbs[0] = tmp
	}
	fmt.Println(s)

}
