package main

import (
	"log"
	"fmt"
)

/*************
** Business **
**************/
type Business struct{
	Cashiers []*Cashier
	Customers []*Customer
}

func NewBusiness() Business{
	return Business{make([]*Cashier, 0), make([]*Customer, 0)}
}

func (c *Business) Checkout(customer *Customer, timestamp float64){
	if len(c.Cashiers) == 0 {
		log.Println(fmt.Sprintf("[DEBUG] No cashiers available for Customer %d at %f", customer.Id, timestamp))
		c.Customers = append(c.Customers, customer)
	} else {
		cashier := c.Cashiers[0]
		c.Cashiers =  c.Cashiers[1:]
		cashier.BeginService(customer, timestamp)
	}
}

func (c *Business) NotifyCashierAvailable(cashier *Cashier, timestamp float64) {
	cashier.TimeStamp = timestamp
	if len(c.Customers) == 0 {
		c.Cashiers = append(c.Cashiers, cashier)
	}else{
		customer := c.Customers[0]
		c.Customers =  c.Customers[1:]
		cashier.BeginService(customer, timestamp)
	}
}