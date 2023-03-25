package main

import (
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	prev  *Node
	next  *Node
	id int
	price int
	size int
	time int64
}

func newNode(id int, price int, size int, time int64) *Node {
	return &Node{
		prev: nil,
		next: nil,
		id: id,
		price: price,
		size: size,
		time: time}
}

type LinkedList struct {
	head   *Node
	length int
}

func newLinkedList() *LinkedList {
	return &LinkedList{
		head:   nil,
		length: 0}
}

func (ll *LinkedList) getHead() *Node {
	return ll.head
}

func (ll *LinkedList) getLength() int {
	return ll.length
}

func (ll *LinkedList) insert(id int, price int, size int, time int64) {
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

func (ll *LinkedList) getNodeById(id int) *Node {
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

func (ll *LinkedList) printValues() {
	if ll.getLength() == 0 {
		fmt.Fprintf(os.Stderr, "empty\n");
		return
	}
	var answer string = ""
	var dummy *Node = ll.head
	for i := 0; i < ll.getLength(); i++ {
		answer += "(" + strconv.Itoa(dummy.id) + " with " + strconv.Itoa(dummy.price) + "@" + strconv.Itoa(dummy.size) + ") "
		dummy = dummy.next
	}
	fmt.Fprintf(os.Stderr, answer + "\n")
}


