package main

import (
	// "fmt"
	// "os"
	// "strconv"
)

type Node struct {
	prev  *Node
	next  *Node
	id uint32
	price uint32
	size uint32
	time int64
	counter uint32
}

func newNode(id uint32, price uint32, size uint32, time int64) *Node {
	return &Node{
		prev: nil,
		next: nil,
		id: id,
		price: price,
		size: size,
		time: time,
		counter: 1}
}

type LinkedList struct {
	mut chan struct{}
	head   *Node
	length int
}

func newLinkedList() *LinkedList {
	return &LinkedList{
		mut: initMutex(),
		head:   nil,
		length: 0}
}

func (ll *LinkedList) lock() {
	lockMutex(ll.mut)
}

func (ll *LinkedList) unlock() {
	unlockMutex(ll.mut)
}

func (ll *LinkedList) getHead() *Node {
	return ll.head
}

func (ll *LinkedList) getLength() int {
	return ll.length
}

func (ll *LinkedList) insert(id uint32, price uint32, size uint32, time int64) {
	var node *Node = newNode(id, price, size, time)
	if ll.length == 0 {
		node.next = nil
		node.prev = nil
		ll.head = node
	} else {
		node.next = ll.head
		node.prev = nil
		ll.head.prev = node
	}
	ll.length += 1
	ll.head = node
}

func (ll *LinkedList) getNodeById(id uint32) *Node { 
	var dummyNode *Node = ll.getHead()
	for i := 0; i < ll.getLength(); i++ {
		if dummyNode.id == id {
			return dummyNode
		}
		dummyNode = dummyNode.next
	}
	return nil
}

func (ll *LinkedList) deleteNode(toBeDeleted *Node) bool {
	if ll.getLength() == 0 {
		return false
	}
	var answer bool = false
	var dummyNode *Node = ll.getHead()
	for i := 0; i < ll.getLength(); i++ {
		if dummyNode == toBeDeleted {
			answer = true
			break
		}
		dummyNode = dummyNode.next
	}
	if answer == true {
		if ll.getHead() == toBeDeleted {
			ll.head = ll.head.next
		}
		if toBeDeleted.next != nil {
			toBeDeleted.next.prev = toBeDeleted.prev
		}
		if toBeDeleted.prev != nil {
			toBeDeleted.prev.next = toBeDeleted.next
		}
		ll.length -= 1
	}
	return answer
}

// func (ll *LinkedList) printValues() {
// 	if ll.getLength() == 0 {
// 		fmt.Fprintf(os.Stderr, "empty\n");
// 		return
// 	}
// 	var answer string = ""
// 	var dummy *Node = ll.head
// 	for i := 0; i < ll.getLength(); i++ {
// 		answer += "(" + strconv.Itoa(dummy.id) + " with " + strconv.Itoa(dummy.price) + "@" + strconv.Itoa(dummy.size) + ") "
// 		dummy = dummy.next
// 	}
// 	fmt.Fprintf(os.Stderr, answer + "\n")
// }


