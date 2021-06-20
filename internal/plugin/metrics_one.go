package plugin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
	"github.com/niftynei/glightning/glightning"
)

//TODO: export the interface in anther method
type Metric interface {
	// Call this method when the close rpc method is called
	OnClose(msg *Msg, lightning *glightning.Lightning) error
	// Call this method to make the status of the metrics persistent
	MakePersistent() error
	// Call this method when you want update all the metrics without
	// some particular event throw from c-lightning
	Update(lightning *glightning.Lightning) error
	// Class this method when you want catch some event from
	// c-lightning and make some operation on the metrics data.
	UpdateWithMsg(message *Msg, lightning *glightning.Lightning) error

	// convert the object into a json
	ToJSON() (string, error)
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
	Channels  int   `json:"channels"`
	Timestamp int64 `json:"timestamp"`
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

// This method is required by the
func NewMetricOne(nodeId string, architecture string) *MetricOne {
	return &MetricOne{id: 1, Name: "metric_one", NodeId: nodeId,
		Architecture: architecture, UpTime: make([]status, 0)}
}

func (instance *MetricOne) Update(lightning *glightning.Lightning) error {
	log.GetInstance().Debug("Update event on metrics on called")
	/*listFunds, err := lightning.ListFunds()
	log.GetInstance().Debug(fmt.Sprintf("%s", listFunds))
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}*/
	instance.UpTime = append(instance.UpTime,
		status{Timestamp: time.Now().Unix(), Channels: 0})
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

//TODO: here the message is not useful, but we keep it
func (instance *MetricOne) OnClose(msg *Msg, lightning *glightning.Lightning) error {
	log.GetInstance().Debug("On close event on metrics called")
	/*listFunds, err := lightning.ListFunds()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	} */
	instance.UpTime = append(instance.UpTime,
		status{Timestamp: time.Now().Unix(), Channels: 0})
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
