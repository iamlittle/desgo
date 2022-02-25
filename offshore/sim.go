package main

import (
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/logutils"
)

func RunSim(config *SimConfig) *Stats {
	var migration = NewMigration()
	var stats = NewStats(&config.Spec.StatsConfig)
	var pendingEventSet = NewPendingEventSet(&stats)

	for i := 0; i < config.Spec.ComponentCount; i++ {
		component := NewComponent(float64(i), &pendingEventSet, &migration, &stats)
		pendingEventSet.scheduleEvent(&component)
	}
	heap.Init(&pendingEventSet)
	for i := 0; i < config.Spec.OnshoreResourceCount; i++ {
		resource := NewResource(1, Onshore, &pendingEventSet, &migration, &stats)
		migration.NotifyResourceAvailable(&resource, float64(i))
	}
	for i := 0; i < config.Spec.OffshoreResourceCount; i++ {
		resource := NewResource(1, Offshore, &pendingEventSet, &migration, &stats)
		migration.NotifyResourceAvailable(&resource, float64(i))
	}
	stats.WarmedUp = true
	for len(pendingEventSet.Events) > 0 {
		event := pendingEventSet.nextEvent()
		event.Transition()
	}

	return &stats
}

func printResults(index int, simConfig *SimConfig, results *Stats) {
	log.Println(fmt.Sprintf("[INFO] ------ Begin Simulation  %s_%d ------", simConfig.Metadata.Name, index))
	codeMigratedMean := results.Mean(results.CodeMigratedTimes)
	reviewMean := results.Mean(results.ReviewTimes)
	conversionMean := results.Mean(results.ConversionTimes)
	unitTestMean := results.Mean(results.UnitTestTimes)
	validatedMean := results.Mean(results.ValidatedTimes)
	cutoverMean := results.Mean(results.CutoverTimes)
	log.Println(fmt.Sprintf("[INFO] Migration finished at %f", results.GlobalTime))
	log.Println(fmt.Sprintf("[INFO] Code Migration Mean %f", codeMigratedMean))
	log.Println(fmt.Sprintf("[INFO] Review Mean %f", reviewMean))
	log.Println(fmt.Sprintf("[INFO] Review StdDev %f", results.StdDev(reviewMean, results.ReviewTimes)))
	log.Println(fmt.Sprintf("[INFO] Conversion StdDev %f", results.StdDev(conversionMean, results.ConversionTimes)))
	log.Println(fmt.Sprintf("[INFO] Conversion Mean %f", conversionMean))
	log.Println(fmt.Sprintf("[INFO] Unit Test StdDev %f", results.StdDev(unitTestMean, results.UnitTestTimes)))
	log.Println(fmt.Sprintf("[INFO] Unit Test Mean %f", unitTestMean))
	log.Println(fmt.Sprintf("[INFO] Validated StdDev %f", results.StdDev(validatedMean, results.ValidatedTimes)))
	log.Println(fmt.Sprintf("[INFO] Validated Mean %f", validatedMean))
	log.Println(fmt.Sprintf("[INFO] Cutover StdDev %f", results.StdDev(cutoverMean, results.CutoverTimes)))
	log.Println(fmt.Sprintf("[INFO] Cutover Mean %f", cutoverMean))
	log.Println(fmt.Sprintf("[INFO] ------ End Simulation %s_%d ------", simConfig.Metadata.Name, index))
}

func main() {
	// Use the Hashicorp logutils to create level capable logging functionality
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	inputFile := flag.String("input", "", "path to a cashier yaml file")
	flag.Parse()

	simConfigs := ReadSimConfig(*inputFile)
	var wg sync.WaitGroup
	for _, simConfig := range simConfigs {
		simConfig.CleanOutputFile()
		for i := 0; i < simConfig.NumberOfRuns; i++ {
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
