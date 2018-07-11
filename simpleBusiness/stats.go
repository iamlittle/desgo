package main

import (
	"math/rand"
	"time"
)

const randomStdDevServiceTime float32 = 2
const randomMeanServiceTime float32 = 8
const randomStdDevShopTime float32 = 2
const randomMeanShopTime float32 = 8

type Stats struct{
	EntityCount int
	CumulativeServiceTime float32
	CumulativeShopTime float32
	CumulativeWaitTime float32
	CompletedJobCount int
	JobCount int
	GlobalTime float32
}

func (*Stats) generateGaussianRandomNumber(stdDev float32, mean float32) float32{
	source := rand.NewSource(time.Now().Unix())
	rnd := rand.New(source)
	return float32(rnd.NormFloat64()) * stdDev + mean
}

func (s *Stats) RecordWaitTime(waitTime float32) {
	if s.JobCount >= 100 {
		if s.CompletedJobCount >= 50 {
			s.CumulativeWaitTime += waitTime
		}
	}else {
			s.CumulativeWaitTime += waitTime
	}
}

func (s *Stats) generateServiceTime() float32{
	return s.generateGaussianRandomNumber(randomStdDevServiceTime, randomMeanServiceTime)
}

func (s *Stats) generateShopTime() float32{
	return s.generateGaussianRandomNumber(randomStdDevShopTime, randomMeanShopTime)
}

func (s *Stats) generateEntityId() int{
	id := s.EntityCount
	s.EntityCount++
	return id
}