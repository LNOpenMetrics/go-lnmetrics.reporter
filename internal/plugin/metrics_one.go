package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
	"github.com/niftynei/glightning/glightning"
)

// Wrap all useful information
type status struct {
	//node_id  string    `json:node_id`
	Event     string `json:"event"`
	Channels  int    `json:"channels"`
	Forwords  int    `json:"forwords"`
	Timestamp int64  `json:"timestamp"`
}

type MetricOne struct {
	id           int    `json:"-"`
	Name         string `json:"metric_name"`
	NodeId       string `json:"node_id"`
	Architecture string `json:"architecture"`
	// TODO: Here we need to store a object?
	// With an object we can add also some custom message
	UpTime []status `json:"up_time"`
	// TODO: missing the check to other channels
}

var MetricsSupported map[int]string

func init() {
	log.GetInstance().Debug("Init metrics map with all the name")
	MetricsSupported = make(map[int]string)
	MetricsSupported[1] = "metric_one"
}

// This method is required by the
func NewMetricOne(nodeId string, architecture string) *MetricOne {
	return &MetricOne{id: 1, Name: MetricsSupported[1], NodeId: nodeId,
		Architecture: architecture, UpTime: make([]status, 0)}
}

func (instance *MetricOne) Update(lightning *glightning.Lightning) error {
	log.GetInstance().Debug("Update event on metrics on called")
	listFunds, err := lightning.ListFunds()
	log.GetInstance().Debug(fmt.Sprintf("%s", listFunds))
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	listForwords, err := lightning.ListForwards()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	instance.UpTime = append(instance.UpTime,
		status{Event: "on_update",
			Timestamp: time.Now().Unix(),
			Channels:  len(listFunds.Channels),
			Forwords:  len(listForwords)})
	return instance.MakePersistent()
}

func (metric *MetricOne) UpdateWithMsg(message *Msg,
	lightning *glightning.Lightning) error {
	return nil
}

func (instance *MetricOne) MakePersistent() error {
	log.GetInstance().Debug(fmt.Sprintf("%s", instance))
	json, err := instance.ToJSON()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("JSON error %s", err))
		return err
	}
	return db.GetInstance().PutValue(instance.Name, json)
}

// here the message is not useful, but we keep it only for future evolution
// or we will remove it from here.
func (instance *MetricOne) OnClose(msg *Msg, lightning *glightning.Lightning) error {
	log.GetInstance().Debug("On close event on metrics called")
	lastValue := 0
	sizeForwords := 0
	if len(instance.UpTime) > 0 {
		lastValue = instance.UpTime[len(instance.UpTime)-1].Channels
		sizeForwords = instance.UpTime[len(instance.UpTime)-1].Forwords
	}
	instance.UpTime = append(instance.UpTime,
		status{Event: "on_close",
			Timestamp: time.Now().Unix(),
			Channels:  lastValue, Forwords: sizeForwords})
	return instance.MakePersistent()
}

func (instance *MetricOne) ToJSON() (string, error) {
	json, err := json.Marshal(&instance)
	if err != nil {
		log.GetInstance().Error(err)
		return "", err
	}
	return string(json), nil
}
