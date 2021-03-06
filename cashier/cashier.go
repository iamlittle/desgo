package main

import (
	"log"
	"fmt"
)

type Cashier struct{
	Id, Index, State int
	Created, TimeStamp, ServiceTime float64
	Customer *Customer
	PendingEventSet *PendingEventSet
	Business *Business
	Stats *Stats
}

func NewCashier(timestamp float64, pendingEventSet *PendingEventSet, customerServer *Business, stats *Stats) Cashier{
	cpu := Cashier{ stats.generateEntityId(),
		0,
		0,
		timestamp,
		timestamp,
		stats.generateServiceTime(),
		nil,
		pendingEventSet,
		customerServer,
		stats,
	}
	return cpu
}

func (c *Cashier) EventInfo() (int, int, float64){
	return c.Id, c.Index, c.TimeStamp
}

func (c *Cashier) Transition() bool{
	log.Println(fmt.Sprintf("[DEBUG] Cashier %d finished with Customer %d at %f", c.Id, c.Customer.Id, c.TimeStamp))
	c.Customer.EndService(c.TimeStamp)
	c.Customer = nil
	c.Business.NotifyCashierAvailable(c, c.TimeStamp)
	return true
}

func (c *Cashier) BeginService(customer *Customer, timestamp float64){
	c.Customer = customer
	log.Println(fmt.Sprintf("[DEBUG] Cashier %d servicing Customer %d at %f", c.Id, customer.Id, timestamp))
	c.Customer.EndWait(timestamp)
	serviceTime := c.Stats.generateServiceTime()
	c.Stats.RecordCashierServiceTime(serviceTime)
	c.Stats.RecordCashierIdleTime(timestamp - c.TimeStamp)
	c.TimeStamp = timestamp + c.ServiceTime
	c.PendingEventSet.scheduleEvent(c)
}

