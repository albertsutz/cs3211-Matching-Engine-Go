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

type Engine struct{
	ob *OrderBook 
}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go e.handleConn(conn)
}

func (e *Engine) handleConn(conn net.Conn) {
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
		switch in.orderType {
		case inputCancel:
			e.ob.process_cancel(in.orderId, clientChan)
			<- clientChan 
		default:
			e.ob.process_order(in.orderType, in.orderId, in.instrument, in.price, in.count, clientChan)
			<- clientChan 
		}
	}
}
