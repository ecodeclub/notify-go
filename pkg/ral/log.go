package ral

import (
	"fmt"
	"github.com/ecodeclub/notify-go/pkg/logger"
	"sync"
	"time"
)

type LogRecord struct {
	LogId             int64
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

func newLogRecord() LogRecord {
	return LogRecord{
		timeCostSpan:  make(map[string]*StaticsItem),
		timeCostPoint: make(map[string]time.Duration),
		field:         make(map[string]any),
	}
}

func NewLogRecord() LogRecord {
	return newLogRecord()
}

func NewLogRecordWithID(id int64) LogRecord {
	l := newLogRecord()
	l.LogId = id
	return l
}

func (lr *LogRecord) PointStart(name string) {
	defer lr.timeCostSpanLock.Unlock()
	lr.timeCostSpanLock.Lock()
	item := new(StaticsItem)
	item.StartPoint = time.Now()
	lr.timeCostSpan[name] = item
}

func (lr *LogRecord) PointStop(name string) {
	defer lr.timeCostSpanLock.Unlock()
	lr.timeCostSpanLock.Lock()
	item, ok := lr.timeCostSpan[name]
	if ok {
		item.StopPoint = time.Now()
	}
}

func (lr *LogRecord) AddTimeCostPoint(name string, d time.Duration) {
	defer lr.timeCostPointLock.Unlock()
	lr.timeCostPointLock.Lock()
	lr.timeCostPoint[name] = d
}

func (lr *LogRecord) AddField(name string, value any) {
	defer lr.fieldLock.Unlock()
	lr.fieldLock.Lock()
	lr.field[name] = value
}

func (s *StaticsItem) GetDuration() time.Duration {
	return s.StopPoint.Sub(s.StartPoint)
}

func (lr *LogRecord) Flush() {
	field := make([]logger.Field, 0, 16)
	field = append(field,
		logger.Int("code", lr.RspCode), logger.String("path", lr.Url),
		logger.String("port", lr.Port), logger.String("host", lr.Host),
		logger.Int("retry", lr.retry), logger.String("protocol", lr.Protocol),
		logger.String("method", lr.Method))

	for name, sItem := range lr.timeCostSpan {
		dura := sItem.GetDuration()
		field = append(field, logger.Duration(name, dura))
	}

	for name, f := range lr.field {
		field = append(field, logger.Any(name, f))
	}

	for name, d := range lr.timeCostPoint {
		field = append(field, logger.Duration(name, d))
	}

	if lr.Error != nil {
		logger.Default().Error(fmt.Sprintf("%v", lr.Error), field...)
	} else {
		logger.Default().Info("", field...)
	}
}
