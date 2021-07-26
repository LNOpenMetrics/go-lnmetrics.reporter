package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"
	"github.com/OpenLNMetrics/go-metrics-reported/pkg/log"

	sysinfo "github.com/elastic/go-sysinfo/types"
	"github.com/niftynei/glightning/glightning"
)

// Information about the Payment forward by the node
type PaymentInfo struct {
	// the payment is received by the channel or is sent to the channel
	Direction string `json:"direction"`
	// the status of the channels is completed, failed or pending
	Status string `json:"status"`
	// The motivation for the failure
	FailureReason string `json:"failure_reason"`
	// The code of the failure
	FailureCode int `json:"failure_code"`
}

// Only a wrapper to pass collected information about the channel
type ChannelInfo struct {
	NodeId    string
	Alias     string
	Color     string
	Direction string
	Forwards  []*PaymentInfo
}

// Wrap all useful information about the own node
type status struct {
	//node_id  string    `json:node_id`
	Event string `json:"event"`
	// how many channels the node have
	Channels int `json:"channels"`
	// how many payments it forwords
	Forwards *PaymentsSummary `json:"forwards"`
	// unix time where the check is made.
	Timestamp int64 `json:"timestamp"`
}

type channelStatus struct {
	Timestamp int64 `json:"timestamp"`
	// Status of the channel
	Status string `json:"status"`
}

// Wrap all the information about the node that the own node
// has some channel open.
type statusChannel struct {
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
	UpTimes []*channelStatus `json:"up_times"`
	// the node is ready to receive payment to share?
	Online bool `json:"online"`
	// last message (channel_update) received from the gossip
	LastUpdate uint `json:"last_update"`
	// the channel is public?
	Public bool `json:"public"`
	// information about the direction of the channel: out, in, mutual.
	Direction string `json:"direction"`
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

type MetricOne struct {
	// Internal id to identify the metric
	id int `json:"-"`
	// Name of the metrics
	Name   string  `json:"metric_name"`
	NodeId string  `json:"node_id"`
	Color  string  `json:"color"`
	OSInfo *osInfo `json:"os_info"`
	// timezone where the node is located
	Timezone string `json:"timezone"`
	// array of the up_time
	UpTime []*status `json:"up_time"`
	// map of informatonof channel information
	ChannelsInfo map[string]*statusChannel `json:"channels_info"`
}

// FIXME: Move this in a separate file like a generic metrics file
var MetricsSupported map[int]string

// FIXME: Move this in a utils file to give the possibility also to other metrics to access to it
// 0 = outcoming
// 1 = incoming
// 2 = mutual.
var ChannelDirections map[int]string

func init() {
	log.GetInstance().Debug("Init metrics map with all the name")
	MetricsSupported = make(map[int]string)
	MetricsSupported[1] = "metric_one"

	ChannelDirections = make(map[int]string)
	ChannelDirections[0] = "OUTCOMING"
	ChannelDirections[1] = "INCOOMING"
	ChannelDirections[2] = "MUTUAL"
}

// This method is required by the
func NewMetricOne(nodeId string, sysInfo sysinfo.HostInfo) *MetricOne {
	return &MetricOne{id: 1, Name: MetricsSupported[1], NodeId: nodeId,
		OSInfo: &osInfo{OS: sysInfo.OS.Name,
			Version:      sysInfo.OS.Version,
			Architecture: sysInfo.Architecture},
		Timezone: sysInfo.Timezone, UpTime: make([]*status, 0),
		ChannelsInfo: make(map[string]*statusChannel), Color: ""}
}

func (instance *MetricOne) OnInit(lightning *glightning.Lightning) error {
	getInfo, err := lightning.GetInfo()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during the OnInit method; %s", err))
		return err
	}

	instance.NodeId = getInfo.Id
	instance.Color = getInfo.Color

	log.GetInstance().Debug("Init event on metrics on called")
	listFunds, err := lightning.ListFunds()
	log.GetInstance().Debug(fmt.Sprintf("%s", listFunds))
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	instance.collectInfoChannels(lightning, listFunds.Channels)

	listForwards, err := lightning.ListForwards()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	statusPayments, err := instance.makePaymentsSummary(lightning, listForwards)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	instance.UpTime = append(instance.UpTime,
		&status{Event: "on_start",
			Timestamp: time.Now().Unix(),
			Channels:  len(listFunds.Channels),
			Forwards:  statusPayments})
	return instance.MakePersistent()

	return nil
}

func (instance *MetricOne) Update(lightning *glightning.Lightning) error {
	log.GetInstance().Debug("Update event on metrics on called")
	listFunds, err := lightning.ListFunds()
	log.GetInstance().Debug(fmt.Sprintf("%s", listFunds))
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	instance.collectInfoChannels(lightning, listFunds.Channels)

	listForwards, err := lightning.ListForwards()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	// TODO: feel payment status
	statusPayments, err := instance.makePaymentsSummary(lightning, listForwards)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}
	instance.UpTime = append(instance.UpTime,
		&status{Event: "on_update",
			Timestamp: time.Now().Unix(),
			Channels:  len(listFunds.Channels),
			Forwards:  statusPayments,
		})
	return instance.MakePersistent()
}

func (metric *MetricOne) UpdateWithMsg(message *Msg,
	lightning *glightning.Lightning) error {
	return nil
}

func (instance *MetricOne) MakePersistent() error {
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
	forwards := &PaymentsSummary{0, 0}
	if len(instance.UpTime) > 0 {
		lastValue = instance.UpTime[len(instance.UpTime)-1].Channels
		forwards = instance.UpTime[len(instance.UpTime)-1].Forwards
	}
	instance.UpTime = append(instance.UpTime,
		&status{Event: "on_close",
			Timestamp: time.Now().Unix(),
			Channels:  lastValue, Forwards: forwards})
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

func (instance *MetricOne) makePaymentsSummary(lightning *glightning.Lightning, forwards []glightning.Forwarding) (*PaymentsSummary, error) {
	statusPayments := PaymentsSummary{Completed: 0, Failed: 0}

	for _, forward := range forwards {
		switch forward.Status {
		case "settled":
			statusPayments.Completed++
		case "failed", "local_failed":
			statusPayments.Failed++
		default:
			return nil, errors.New(fmt.Sprintf("Status %s unexpected", forward.Status))
		}
	}

	return &statusPayments, nil
}

// private method of the module
func (instance *MetricOne) collectInfoChannels(lightning *glightning.Lightning, channels []*glightning.FundingChannel) error {
	cache := make(map[string]bool)
	for _, channel := range channels {
		err := instance.collectInfoChannel(lightning, channel)
		if err != nil {
			// void returning error here? We can continue to make the analysis over the channels
			log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			return err
		}
		err = instance.collectInfoChannel(lightning, channel)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			return nil
		}
		cache[channel.ShortChannelId] = true
	}

	// make intersection of the channels in the cache and a
	// channels in the metrics plugin
	// this is useful to remove the metrics over closed channels
	// in the metrics one we are not interested to have a story of
	// of the old channels (for the moments).
	for key, _ := range instance.ChannelsInfo {
		_, found := cache[key]
		if !found {
			delete(instance.ChannelsInfo, key)
		}
	}

	return nil
}

func (instance *MetricOne) collectInfoChannel(lightning *glightning.Lightning, channel *glightning.FundingChannel) error {

	shortChannelId := channel.ShortChannelId
	infoChannel, found := instance.ChannelsInfo[shortChannelId]
	var timestamp int64 = 0
	if instance.pingNode(lightning, channel.Id) {
		timestamp = time.Now().Unix()
	}
	info, err := instance.getChannelInfo(lightning, channel)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during get the information about the channel: %s", err))
		return err
	}

	// A new channels found
	channelStat := channelStatus{timestamp, channel.State}
	if !found {
		upTimes := make([]*channelStatus, 1)
		upTimes[0] = &channelStat
		// TODO: Could be good to have a information about the direction of the channel
		newInfoChannel := statusChannel{NodeId: info.NodeId, NodeAlias: info.Alias, Color: info.Color,
			Capacity: channel.ChannelSatoshi, Forwards: info.Forwards,
			UpTimes: upTimes, Online: channel.Connected}
		instance.ChannelsInfo[shortChannelId] = &newInfoChannel
	} else {
		infoChannel.Capacity = channel.ChannelSatoshi
		infoChannel.UpTimes = append(infoChannel.UpTimes, &channelStat)
		infoChannel.Color = info.Color
		infoChannel.Online = channel.Connected
	}
	return nil
}

func (instance *MetricOne) pingNode(lightning *glightning.Lightning, nodeId string) bool {
	_, err := lightning.Ping(nodeId)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during pinging node: %s", err))
		return false
	}
	return true
}

func (instance *MetricOne) getChannelInfo(lightning *glightning.Lightning, channel *glightning.FundingChannel) (*ChannelInfo, error) {

	nodeInfo, err := lightning.GetNode(channel.Id)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during the call listNodes: %s", err))
		return nil, err
	}

	channelInfo := ChannelInfo{NodeId: channel.Id, Alias: nodeInfo.Alias, Color: nodeInfo.Color, Direction: "unknown"}

	listForwards, err := lightning.ListForwards()

	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during the listForwards call: %s", err))
		return nil, err
	}

	for _, forward := range listForwards {
		if channel.ShortChannelId == forward.InChannel {
			channelInfo.Forwards = append(channelInfo.Forwards, &PaymentInfo{Direction: ChannelDirections[1], Status: forward.Status})
		} else if channel.ShortChannelId == forward.OutChannel {
			channelInfo.Forwards = append(channelInfo.Forwards, &PaymentInfo{Direction: ChannelDirections[0], Status: forward.Status})
		}

		switch forward.Status {
		case "settled", "failed":
			// do nothings
			continue
		case "local_failed":
			// store the information about the failure
			paymentInfo := channelInfo.Forwards[len(channelInfo.Forwards)-1]
			paymentInfo.FailureReason = forward.FailReason
			paymentInfo.FailureCode = forward.FailCode
		default:
			return nil, errors.New(fmt.Sprintf("Status %s unexpected", forward.Status))
		}
	}
	//TODO Adding support for the dual founding channels.
	return &channelInfo, nil
}
