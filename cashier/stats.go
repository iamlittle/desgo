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
	StatsConfig *StatsConfig
	WarmedUp bool
}

type StatsConfig struct{
	ServiceTimeVariance float64
	ServiceTimeMean float64
	ShopTimeVariance float64
	ShopTimeMean float64
	EntryTimeVariance float64
	EntryTimeMean float64
}

var source = rand.NewSource(time.Now().Unix())

func NewStats(config *StatsConfig) Stats{
	return Stats{
		0, 0, 0, 0,
		make([]float64, 0), make([]float64, 0),
		make([]float64, 0),make([]float64, 0), config, false,
	}
}

func (*Stats) generateGaussianRandomNumber(variance float64, mean float64) float64{
	rnd := rand.New(source)
	return rnd.NormFloat64() * math.Sqrt(variance) + mean
}

func (s *Stats) generateLogNormalRandomNumber(variance float64, mean float64) float64{
	//ensure the value is positive
	gaussian :=  math.Abs(s.generateGaussianRandomNumber(1, 0))
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	return math.Exp(s.Mu(variance, mean) + s.Sigma(variance, mean) * gaussian)
}

func (s *Stats) generateExponentialRandomNumber(variance float64, mean float64, rate float64) float64{
	gaussian := math.Abs(s.generateGaussianRandomNumber(1, 0))
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	exponential :=  math.Log(1-gaussian) / -rate
	return exponential * math.Sqrt(variance) + mean
}

func (s *Stats) generateServiceTime() float64{
	return math.Abs(s.generateGaussianRandomNumber(s.StatsConfig.ServiceTimeVariance, s.StatsConfig.ServiceTimeMean))
}

func (s *Stats) generateShopTime() float64{
	return s.generateLogNormalRandomNumber(s.StatsConfig.ShopTimeVariance, s.StatsConfig.ShopTimeMean)
}

func (s *Stats) generateEntryTime() float64{
	entryTime := s.generateLogNormalRandomNumber(s.StatsConfig.EntryTimeVariance, s.StatsConfig.EntryTimeMean)
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
	return math.Log(mean / math.Sqrt(1+math.Pow(variance, 2)/math.Pow(mean, 2)))
}

func (s *Stats) Sigma(variance float64, mean float64) float64{
	return math.Log(1+math.Pow(variance, 2)/math.Pow(mean, 2))
}
