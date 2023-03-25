package main

import (
)

type PairInsOrd struct {
	instrument string 
	orderType inputType
}

type OrderBook struct{
	mut *Mutex
	m_instrument_map map[string](chan Order)
	m_id_map map[uint32]PairInsOrd
	done <-chan struct{}
}

func newOrderBook(done <-chan struct{}) *OrderBook{
	return &OrderBook {
		mut: initMutex(),
		m_instrument_map: make(map[string](chan Order)),
		m_id_map: make(map[uint32]PairInsOrd),
		done: done}
}

func (o *OrderBook) process_order(orderType inputType, id uint32, instr string, price uint32, size uint32, clientChan chan struct{}) {
	o.mut.lock()
	if !o.is_exist_instr(instr) {
		o.m_instrument_map[instr] = instrumentFunc(o.done, instr)
	}
	o.m_id_map[id] = PairInsOrd{instrument: instr, orderType: orderType}
	o.mut.unlock()

	o.m_instrument_map[instr] <- Order{orderType: orderType, id: id, price: price, size: size, clientChan: clientChan}
}

func (o *OrderBook) process_cancel(id uint32, clientChan chan struct{}) {
	o.mut.lock()
    pair, ok := o.m_id_map[id] 
	o.mut.unlock()
	if !ok {
		// fmt.Fprintf(os.Stderr, "Declined Cancel with ID %v at %v", id, GetCurrentTimestamp())
		outputOrderDeleted(input{orderId: id}, false, GetCurrentTimestamp())
		clientChan <- struct{}{}
		return  
	}

	o.m_instrument_map[pair.instrument] <- Order{orderType: inputCancel, id: id, clientChan: clientChan}
}

func (o *OrderBook) is_exist_instr(instr string) bool {
	_, ok := o.m_instrument_map[instr]
	return ok
}

func (o *OrderBook) is_exist_id(id uint32) bool {
	_, ok := o.m_id_map[id]
	return ok
}