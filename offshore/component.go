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
type ComponentType string

const (
	Spark     ComponentType = "spark"
	DataRobot               = "datarobot"
	Inline                  = "inline"
	ETL                     = "etl"
)

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
8 - Complete
*************/

type ComponentState int64

const (
	CodeMigrated ComponentState = iota
	ReviewArtifacts
	RoughConversion
	DDLAvailable
	UnitTest
	DataAvailable
	Validate
	Release
	Complete
)

type Component struct {
	Id, Index, State                                                int
	CodeMigrateDuration, ReviewDuration, ConvertDuration            float64
	UnitTestDuration, ValidateDuration, CutoverDuration             float64
	Created, Timestamp, Lifespan, CodeMigrated, Reviewed, Converted float64
	DDLAvailable, UnitTested, DataAvailable, Validated, Cutover     float64
	ComponentType                                                   ComponentType
	PendingEventSet                                                 *PendingEventSet
	Migration                                                       *Migration
	Stats                                                           *Stats
}

func NewComponent(
	timestamp float64,
	componentType ComponentType,
	pendingEventSet *PendingEventSet,
	migration *Migration,
	stats *Stats) Component {
	customer := Component{
		stats.generateEntityId(), //Id
		0,                        //Index
		0,                        //State
		stats.generateCodeMigratedDuration(componentType), //CodeMigrateDuration
		stats.generateReviewDuration(componentType),       //ReviewDuration
		stats.generateConvertDuration(componentType),      //ConvertDuration
		stats.generateUnitTestDuration(componentType),     //UnitTestDuration
		stats.generateValidateDuration(componentType),     //ValidateDuration
		stats.generateCutoverDuration(componentType),      //CutoverDuration
		timestamp,       //Created
		timestamp,       //Timestamp
		0,               //Lifespan
		0,               //CodeMigrated
		0,               //Reviewed
		0,               //Converted
		0,               //DDLAvailable
		0,               //UnitTested
		0,               //DataAvailable
		0,               //Validated
		0,               //Released
		componentType,   //ComponentType
		pendingEventSet, //PendingEventSet
		migration,       //Migration
		stats,           //Component
	}
	return customer
}

func (c *Component) EventInfo() (int, int, float64) {
	return c.Id, c.Index, c.Timestamp
}

func (c *Component) Transition() bool {
	if c.State == int(Complete) {
		c.EndService(c.Timestamp)
	} else {
		c.Migration.Process(c, c.Timestamp)
	}
	return true
}

func (c *Component) EndService(timestamp float64) {
	c.Lifespan = timestamp - c.Created
	log.Println(fmt.Sprintf("[DEBUG] Component %d spent %f in migration", c.Id, c.Lifespan))
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

// func (c *ComponentGenerator) Transition() bool {
// 	customer := NewComponent(c.Stats.generateEntryTime(), c.PendingEventSet, c.Migration, c.Stats)
// 	c.PendingEventSet.scheduleEvent(&customer)
// 	log.Println(fmt.Sprintf("[DEBUG] Component %d generated at %f", customer.Id, c.Timestamp))
// 	return false
// }
