package main

import "C"
import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	// "time"
)

type Request struct {
	in input
	clientChan chan struct{}
}

type Engine struct{
	reqChan chan Request
}

func (e *Engine) processRequest(done chan struct{}) {
	ob := newOrderBook(done) 

	for {
		req := <-e.reqChan
		in := req.in 
		switch in.orderType {
		case inputCancel:
			ob.process_cancel(in.orderId, req.clientChan)
		default:
			ob.process_order(in.orderType, in.orderId, in.instrument, in.price, in.count, req.clientChan)
		}
	}
}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go handleConn(conn, e.reqChan)
}

func handleConn(conn net.Conn, reqChan chan Request) {
	defer conn.Close()
	clientChan := make(chan struct{})
	for {
		in, err := readInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		request := Request{in: in, clientChan: clientChan}
		reqChan <- request
		<- clientChan 
	}
}
