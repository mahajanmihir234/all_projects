package main

type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

type LRUCache struct {
	store    map[int]*Node
	capacity int
	head     *Node
	tail     *Node
}

func Constructor(capacity int) LRUCache {
	dummyHead := &Node{}
	dummyTail := &Node{}

	dummyHead.next = dummyTail
	dummyTail.prev = dummyHead

	return LRUCache{
		capacity: capacity,
		store:    map[int]*Node{},
		head:     dummyHead,
		tail:     dummyTail,
	}
}

func (this *LRUCache) bringToFront(node *Node) {
	currentPrev := node.prev
	currentNext := node.next
	currentPrev.next = currentNext
	currentNext.prev = currentPrev

	this.putInFront(node)
}

func (this *LRUCache) putInFront(node *Node) {
	currentFirst := this.head.next
	this.head.next = node
	node.prev = this.head
	node.next = currentFirst
	currentFirst.prev = node
}

func (this *LRUCache) Get(key int) int {
	node, ok := this.store[key]
	if !ok {
		return -1
	}
	this.bringToFront(node)
	return node.value
}

func (this *LRUCache) Put(key int, value int) {
	if node, ok := this.store[key]; ok {
		this.bringToFront(node)
		node.value = value
		return
	}

	if len(this.store) >= this.capacity {
		currentLast := this.tail.prev
		delete(this.store, currentLast.key)
		currentLast.prev.next = this.tail
		this.tail.prev = currentLast.prev
	}

	node := &Node{key: key, value: value}
	this.putInFront(node)
	this.store[key] = node
}

func main() {}
