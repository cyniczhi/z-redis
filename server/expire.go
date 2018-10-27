package server

import (
	"log"
	"fmt"
	"github.com/cyniczhi/z-redis/server/baseds"
)

const (
	maxKeys = 5
)

type LRUDict struct {
	Dict baseds.ExpireDict
	max  int32 // max keys in dict
	head *baseds.Node // the Next key need to be expired
	tail *baseds.Node // recently used key
	len  int32
}
// delete Key
func (d *LRUDict) Del(key string, db *Database) {
	if n, ok := d.Dict[key]; ok {
		if n == d.head {
			// delete head
			d.head = n.Next
			d.head.Prev = nil
		} else if n == d.tail {
			// delete tail
			d.tail = n.Prev
			d.tail.Next = nil
		} else {
			// delete internal
			n.Prev.Next = n.Next
			n.Next.Prev = n.Prev
		}
		// release mem
		delete(db.Dict, key)
		delete(d.Dict, n.Key)
		n = nil
		d.len--
	} else {
		log.Println(fmt.Sprintf("Query key %s from baseds dict error", key))
	}
}

// Renew key: move node to tail of the link
func (d *LRUDict) Renew(key string) (ok bool){
	if n, ok := d.Dict[key]; ok {
		if n == d.head {
			d.tail.Next = n
			n.Next.Prev = nil
			d.head = n.Next

			n.Prev = d.tail
			n.Next = nil
			d.tail = n
		} else if n == d.tail {
			// do nothing
			return true
		} else {
			d.tail.Next = n
			n.Prev.Next = n.Next
			n.Next.Prev = n.Prev
			n.Prev = d.tail
			n.Next = nil
			d.tail = n
		}
	}
	return true
}

// Insert key: insert key, if num of keys > max num of the keys need to cache defined, expire the head.
func (d *LRUDict) Insert(key string, db *Database) (ok bool) {
	// TODO: err handling
	node := new(baseds.Node)
	node.Key = key

	if d.len == 0 {
		// add 1st node
		d.head = node
		d.tail = node
	} else {
		// insert new node to the tail
		d.tail.Next = node
		node.Prev = d.tail
		d.tail = node
	}
	d.Dict[key] = node
	d.len++

	if d.len > d.max {
		// baseds the head node
		d.Del(d.head.Key, db)
	}
	return true
}

// Check if key in expire dict
func (d *LRUDict) Has(key string) bool{
	if _, ok := d.Dict[key]; ok {
		return true
	}
	return false
}
