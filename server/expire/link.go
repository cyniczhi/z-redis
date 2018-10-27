package expire

type Node struct {
	prev       *Node
	next       *Node
	key        string
	expireTime int64
}

type dict map[string]*Node
