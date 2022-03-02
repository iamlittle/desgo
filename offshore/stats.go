package main

import (
	"math"
	"math/rand"
	"time"
)

type Stats struct {
	EntityCount       int
	CompletedJobCount int
	JobCount          int
	GlobalTime        float64
	CodeMigratedTimes []float64
	ReviewTimes       []float64
	ConversionTimes   []float64
	UnitTestTimes     []float64
	ValidatedTimes    []float64
	CutoverTimes      []float64
	StatsConfig       map[string]*StatsConfig
	WarmedUp          bool
	Source            rand.Source
}

func NewStats(configs []StatsConfig) Stats {
	m := make(map[string]*StatsConfig, len(configs))
	for i := 0; i < len(configs); i++ {
		m[configs[i].Name] = &configs[i]
	}
	return Stats{
		0, 0, 0, 0,
		make([]float64, 0), make([]float64, 0), make([]float64, 0),
		make([]float64, 0), make([]float64, 0), make([]float64, 0),
		m, false, rand.NewSource(time.Now().UnixNano()),
	}
}

func (s *Stats) generateGaussianRandomNumber(stdDev float64, mean float64) float64 {
	rnd := rand.New(s.Source)
	return rnd.NormFloat64()*stdDev + mean
}

func (s *Stats) generateLogNormalRandomNumber(stdDev float64, mean float64) float64 {
	//ensure the value is positive
	mu := s.Mu(math.Pow(stdDev, 2), mean)
	sigma := s.Sigma(math.Pow(stdDev, 2), mean)
	gaussian := math.Abs(s.generateGaussianRandomNumber(sigma, mu))
	return math.Exp(gaussian)
}

func (s *Stats) generateExponentialRandomNumber(stdDev float64, mean float64) float64 {
	rnd := rand.New(s.Source)
	return rnd.ExpFloat64()*stdDev + (mean - stdDev)
}

func (s *Stats) generateCodeMigratedDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].CodeMigratedStdDev
	mean := s.StatsConfig[string(componentType)].CodeMigratedMean
	return math.Abs(s.generateExponentialRandomNumber(std_dev, mean))
}

func (s *Stats) generateReviewDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].ReviewStdDev
	mean := s.StatsConfig[string(componentType)].ReviewMean
	return s.generateExponentialRandomNumber(std_dev, mean)
}

func (s *Stats) generateConvertDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].ConvertStdDev
	mean := s.StatsConfig[string(componentType)].ConvertMean
	return s.generateExponentialRandomNumber(std_dev, mean)
}

func (s *Stats) generateUnitTestDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].UnitTestStdDev
	mean := s.StatsConfig[string(componentType)].UnitTestMean
	return s.generateExponentialRandomNumber(std_dev, mean)
}

func (s *Stats) generateValidateDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].ValidateStdDev
	mean := s.StatsConfig[string(componentType)].ValidateMean
	return s.generateExponentialRandomNumber(std_dev, mean)
}

func (s *Stats) generateCutoverDuration(componentType ComponentType) float64 {
	std_dev := s.StatsConfig[string(componentType)].CutoverStdDev
	mean := s.StatsConfig[string(componentType)].CutoverMean
	return s.generateExponentialRandomNumber(std_dev, mean)
}

// func (s *Stats) generateEntryTime() float64 {
// 	entryTime := s.generateLogNormalRandomNumber(s.StatsConfig.CutoverStdDev, s.StatsConfig.CutoverStdDev)
// 	if entryTime < s.GlobalTime {
// 		entryTime += s.GlobalTime
// 	}
// 	return entryTime
// }

func (s *Stats) generateEntityId() int {
	id := s.EntityCount
	s.EntityCount++
	return id
}

func (s *Stats) RecordComponentCodeMigrationTime(component *Component) {
	if s.WarmedUp {
		s.CodeMigratedTimes = append(s.CodeMigratedTimes, component.CodeMigrateDuration)
	}
}

func (s *Stats) RecordComponentReviewTime(component *Component) {
	if s.WarmedUp {
		s.ReviewTimes = append(s.ReviewTimes, component.ReviewDuration)
	}
}

func (s *Stats) RecordComponentRoughConversionTime(component *Component) {
	if s.WarmedUp {
		s.ConversionTimes = append(s.ConversionTimes, component.ConvertDuration)
	}
}

func (s *Stats) RecordComponentUnitTestTime(component *Component) {
	if s.WarmedUp {
		s.UnitTestTimes = append(s.UnitTestTimes, component.UnitTestDuration)
	}
}

func (s *Stats) RecordComponentValidateTime(component *Component) {
	if s.WarmedUp {
		s.ValidatedTimes = append(s.ValidatedTimes, component.ValidateDuration)
	}
}

func (s *Stats) RecordComponentCutoverTime(component *Component) {
	if s.WarmedUp {
		s.CutoverTimes = append(s.CutoverTimes, component.CutoverDuration)
	}
}

func (s *Stats) Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	} else {
		var sum float64 = 0
		for _, value := range values {
			sum += value
		}
		return sum / float64(len(values))
	}
}

func (s *Stats) Variance(mean float64, values []float64) float64 {
	if len(values) == 0 {
		return 0
	} else {
		var variance float64 = 0
		for _, value := range values {
			variance += math.Pow(value-mean, 2)
		}
		return variance
	}
}

func (s *Stats) StdDev(mean float64, values []float64) float64 {
	return math.Sqrt(s.Variance(mean, values) / float64(len(values)))
}

func (s *Stats) Mu(variance float64, mean float64) float64 {
	return math.Log(math.Pow(mean, 2) / math.Sqrt(variance+math.Pow(mean, 2)))
}

func (s *Stats) Sigma(variance float64, mean float64) float64 {
	return math.Sqrt(math.Log(1 + variance/math.Pow(mean, 2)))
}
