package main

import (
	"log"
	"fmt"
)

/*************
** Customer **
0 - enter store
1 - shop
2 - checkout
3 - wait
4 - get served
5 - leave
*************/

type Customer struct {
    Id, Index, State int
	Created, Timestamp, EnterQueue, Lifespan, WaitTime, ShopTime float64
	PendingEventSet *PendingEventSet
	Business *Business
	Stats *Stats
}

func NewCustomer(timestamp float64,
				pendingEventSet *PendingEventSet,
				business *Business,
				stats *Stats) *Customer{
	shopTime := stats.generateShopTime()
	stats.RecordShopTime(shopTime)
	customer := &Customer{ stats.generateEntityId(),
	    0,
	    0,
		timestamp,
		timestamp,
		0,
		0,
		0,
		shopTime,
		pendingEventSet,
		business,
		stats,
	}
	return customer
}

func (c *Customer) EventInfo() (int, int, float64){
	return c.Id, c.Index, c.Timestamp
}

func (c *Customer) Transition() bool {
	switch c.State {
	case 0:
		c.State++
		log.Println(fmt.Sprintf("[DEBUG] Customer %d entered store at %f", c.Id, c.Timestamp))
		c.Timestamp += c.ShopTime
		c.PendingEventSet.scheduleEvent(c)
	case 1:
		c.EnterQueue = c.Timestamp
		log.Println(fmt.Sprintf("[DEBUG] Customer %d finished shopping at %f", c.Id, c.Timestamp))
		c.Business.Checkout(c, c.Timestamp)
		c.State++
	}
	return true
}

func (c *Customer) EndWait(timestamp float64){
	c.WaitTime = timestamp - c.EnterQueue
	log.Println(fmt.Sprintf("[DEBUG] Customer %d waited in line for %f", c.Id, c.WaitTime))
	c.Stats.RecordWaitTime(c.WaitTime)
	c.State++
}

func (c *Customer) EndService(timestamp float64){
	c.Lifespan = timestamp - c.Created
	log.Println(fmt.Sprintf("[DEBUG] Customer %d was at business for %f", c.Id, c.Lifespan))
	c.State++
	c.Stats.CompletedJobCount++
}

/***********************
** Customer Generator **
************************/
type CustomerGenerator struct {
	Id, Index int
	Created, Timestamp float64
	PendingEventSet *PendingEventSet
	Business *Business
	Stats *Stats
}

func NewCustomerGenerator(timestamp float64,
						  pendingEventSet *PendingEventSet,
	  					  business *Business,
						  stats *Stats) CustomerGenerator{
	customerGen := CustomerGenerator{ stats.generateEntityId(),
		0,
		timestamp,
		timestamp,
		pendingEventSet,
		business,
		stats,
	}
	return customerGen
}

func (c *CustomerGenerator) EventInfo() (int, int, float64){
	return c.Id, c.Index, c.Timestamp
}

func (c *CustomerGenerator) Transition() bool {
	customer := NewCustomer(c.Stats.generateEntryTime(), c.PendingEventSet, c.Business, c.Stats)
	c.PendingEventSet.scheduleEvent(customer)
	log.Println(fmt.Sprintf("[DEBUG] Customer %d generated at %f", customer.Id, c.Timestamp))
	return false
}