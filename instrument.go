package main

import (
	"fmt"
	"os"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Instrument struct {
	buyBook  *LinkedList
	sellBook *LinkedList
}

func newInstrument() *Instrument {
	return &Instrument{
		buyBook:  newLinkedList(),
		sellBook: newLinkedList()}
}

// process buy
func (i *Instrument) processBuy(id int, price int, size int) {
	var bestSell *Node
	for true {
		bestSell = i.getBestSell()
		if bestSell == nil {
			break
		}
		if price < bestSell.price {
			break
		}
		//quantity traded
		var quantity = min(size, bestSell.size)
		size -= quantity
		bestSell.size -= quantity
		fmt.Fprintf(os.Stderr, "TRADE BUY %v with %v %v@%v : %v\n", id, bestSell.id, quantity, bestSell.price, GetCurrentTimestamp())

		//if bestSell.size == 0 -> remove from instrument
		if bestSell.size == 0 {
			i.sellBook.deleteNode(bestSell)
		}
		if size == 0 {
			break
		}
	}
	if size > 0 {
		i.addBuy(id, price, size)
	}
}

func (i *Instrument) processSell(id int, price int, size int) {
	var bestBuy *Node
	for true {
		bestBuy = i.getBestBuy()
		if bestBuy == nil {
			break
		}
		if price > bestBuy.price {
			break
		}
		//quantity traded
		var quantity = min(size, bestBuy.size)
		size -= quantity
		bestBuy.size -= quantity
		fmt.Fprintf(os.Stderr, "TRADE SELL %v with %v %v@%v : %v\n",id, bestBuy.id, quantity, bestBuy.price, GetCurrentTimestamp())

		//if bestBuy.size == 0 -> remove from instrument
		if bestBuy.size == 0 {
			i.buyBook.deleteNode(bestBuy)
		}
		if size == 0 {
			break
		}
	}
	if size > 0 {
		i.addSell(id, price, size)
	}
}

func (i *Instrument) getBestBuy() *Node {
	if i.buyBook.getLength() == 0 {
		return nil
	}
	var answer *Node = i.buyBook.getHead()

	// traverse through the buyBook
	var currNode *Node = i.buyBook.getHead()
	for a := 0; a < i.buyBook.getLength(); a++ {
		if currNode.price > answer.price || currNode.price == answer.price && currNode.time < answer.time {
			answer = currNode
		}
		currNode = currNode.next
	}
	return answer
}

func (i *Instrument) getBestSell() *Node {
	if i.sellBook.getLength() == 0 {
		return nil
	}
	var answer *Node = i.sellBook.getHead()

	// traverse through the sellBook
	var currNode *Node = i.sellBook.getHead()
	for a := 0; a < i.sellBook.getLength(); a++ {
		if currNode.price < answer.price || currNode.price == answer.price && currNode.time < answer.time {
			answer = currNode
		}
		currNode = currNode.next
	}
	return answer
}

func (i *Instrument) addBuy(id int, price int, size int) {
	currentTime := GetCurrentTimestamp()
	fmt.Fprintf(os.Stderr, "Buy %v added with %v@%v : %v\n", id, size, price, currentTime)
	i.buyBook.insert(id, price, size, currentTime)
}

func (i *Instrument) addSell(id int, price int, size int) {
	currentTime := GetCurrentTimestamp()
	fmt.Fprintf(os.Stderr, "Sell %v added with %v@%v : %v\n", id, size, price, currentTime)
	i.sellBook.insert(id, price, size, currentTime)
}


func (i *Instrument) deleteNode(id int) {
	node := i.buyBook.getNodeById(id)
	if node != nil {
		i.buyBook.deleteNode(node)
		currentTime := GetCurrentTimestamp()
		fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
	}
	node2 := i.sellBook.getNodeById(id)
	if node2 != nil {
		i.sellBook.deleteNode(node2)
		currentTime := GetCurrentTimestamp()
		fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
	}
	currentTime := GetCurrentTimestamp()
	fmt.Fprintf(os.Stderr, "Declined Cancel with ID %v at %v", id, currentTime)
}