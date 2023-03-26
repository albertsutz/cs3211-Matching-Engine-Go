package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func handleSigs(cancel func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <socket path>\n", os.Args[0])
		return
	}

	socketPath := os.Args[1]
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatal("remove existing sock error: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		handleSigs(cancel)
	}()

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("listen error: ", err)
	}
	go func() {
		<-ctx.Done()
		if err := l.Close(); err != nil {
			log.Fatal("close listener error: ", err)
		}
	}()

	

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(done)
	}()

	var e = Engine{reqChan: make(chan Request)}
	go e.processRequest(done)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error: ", err)
		}

		e.accept(ctx, conn)
	}
}


// package main

// import(
// 	// "fmt"
// 	// "os"
// 	// "strconv"
// 	// "time"
// )

// func main() {
// 	testInstrument()
// }

// func testInstrument() {
// 	done := make(chan struct{})
// 	dummy := instrumentFunc(done, "testing")
// 	clientChan := make(chan struct{})
// 	dummy <- Order{orderType: 'B', id: 1, price: 100, size: 100, clientChan: clientChan}
// 	<-clientChan
	// dummy <- Order{orderType: 'S', id: 2, price: 200, size: 100, clientChan: clientChan}
	// <-clientChan

	// dummy <- Order{orderType: 'B', id: 3, price: 200, size: 50, clientChan: clientChan}
	// <-clientChan

	// dummy <- Order{orderType: 'B', id: 4, price: 200, size: 25, clientChan: clientChan}
	// <-clientChan

	// dummy <- Order{orderType: 'B', id: 5, price: 200, size: 25, clientChan: clientChan}
	// <-clientChan

	// fmt.Fprintf(os.Stderr, "OY1\n")
	// fmt.Fprintf(os.Stderr, "%v\n", clientChan)
	// dummy <- Order{orderType: 'C', id: 1, clientChan: clientChan}
	// <-clientChan
	// fmt.Fprintf(os.Stderr, "OY2\n")


	// 	var ob *OrderBook = newOrderBook()
// 	// orderType string, id int, instr string, price int, size int
// 	ob.process_order(inputBuy, 301, "ABC", 3250, 500)
// 	// ob.process_order("S", 302, "ABC", 3250, 500)
// 	// ob.process_cancel(301)
// 	ob.process_order(inputSell, 302, "ABC", 3200, 1)
// 	ob.process_order(inputSell, 303, "ABC", 3200, 5)
// 	ob.process_order(inputSell, 304, "ABC", 3200, 5)
// 	time.Sleep(8 * time.Second)
// 	// var instrument *Instrument = newInstrument()
// 	// instrument.processBuy(301, 3250, 500)
// 	// instrument.processSell(302, 3200, 1)
// 	// instrument.processSell(303, 3200, 5)
// 	// instrument.processSell(304, 3200, 5)
// 	// instrument.buyBook.printValues()
// 	// instrument.sellBook.printValues()
// }