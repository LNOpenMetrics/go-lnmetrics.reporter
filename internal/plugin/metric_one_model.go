package plugin

import (
	"encoding/json"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"strings"
)

// PaymentInfo Information about the Payment forward by the node
type PaymentInfo struct {
	// the payment is received by the channel or is sent to the channel
	Direction string `json:"direction"`
	// the status of the channels is completed, failed or pending
	Status string `json:"status"`
	// The motivation for the failure
	FailureReason string `json:"failure_reason,omitempty"`
	// The code of the failure
	FailureCode int `json:"failure_code,omitempty"`
	// instance where the payment is started
	Timestamp int64 `json:"timestamp"`
}

// ChannelInfo Only a wrapper to pass collected information about the channel
// not used inside the metrics.
type ChannelInfo struct {
	NodeId     string
	Alias      string
	Color      string
	Direction  string
	LastUpdate uint
	Forwards   []*PaymentInfo
	// information about the channels fee
	Fee *ChannelFee `json:"fee"`
	// HTLC limit of the node where we have a channel with
	Limits *ChannelLimits `json:"limits"`
}

type ChannelSummary struct {
	NodeId    string `json:"node_id"`
	Alias     string `json:"alias"`
	Color     string `json:"color"`
	ChannelId string `json:"channel_id"`
	State     string `json:"state"`
}

type ChannelsSummary struct {
	TotChannels uint64            `json:"tot_channels"`
	Summary     []*ChannelSummary `json:"summary"`
}

// Wrap all useful information about the own node
type status struct {
	//node_id  string    `json:node_id`
	Event string `json:"event"`
	// how many channels the node have
	Channels *ChannelsSummary `json:"channels"`
	// how many payments it forwords
	Forwards *PaymentsSummary `json:"forwards"`
	// unix time where the check is made.
	Timestamp int64 `json:"timestamp"`
	// Node fee settings
	Fee *ChannelFee `json:"fee"`
	// Node htlc limits information
	Limits *ChannelLimits `json:"limits"`
}

type channelStatus struct {
	// the event that originate this status check
	Event string `json:"event"`
	// Timestamp when the check is made
	Timestamp int64 `json:"timestamp"`
	// Status of the channel
	Status string `json:"status"`
}

// ChannelLimits Container of the htlc limit information
// of the node where we have a channel with.
//
// It is used in the statusChannel struct
type ChannelLimits struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// ChannelFee Container of the fee information related
// to the channel, used in statusChannel struct
type ChannelFee struct {
	Base    uint64 `json:"base"`
	PerMSat uint64 `json:"per_msat"`
}

// Wrap all the information about the node that the own node
// has some channel open.
//
// Used in the MetricOne struct
type statusChannel struct {
	// short channel id
	ChannelId string `json:"channel_id"`
	// node id
	NodeId string `json:"node_id"`
	// label of the node
	NodeAlias string `json:"node_alias"`
	// Color of the node
	Color string `json:"color"`
	// the capacity of the channel
	Capacity uint64 `json:"capacity"`
	// how payment the channel forwords
	Forwards []*PaymentInfo `json:"forwards"`
	// The node answer from the ping operation
	UpTimes []*channelStatus `json:"up_time"`
	// the node is ready to receive payment to share?
	Online bool `json:"online"`
	// last message (channel_update) received from the gossip
	LastUpdate uint `json:"last_update"`
	// information about the direction of the channel: out, in, mutual.
	Direction string `json:"direction"`
	// information about the channels fee
	Fee *ChannelFee `json:"fee"`
	// HTLC limit of the node where we have a channel with
	Limits *ChannelLimits `json:"limits"`
}

type osInfo struct {
	// Operating system name
	OS string `json:"os"`
	// Version of the Operating System
	Version string `json:"version"`
	// architecture of the system where the node is running
	Architecture string `json:"architecture"`
}

type PaymentsSummary struct {
	Completed uint64 `json:"completed"`
	Failed    uint64 `json:"failed"`
}

// NodeInfo Contains the info about the ln node.
type NodeInfo struct {
	Implementation string `json:"implementation"`
	Version        string `json:"version"`
}

type NodeAddress struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port uint   `json:"port"`
}

// MetricOne Main data structure that it is filled by the collection data phase.
type MetricOne struct {
	// Internal id to identify the metric
	id int `json:"-"`

	// Version of metrics format, it is used to migrate the
	// JSON payload from previous version of plugin.
	Version int `json:"version"`

	// Name of the metrics
	Name string `json:"metric_name"`

	// Public Key of the Node
	NodeID string `json:"node_id"`

	// Node Alias on the network
	NodeAlias string `json:"node_alias"`

	// Color of the node
	Color string `json:"color"`

	// Network where the node it is running
	Network string `json:"network"`

	// OS host information
	OSInfo *osInfo `json:"os_info"`

	// Node information, like version/implementation
	NodeInfo *NodeInfo `json:"node_info"`

	// Node address, where the node will be reachable by other node
	Address []*NodeAddress `json:"address"`

	// timezone where the node is located
	Timezone string `json:"timezone"`

	// array of the up_time
	UpTime []*status `json:"up_time"`

	// map of information of channel information
	ChannelsInfo map[string]*statusChannel `json:"-"`

	// Last check of the plugin, useful to store the data
	// in the db by timestamp
	lastCheck int64 `json:"-"`

	// Storage reference
	Storage db.PluginDatabase `json:"-"`
}

func (instance MetricOne) MarshalJSON() ([]byte, error) {
	// Declare a new type using the definition of MetricOne,
	// the result of this is that M will have the same structure
	// as MetricOne but none of its methods (this avoids recursive
	// calls to MarshalJSON).
	//
	// Also, because M and MetricOne have the same structure you can
	// easily convert between those two. e.g. M(MetricOne{}) and
	// MetricOne(M{}) are valid expressions.
	type M MetricOne

	// Declare a new type that has a field of the "desired" type and
	// also **embeds** the M type. Embedding promotes M's fields to T
	// and encoding/json will marshal those fields unnested/flattened,
	// i.e. at the same level as the channels_info field.
	type T struct {
		M
		ChannelsInfo []*statusChannel `json:"channels_info"`
	}

	// move map elements to slice
	channels := make([]*statusChannel, 0, len(instance.ChannelsInfo))
	for _, channel := range instance.ChannelsInfo {
		channels = append(channels, channel)
	}

	// Pass in an instance of the new type T to json.Marshal.
	// For the embedded M field use a converted instance of the receiver.
	// For the ChannelsInfo field use the channels slice.
	return json.Marshal(T{
		M:            M(instance),
		ChannelsInfo: channels,
	})
}

// UnmarshalJSON Same as MarshalJSON but in reverse.
func (instance *MetricOne) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]any
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return err
	}

	if err := instance.Migrate(jsonMap); err != nil {
		return err
	}

	data, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}
	type M MetricOne
	type T struct {
		*M
		ChannelsInfo []*statusChannel `json:"channels_info"`
	}
	t := T{M: (*M)(instance)}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	instance.ChannelsInfo = make(map[string]*statusChannel, len(t.ChannelsInfo))
	for _, channel := range t.ChannelsInfo {
		key := strings.Join([]string{channel.ChannelId, channel.Direction}, "_")
		instance.ChannelsInfo[key] = channel
	}

	return nil
}
