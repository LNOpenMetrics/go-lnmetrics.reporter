package plugin

import (
	"github.com/niftynei/glightning/glightning"
)

// All the metrics need to respect this interface
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

// Message struct to pass from the plugin to the metric
type Msg struct {
	// The message is from a command? if not it is nil
	cmd string
	// the map of parameter that the plugin need to feel.
	params map[string]interface{}
}
