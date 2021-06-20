package plugin

import (
	"time"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
	"github.com/niftynei/glightning/glightning"
)

//TODO: export the interface in anther method
type Metric interface {
	// Call this method when the close rpc method is called
	OnClose() error
	// Call this method to make the status of the metrics persistent
	MakePersistent() error
	// Call this method when you want update all the metrics without
	// some particular event throw from c-lightning
	Update(lightning *glightning.Lightning) error
	// Class this method when you want catch some event from
	// c-lightning and make some operation on the metrics data.
	UpdateWithMsg(message *Msg, lightning *glightning.Lightning) error
}

//TODO move also in a common place
type Msg struct {
	// The message is from a command? if not it is nil
	cmd    string
	params map[string]interface{}
}

// Wrap all useful information
type status struct {
	//node_id  string    `json:node_id`
	timestamp time.Time `json:timestamp`
}

type MetricOne struct {
	Metric
	id           int
	name         string `json:metric_name`
	nodeId       string `json:node_id`
	architecture string `json:architecture`
	// TODO: Here we need to store a object?
	// With an object we can add also some custom message
	upTime []status `json:up_time`
	// TODO: missing the check to other channels
}

// This method is required by the
func New(nodeId string, architecture string) *MetricOne {
	return &MetricOne{id: 1, name: "metric_one", nodeId: nodeId,
		architecture: architecture, upTime: make([]status, 0)}
}

func (instance *MetricOne) Update(lightning *glightning.Lightning) error {
	log.GetInstance().Debug("On close event on metrics on called")
	instance.upTime = append(instance.upTime,
		status{timestamp: time.Now()})
	return nil
}

func (metric *MetricOne) UpdateWithMsg(message *Msg,
	lightning *glightning.Lightning) error {
	return nil
}

func (instance *MetricOne) MakePersistent() error {
	return db.GetInstance().PutValue(instance.name, instance)
}

//TODO: here the message is not useful, but we keep it
func (instance *MetricOne) OnClose(msg *Msg) error {
	log.GetInstance().Debug("On close event on metrics on called")
	instance.upTime = append(instance.upTime,
		status{timestamp: time.Now()})
	return instance.MakePersistent()
}
