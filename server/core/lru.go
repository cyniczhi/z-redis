package core

import (
	"log"
	"fmt"
)

type LRUDict struct {
	Dict ExpireDict
	Max  int32 // Max keys in dict
	Head *Node // the Next key need to be expired
	Tail *Node // recently used key
	Len  int32
}

// delete Key
func (d *LRUDict) Del(key string, db *Database) {
	if n, ok := d.Dict[key]; ok {
		if n == d.Head {
			// delete Head
			d.Head = n.Next
			d.Head.Prev = nil
		} else if n == d.Tail {
			// delete Tail
			d.Tail = n.Prev
			d.Tail.Next = nil
		} else {
			// delete internal
			n.Prev.Next = n.Next
			n.Next.Prev = n.Prev
		}
		// release mem
		delete(db.Dict, key)
		delete(d.Dict, n.Key)
		n = nil
		d.Len--
	} else {
		log.Println(fmt.Sprintf("Query key %s from core dict error", key))
	}
}

// Renew key: move node to Tail of the link
func (d *LRUDict) Renew(key string) (ok bool) {
	if n, ok := d.Dict[key]; ok {
		if n == d.Head {
			d.Tail.Next = n
			n.Next.Prev = nil
			d.Head = n.Next

			n.Prev = d.Tail
			n.Next = nil
			d.Tail = n
		} else if n == d.Tail {
			// do nothing
			return true
		} else {
			d.Tail.Next = n
			n.Prev.Next = n.Next
			n.Next.Prev = n.Prev
			n.Prev = d.Tail
			n.Next = nil
			d.Tail = n
		}
	}
	return true
}

// Insert key: insert key, if num of keys > Max num of the keys need to cache defined, expire the Head.
func (d *LRUDict) Insert(key string, db *Database) (ok bool) {
	// TODO: err handling
	node := new(Node)
	node.Key = key

	if d.Len == 0 {
		// add 1st node
		d.Head = node
		d.Tail = node
	} else {
		// insert new node to the Tail
		d.Tail.Next = node
		node.Prev = d.Tail
		d.Tail = node
	}
	d.Dict[key] = node
	d.Len++

	if d.Len > d.Max {
		// core the Head node
		d.Del(d.Head.Key, db)
	}
	return true
}

// Check if key in expire dict
func (d *LRUDict) Has(key string) bool {
	if _, ok := d.Dict[key]; ok {
		return true
	}
	return false
}
