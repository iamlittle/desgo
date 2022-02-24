package main

import (
	"fmt"
	"log"
)

type OffshoreResource struct {
	Id, Index, State                int
	Created, TimeStamp, ServiceTime float64
	Customer                        *Customer
	PendingEventSet                 *PendingEventSet
	Business                        *Business
	Stats                           *Stats
}

func NewOffshoreResource(timestamp float64, pendingEventSet *PendingEventSet, customerServer *Business, stats *Stats) OffshoreResource {
	cpu := OffshoreResource{stats.generateEntityId(),
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

func (c *OffshoreResource) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.TimeStamp
}

func (c *OffshoreResource) Transition() bool {
	log.Println(fmt.Sprintf("[DEBUG] OffshoreResource %d finished with Customer %d at %f", c.Id, c.Customer.Id, c.TimeStamp))
	c.Customer.EndService(c.TimeStamp)
	c.Customer = nil
	c.Business.NotifyOffshoreResourceAvailable(c, c.TimeStamp)
	return true
}

func (c *OffshoreResource) BeginService(customer *Customer, timestamp float64) {
	c.Customer = customer
	log.Println(fmt.Sprintf("[DEBUG] OffshoreResource %d servicing Customer %d at %f", c.Id, customer.Id, timestamp))
	c.Customer.EndWait(timestamp)
	serviceTime := c.Stats.generateServiceTime()
	c.Stats.RecordOffshoreResourceServiceTime(serviceTime)
	c.Stats.RecordOffshoreResourceIdleTime(timestamp - c.TimeStamp)
	c.TimeStamp = timestamp + c.ServiceTime
	c.PendingEventSet.scheduleEvent(c)
}
