package main

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"container/heap"
	"fmt"
)

func main(){
	// Use the Hashicorp logutils to create level capable logging functionality
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("DEBUG"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)

	var business = NewBusiness()
	var stats = NewStats()
	var pendingEventSet = NewPendingEventSet(&stats)

	customerCount := 10
	cashierCount := 2
	for i:=0; i < customerCount; i++ {
		customerGenerator := NewCustomerGenerator(float32(i), &pendingEventSet, &business, &stats)
		pendingEventSet.scheduleEvent(&customerGenerator)
	}
	heap.Init(&pendingEventSet)
	for i:=0; i < cashierCount; i++ {
		cashier := NewCashier(float32(i), &pendingEventSet, &business, &stats)
		business.NotifyCashierAvailable(&cashier, float32(i))
	}

	for len(pendingEventSet.Events) > 0 {
		pendingEventSet.nextEvent().Transition()
	}

	log.Println(fmt.Sprintf("[INFO] Business closed at %f", stats.GlobalTime))
	log.Println(fmt.Sprintf("[INFO] Mean Idle Time %f", stats.Mean(stats.IdleTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Service Time %f", stats.Mean(stats.ServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Shop Time %f", stats.Mean(stats.ShopTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Wait Time %f", stats.Mean(stats.WaitTimes)))


}