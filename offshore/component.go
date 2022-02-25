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
0 - Code migrated
1 - Review Artifacts
2 - Rough Conversion
3 - DDL Available
4 - Unit Test
5 - Data Available
6 - Execute and Validate
7 - Release and Production Cutover
*************/

type Component struct {
	Id, Index, State                                                int
	CodeMigrateDuration, ReviewDuration, ConvertDuration            float64
	UnitTestDuration, ValidateDuration, CutoverDuration             float64
	Created, Timestamp, Lifespan, CodeMigrated, Reviewed, Converted float64
	DDLAvailable, UnitTested, DataAvailable, Validated, Cutover     float64
	PendingEventSet                                                 *PendingEventSet
	Migration                                                       *Migration
	Stats                                                           *Stats
}

func NewComponent(timestamp float64,
	pendingEventSet *PendingEventSet,
	migration *Migration,
	stats *Stats) Component {
	customer := Component{
		stats.generateEntityId(),             //Id
		0,                                    //Index
		0,                                    //State
		stats.generateCodeMigratedDuration(), //CodeMigrateDuration
		stats.generateReviewDuration(),       //ReviewDuration
		stats.generateConvertDuration(),      //ConvertDuration
		stats.generateUnitTestDuration(),     //UnitTestDuration
		stats.generateValidateDuration(),     //ValidateDuration
		stats.generateCutoverDuration(),      //CutoverDuration
		timestamp,                            //Created
		timestamp,                            //Timestamp
		0,                                    //Lifespan
		0,                                    //CodeMigrated
		0,                                    //Reviewed
		0,                                    //Converted
		0,                                    //DDLAvailable
		0,                                    //UnitTested
		0,                                    //DataAvailable
		0,                                    //Validated
		0,                                    //Released
		pendingEventSet,                      //PendingEventSet
		migration,                            //Migration
		stats,                                //Component
	}
	return customer
}

func (c *Component) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.Timestamp
}

func (c *Component) Transition() bool {
	switch c.State {
	case 0:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentCodeMigrationTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d code is migrated at %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 1:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentReviewTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d artifacts reviewed at %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 2:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentRoughConversionTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d rough conversion %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 3:
		c.Migration.Process(c, c.Timestamp)
		log.Println(fmt.Sprintf("[DEBUG] Component %d DDL available at %f, processed at %f", c.Id, c.DDLAvailable, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 4:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentUnitTestTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d unit tested at %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 5:
		c.Migration.Process(c, c.Timestamp)
		log.Println(fmt.Sprintf("[DEBUG] Component %d Data available at %f, processed at %f", c.Id, c.DDLAvailable, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 6:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentValidateTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d executed and validated at %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	case 7:
		c.Migration.Process(c, c.Timestamp)
		c.Stats.RecordComponentCutoverTime(c)
		log.Println(fmt.Sprintf("[DEBUG] Component %d production cutover at %f", c.Id, c.Timestamp))
		c.State++
		c.PendingEventSet.scheduleEvent(c)
	}
	return true
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
	c.PendingEventSet.scheduleEvent(&customer)
	log.Println(fmt.Sprintf("[DEBUG] Component %d generated at %f", customer.Id, c.Timestamp))
	return false
}
