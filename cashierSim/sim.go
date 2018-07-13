package main

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"container/heap"
	"fmt"
	"math"
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
	var stats = NewStats(&StatsConfig{
		math.Pow(2,2),
		2,
		math.Pow(5,2),
		8,
		math.Pow(10,2),
		30,
	})
	var pendingEventSet = NewPendingEventSet(&stats)

	customerCount := 100
	cashierCount := 2
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
		event := pendingEventSet.nextEvent()
		_, ok := event.(*CustomerGenerator)
		if !ok && !stats.WarmedUp {
			stats.WarmedUp = true
		}
		event.Transition()
	}

	serviceMean := stats.Mean(stats.CashierServiceTimes)
	shopMean := stats.Mean(stats.CustomerShopTimes)
	waitMean := stats.Mean(stats.CustomerWaitTimes)
	log.Println(fmt.Sprintf("[INFO] Business closed at %f", stats.GlobalTime))
	log.Println(fmt.Sprintf("[INFO] Idle Time Mean %f", stats.Mean(stats.CashierIdleTimes)))
	log.Println(fmt.Sprintf("[INFO] Service Time Mean %f", stats.Mean(stats.CashierServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Service Time StdDev %f", stats.StdDev(serviceMean, stats.CashierServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Shop Time StdDev %f", stats.StdDev(shopMean, stats.CustomerShopTimes)))
	log.Println(fmt.Sprintf("[INFO] Shop Time Mean %f", shopMean))
	log.Println(fmt.Sprintf("[INFO] Wait Time StdDev %f", stats.StdDev(waitMean, stats.CustomerWaitTimes)))
	log.Println(fmt.Sprintf("[INFO] Wait Time Mean %f", waitMean))

}