package main

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"container/heap"
	"fmt"
)

func RunSim(config SimConfig) *Stats{
	var business = NewBusiness()
	var stats = NewStats(&config.Spec.StatsConfig)
	var pendingEventSet = NewPendingEventSet(&stats)

	for i:=0; i < config.Spec.CustomerCount; i++ {
		customerGenerator := NewCustomerGenerator(float64(i), &pendingEventSet, &business, &stats)
		pendingEventSet.scheduleEvent(&customerGenerator)
	}
	heap.Init(&pendingEventSet)
	for i:=0; i < config.Spec.CashierCount; i++ {
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
	return &stats
}

func main(){
	// Use the Hashicorp logutils to create level capable logging functionality
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)

	simConfigs := ReadSimConfig("./examples/cashier.yaml")

	for _, simConfig := range simConfigs{
		stats := RunSim(simConfig)

		log.Println(fmt.Sprintf("[INFO] ------ Simulation  %s ------", simConfig.Metadata.Name))

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

		log.Println("[INFO] ---------- END Simulation ----------")
	}


}