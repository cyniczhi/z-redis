package core

type Node struct {
	Prev       *Node
	Next       *Node
	Key        string
	ExpireTime int64
}

type ExpireDict map[string]*Node
