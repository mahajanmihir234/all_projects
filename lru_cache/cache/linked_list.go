package cache

type Node struct {
	key   interface{}
	value interface{}
	prev  *Node
	next  *Node
}
