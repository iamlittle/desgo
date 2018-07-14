package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
	"fmt"
	"path"
	"strconv"
	"sync"
)

type SimConfig struct {
	Kind string `yaml:"kind"`
	NumberOfRuns int `yaml:"number_of_runs"`
	Metadata struct {
		Name string
	} `yaml:"metadata"`
	Spec struct {
		StatsConfig
		CustomerCount int `yaml:"customer_count"`
		CashierCount int `yaml:"cashier_count"`
	} `yaml:"spec"`
	Output struct {
		Path string `yaml:"path"`
	} `yaml:"output"`
	sync.Mutex
}

type StatsConfig struct {
	ServiceTimeVariance float64 `yaml:"service_time_variance"`
	ServiceTimeMean float64 `yaml:"service_time_mean"`
	ShopTimeVariance float64 `yaml:"shop_time_variance"`
	ShopTimeMean float64 `yaml:"shop_time_mean"`
	EntryTimeVariance float64 `yaml:"entry_time_variance"`
	EntryTimeMean float64 `yaml:"entry_time_mean"`
}

func ReadSimConfig(filename string) []SimConfig {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		log.Printf("[ERROR] Could not read file: %s", filename)
		os.Exit(1)
	}
	configs := strings.Split(string(dat), "---")
	simConfigs := make([]SimConfig, len(configs))
	for i, config := range configs {
		simConfig := SimConfig{}
		err = yaml.Unmarshal([]byte(config), &simConfig)
		simConfigs[i] = simConfig
	}
	return simConfigs
}

func (*SimConfig) CleanOutputFile(sc SimConfig){
	_ = os.Remove(sc.Output.Path)
}

func (s *SimConfig) WriteResults(runIndex int, sc SimConfig, stats *Stats){
	s.Lock()
	_ = os.Mkdir(path.Dir(sc.Output.Path), 0755)

	f, _ := os.OpenFile(sc.Output.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	//only write the input / configuration section once
	if runIndex == 0{
		f.WriteString(fmt.Sprintf("SIM_NAME\t%s\n", sc.Metadata.Name))
		f.WriteString(fmt.Sprintf("SIM_KIND\t%s\n", sc.Kind))
		f.WriteString(fmt.Sprintf("SIM_NUMBER_OF_RUNS\t%d\n", sc.NumberOfRuns))
		f.WriteString(fmt.Sprintf("INP_CASHIER_COUNT\t%d\n", sc.Spec.CashierCount))
		f.WriteString(fmt.Sprintf("INP_CASHIER_COUNT\t%d\n", sc.Spec.CashierCount))
		f.WriteString(fmt.Sprintf("INP_CUSTOMER_COUNT\t%d\n", sc.Spec.CustomerCount))
		f.WriteString(fmt.Sprintf("INP_ENTRY_TIME_VARIANCE\t%f\n", sc.Spec.EntryTimeVariance))
		f.WriteString(fmt.Sprintf("INP_ENTRY_TIME_MEAN\t%f\n", sc.Spec.EntryTimeMean))
		f.WriteString(fmt.Sprintf("INP_SHOP_TIME_VARIANCE\t%f\n", sc.Spec.ShopTimeVariance))
		f.WriteString(fmt.Sprintf("INP_SHOP_TIME_MEAN\t%f\n", sc.Spec.ShopTimeMean))
		f.WriteString(fmt.Sprintf("INP_SERVICE_TIME_VARIANCE\t%f\n", sc.Spec.ServiceTimeVariance))
		f.WriteString(fmt.Sprintf("INP_SERVICE_TIME_MEAN\t%f\n", sc.Spec.ServiceTimeMean))
		f.WriteString("\n")
	}

	f.WriteString(fmt.Sprintf("OUT_ENTITY_COUNT_%d\t%d\n", runIndex, stats.EntityCount))
	f.WriteString(fmt.Sprintf("OUT_TERMINATION_TIME_%d\t%f\n", runIndex, stats.GlobalTime))
	writeFloatSlice(f,fmt.Sprintf("OUT_CUSTOMER_WAIT_TIMES_%d", runIndex), stats.CustomerWaitTimes)
	writeFloatSlice(f,fmt.Sprintf("OUT_CUSTOMER_SHOP_TIMES_%d", runIndex), stats.CustomerShopTimes)
	writeFloatSlice(f,fmt.Sprintf("OUT_CASHIER_IDLE_TIMES_%d", runIndex), stats.CashierIdleTimes)
	writeFloatSlice(f,fmt.Sprintf("OUT_CASHIER_SERVICE_TIMES_%d", runIndex), stats.CashierServiceTimes)
	f.WriteString("\n")
	f.Sync()
	s.Unlock()
}

func writeFloatSlice (f *os.File, fieldName string, values []float64){
	f.WriteString(fieldName)
	f.WriteString("\t")
	for _, val := range values{
		f.WriteString(strconv.FormatFloat(val, 'f', 4, 64))
		f.WriteString("\t")
	}
	f.WriteString("\n")
}