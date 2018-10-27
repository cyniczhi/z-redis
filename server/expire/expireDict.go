package expire

const (
	maxKeys = 5
)

type LRUDict struct {
	Dict dict
	max  int32 // max keys in dict
	head *Node // the next key need to be expired
	tail *Node // recently used key
	len  int32
}

// update expire dict when call set(key, v) or get(key)
func (d *LRUDict) update(key string) {
	if node, err := d.Dict[key]; !err && node != nil {
		// key exists
		// move Node to tail
		node.prev.next = node.next
		node.next = nil
		node.prev = d.tail
		d.tail = node
	} else {
		if node == nil {

		}
		// key not exists
		// insert key into dict and check if d.link
		newNode := new(Node)
		newNode.prev = nil
		newNode.next = nil
		newNode.key = key
		d.tail = newNode
		d.len++
		if d.len > d.max {
			// remove head
			d.del(d.head.key)
		}
	}

}

// update expire dict when call del(key)
func (d *LRUDict) del(key string) {

}
