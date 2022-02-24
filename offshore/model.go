package main

import (
	"fmt"
	"log"
)

/*************
** Component Type **
1 - Batch Model
2 - Real Time Model
3 - ETL
4 - Bundle
5 - QES
*************/
type ComponentType struct {
	Id int
}

/*************
** Component State **
0 - Not Started
1 - Code migrated
2 - Review Artifacts
3 - Rough Conversion
4 - DDL Available
5 - Unit Test
6 - Data Available
7 - Execute and Validate
8 - Release and Production Cutover
*************/

type Component struct {
	Id, Index, State                                                       int
	Review_Time, Convert_Time, Unit_Test_Time, Validate_Time, Release_Time float64
	Created, Timestamp, Lifespan, Code_Migrated, Reviewed, Converted       float64
	DDL_Available, Unit_Tested, Data_Available, Validated, Released        float64
	PendingEventSet                                                        *PendingEventSet
	Migration                                                              *Migration
	Stats                                                                  *Stats
}

func NewComponent(timestamp float64,
	pendingEventSet *PendingEventSet,
	migration *Migration,
	stats *Stats) *Component {
	customer := &Component{
		stats.generateEntityId(),     //Id
		0,                            //Index
		0,                            //State
		stats.generateReviewTime(),   //Review_Time
		stats.generateConvertTime(),  //Convert_Time
		stats.generateUnitTestTime(), //Unit_Test_Time
		stats.generateValidateTime(), //Validate_Time
		stats.generateReleaseTime(),  //Release_Time
		timestamp,                    //Created
		timestamp,                    //Timestamp
		0,                            //Lifespan
		0,                            //Code_Migrated
		0,                            //Reviewed
		0,                            //Converted
		0,                            //DDL_Available
		0,                            //Unit_Tested
		0,                            //Data_Available
		0,                            //Validated
		0,                            //Released
		pendingEventSet,              //PendingEventSet
		migration,                    //Migration
		stats,                        //Component
	}
	return customer
}

func (c *Component) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.Timestamp
}

func (c *Component) Transition() bool {
	switch c.State {
	case 0:
		c.State++
		c.Stats.RecordComponentEntryTime(c.Timestamp)
		c.Stats.RecordComponentShopTime(c.ShopTime)
		log.Println(fmt.Sprintf("[DEBUG] Component %d entered store at %f", c.Id, c.Timestamp))
		c.Timestamp += c.ShopTime
		c.PendingEventSet.scheduleEvent(c)
	case 1:
		c.EnterQueue = c.Timestamp
		log.Println(fmt.Sprintf("[DEBUG] Component %d finished shopping at %f", c.Id, c.Timestamp))
		c.Migration.Checkout(c, c.Timestamp)
		c.State++
	}
	return true
}

func (c *Component) EndWait(timestamp float64) {
	c.WaitTime = timestamp - c.EnterQueue
	log.Println(fmt.Sprintf("[DEBUG] Component %d waited in line for %f", c.Id, c.WaitTime))
	c.Stats.RecordComponentWaitTime(c.WaitTime)
	c.State++
}

func (c *Component) EndService(timestamp float64) {
	c.Lifespan = timestamp - c.Created
	log.Println(fmt.Sprintf("[DEBUG] Component %d was at migration for %f", c.Id, c.Lifespan))
	c.State++
	c.Stats.CompletedJobCount++
}

/***********************
** Component Generator **
************************/
type ComponentGenerator struct {
	Id, Index          int
	Created, Timestamp float64
	PendingEventSet    *PendingEventSet
	Migration          *Migration
	Stats              *Stats
}

func NewComponentGenerator(timestamp float64,
	pendingEventSet *PendingEventSet,
	migration *Migration,
	stats *Stats) ComponentGenerator {
	customerGen := ComponentGenerator{stats.generateEntityId(),
		0,
		timestamp,
		timestamp,
		pendingEventSet,
		migration,
		stats,
	}
	return customerGen
}

func (c *ComponentGenerator) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.Timestamp
}

func (c *ComponentGenerator) Transition() bool {
	customer := NewComponent(c.Stats.generateEntryTime(), c.PendingEventSet, c.Migration, c.Stats)
	c.PendingEventSet.scheduleEvent(customer)
	log.Println(fmt.Sprintf("[DEBUG] Component %d generated at %f", customer.Id, c.Timestamp))
	return false
}
