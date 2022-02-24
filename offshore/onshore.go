package main

import (
	"fmt"
	"log"
)

type OnshoreResource struct {
	Id, Index, State                int
	Created, TimeStamp, ServiceTime float64
	Customer                        *Customer
	PendingEventSet                 *PendingEventSet
	Business                        *Business
	Stats                           *Stats
}

func NewOnshoreResource(timestamp float64, pendingEventSet *PendingEventSet, customerServer *Business, stats *Stats) OnshoreResource {
	cpu := OnshoreResource{stats.generateEntityId(),
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

func (c *OnshoreResource) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.TimeStamp
}

func (c *OnshoreResource) Transition() bool {
	log.Println(fmt.Sprintf("[DEBUG] OnshoreResource %d finished with Customer %d at %f", c.Id, c.Customer.Id, c.TimeStamp))
	c.Customer.EndService(c.TimeStamp)
	c.Customer = nil
	c.Business.NotifyOnshoreResourceAvailable(c, c.TimeStamp)
	return true
}

func (c *OnshoreResource) BeginService(customer *Customer, timestamp float64) {
	c.Customer = customer
	log.Println(fmt.Sprintf("[DEBUG] OnshoreResource %d servicing Customer %d at %f", c.Id, customer.Id, timestamp))
	c.Customer.EndWait(timestamp)
	serviceTime := c.Stats.generateServiceTime()
	c.Stats.RecordOnshoreResourceServiceTime(serviceTime)
	c.Stats.RecordOnshoreResourceIdleTime(timestamp - c.TimeStamp)
	c.TimeStamp = timestamp + c.ServiceTime
	c.PendingEventSet.scheduleEvent(c)
}
