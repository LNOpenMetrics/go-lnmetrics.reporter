package plugin

import (
	"time"

	"github.com/niftynei/glightning/glightning"
)

//TODO: export the interface in anther method
type Metric interface {
	MakePersistent() error
	Update(lightning glightning.Lightning) error
	UpdateWithMsg(message *Msg, lightning glightning.Lightning) error
}

//TODO move also in a common place
type Msg struct {
	// The message is from a command? if not it is nil
	cmd    *string
	params map[string]interface{}
}

type MetricOne struct {
	Metric
	nodeId       string `json:node_id`
	architecture string `json:architecture`
	// TODO: Here we need to store a object?
	// With an object we can add also some custom message
	upTime []time.Time `json:up_time`
	// TODO: missing the check to other channels
}

func New(nodeId string, architecture *string) *MetricOne {
	return nil
}

func (metric *MetricOne) Update(lightning glightning.Lightning) error {
	return nil
}

func (metric *MetricOne) UpdateWithMsg(message *Msg, lightning glightning.Lightning) error {
	return nil
}

func (metric *MetricOne) MakePersistent() error {
	return nil
}
