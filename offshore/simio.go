package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type SimConfig struct {
	Kind         string `yaml:"kind"`
	NumberOfRuns int    `yaml:"number_of_runs"`
	Metadata     struct {
		Name string
	} `yaml:"metadata"`
	Spec struct {
		Configs               []StatsConfig `yaml:"statsconfig"`
		OnshoreResourceCount  int           `yaml:"onshore_count"`
		OffshoreResourceCount int           `yaml:"offshore_count"`
	} `yaml:"spec"`
	Output struct {
		Path string `yaml:"path"`
	} `yaml:"output"`
	sync.Mutex
}

type StatsConfig struct {
	Name               string  `yaml:"name"`
	CodeMigratedStdDev float64 `yaml:"code_migrated_std_dev"`
	CodeMigratedMean   float64 `yaml:"code_migrated_mean"`
	ReviewStdDev       float64 `yaml:"review_std_dev"`
	ReviewMean         float64 `yaml:"review_mean"`
	ConvertStdDev      float64 `yaml:"convert_std_dev"`
	ConvertMean        float64 `yaml:"convert_mean"`
	UnitTestStdDev     float64 `yaml:"unit_test_std_dev"`
	UnitTestMean       float64 `yaml:"unit_test_mean"`
	ValidateStdDev     float64 `yaml:"validate_std_dev"`
	ValidateMean       float64 `yaml:"validate_mean"`
	CutoverStdDev      float64 `yaml:"cutover_std_dev"`
	CutoverMean        float64 `yaml:"cutover_mean"`
	ComponentCount     int     `yaml:"component_count"`
}

func ReadSimConfig(filename string) []*SimConfig {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		log.Printf("[ERROR] Could not read file: %s", filename)
		os.Exit(1)
	}
	configs := strings.Split(string(dat), "---")
	simConfigs := make([]*SimConfig, len(configs))
	for i, config := range configs {
		simConfig := &SimConfig{}
		err = yaml.Unmarshal([]byte(config), &simConfig)
		simConfigs[i] = simConfig
	}
	return simConfigs
}

func (s *SimConfig) CleanOutputFile() {
	_ = os.Remove(s.Output.Path)
}

func (s *SimConfig) WriteResults(runIndex int, stats *Stats) {
	s.Lock()
	defer s.Unlock()
	_ = os.Mkdir(path.Dir(s.Output.Path), 0755)

	f, _ := os.OpenFile(s.Output.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	//only write the input / configuration section once
	if runIndex == 0 {
		f.WriteString(fmt.Sprintf("SIM_NAME\t%s\n", s.Metadata.Name))
		f.WriteString(fmt.Sprintf("SIM_KIND\t%s\n", s.Kind))
		f.WriteString(fmt.Sprintf("SIM_NUMBER_OF_RUNS\t%d\n", s.NumberOfRuns))
		// f.WriteString(fmt.Sprintf("INP_CASHIER_COUNT\t%d\n", s.Spec.CashierCount))
		// f.WriteString(fmt.Sprintf("INP_CASHIER_COUNT\t%d\n", s.Spec.CashierCount))
		// f.WriteString(fmt.Sprintf("INP_CUSTOMER_COUNT\t%d\n", s.Spec.CustomerCount))
		// f.WriteString(fmt.Sprintf("INP_ENTRY_TIME_STD_DEV\t%f\n", s.Spec.EntryTimeStdDev))
		// f.WriteString(fmt.Sprintf("INP_ENTRY_TIME_MEAN\t%f\n", s.Spec.EntryTimeMean))
		// f.WriteString(fmt.Sprintf("INP_SHOP_TIME_STD_DEV\t%f\n", s.Spec.ShopTimeStdDev))
		// f.WriteString(fmt.Sprintf("INP_SHOP_TIME_MEAN\t%f\n", s.Spec.ShopTimeMean))
		// f.WriteString(fmt.Sprintf("INP_SERVICE_TIME_STD_DEV\t%f\n", s.Spec.ServiceTimeStdDev))
		// f.WriteString(fmt.Sprintf("INP_SERVICE_TIME_MEAN\t%f\n", s.Spec.ServiceTimeMean))
		f.WriteString("\n")
	}

	f.WriteString(fmt.Sprintf("OUT_ENTITY_COUNT_%d\t%d\n", runIndex, stats.EntityCount))
	f.WriteString(fmt.Sprintf("OUT_TERMINATION_TIME_%d\t%f\n", runIndex, stats.GlobalTime))
	// writeFloatSlice(f, fmt.Sprintf("OUT_CUSTOMER_WAIT_TIMES_%d", runIndex), stats.CustomerWaitTimes)
	// writeFloatSlice(f, fmt.Sprintf("OUT_CUSTOMER_ENTRY_TIMES_%d", runIndex), stats.CustomerEntryTimes)
	// writeFloatSlice(f, fmt.Sprintf("OUT_CUSTOMER_SHOP_TIMES_%d", runIndex), stats.CustomerShopTimes)
	// writeFloatSlice(f, fmt.Sprintf("OUT_CASHIER_IDLE_TIMES_%d", runIndex), stats.CashierIdleTimes)
	// writeFloatSlice(f, fmt.Sprintf("OUT_CASHIER_SERVICE_TIMES_%d", runIndex), stats.CashierServiceTimes)
	f.WriteString("\n")
	f.Sync()
}

func writeFloatSlice(f *os.File, fieldName string, values []float64) {
	f.WriteString(fieldName)
	f.WriteString("\t")
	for _, val := range values {
		f.WriteString(strconv.FormatFloat(val, 'f', 4, 64))
		f.WriteString("\t")
	}
	f.WriteString("\n")
}
