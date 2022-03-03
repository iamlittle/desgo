package main

import (
	"fmt"
	"log"
	"math"
)

type Resource struct {
	Id, Index, State   int
	Type               ResourceType
	Created, TimeStamp float64
	Component          *Component
	TimeOff            bool
	PendingEventSet    *PendingEventSet
	Migration          *Migration
	Stats              *Stats
}

func NewResource(timestamp float64, _type ResourceType, time_off bool, pendingEventSet *PendingEventSet, migration *Migration, stats *Stats) Resource {
	cpu := Resource{
		stats.generateEntityId(), //Id
		0,                        //Index
		0,                        //State
		_type,                    //Type
		timestamp,                //Created
		timestamp,                //TimeStamp
		nil,                      //Component
		time_off,                 //TimeOff
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
	log.Println(fmt.Sprintf("[DEBUG] Resource %d finished with Component %d at %s", c.Id, c.Component.Id, c.formatTime(c.TimeStamp)))
	c.Component = nil
	c.Migration.NotifyResourceAvailable(c, c.TimeStamp)
	return true
}

func (c *Resource) Process(component *Component, timestamp float64) {
	c.Component = component
	c.TimeStamp = timestamp
	switch component.State {
	case 1:
		if c.TimeOff {
			component.Reviewed = c.calculateDurationIncludingTimeOff(component.ReviewDuration, timestamp)
		} else {
			component.Reviewed = timestamp + component.ReviewDuration
		}
		component.Timestamp = component.Reviewed
		c.TimeStamp = component.Reviewed
		c.Stats.RecordComponentReviewTime(component)
		log.Println(fmt.Sprintf("[DEBUG] Component %d artifacts reviewed at %s", component.Id, c.formatTime(component.Timestamp)))
	case 2:
		if c.TimeOff {
			component.Converted = c.calculateDurationIncludingTimeOff(component.ConvertDuration, timestamp)
		} else {
			component.Converted = timestamp + component.ConvertDuration
		}
		component.Timestamp = component.Converted
		c.TimeStamp = component.Converted
		c.Stats.RecordComponentRoughConversionTime(component)
		log.Println(fmt.Sprintf("[DEBUG] Component %d rough conversion %s", component.Id, c.formatTime(component.Timestamp)))
	case 4:
		if c.TimeOff {
			component.UnitTested = c.calculateDurationIncludingTimeOff(component.UnitTestDuration, timestamp)
		} else {
			component.UnitTested = timestamp + component.UnitTestDuration
		}
		component.UnitTested = timestamp + component.UnitTestDuration
		component.Timestamp = component.UnitTested
		c.TimeStamp = component.UnitTested
		c.Stats.RecordComponentUnitTestTime(component)
		log.Println(fmt.Sprintf("[DEBUG] Component %d unit tested at %s", component.Id, c.formatTime(component.Timestamp)))
	case 6:
		if c.TimeOff {
			component.Validated = c.calculateDurationIncludingTimeOff(component.ValidateDuration, timestamp)
		} else {
			component.Validated = timestamp + component.ValidateDuration
		}
		component.Validated = timestamp + component.ValidateDuration
		component.Timestamp = component.Validated
		c.TimeStamp = component.Validated
		c.Stats.RecordComponentValidateTime(component)
		log.Println(fmt.Sprintf("[DEBUG] Component %d executed and validated at %s", component.Id, c.formatTime(component.Timestamp)))
	case 7:
		if c.TimeOff {
			component.Cutover = c.calculateDurationIncludingTimeOff(component.CutoverDuration, timestamp)
		} else {
			component.Cutover = timestamp + component.CutoverDuration
		}
		component.Cutover = timestamp + component.CutoverDuration
		component.Timestamp = component.Cutover
		c.TimeStamp = component.Cutover
		c.Stats.RecordComponentCutoverTime(component)
		log.Println(fmt.Sprintf("[DEBUG] Component %d production cutover at %s", component.Id, c.formatTime(component.Timestamp)))
	}
	c.PendingEventSet.scheduleEvent(c)
}

func (c *Resource) calculateDurationIncludingTimeOff(duration float64, timestamp float64) float64 {
	remainingDuration := duration
	currentDayTime := math.Mod(timestamp, 24)
	currentWeekDay := math.Mod(timestamp, 24*7)
	var clockIn float64
	var clockOut float64
	if c.Type == Offshore {
		clockIn = 0
		clockOut = 8
	} else {
		clockIn = 8
		clockOut = 16
	}
	if currentDayTime < clockOut &&
		currentDayTime <= clockIn &&
		currentWeekDay < (24*5) {
		//On the clock, calculate avaliable bandwidth for the day
		remainingCapacity := clockOut - currentDayTime
		if remainingCapacity > remainingDuration {
			return duration + timestamp
		}
		remainingDuration -= remainingCapacity
	}

	last_leg := math.Mod(remainingDuration, 8)
	daysDuration := int(remainingDuration) / 8
	proposedEnd := timestamp + 24 + float64(24*daysDuration) + last_leg
	if int(proposedEnd)/(24*5) > 0 {
		//then we ran into the weekend, add 2 days
		proposedEnd += 24 * 2
	}
	return proposedEnd

	// currentWorkingDay := math.Mod(timestamp, 24*7)
	// var nextAvaliableWindow float64
	// if currentWorkingDay > (24 * 4) {
	// 	// then the next working day will be Monday
	// 	nextWeek := int(timestamp)/(24*7) + 24*7
	// 	nextAvaliableWindow = float64(nextWeek) + clockIn
	// } else {
	// 	tomorrow := int(timestamp)/24 + 24
	// 	nextAvaliableWindow = float64(tomorrow) + clockIn
	// }
	// last_leg := math.Mod(remainingDuration, 8)

	// return nextAvaliableWindow + float64(24*daysDuration) + last_leg

}

func (c *Resource) formatTime(timestamp float64) string {
	day := int(timestamp) / (24)
	week := int(timestamp) / (24 * 7)
	return fmt.Sprintf("Time %f, Day %d, Week %d", timestamp, day, week)
}
