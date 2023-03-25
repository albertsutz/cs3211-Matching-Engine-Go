package main

import (
	"time"
	"fmt"
	"os"
// )
)

type Order struct {
	orderType inputType
	clientChan chan struct{} 
	id uint32
	price uint32
	size uint32
}

type Instrument struct {
	name string
	em *Mutex 
	buyBook  *LinkedList
	sellBook *LinkedList

	counter int
	currentType inputType
	turnMut *Mutex
	turnChan chan struct{}
}

func newInstrument(name string) *Instrument {
	var i = Instrument{
		name: name, 
		em: initMutex(),
		buyBook:  newLinkedList(),
		sellBook: newLinkedList(),
		counter: 0,
		currentType: 'X',
		turnMut: initMutex(),
		turnChan: make(chan struct{}, 1)}
	i.turnChan <- struct{}{}
	return &i
}

func instrumentFunc(done <-chan struct{}, name string) chan Order {
	orderChan := make(chan Order)
	go func() {
		instrument := newInstrument(name)
		for {
			select {
			case order:= <- orderChan:
				fmt.Fprintf(os.Stderr, "GOT ORDER\n")
				if order.orderType != instrument.currentType {
					<-instrument.turnChan
				}
				fmt.Fprintf(os.Stderr, "DOING ORDER\n")

				instrument.turnMut.lock()
				instrument.counter++
				instrument.currentType = order.orderType
				instrument.turnMut.unlock()

				if order.orderType == inputBuy {
					go instrument.processBuy(order.id, order.price, order.size, order.clientChan)
				} else if order.orderType == inputSell {
					go instrument.processSell(order.id, order.price, order.size, order.clientChan)
				} else {
					go instrument.processCancel(order.id, order.clientChan)
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

// process buy
func (i *Instrument) processBuy(id uint32, price uint32, size uint32, clientChan chan struct {}) {
	var bestSell *Node
	bestSell = i.getBestSell()

	if bestSell == nil || price < bestSell.price {
		i.addBuy(id, price, size)
	} else {
		i.executeBuy(id, price, size, clientChan)
	}

	i.turnMut.lock()
	defer i.turnMut.unlock()
	i.counter--
	if(i.counter == 0) {
		i.currentType = 'X'
		i.turnChan <- struct{}{}
	}	

	fmt.Fprintf(os.Stderr, "DONE ORDER\n")
	clientChan <- struct{}{}
}

func (i *Instrument) executeBuy(id uint32, price uint32, size uint32, clientChan chan struct {}) {
	// acquire execute lock
	i.em.lock()
	defer i.em.unlock() 

	var bestSell *Node
	for true {
		bestSell = i.getBestSell()

		if bestSell == nil || price < bestSell.price {
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
			i.sellBook.lock()
			i.sellBook.deleteNode(bestSell)
			i.sellBook.unlock()
		}
		if size == 0 {
			break
		}
	}
	if size > 0 {
		i.addBuy(id, price, size)
	}
}

func (i *Instrument) processSell(id uint32, price uint32, size uint32, clientChan chan struct{}) {
	bestBuy := i.getBestBuy()
	if bestBuy == nil || price > bestBuy.price {
		i.addSell(id, price, size)
	} else {
		i.executeSell(id, price, size, clientChan)
	}

	i.turnMut.lock()
	defer i.turnMut.unlock()
	i.counter--
	if(i.counter == 0) {
		i.currentType = 'X'
		i.turnChan <- struct{}{}
	}

	fmt.Fprintf(os.Stderr, "DONE ORDER\n")
	clientChan <- struct{}{}
}

func (i *Instrument) executeSell(id uint32, price uint32, size uint32, clientChan chan struct{}) {
	// acquire mutex here 
	i.em.lock()
	defer i.em.unlock()

	var bestBuy *Node
	for true {
		bestBuy = i.getBestBuy()
		if bestBuy == nil || price > bestBuy.price {
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
			i.buyBook.lock()
			defer i.buyBook.unlock()
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
	i.buyBook.lock()
	defer i.buyBook.unlock()

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
	i.sellBook.lock()
	defer i.sellBook.unlock()

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
	i.buyBook.lock()
	defer i.buyBook.unlock()

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
	i.sellBook.lock()
	defer i.sellBook.unlock()
	

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
	i.sellBook.lock()
	i.buyBook.lock()
	defer i.sellBook.unlock()
	defer i.buyBook.unlock()

	in := input {orderId: id} 
	node := i.buyBook.getNodeById(id)
	if node != nil {
		i.buyBook.deleteNode(node)
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		// fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
		i.turnMut.lock()
		defer i.turnMut.unlock()
		i.counter--
		if(i.counter == 0) {
			i.currentType = 'X'
			i.turnChan <- struct{}{}
		}
		fmt.Fprintf(os.Stderr, "DONE ORDER\n")
		clientChan <- struct{}{}
		return
	}
	node2 := i.sellBook.getNodeById(id)
	if node2 != nil {
		i.sellBook.deleteNode(node2)
		outputOrderDeleted(in, true, GetCurrentTimestamp())
		// fmt.Fprintf(os.Stderr, "Accepted Cancel with ID %v at %v", id, currentTime)
		i.turnMut.lock()
		defer i.turnMut.unlock()
		i.counter--
		if(i.counter == 0) {
			i.currentType = 'X'
			i.turnChan <- struct{}{}
		}
		fmt.Fprintf(os.Stderr, "DONE ORDER\n")
		clientChan <- struct{}{}
		return
	}
	outputOrderDeleted(in, false, GetCurrentTimestamp())

	i.turnMut.lock()
	defer i.turnMut.unlock()
	i.counter--
	if(i.counter == 0) {
		i.currentType = 'X'
		i.turnChan <- struct{}{}
	}
	fmt.Fprintf(os.Stderr, "DONE ORDER\n")
	clientChan <- struct{}{}
}
