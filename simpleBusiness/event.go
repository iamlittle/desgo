package main

import (
	"container/heap"
	"log"
	"fmt"
	"os"
)

type Event interface {
	//return ID, Index, TimeStamp
	EventInfo() (int, int, float32)
	Transition() bool
}

type PendingEventSet struct {
	Events []Event
	Indices map[int]int
	Stats *Stats
}

func NewPendingEventSet(stats *Stats) PendingEventSet{
	return PendingEventSet{make([]Event, 0), make(map[int]int), stats}
}

func (pes PendingEventSet) Len() int { return len(pes.Events) }

func (pes PendingEventSet) Less(i, j int) bool {
	_, _, iTimestamp := pes.Events[i].EventInfo()
	_, _, jTimestamp := pes.Events[j].EventInfo()
	return  iTimestamp < jTimestamp
}

func (pes PendingEventSet) Swap(i, j int) {
	iId, _, _ := pes.Events[i].EventInfo()
	jId, _, _ := pes.Events[j].EventInfo()
	pes.Events[i], pes.Events[j] = pes.Events[j], pes.Events[i]
	pes.Indices[iId] = i
	pes.Indices[jId] = j
}

func (pes *PendingEventSet) Push(x interface{}) {
	n := len(pes.Events)
	event := x.(Event)
	id, _, _ := event.EventInfo()
	pes.Indices[id] = n
	pes.Events = append(pes.Events, event)
}

func (pes *PendingEventSet) Pop() interface{} {
	n := len(pes.Events)
	event := pes.Events[n-1]
	id, _, _ := event.EventInfo()
	pes.Indices[id] = -1
	pes.Events = pes.Events[0 : n-1]
	return event
}

func (pes *PendingEventSet) scheduleEvent(e Event){
	heap.Push(pes, e)
}

func (pes *PendingEventSet) nextEvent() Event{
	event := heap.Pop(pes).(Event)
	_, _, timestamp := event.EventInfo()
	if pes.Stats.GlobalTime > timestamp{
		log.Println(fmt.Sprintf("[ERROR] next event in queue has a timestamp (%f) less than the Global Clock (%f). " +
			"The simulation is corrupt!", timestamp, pes.Stats.GlobalTime))
		os.Exit(1)
	}else{
		pes.Stats.GlobalTime = timestamp
	}
	return event
}
