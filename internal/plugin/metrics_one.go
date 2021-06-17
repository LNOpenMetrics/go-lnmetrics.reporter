package plugin

import (
	"time"
)

//TODO: export the interface in anther method
type Metrics interface {
	MakePersistent() error
}

type MetricsOne struct {
	Metrics
	nodeId string
	upTime []time.Time
}

func (metrics *MetricsOne) MakePersistent() error {
	return nil
}
