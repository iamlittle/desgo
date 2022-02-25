package main

import (
	"fmt"
	"log"
)

/*************
** Migration **
**************/
type Migration struct {
	OnshoreResources   []*Resource
	OffshoreResources  []*Resource
	OnshoreComponents  []*Component
	OffshoreComponents []*Component
}

func NewMigration() Migration {
	return Migration{make([]*Resource, 0), make([]*Resource, 0), make([]*Component, 0), make([]*Component, 0)}
}

type ResourceType int64

const (
	Offshore ResourceType = iota
	Onshore
	Other
)

func (c *Migration) getResourceType(component *Component) ResourceType {
	switch component.State {
	case
		1,
		2,
		4:
		return Offshore
	case
		6,
		7:
		return Onshore
	}
	return Other
}

func (c *Migration) Process(component *Component, timestamp float64) {
	t := c.getResourceType(component)
	if t == Offshore {
		if len(c.OffshoreResources) == 0 {
			log.Println(fmt.Sprintf("[DEBUG] No offshore resources available for Component %d at %f", component.Id, timestamp))
			c.OffshoreComponents = append(c.OffshoreComponents, component)
		} else {
			resource := c.OffshoreResources[0]
			c.OffshoreResources = c.OffshoreResources[1:]
			resource.Process(component, timestamp)
		}
	} else if t == Onshore {
		if len(c.OnshoreResources) == 0 {
			log.Println(fmt.Sprintf("[DEBUG] No onshore resources available for Component %d at %f", component.Id, timestamp))
			c.OffshoreComponents = append(c.OffshoreComponents, component)
		} else {
			resource := c.OnshoreResources[0]
			c.OnshoreResources = c.OnshoreResources[1:]
			resource.Process(component, timestamp)
		}
	} else {
		switch component.State {
		case 0:
			component.CodeMigrated = component.CodeMigrateDuration
			component.Timestamp = timestamp
		case 3:
			if component.DDLAvailable > timestamp {
				component.Timestamp = component.DDLAvailable
			} else {
				component.Timestamp = timestamp
			}
		case 5:
			if component.DataAvailable > timestamp {
				component.Timestamp = component.DataAvailable
			} else {
				component.Timestamp = timestamp
			}
		}
	}
}

func (c *Migration) NotifyResourceAvailable(resource *Resource, timestamp float64) {
	resource.TimeStamp = timestamp
	if resource.Type == 0 && len(c.OffshoreResources) == 0 {
		c.OffshoreResources = append(c.OffshoreResources, resource)
	} else if resource.Type == 0 {
		component := c.OffshoreComponents[0]
		c.OffshoreComponents = c.OffshoreComponents[1:]
		resource.Process(component, timestamp)
	} else if len(c.OnshoreResources) == 0 {
		c.OnshoreResources = append(c.OnshoreResources, resource)
	} else {
		component := c.OnshoreComponents[0]
		c.OnshoreComponents = c.OnshoreComponents[1:]
		resource.Process(component, timestamp)
	}
}
