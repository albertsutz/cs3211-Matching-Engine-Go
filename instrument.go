package main

import (
	"time"
	// "fmt"
	// "os"
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
	em chan struct{}
	// buyBook  *LinkedList
	// sellBook *LinkedList

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
	// getBestBuy, insert, delete, updateNode
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
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				bestBuy := getBestBuy(buyBook)
				if bestBuy != nil {
					output := InstrumentOutput {
						requestType: request.requestType,
						status: true,
						id: bestBuy.id,
						price: bestBuy.price,
						size: bestBuy.size,
						counter: bestBuy.counter}
					//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
					request.outputChan <- output
				} else {
					output := InstrumentOutput {requestType: request.requestType, status: false}
					//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
					request.outputChan <- output
				}
			} else if (request.requestType == "insert") {
				// sus
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				currentTime:= GetCurrentTimestamp()
				buyBook.insert(request.id, request.price, request.size, currentTime)
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true,
					time: currentTime}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
				request.outputChan <- output
			} else if (request.requestType == "delete") {
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				node := buyBook.getNodeById(request.id)
				if node != nil {
					buyBook.deleteNode(node)
				}
				output := InstrumentOutput {
					requestType: request.requestType,
					status: node != nil,
					time: GetCurrentTimestamp()}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
				request.outputChan <- output
			} else if (request.requestType == "update") {
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				node := buyBook.getNodeById(request.id)
				node.size = request.size
				node.counter++
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
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
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				bestSell := getBestSell(sellBook)
				if bestSell != nil {
					output := InstrumentOutput {
						requestType: request.requestType,
						status: true,
						id: bestSell.id,
						price: bestSell.price,
						size: bestSell.size,
						counter: bestSell.counter}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
				request.outputChan <- output
				} else {
					output := InstrumentOutput {requestType: request.requestType, status: false}
					//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
					request.outputChan <- output
				}
			} else if (request.requestType == "insert") {
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				// sus
				currentTime:= GetCurrentTimestamp()
				sellBook.insert(request.id, request.price, request.size, currentTime)
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true,
					time: currentTime}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
				request.outputChan <- output
			} else if (request.requestType == "delete") {
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				node := sellBook.getNodeById(request.id)
				if node != nil {
					sellBook.deleteNode(node)
				}
				output := InstrumentOutput {
					requestType: request.requestType,
					status: node != nil,
					time: GetCurrentTimestamp()}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
				request.outputChan <- output
			} else if (request.requestType == "update") {
				//fmt.Fprintf(os.Stderr, "Request: %+v\n", request)
				node := sellBook.getNodeById(request.id)
				node.size = request.size
				node.counter++
				output := InstrumentOutput {
					requestType: request.requestType,
					status: true}
				//fmt.Fprintf(os.Stderr, "Answer: %+v\n", output)
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
			// fmt.Fprintf(os.Stderr, "Get Request: %s\n", req.requestType)
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
					// fmt.Fprintf(os.Stderr, "Hit zero!\n")
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
					go instrument.processBuy(order.id, order.price, order.size, order.clientChan, inputBBH, inputSBH, counterChan)
				} else if order.orderType == inputSell {
					go instrument.processSell(order.id, order.price, order.size, order.clientChan, inputBBH, inputSBH, counterChan)
				} else {
					go instrument.processCancel(order.id, order.clientChan, inputBBH, inputSBH, counterChan)
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

func (i *Instrument)processBuy(id uint32, price uint32, size uint32, clientChan chan struct {}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest) {

	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	// ask for bestSell
	inputSBH <- InstrumentRequest {requestType: "getBestSell", outputChan: outputChan}
	output = <- outputChan

	if output.status == false || price < output.price {
		addBuy(i.name, id, price, size, inputBBH, inputSBH)
	} else {
		i.executeBuy(id, price, size, inputBBH, inputSBH)
	}

	resChan := make(chan CounterResponse)
	counterChan <- CounterRequest{requestType: "decrement", curType: 'B', outputChan: resChan}
	<-resChan

	clientChan <- struct{}{}
}

func (i *Instrument)executeBuy(id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest) {
	// acquire execute lock
	lockMutex(i.em)
	defer unlockMutex(i.em)

	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	for true {
		// bestSell = i.getBestSell()
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
		// bestSell.counter++

		if output.size == 0 {
			// i.sellBook.deleteNode(bestSell)
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
		addBuy(i.name, id, price, size, inputBBH, inputSBH)
	}
}

func (i *Instrument) processSell(id uint32, price uint32, size uint32, clientChan chan struct{}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest) {
	// bestBuy := i.getBestBuy()
	outputChan:= make(chan InstrumentOutput, 1)
	var output InstrumentOutput
	// ask for bestBuy
	inputBBH <- InstrumentRequest {requestType: "getBestBuy", outputChan: outputChan}
	output = <- outputChan

	if output.status == false || price > output.price {
		addSell(i.name, id, price, size, inputBBH, inputSBH)
	} else {
		i.executeSell(id, price, size, inputBBH, inputSBH)
	}

	resChan := make(chan CounterResponse)
	counterChan <- CounterRequest{requestType: "decrement", curType: 'S', outputChan: resChan}
	<-resChan
	clientChan <- struct{}{}
}

func (i *Instrument)executeSell(id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest) {
	// acquire mutex here 
	lockMutex(i.em)
	defer unlockMutex(i.em)

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
		addSell(i.name, id, price, size, inputBBH, inputSBH)
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
	// i.buyBook.insert(id, price, size, currentTime)

	outputOrderAdded(in, output.time)
}

func addSell(name string, id uint32, price uint32, size uint32, inputBBH chan InstrumentRequest, inputSBH chan InstrumentRequest) {
	// i.sellBook.lock()
	// defer i.sellBook.unlock()

	var in = input {
		orderType: inputSell,
		orderId: id, 
		price: price, 
		count: size,
		instrument: name}
	outputChan:= make(chan InstrumentOutput, 1)
	inputSBH <- InstrumentRequest {requestType: "insert", id: id, price: price, size: size, outputChan: outputChan}
	output := <- outputChan
	// i.sellBook.insert(id, price, size, currentTime)
	outputOrderAdded(in, output.time)
}


func (i *Instrument) processCancel(id uint32, clientChan chan struct{}, inputBBH chan InstrumentRequest, 
	inputSBH chan InstrumentRequest, counterChan chan CounterRequest) {
	// i.sellBook.lock()
	// defer i.sellBook.unlock()
	// i.buyBook.lock()
	// defer i.buyBook.unlock()
	var output InstrumentOutput
	in := input {orderId: id} 
	outputChan := make(chan InstrumentOutput, 1)

	// node := i.buyBook.getNodeById(id)
	inputBBH <- InstrumentRequest{requestType: "delete", id: id, outputChan: outputChan}
	output = <- outputChan
	if output.status != false {
		// i.buyBook.deleteNode(node)
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
		// i.sellBook.deleteNode(node2)
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
	
	//fmt.Fprintf(os.Stderr, "%v\n", clientChan)
	clientChan <- struct{}{}
}
