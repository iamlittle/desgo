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

func printResults(results map[SimConfig]*Stats){
	for sc, stats := range results{
		log.Println(fmt.Sprintf("[INFO] ------ Begin Simulation  %s ------", sc.Metadata.Name))
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
		log.Println(fmt.Sprintf("[INFO] ------ End Simulation %s ------", sc.Metadata.Name))
	}
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
	results := make(map[SimConfig]*Stats)
	var wg sync.WaitGroup
	for _, simConfig := range simConfigs{
		wg.Add(1)
		go func(sc SimConfig){
			results[sc] = RunSim(sc)
			WriteResults(sc, results[sc])
			wg.Done()
		}(simConfig)
	}
	wg.Wait()
	printResults(results)
}