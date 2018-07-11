package main

import (
	"math/rand"
	"time"
	"math"
)

const randomStdDevServiceTime float64 = 1
const randomMeanServiceTime float64 = 2
const randomStdDevShopTime float64 = 1
const randomMeanShopTime float64 = 8
const randomStdDevEntryTime float64 = 10
const randomMeanEntryTime float64 = 30

type Stats struct{
	EntityCount int
	CumulativeServiceTime float32
	CumulativeShopTime float32
	CompletedJobCount int
	JobCount int
	GlobalTime float32
	WaitTimes []float32
}

var source = rand.NewSource(time.Now().Unix())

func NewStats() Stats{
	return Stats{
		0, 0, 0, 0,
		0, 0, make([]float32, 0),
	}
}

func (*Stats) generateGaussianRandomNumber(stdDev float64, mean float64) float32{
	rnd := rand.New(source)
	return float32(math.Abs(rnd.NormFloat64() * stdDev + mean))
}

func (s *Stats) RecordWaitTime(waitTime float32) {
	s.WaitTimes = append(s.WaitTimes, waitTime)
}

func (s *Stats) generateServiceTime() float32{
	return s.generateGaussianRandomNumber(randomStdDevServiceTime, randomMeanServiceTime)
}

func (s *Stats) generateShopTime() float32{
	return s.generateGaussianRandomNumber(randomStdDevShopTime, randomMeanShopTime)
}

func (s *Stats) generateEntryTime() float32{
	entryTime := s.generateGaussianRandomNumber(randomStdDevEntryTime, randomMeanEntryTime)
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