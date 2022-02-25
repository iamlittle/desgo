package main

import (
	"fmt"
	"log"
)

type Resource struct {
	Id, Index, State   int
	Type               ResourceType
	Created, TimeStamp float64
	Component          *Component
	PendingEventSet    *PendingEventSet
	Migration          *Migration
	Stats              *Stats
}

func NewResource(timestamp float64, _type ResourceType, pendingEventSet *PendingEventSet, migration *Migration, stats *Stats) Resource {
	cpu := Resource{
		stats.generateEntityId(), //Id
		0,                        //Index
		0,                        //State
		_type,                    //Type
		timestamp,                //Created
		timestamp,                //TimeStamp
		nil,                      //Component
		pendingEventSet,          //PendingEventSet
		migration,                //Migration
		stats,                    //Stats
	}
	return cpu
}

func (c *Resource) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.TimeStamp
}

func (c *Resource) Transition() bool {
	log.Println(fmt.Sprintf("[DEBUG] Resource %d finished with Component %d at %f", c.Id, c.Component.Id, c.TimeStamp))
	c.Component.EndService(c.TimeStamp)
	c.Component = nil
	c.Migration.NotifyResourceAvailable(c, c.TimeStamp)
	return true
}

func (c *Resource) Process(component *Component, timestamp float64) {
	c.Component = component
	c.TimeStamp = timestamp
	log.Println(fmt.Sprintf("[DEBUG] Resource %d processing Component %d at %f, state %d", c.Id, component.Id, timestamp, component.State))
	switch component.State {
	case 1:
		component.Reviewed = timestamp + component.ReviewDuration
		component.Timestamp = component.Reviewed
		c.TimeStamp = component.Reviewed
	case 2:
		component.Converted = timestamp + component.ConvertDuration
		component.Timestamp = component.Converted
		c.TimeStamp = component.Converted
	case 4:
		component.UnitTested = timestamp + component.UnitTestDuration
		component.Timestamp = component.UnitTested
		c.TimeStamp = component.UnitTested
	case 6:
		component.Validated = timestamp + component.ValidateDuration
		component.Timestamp = component.Validated
		c.TimeStamp = component.Validated
	case 7:
		component.Cutover = timestamp + component.CutoverDuration
		component.Timestamp = component.Cutover
		c.TimeStamp = component.Cutover
	}
	c.PendingEventSet.scheduleEvent(c)
}
