package main

import (
	"time"
)

type Order struct {
	orderType inputType
	clientChan chan struct{} 
	id uint32
	price uint32
	size uint32
}


func instrumentFunc(done <-chan struct{}, name string) chan Order {
	orderChan := make(chan Order)
	go func() {
		instrument := newInstrument(name)
		for {
			select {
			case order:= <- orderChan:
				if order.orderType == inputBuy {
					instrument.processBuy(order.id, order.price, order.size, order.clientChan)
				} else if order.orderType == inputSell {
					instrument.processSell(order.id, order.price, order.size, order.clientChan)
				} else {
					instrument.processCancel(order.id, order.clientChan)
				}
			case <- done:
				return
			}
		}
	}()
	return orderChan
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

type Instrument struct {
	name string
	buyBook  *LinkedList
	sellBook *LinkedList
}

func newInstrument(name string) *Instrument {
	return &Instrument{
		name: name, 
		buyBook:  newLinkedList(),
		sellBook: newLinkedList()}
}

// process buy
func (i *Instrument) processBuy(id uint32, price uint32, size uint32, clientChan chan struct {}) {
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
		// fmt.Fprintf(os.Stderr, "TRADE BUY %v with %v %v@%v : %v counter %v\n", id, bestSell.id, quantity, bestSell.price, GetCurrentTimestamp(), bestSell.counter)
		outputOrderExecuted(bestSell.id, id, bestSell.counter, bestSell.price, quantity, GetCurrentTimestamp())
		bestSell.counter++

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

	clientChan <- struct{}{}
}

func (i *Instrument) processSell(id uint32, price uint32, size uint32, clientChan chan struct{}) {
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
		// fmt.Fprintf(os.Stderr, "TRADE SELL %v with %v %v@%v : %v counter %v\n",id, bestBuy.id, quantity, bestBuy.price, GetCurrentTimestamp(), bestBuy.counter)
		outputOrderExecuted(bestBuy.id, id, bestBuy.counter, bestBuy.price, quantity, GetCurrentTimestamp())
		bestBuy.counter++

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

	clientChan <- struct{}{}
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

func (i *Instrument) addBuy(id uint32, price uint32, size uint32) {
	currentTime := GetCurrentTimestamp()
	// fmt.Fprintf(os.Stderr, "Buy %v added with %v@%v : %v\n", id, size, price, currentTime)
	var in = input {
		orderType: inputBuy,
		orderId: id, 
		price: price, 
		count: size,
		instrument: i.name}
	i.buyBook.insert(id, price, size, currentTime)
	outputOrderAdded(in, currentTime)
}

func (i *Instrument) addSell(id uint32, price uint32, size uint32) {
	currentTime := GetCurrentTimestamp()
	// fmt.Fprintf(os.Stderr, "Sell %v added with %v@%v : %v\n", id, size, price, currentTime)
	var in = input {
		orderType: inputSell,
		orderId: id, 
		price: price, 
		count: size,
		instrument: i.name}
	i.sellBook.insert(id, price, size, currentTime)
	outputOrderAdded(in, currentTime)
}


func (i *Instrument) processCancel(id uint32, clientChan chan struct{}) {
	in := input {orderId: id} 
	node := i.buyBook.getNodeById(id)
	if node != nil {
		i.buyBook.deleteNode(node)
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		// fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
		clientChan <- struct{}{}
		return
	}
	node2 := i.sellBook.getNodeById(id)
	if node2 != nil {
		i.sellBook.deleteNode(node2)
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		// fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
		clientChan <- struct{}{}
		return
	}
	outputOrderDeleted(in, false, GetCurrentTimestamp())
	clientChan <- struct{}{}
}