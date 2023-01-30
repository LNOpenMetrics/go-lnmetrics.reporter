// Package plugin implement all the necessary building blocks to implement
// an open source metrics.
package metrics

import (
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/vincenzopalazzo/cln4go/client"
)

const (
	RawLocalScoreID int = 1
)

// MetricsSupported mapping the internal id with the name of the metrics.
// the id is passed by the plugin RPC name.
var MetricsSupported map[int]string

// 0 = outcoming
// 1 = incoming
var ChannelDirections map[int]string

func init() {
	MetricsSupported = make(map[int]string)
	MetricsSupported[1] = "metric_one"

	ChannelDirections = make(map[int]string)
	ChannelDirections[0] = "OUTCOMING"
	ChannelDirections[1] = "INCOOMING"
}

// Metric All the metrics need to respect this interface
type Metric interface {
	// MetricName return the name of the metric
	MetricName() *string

	// OnInit initialize the method with node information
	OnInit(lightning client.Client) error

	// OnStop commit the actual information before exit
	OnStop(msg *Msg, lightning client.Client) error

	// MakePersistent make the metric persistent
	MakePersistent() error

	// UploadOnRepo Commit the metric on remote server
	UploadOnRepo(client *graphql.Client, lightning client.Client) error

	// InitOnRepo Init metric on the remote server.
	InitOnRepo(client *graphql.Client, lightning client.Client) error

	// Update the metric with the last information of the node
	Update(lightning client.Client) error

	// UpdateWithMsg update the metric with the last information fo the node with some msg info
	UpdateWithMsg(message *Msg, lightning client.Client) error

	// ToMap convert the Metric object into a map
	ToMap() (map[string]any, error)

	// ToJSON convert the object into a json
	ToJSON() (string, error)

	// Migrate to a new version of the metrics, some new version of the plugin
	// ca contains some format change and for this reason, this method help
	// to migrate the old version to a new version.
	Migrate(payload map[string]any) error
}

// Msg Message struct to pass from the plugin to the metric
type Msg struct {
	// The message is from a command? if not it is nil
	cmd string
	// the map of parameter that the plugin need to feel.
	params map[string]any
}

func NewMsg(cmd string, params map[string]any) Msg {
	return Msg{
		cmd:    cmd,
		params: params,
	}
}
