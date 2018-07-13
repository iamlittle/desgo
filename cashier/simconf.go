package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"gopkg.in/yaml.v2"
)

type SimConfig struct {
	Kind string `yaml:"kind"`
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