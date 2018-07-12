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
		MinLevel: logutils.LogLevel("INFO"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)

	var business = NewBusiness()
	var stats = NewStats()
	var pendingEventSet = NewPendingEventSet(&stats)

	customerCount := 100
	cashierCount := 4
	for i:=0; i < customerCount; i++ {
		customerGenerator := NewCustomerGenerator(float64(i), &pendingEventSet, &business, &stats)
		pendingEventSet.scheduleEvent(&customerGenerator)
	}
	heap.Init(&pendingEventSet)
	for i:=0; i < cashierCount; i++ {
		cashier := NewCashier(float64(i), &pendingEventSet, &business, &stats)
		business.NotifyCashierAvailable(&cashier, float64(i))
	}

	for len(pendingEventSet.Events) > 0 {
		pendingEventSet.nextEvent().Transition()
	}

	shopMean := stats.Mean(stats.ShopTimes)
	serviceMean := stats.Mean(stats.ServiceTimes)
	log.Println(fmt.Sprintf("[INFO] Business closed at %f", stats.GlobalTime))
	log.Println(fmt.Sprintf("[INFO] Mean Idle Time %f", stats.Mean(stats.IdleTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Service Time %f", stats.Mean(stats.ServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] StdDev Service Time %f", stats.StdDev(serviceMean, stats.ServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Shop Time %f", shopMean))
	log.Println(fmt.Sprintf("[INFO] StdDev Shop Time %f", stats.StdDev(shopMean, stats.ShopTimes)))
	log.Println(fmt.Sprintf("[INFO] Mean Wait Time %f", stats.Mean(stats.WaitTimes)))


}