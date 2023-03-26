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

type Instrument struct {
	name string
	em chan struct{}

	counter int
	currentType inputType
	turnMut chan struct{}
	turnChan chan struct{}
}

func newInstrument(name string) *Instrument {
	var i = Instrument{
		name: name, 
		em: initMutex(),
		counter: 0,
		currentType: 'X',
		turnMut: initMutex(),
		turnChan: make(chan struct{}, 1)}
	i.turnChan <- struct{}{}
	return &i
}

type InstrumentRequest struct {
	requestType string
	id uint32
	price uint32
	size uint32
	outputChan chan InstrumentOutput
}

type InstrumentOutput struct {
	requestType string
	status bool
	id uint32
	price uint32
	size uint32
	time int64
	counter uint32
}

func handleBuyBook(incomingChan chan InstrumentRequest, done <-chan struct{}) {
	buyBook := newLinkedList()
	for {
		select {
		case request, _ := <- incomingChan:
			if (request.requestType == "getBestBuy") {
				bestBuy := getBestBuy(buyBook)
				if bestBuy != nil {
					output := InstrumentOutput {
						requestType: request.requestType,
						status: true,
						id: bestBuy.id,
						price: bestBuy.price,
						size: bestBuy.size,
						counter: bestBuy.counter}
					request.outputChan <- output
				} else {
					output := InstrumentOutput {requestType: request.requestType, status: false}
					request.outputChan <- output
				}
			} else if (request.requestType == "insert") {
				currentTime:= GetCurrentTimestamp()
				buyBook.insert(request.id, request.price, request.size, currentTime)
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true,
					time: currentTime}

				request.outputChan <- output
			} else if (request.requestType == "delete") {
				node := buyBook.getNodeById(request.id)
				if node != nil {
					buyBook.deleteNode(node)
				}
				output := InstrumentOutput {
					requestType: request.requestType,
					status: node != nil,
					time: GetCurrentTimestamp()}
				request.outputChan <- output
			} else if (request.requestType == "update") {
				node := buyBook.getNodeById(request.id)
				node.size = request.size
				node.counter++
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true}
				request.outputChan <- output
			}
		case <- done:
			return
		}
		
	}
}

func handleSellBook(incomingChan chan InstrumentRequest, done <-chan struct{}) {
	sellBook := newLinkedList()
	for {
		select {
		case request, _ := <- incomingChan:
			if (request.requestType == "getBestSell") {
				bestSell := getBestSell(sellBook)
				if bestSell != nil {
					output := InstrumentOutput {
						requestType: request.requestType,
						status: true,
						id: bestSell.id,
						price: bestSell.price,
						size: bestSell.size,
						counter: bestSell.counter}
				request.outputChan <- output
				} else {
					output := InstrumentOutput {requestType: request.requestType, status: false}
					request.outputChan <- output
				}
			} else if (request.requestType == "insert") {
				currentTime:= GetCurrentTimestamp()
				sellBook.insert(request.id, request.price, request.size, currentTime)
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true,
					time: currentTime}
				request.outputChan <- output
			} else if (request.requestType == "delete") {
				node := sellBook.getNodeById(request.id)
				if node != nil {
					sellBook.deleteNode(node)
				}
				output := InstrumentOutput {
					requestType: request.requestType,
					status: node != nil,
					time: GetCurrentTimestamp()}
				request.outputChan <- output
			} else if (request.requestType == "update") {
				node := sellBook.getNodeById(request.id)
				node.size = request.size
				node.counter++
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true}
				request.outputChan <- output
			}
		case <-done:
			return
		}
	}
}

func getBestBuy(buyBook *LinkedList) *Node {
	if buyBook.getLength() == 0 {
		return nil
	}
	var answer *Node = buyBook.getHead()
	// traverse through the buyBook
	var currNode *Node = buyBook.getHead()
	for a := 0; a < buyBook.getLength(); a++ {
		if currNode.price > answer.price || currNode.price == answer.price && currNode.time < answer.time {
			answer = currNode
		}
		currNode = currNode.next
	}
	return answer
}

func getBestSell(sellBook *LinkedList) *Node {
	if sellBook.getLength() == 0 {
		return nil
	}
	var answer *Node = sellBook.getHead()
	// traverse through the buyBook
	var currNode *Node = sellBook.getHead()
	for a := 0; a < sellBook.getLength(); a++ {
		if currNode.price < answer.price || currNode.price == answer.price && currNode.time < answer.time {
			answer = currNode
		}
		currNode = currNode.next
	}
	return answer
}

type CounterRequest struct {
	// increment, decrement, getType
	requestType string
	curType inputType 
	outputChan chan CounterResponse
}

type CounterResponse struct {
	curType inputType 
}

func counterHandler(incomingChan chan CounterRequest, turnChan chan struct{}, done <-chan struct{}) {
	counter := 0
	var currentType inputType = 'X'
	
	for {
		select {
		case req := <- incomingChan:
			switch req.requestType {
			case "increment":
				select {
				case <-turnChan:
				default:
				}
				counter++
				currentType = req.curType 
				req.outputChan <- CounterResponse{}
			case "decrement":
				counter--
				if counter == 0 {
					currentType = 'X'
					turnChan <- struct{}{}
				}
				req.outputChan <- CounterResponse{}
			case "getType":
				req.outputChan <- CounterResponse{curType: currentType}
			}
		case <- done:
			return 
		}
	}
}

func instrumentFunc(done <-chan struct{}, name string) chan Order {
	orderChan := make(chan Order)
	go func() {
		instrument := newInstrument(name)
		inputBBH := make(chan InstrumentRequest)
		inputSBH := make(chan InstrumentRequest)
		counterChan := make(chan CounterRequest)
		go handleBuyBook(inputBBH, done)
		go handleSellBook(inputSBH, done)
		go counterHandler(counterChan,instrument.turnChan, done)
		for {
			select {
			case order:= <- orderChan:
				resChan := make(chan CounterResponse)
				counterChan <- CounterRequest{requestType: "getType", outputChan: resChan}
				typ := <-resChan

				if order.orderType != typ.curType {
					<-instrument.turnChan
				}

				resChan = make(chan CounterResponse)
				counterChan <- CounterRequest{requestType: "increment", curType: order.orderType, outputChan: resChan}
				<-resChan

				if order.orderType == inputBuy {
					go processBuy(instrument.name, order.id, order.price, order.size, order.clientChan, inputBBH, inputSBH, counterChan, instrument.em)
				} else if order.orderType == inputSell {
					go processSell(instrument.name, order.id, order.price, order.size, order.clientChan, inputBBH, inputSBH, counterChan, instrument.em)
				} else {
					go processCancel(order.id, order.clientChan, inputBBH, inputSBH, counterChan)
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

func processBuy(name string, id uint32, price uint32, size uint32, clientChan chan struct {}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest, em chan struct{}) {

	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	// ask for bestSell
	inputSBH <- InstrumentRequest {requestType: "getBestSell", outputChan: outputChan}
	output = <- outputChan

	if output.status == false || price < output.price {
		addBuy(name, id, price, size, inputBBH, inputSBH)
	} else {
		executeBuy(name, id, price, size, inputBBH, inputSBH, em)
	}

	resChan := make(chan CounterResponse)
	counterChan <- CounterRequest{requestType: "decrement", curType: 'B', outputChan: resChan}
	<-resChan

	clientChan <- struct{}{}
}

func executeBuy(name string, id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest,
	em chan struct{}) {
	// acquire execute lock
	lockMutex(em)
	defer unlockMutex(em)

	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	for true {
		inputSBH <- InstrumentRequest {requestType: "getBestSell", outputChan: outputChan}
		output = <- outputChan

		if output.status == false || output.price < output.price {
			break
		} 

		//quantity traded
		var quantity = min(size, output.size)
		size -= quantity
		output.size -= quantity
		outputOrderExecuted(output.id, id, output.counter, output.price, quantity, GetCurrentTimestamp())

		if output.size == 0 {
			inputSBH <- InstrumentRequest {requestType: "delete", id: output.id, outputChan: outputChan}
			output = <- outputChan
		} else {
			//update
			inputSBH <- InstrumentRequest {requestType: "update", id: output.id, size: output.size, outputChan: outputChan}
			output = <- outputChan
		}
		if size == 0 {
			break
		}
	}
	if size > 0 {
		addBuy(name, id, price, size, inputBBH, inputSBH)
	}
}

func processSell(name string, id uint32, price uint32, size uint32, clientChan chan struct{}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest, em chan struct{}) {
	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	// ask for bestBuy
	inputBBH <- InstrumentRequest {requestType: "getBestBuy", outputChan: outputChan}
	output = <- outputChan

	if output.status == false || price > output.price {
		addSell(name, id, price, size, inputBBH, inputSBH)
	} else {
		executeSell(name, id, price, size, inputBBH, inputSBH, em)
	}

	resChan := make(chan CounterResponse)
	counterChan <- CounterRequest{requestType: "decrement", curType: 'S', outputChan: resChan}
	<-resChan
	clientChan <- struct{}{}
}

func executeSell(name string, id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest,
	em chan struct{}) {
	// acquire mutex here 
	lockMutex(em)
	defer unlockMutex(em)

	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	for true {
		inputBBH <- InstrumentRequest {requestType: "getBestBuy", outputChan: outputChan}
		output = <- outputChan
		if output.status == false || price > output.price {
			break
		}
		//quantity traded
		var quantity = min(size, output.size)
		size -= quantity
		output.size -= quantity
		outputOrderExecuted(output.id, id, output.counter, output.price, quantity, GetCurrentTimestamp())

		//if bestBuy.size == 0 -> remove from instrument
		if output.size == 0 {
			inputBBH <- InstrumentRequest {requestType: "delete", id: output.id, outputChan: outputChan}
			output = <- outputChan
		} else {
			inputBBH <- InstrumentRequest {requestType: "update", id: output.id, size: output.size, outputChan: outputChan}
			output = <- outputChan			
		}
		if size == 0 {
			break
		}
	}
	if size > 0 {
		addSell(name, id, price, size, inputBBH, inputSBH)
	}
}

func addBuy(name string, id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest) {
	var in = input {
		orderType: inputBuy,
		orderId: id, 
		price: price, 
		count: size,
		instrument: name}

	outputChan:= make(chan InstrumentOutput, 1)
	inputBBH <- InstrumentRequest {requestType: "insert", id: id, price: price, size: size, outputChan: outputChan}
	output := <- outputChan

	outputOrderAdded(in, output.time)
}

func addSell(name string, id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest) {

	var in = input {
		orderType: inputSell,
		orderId: id, 
		price: price, 
		count: size,
		instrument: name}
	outputChan:= make(chan InstrumentOutput, 1)
	inputSBH <- InstrumentRequest {requestType: "insert", id: id, price: price, size: size, outputChan: outputChan}
	output := <- outputChan
	outputOrderAdded(in, output.time)
}


func processCancel(id uint32, clientChan chan struct{}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest) {
	var output InstrumentOutput
	in := input {orderId: id} 
	outputChan := make(chan InstrumentOutput, 1)

	inputBBH <- InstrumentRequest{requestType: "delete", id: id, outputChan: outputChan}
	output = <- outputChan
	if output.status != false {
		outputOrderDeleted(in, true, output.time)
		resChan := make(chan CounterResponse)
		counterChan <- CounterRequest{requestType: "decrement", curType: 'C', outputChan: resChan}
		<-resChan
		clientChan <- struct{}{}
		return
	}
	inputSBH <- InstrumentRequest{requestType: "delete", id: id, outputChan: outputChan}
	output = <- outputChan
	if output.status != false {
		outputOrderDeleted(in, true, output.time)
		resChan := make(chan CounterResponse)
		counterChan <- CounterRequest{requestType: "decrement", curType: 'C', outputChan: resChan}
		<-resChan
		clientChan <- struct{}{}
		return
	}
	outputOrderDeleted(in, false, GetCurrentTimestamp())
	resChan := make(chan CounterResponse)
	counterChan <- CounterRequest{requestType: "decrement", curType: 'C', outputChan: resChan}
	<-resChan
	
	clientChan <- struct{}{}
}
