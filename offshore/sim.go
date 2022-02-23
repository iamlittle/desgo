package main

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"container/heap"
	"sync"
	"fmt"
	"flag"
)

func RunSim(config *SimConfig) *Stats{
	var business = NewBusiness()
	var stats = NewStats(&config.Spec.StatsConfig)
	var pendingEventSet = NewPendingEventSet(&stats)

	for i:=0; i < config.Spec.CustomerCount; i++ {
		customerGenerator := NewCustomerGenerator(1, &pendingEventSet, &business, &stats)
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

func printResults(index int, simConfig *SimConfig, results *Stats){
	log.Println(fmt.Sprintf("[INFO] ------ Begin Simulation  %s_%d ------", simConfig.Metadata.Name, index))
	serviceMean := results.Mean(results.CashierServiceTimes)
	shopMean := results.Mean(results.CustomerShopTimes)
	waitMean := results.Mean(results.CustomerWaitTimes)
	entryMean := results.Mean(results.CustomerEntryTimes)
	log.Println(fmt.Sprintf("[INFO] Business closed at %f", results.GlobalTime))
	log.Println(fmt.Sprintf("[INFO] Idle Time Mean %f", results.Mean(results.CashierIdleTimes)))
	log.Println(fmt.Sprintf("[INFO] Service Time Mean %f", results.Mean(results.CashierServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Service Time StdDev %f", results.StdDev(serviceMean, results.CashierServiceTimes)))
	log.Println(fmt.Sprintf("[INFO] Shop Time StdDev %f", results.StdDev(shopMean, results.CustomerShopTimes)))
	log.Println(fmt.Sprintf("[INFO] Shop Time Mean %f", shopMean))
	log.Println(fmt.Sprintf("[INFO] Wait Time StdDev %f", results.StdDev(waitMean, results.CustomerWaitTimes)))
	log.Println(fmt.Sprintf("[INFO] Wait Time Mean %f", waitMean))
	log.Println(fmt.Sprintf("[INFO] Entry Time StdDev %f", results.StdDev(entryMean, results.CustomerEntryTimes)))
	log.Println(fmt.Sprintf("[INFO] Entry Time Mean %f", entryMean))
	log.Println(fmt.Sprintf("[INFO] ------ End Simulation %s_%d ------", simConfig.Metadata.Name, index))
}

func main(){
	// Use the Hashicorp logutils to create level capable logging functionality
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer: os.Stderr,
	}
	log.SetOutput(filter)

	inputFile := flag.String("input", "", "path to a cashier yaml file")
	flag.Parse()

	simConfigs := ReadSimConfig(*inputFile)
	var wg sync.WaitGroup
	for _, simConfig := range simConfigs{
		simConfig.CleanOutputFile()
		for i:=0; i<simConfig.NumberOfRuns; i++ {
			wg.Add(1)
			go func(index int, sc *SimConfig) {
				results := RunSim(sc)
				sc.WriteResults(index, results)
				printResults(index, sc, results)
				wg.Done()
			}(i, simConfig)
		}
	}
	wg.Wait()
}