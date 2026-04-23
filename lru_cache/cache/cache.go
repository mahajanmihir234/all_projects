package cache

import (
	"fmt"
	"sync"
)

type LRUCache struct {
	size           int
	hashMap        map[interface{}]*Node
	linkedListHead *Node
	linkedListTail *Node
	mutex          sync.Mutex
}

func NewLRUCache(size int) LRUCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return LRUCache{
		size:           size,
		hashMap:        map[interface{}]*Node{},
		linkedListHead: head,
		linkedListTail: tail,
	}
}

func (c *LRUCache) Get(key interface{}) interface{} {
	c.mutex.Lock()
	node, ok := c.hashMap[key]
	if !ok || node == nil {
		return nil
	}

	c.moveToFirst(node)
	c.mutex.Unlock()
	return node.value
}

func (c *LRUCache) moveToFirst(node *Node) {
	previousElement := node.prev
	nextElement := node.next

	previousElement.next = nextElement
	nextElement.prev = previousElement

	// Setting new prev and next pointers for current node
	node.prev = c.linkedListHead
	node.next = c.linkedListHead.next

	// Setting pointers for the new position
	firstElement := c.linkedListHead.next
	firstElement.prev = node
	c.linkedListHead.next = node
}

func (c *LRUCache) Put(key interface{}, value interface{}) {
	c.mutex.Lock()
	fmt.Printf("Putting %v key with value %v \n", key, value)
	if node, ok := c.hashMap[key]; !ok {
		if c.size == len(c.hashMap) {
			c.removeLast()
		}
		// Setting value and next and prev
		newNode := Node{key: key, value: value, next: c.linkedListHead.next, prev: c.linkedListHead}

		// Inserting in the hashMap
		c.hashMap[key] = &newNode

		// Setting pointers for the new position
		firstElement := c.linkedListHead.next
		firstElement.prev = &newNode
		c.linkedListHead.next = &newNode
	} else {
		// Setting value
		node.value = value
		c.moveToFirst(node)
	}
	c.mutex.Unlock()
}

func (c *LRUCache) removeLast() {
	lastElement := c.linkedListTail.prev
	newLastElement := lastElement.prev

	newLastElement.next = c.linkedListTail
	c.linkedListTail.prev = newLastElement

	lastElement.next = nil
	lastElement.prev = nil
	delete(c.hashMap, lastElement.key)
}

func (c *LRUCache) Print() {
	head := c.linkedListHead.next
	for head != c.linkedListTail {
		fmt.Println(head.key, head.value)
		head = head.next
	}
}
