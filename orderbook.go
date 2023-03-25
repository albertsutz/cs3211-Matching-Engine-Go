package main

import (

)

type PairInsOrd struct {
	id string
	orderType string
}

type OrderBook struct{
	m_instrument_map map[string]*Instrument
	m_id_map map[int]PairInsOrd
}

func (o *OrderBook) process_order(orderType string, id int, instr string, price int, size int) {
	// if !is_exist_instr {

	// }
	// var instrument_object


    // if (!is_exist_instr(order.instrument)) {
    //     m_instrument_map.emplace(std::piecewise_construct, std::forward_as_tuple(order.instrument), std::forward_as_tuple());
    // }
    // auto& instruction_object = m_instrument_map.at(order.instrument);
    // m_id_map[order.order_id] = std::make_pair(order.instrument, order.order_type);

    // auto result = instruction_object.process_order(order); 
}

func (o *OrderBook) process_cancel(id int) {
    // pair_name_type, ok = o.m_id_map[id] 

    // auto& instruction_object = m_instrument_map.at(pair_name_type.first);

    // auto res = instruction_object.process_cancel(order, pair_name_type.second); 
}

func (o *OrderBook) is_exist_instr(instr string) bool {
	_, ok := o.m_instrument_map[instr]
	return ok
}

func (o *OrderBook) is_exist_id(id int) bool {
	_, ok := o.m_id_map[id]
	return ok
}