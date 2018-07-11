package main

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"container/heap"
	"fmt"
)

func main(){
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)
	var business = &Business{make([]*Cashier, 0), make([]*Customer, 0)}
	var stats = &Stats{0, 0, 0, 0,
				0, 0, 0}
	var pendingEventSet = &PendingEventSet{make([]Event, 0), make(map[int]int), stats}

	customerCount := 10
	cashierCount := 2
	for i:=0; i < customerCount; i++ {
		pendingEventSet.scheduleEvent(NewCustomerGenerator(float32(i), pendingEventSet, business, stats))
	}
	heap.Init(pendingEventSet)
	for i:=0; i < cashierCount; i++ {
		business.NotifyCashierAvailable(NewCashier(float32(i), pendingEventSet, business, stats), float32(i))
	}

	for len(pendingEventSet.Events) > 0 {
		pendingEventSet.nextEvent().Transition()
	}

	log.Println(fmt.Sprintf("[INFO] Business closed at %f", stats.GlobalTime))

}