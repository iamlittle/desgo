package main

import (
	"math/rand"
	"time"
	"math"
)

type Stats struct{
	EntityCount int
	CompletedJobCount int
	JobCount int
	GlobalTime float64
	CustomerShopTimes []float64
	CustomerWaitTimes []float64
	CashierIdleTimes []float64
	CashierServiceTimes []float64
	CustomerEntryTimes []float64
	StatsConfig *StatsConfig
	WarmedUp bool
	Source rand.Source
}


func NewStats(config *StatsConfig) Stats{
	return Stats{
		0, 0, 0, 0,
		make([]float64, 0), make([]float64, 0),
		make([]float64, 0),make([]float64, 0),make([]float64, 0),
		config, false, rand.NewSource(time.Now().UnixNano()),
	}
}

func (s *Stats) generateGaussianRandomNumber(stdDev float64, mean float64) float64{
	rnd := rand.New(s.Source)
	return rnd.NormFloat64() * stdDev + mean
}

func (s *Stats) generateLogNormalRandomNumber(stdDev float64, mean float64) float64{
	//ensure the value is positive
	mu := s.Mu(math.Pow(stdDev, 2), mean)
	sigma := s.Sigma(math.Pow(stdDev, 2), mean)
	gaussian :=  math.Abs(s.generateGaussianRandomNumber(sigma, mu))
	return math.Exp(gaussian)
}

func (s *Stats) generateExponentialRandomNumber(stdDev float64, mean float64, rate float64) float64{
	gaussian := math.Abs(s.generateGaussianRandomNumber(1, 0))
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	exponential :=  math.Log(1-gaussian) / -rate
	return exponential * stdDev + mean
}

func (s *Stats) generateServiceTime() float64{
	return math.Abs(s.generateGaussianRandomNumber(s.StatsConfig.ServiceTimeStdDev, s.StatsConfig.ServiceTimeMean))
}

func (s *Stats) generateShopTime() float64{
	return s.generateLogNormalRandomNumber(s.StatsConfig.ShopTimeStdDev, s.StatsConfig.ShopTimeMean)
}

func (s *Stats) generateEntryTime() float64{
	entryTime := s.generateLogNormalRandomNumber(s.StatsConfig.EntryTimeStdDev, s.StatsConfig.EntryTimeMean)
	if entryTime < s.GlobalTime{
		entryTime += s.GlobalTime
	}
	return entryTime
}

func (s *Stats) generateEntityId() int{
	id := s.EntityCount
	s.EntityCount++
	return id
}

func (s *Stats) RecordCashierIdleTime(idleTime float64) {
	if s.WarmedUp{
		s.CashierIdleTimes = append(s.CashierIdleTimes, idleTime)
	}
}

func (s *Stats) RecordCashierServiceTime(serviceTime float64) {
	if s.WarmedUp {
		s.CashierServiceTimes = append(s.CashierServiceTimes, serviceTime)
	}
}

func (s *Stats) RecordCustomerShopTime(shopTime float64) {
	if s.WarmedUp{
		s.CustomerShopTimes = append(s.CustomerShopTimes, shopTime)
	}
}

func (s *Stats) RecordCustomerWaitTime(waitTime float64) {
	if s.WarmedUp{
		s.CustomerWaitTimes = append(s.CustomerWaitTimes, waitTime)
	}
}

func (s *Stats) RecordCustomerEntryTime(waitTime float64) {
	if s.WarmedUp{
		s.CustomerEntryTimes = append(s.CustomerEntryTimes, waitTime)
	}
}

func (s *Stats) Mean(values []float64) float64{
	if len(values) == 0 {
		return 0
	}else{
		var sum float64 = 0
		for _, value := range values { sum += value }
		return sum / float64(len(values))
	}
}

func (s *Stats) Variance (mean float64, values []float64) float64{
	if len(values) == 0 {
		return 0
	}else{
		var variance float64 = 0
		for _, value := range values { variance += math.Pow( value - mean, 2) }
		return variance
	}
}

func (s *Stats) StdDev(mean float64, values []float64) float64{
	return math.Sqrt(s.Variance(mean, values) / float64(len(values)))
}

func (s *Stats) Mu(variance float64, mean float64) float64{
	return math.Log(math.Pow(mean, 2) / math.Sqrt(variance + math.Pow(mean, 2)))
}

func (s *Stats) Sigma(variance float64, mean float64) float64{
	return math.Sqrt(math.Log(1+variance/math.Pow(mean, 2)))
}
