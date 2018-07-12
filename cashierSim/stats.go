package main

import (
	"math/rand"
	"time"
	"math"
)

const serviceTimeStdDev float64 = 2
const serviceTimeMean float64 = 2
const shopTimeVariance float64 = 512
const shopTimeMean float64 = 8
const entryTimeVariance float64 = 100
const entryTimeMean float64 = 30

type Stats struct{
	EntityCount int
	CompletedJobCount int
	JobCount int
	GlobalTime float64
	WaitTimes []float64
	ServiceTimes []float64
	ShopTimes []float64
	IdleTimes []float64
}

var source = rand.NewSource(time.Now().Unix())

func NewStats() Stats{
	return Stats{
		0, 0, 0, 0,
		make([]float64, 0), make([]float64, 0),
		make([]float64, 0),make([]float64, 0),
	}
}

func (*Stats) generateGaussianRandomNumber(stdDev float64, mean float64) float64{
	rnd := rand.New(source)
	return rnd.NormFloat64() * stdDev + mean
}

func (s *Stats) generateLogNormalRandomNumber(variance float64, mean float64) float64{
	//ensure the value is positive
	gaussian :=  math.Abs(s.generateGaussianRandomNumber(1, 0))
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	return math.Exp(s.Mu(variance, mean) + s.Sigma(variance, mean) * gaussian)
}

func (s *Stats) generateExponentialRandomNumber(stdDev float64, mean float64, rate float64) float64{
	gaussian := math.Abs(s.generateGaussianRandomNumber(1, 0))
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	exponential :=  math.Log(1-gaussian) / -rate
	return exponential * stdDev + mean
}

func (s *Stats) generateServiceTime() float64{
	return math.Abs(s.generateGaussianRandomNumber(serviceTimeStdDev, serviceTimeMean))
}

func (s *Stats) generateShopTime() float64{
	return s.generateLogNormalRandomNumber(shopTimeVariance, shopTimeMean)
}

func (s *Stats) generateEntryTime() float64{
	entryTime := s.generateLogNormalRandomNumber(entryTimeVariance, entryTimeMean)
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

func (s *Stats) RecordIdleTime(idleTime float64) {
	s.IdleTimes = append(s.IdleTimes, idleTime)
}

func (s *Stats) RecordServiceTime(serviceTime float64) {
	s.ServiceTimes = append(s.ServiceTimes, serviceTime)
}

func (s *Stats) RecordShopTime(shopTime float64) {
	s.ShopTimes = append(s.ShopTimes, shopTime)
}

func (s *Stats) RecordWaitTime(waitTime float64) {
	s.WaitTimes = append(s.WaitTimes, waitTime)
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
	return math.Log(mean / math.Sqrt(1+variance/math.Pow(mean, 2)))
}

func (s *Stats) Sigma(variance float64, mean float64) float64{
	return math.Log(1+variance/math.Pow(mean, 2))
}
