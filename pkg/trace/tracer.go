package trace

import (
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	"github.com/vincenzopalazzo/cln4go/comm/tracer"
)

type Tracer struct{}

func (self *Tracer) Log(lebel tracer.TracerLevel, msg string) {}

func (self *Tracer) Logf(level tracer.TracerLevel, msg string, args ...any) {}

func (self *Tracer) Info(msg string) {
	log.GetInstance().Info(msg)
}

func (self *Tracer) Infof(msg string, args ...any) {
	log.GetInstance().Infof(msg, args...)
}

func (self *Tracer) Trace(msg string) {
	log.GetInstance().Error(msg)
}

func (self *Tracer) Tracef(msg string, args ...any) {
	log.GetInstance().Errorf(msg, args...)
}
