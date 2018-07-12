package main

import (
	"math/rand"
	"time"
	"math"
)

const randomStdDevServiceTime float64 = 2
const randomMeanServiceTime float64 = 2
const randomStdDevShopTime float64 = 7
const randomMeanShopTime float64 = 8
const randomStdDevEntryTime float64 = 10
const randomMeanEntryTime float64 = 30

type Stats struct{
	EntityCount int
	CompletedJobCount int
	JobCount int
	GlobalTime float32
	WaitTimes []float32
	ServiceTimes []float32
	ShopTimes []float32
	IdleTimes []float32
}

var source = rand.NewSource(time.Now().Unix())

func NewStats() Stats{
	return Stats{
		0, 0, 0, 0,
		make([]float32, 0), make([]float32, 0),
		make([]float32, 0),make([]float32, 0),
	}
}

func (*Stats) generateGaussianRandomNumber(stdDev float64, mean float64) float64{
	rnd := rand.New(source)
	//ensure the value is positive (x^2)^(1/2)
	return math.Sqrt(math.Pow(rnd.NormFloat64(), 2)) * stdDev + mean
}

func (s *Stats) generateLogNormalRandomNumber(stdDev float64, mean float64) float64{
	gaussian := s.generateGaussianRandomNumber(1, 0)
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	return math.Exp(gaussian) * stdDev + mean
}

func (s *Stats) generateExponentialRandomNumber(stdDev float64, mean float64, rate float64) float64{
	gaussian := s.generateGaussianRandomNumber(1, 0)
	//ensures the value is between 0 and 1
	gaussian = gaussian - float64(int(gaussian))
	exponential :=  math.Log(1-gaussian) / -rate
	return exponential * stdDev + mean
}

func (s *Stats) generateServiceTime() float32{
	return float32(s.generateGaussianRandomNumber(randomStdDevServiceTime, randomMeanServiceTime))
}

func (s *Stats) generateShopTime() float32{
	return float32(s.generateLogNormalRandomNumber(randomStdDevShopTime, randomMeanShopTime))
}

func (s *Stats) generateEntryTime() float32{
	entryTime := float32(s.generateLogNormalRandomNumber(randomStdDevEntryTime, randomMeanEntryTime))
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

func (s *Stats) RecordIdleTime(idleTime float32) {
	s.IdleTimes = append(s.IdleTimes, idleTime)
}

func (s *Stats) RecordServiceTime(serviceTime float32) {
	s.ServiceTimes = append(s.ServiceTimes, serviceTime)
}

func (s *Stats) RecordShopTime(shopTime float32) {
	s.ShopTimes = append(s.ShopTimes, shopTime)
}

func (s *Stats) RecordWaitTime(waitTime float32) {
	s.WaitTimes = append(s.WaitTimes, waitTime)
}

func (s *Stats) Mean(values []float32) float32{
	if len(values) == 0 {
		return 0
	}else{
		var sum float32 = 0
		for _, value := range values { sum += value }
		return sum / float32(len(values))
	}
}

