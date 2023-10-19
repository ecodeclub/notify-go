package ral

import (
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type Record struct {
	RspCode           int
	Protocol          string
	Method            string
	Url               string
	Port              string
	Host              string
	Error             any
	retry             int
	timeCostSpan      map[string]*StaticsItem
	timeCostSpanLock  sync.Mutex
	timeCostPoint     map[string]time.Duration
	timeCostPointLock sync.Mutex
	field             map[string]any
	fieldLock         sync.Mutex
}

type StaticsItem struct {
	StartPoint time.Time
	StopPoint  time.Time
}

func newLogRecord() Record {
	return Record{
		timeCostSpan:  make(map[string]*StaticsItem),
		timeCostPoint: make(map[string]time.Duration),
		field:         make(map[string]any),
	}
}

func NewLogRecord() Record {
	return newLogRecord()
}

func (lr *Record) PointStart(name string) {
	defer lr.timeCostSpanLock.Unlock()
	lr.timeCostSpanLock.Lock()
	item := new(StaticsItem)
	item.StartPoint = time.Now()
	lr.timeCostSpan[name] = item
}

func (lr *Record) PointStop(name string) {
	defer lr.timeCostSpanLock.Unlock()
	lr.timeCostSpanLock.Lock()
	item, ok := lr.timeCostSpan[name]
	if ok {
		item.StopPoint = time.Now()
	}
}

func (lr *Record) AddTimeCostPoint(name string, d time.Duration) {
	defer lr.timeCostPointLock.Unlock()
	lr.timeCostPointLock.Lock()
	lr.timeCostPoint[name] = d
}

func (lr *Record) AddField(name string, value any) {
	defer lr.fieldLock.Unlock()
	lr.fieldLock.Lock()
	lr.field[name] = value
}

func (s *StaticsItem) GetDuration() time.Duration {
	return s.StopPoint.Sub(s.StartPoint)
}

func (lr *Record) Flush(l *slog.Logger) {
	field := make([]any, 0, 16)
	field = append(field,
		"code", lr.RspCode,
		"path", lr.Url,
		"port", lr.Port,
		"host", lr.Host,
		"retry", lr.retry,
		"protocol", lr.Protocol,
		"method", lr.Method)

	for name, sItem := range lr.timeCostSpan {
		dura := sItem.GetDuration()
		field = append(field, slog.Duration(name, dura))
	}

	for name, f := range lr.field {
		field = append(field, slog.Any(name, f))
	}

	for name, d := range lr.timeCostPoint {
		field = append(field, slog.Duration(name, d))
	}

	if lr.Error != nil {
		l.Error(fmt.Sprintf("%v", lr.Error), field...)
	} else {
		l.Info("success", field...)
	}
}
