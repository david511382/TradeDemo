package flow

import (
	"time"
	errUtil "zerologix-homework/src/pkg/util/error"
)

type IStep interface {
	Name() string
	Run() (records StepRecords, resultErrInfo errUtil.IError)
}

type IDeferStep interface {
	HasDefer() bool
	DeferRun(resultErrInfo errUtil.IError) (records StepRecords)
}

type Step struct {
	StepName string
	DeferFun func(resultErrInfo errUtil.IError)
	Fun      func() (resultErrInfo errUtil.IError)
	StepsFun func() Steps
}

func (s Step) Name() string {
	return s.StepName
}

func (s Step) Run() (records StepRecords, resultErrInfo errUtil.IError) {
	if s.Fun == nil {
		return s.StepsFun().SetFeatureName(s.Name()).Run()
	}
	startTime := time.Now()
	resultErrInfo = s.Fun()
	durationTime := time.Since(startTime)
	record := &StepRecord{
		Name:         s.Name(),
		DurationTime: durationTime,
	}
	records = append(records, record)
	return
}

func (s Step) HasDefer() bool {
	return s.DeferFun != nil
}

func (s Step) DeferRun(resultErrInfo errUtil.IError) (records StepRecords) {
	startTime := time.Now()
	s.DeferFun(resultErrInfo)
	durationTime := time.Since(startTime)
	record := &StepRecord{
		Name:         s.Name(),
		DurationTime: durationTime,
	}
	records = append(records, record)
	return
}

type Steps []IStep

func (ss Steps) SetFeatureName(name string) flowType {
	return Flow(name, ss...)
}
