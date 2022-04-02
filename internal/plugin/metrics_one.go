package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/cache"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/db"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"

	"github.com/LNOpenMetrics/lnmetrics.utils/hash/sha256"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
	"github.com/LNOpenMetrics/lnmetrics.utils/utime"

	sysinfo "github.com/elastic/go-sysinfo/types"
	"github.com/vincenzopalazzo/glightning/glightning"
)

// Information about the Payment forward by the node
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

// Only a wrapper to pass collected information about the channel
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
	// Node htlc limits informations
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

// Container of the htlc limit information
// of the node where we have a channel with.
//
// It is used in the statusChannel struct
type ChannelLimits struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

// Container of the fee information related
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

	// map of informaton of channel information
	// TODO: managing the dualfunding channels
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
	// Also because M and MetricOne have the same structure you can
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
	var jsonMap map[string]interface{}
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

// NewMetricOne This method is required by the interface
func NewMetricOne(nodeId string, sysInfo sysinfo.HostInfo, storage db.PluginDatabase) *MetricOne {
	return &MetricOne{
		id:        1,
		Version:   4,
		Name:      MetricsSupported[1],
		NodeID:    nodeId,
		NodeAlias: "unknown",
		Network:   "unknown",
		OSInfo: &osInfo{
			OS:           sysInfo.OS.Name,
			Version:      sysInfo.OS.Version,
			Architecture: sysInfo.Architecture,
		},
		NodeInfo: &NodeInfo{
			Implementation: "unknown",
			Version:        "unknown",
		},
		Address:      make([]*NodeAddress, 0),
		Timezone:     sysInfo.Timezone,
		UpTime:       make([]*status, 0),
		ChannelsInfo: make(map[string]*statusChannel),
		Color:        "",
		Storage:      storage,
	}
}

func (instance *MetricOne) MetricName() *string {
	metricName := MetricsSupported[1]
	return &metricName
}

// Migrate from a payload format to another, with the help of the version number.
// Note that it is required implementing a required strategy only if the same
// properties will change during the time, if something's it is only add, we
// don't have anything's to migrate.
func (instance *MetricOne) Migrate(payload map[string]interface{}) error {
	version, found := payload["version"]

	if !found || int(version.(float64)) < 1 {
		log.GetInstance().Info("Migrate channels_info from version 0 to version 1")
		channelsInfoMap, found := payload["channels_info"]
		if !found {
			log.GetInstance().Error("Error: channels_info is not in the payload for migration")
			return errors.New("Error: channels_info is not in the payload for migration")
		}
		if reflect.ValueOf(channelsInfoMap).Kind() == reflect.Map {
			channelsInfoList := make([]interface{}, 0)
			for _, value := range channelsInfoMap.(map[string]interface{}) {
				channelsInfoList = append(channelsInfoList, value)
			}
			payload["channels_info"] = channelsInfoList
			payload["version"] = 1
		}
	}
	payload["version"] = 4
	return nil
}

// Generic Plugin callback that it is run each time that the plugin need to recording a new event.
func (instance *MetricOne) onEvent(nameEvent string, lightning *glightning.Lightning) (*status, error) {
	listFunds, err := lightning.ListFunds()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, err
	}
	if err := instance.collectInfoChannels(lightning, listFunds.Channels, nameEvent); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		// We admit this error here, we print only some log information.
	}

	listForwards, err := lightning.ListForwards()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, err
	}
	statusPayments, err := instance.makePaymentsSummary(lightning, listForwards)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, err
	}

	channelsSummary, err := instance.makeChannelsSummary(lightning, listFunds.Channels)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		// We admit this error here, we print only some log information.
		// In the call that cause this error we make a call to getListNodes, but
		// the node with that we have the channels with can be offiline for a while
		// and this mean that can be out of the gossip map.
	}

	listConfig, err := lightning.ListConfigs()
	if err != nil {
		log.GetInstance().Errorf("Error during the list config rpc command: %s", err)
		return nil, err
	}

	nodeLimits := &ChannelLimits{
		Min: int64(listConfig["min-capacity-sat"].(float64)),
		Max: 0, // TODO: Where is it the max? there is no max so I can put 0 here?
	}

	nodeFee := &ChannelFee{
		Base:    uint64(listConfig["fee-base"].(float64)),
		PerMSat: uint64(listConfig["fee-per-satoshi"].(float64)),
	}

	status := &status{
		Event:     nameEvent,
		Timestamp: time.Now().Unix(),
		Channels:  channelsSummary,
		Forwards:  statusPayments,
		Fee:       nodeFee,
		Limits:    nodeLimits,
	}

	return status, nil
}

// OnInit One time callback called from the lightning implementation
func (instance *MetricOne) OnInit(lightning *glightning.Lightning) error {
	getInfo, err := lightning.GetInfo()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during the OnInit method; %s", err))
		return err
	}

	instance.NodeID = getInfo.Id
	instance.Color = getInfo.Color
	instance.NodeAlias = getInfo.Alias
	instance.Network = getInfo.Network
	instance.NodeInfo = &NodeInfo{
		Implementation: "c-lightning", // It is easy, it is coupled with c-lightning plugin now
		Version:        getInfo.Version,
	}
	status, err := instance.onEvent("on_start", lightning)
	if err != nil {
		return err
	}
	instance.UpTime = append(instance.UpTime, status)
	instance.lastCheck = time.Now().Unix()
	if status.Timestamp > 0 {
		instance.lastCheck = status.Timestamp
	}

	//FIXME: We could use a set datastructures
	instance.Address = make([]*NodeAddress, 0)
	for _, address := range getInfo.Addresses {
		nodeAddress := &NodeAddress{
			Type: address.Type,
			Host: address.Addr,
			Port: uint(address.Port),
		}
		instance.Address = append(instance.Address, nodeAddress)
	}
	log.GetInstance().Info("Plugin initialized with OnInit event")
	return instance.MakePersistent()
}

func (instance *MetricOne) Update(lightning *glightning.Lightning) error {
	status, err := instance.onEvent("on_update", lightning)
	if err != nil {
		return err
	}
	instance.UpTime = append(instance.UpTime, status)
	instance.lastCheck = time.Now().Unix()
	if status.Timestamp > 0 {
		instance.lastCheck = status.Timestamp
	}
	return instance.MakePersistent()
}

func (instance *MetricOne) UpdateWithMsg(message *Msg,
	lightning *glightning.Lightning) error {
	return fmt.Errorf("method not supported")
}

func (instance *MetricOne) MakePersistent() error {
	instanceJson, err := instance.ToJSON()
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("JSON error %s", err))
		return err
	}
	return instance.Storage.StoreMetricOneSnapshot(instance.lastCheck, &instanceJson)
}

func (instance *MetricOne) OnStop(msg *Msg, lightning *glightning.Lightning) error {
	log.GetInstance().Debug("On close event on metrics called")
	//TODO: Check if the values are empty, if yes, try a solution
	// to avoid to push empty payload.
	var lastMetric MetricOne
	jsonLast, err := instance.Storage.LoadLastMetricOne()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(*jsonLast), &lastMetric); err != nil {
		return err
	}
	now := time.Now().Unix()
	lastStatus := lastMetric.UpTime[len(lastMetric.UpTime)-1]
	statusItem := &status{
		Event:     "on_close",
		Timestamp: now,
		Channels:  lastStatus.Channels,
		Forwards:  lastStatus.Forwards,
		Fee:       lastStatus.Fee,
		Limits:    lastStatus.Limits,
	}
	instance.UpTime = append(instance.UpTime, statusItem)
	instance.lastCheck = now
	if err := instance.MakePersistent(); err != nil {
		return err
	}
	return nil
}

// ToJSON Convert the MetricOne structure to a JSON string.
func (instance *MetricOne) ToJSON() (string, error) {
	json, err := json.Marshal(&instance)
	if err != nil {
		log.GetInstance().Error(err)
		return "", err
	}
	return string(json), nil
}

// InitOnRepo Contact the server and make an init the node.
func (instance *MetricOne) InitOnRepo(client *graphql.Client, lightning *glightning.Lightning) error {
	log.GetInstance().Info("Init plugin on repository")
	err := client.GetNodeMetadata(instance.NodeID, instance.Network)
	if err != nil {
		// If we received an error from the find method, maybe
		// the node it is not initialized on the server, and we
		// can try to init it.

		payload, err := instance.ToJSON()
		if err != nil {
			return err
		}
		// A restart of the plugin it is also caused from an update of it
		// and, so we supported only the migration of the previous version
		// for the moment.
		oldData, found := instance.Storage.GetOldData("metric_one", true)
		if found {
			log.GetInstance().Info("Found old data from db migration")
			payload = *oldData
		}

		toSign := sha256.SHA256(&payload)
		log.GetInstance().Info(fmt.Sprintf("Hash of the paylad: %s", toSign))
		signPayload, err := lightning.SignMessage(toSign)
		if err != nil {
			return err
		}

		if err := client.InitMetric(instance.NodeID, &payload, signPayload.ZBase); err != nil {
			return err
		}

		now := time.Now()
		log.GetInstance().Info(fmt.Sprintf("Metric One:Initialized on server at %s", now.Format(time.RFC850)))
		return nil
	} else {
		log.GetInstance().Info("Metric One: No initialization need, we simple tell to the server that we are back!")
		return instance.UploadOnRepo(client, lightning)
	}
}

// UploadOnRepo Contact the server and make an update request
func (instance *MetricOne) UploadOnRepo(client *graphql.Client, lightning *glightning.Lightning) error {
	payload, err := instance.ToJSON()
	if err != nil {
		return err
	}
	toSign := sha256.SHA256(&payload)
	log.GetInstance().Info(fmt.Sprintf("Hash of the paylad: %s", toSign))
	signPayload, err := lightning.SignMessage(toSign)
	if err != nil {
		return err
	}
	if err := client.UploadMetric(instance.NodeID, &payload, signPayload.ZBase); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error %s: ", err))
		return err
	}

	instance.UpTime = make([]*status, 0)
	instance.ChannelsInfo = make(map[string]*statusChannel)

	// Refactored this method in a utils functions
	t := time.Now()
	log.GetInstance().Info(fmt.Sprintf("Metric One Upload at %s", t.Format(time.RFC850)))
	return nil
}

func (instance *MetricOne) checkChannelInCache(lightning *glightning.Lightning, channelID string) (*cache.NodeInfoCache, error) {
	var nodeInfo cache.NodeInfoCache
	inCache := false
	if cache.GetInstance().IsInCache(channelID) {
		bytes, err := cache.GetInstance().GetFromCache(channelID)
		if err != nil {
			log.GetInstance().Errorf("Error %s:", err)
			return nil, err
		}
		if err := json.Unmarshal(bytes, &nodeInfo); err != nil {
			log.GetInstance().Errorf("Error %s", err)
			return nil, err
		}
		inCache = true
	}

	if !inCache {
		node, err := lightning.GetNode(channelID)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error in command listNodes in makeChannelsSummary: %s", err))
			// We admit this error, a node can be forgotten by the gossip if it is offline for long time.
			return nil, err
		}
		nodeInfo = cache.NodeInfoCache{
			ID:       node.Id,
			Alias:    node.Alias,
			Color:    node.Color,
			Features: node.Features,
		}
		if err := cache.GetInstance().PutToCache(nodeInfo.ID, nodeInfo); err != nil {
			log.GetInstance().Errorf("%s", err)
		}
	}
	return &nodeInfo, nil
}

// Make a summary of all the channels information that the node have a channels with.
func (instance *MetricOne) makeChannelsSummary(lightning *glightning.Lightning, channels []*glightning.FundingChannel) (*ChannelsSummary, error) {
	channelsSummary := &ChannelsSummary{
		TotChannels: 0,
		Summary:     make([]*ChannelSummary, 0),
	}

	if len(channels) > 0 {
		summary := make([]*ChannelSummary, 0)
		for _, channel := range channels {

			if channel.State == "ONCHAIN" {
				// When the channel is on chain, it is not longer a channel,
				// it stay in the listfunds for 100 block (bitcoin time) after the closing commitment
				log.GetInstance().Debug(fmt.Sprintf("The channel with ID %s has ON_CHAIN status", channel.Id))
				continue
			}

			channelSummary := &ChannelSummary{
				NodeId:    channel.Id,
				ChannelId: channel.ShortChannelId,
				State:     channel.State,
			}

			nodeInfo, err := instance.checkChannelInCache(lightning, channel.Id)
			if err != nil {
				// THIS should never happen, because can corrupt
				// the metrics
				return nil, err
			}
			channelsSummary.TotChannels++
			channelSummary.Alias = nodeInfo.Alias
			channelSummary.Color = nodeInfo.Color
			summary = append(summary, channelSummary)
		}
		channelsSummary.Summary = summary
	}

	return channelsSummary, nil
}

func (instance *MetricOne) makePaymentsSummary(lightning *glightning.Lightning, forwards []glightning.Forwarding) (*PaymentsSummary, error) {
	statusPayments := PaymentsSummary{
		Completed: 0,
		Failed:    0,
	}

	for _, forward := range forwards {
		switch forward.Status {
		case "settled", "offered":
			statusPayments.Completed++
		case "failed", "local_failed":
			statusPayments.Failed++
		default:
			return nil, fmt.Errorf("Status %s unexpected", forward.Status)
		}
	}

	return &statusPayments, nil
}

// private method of the module
func (instance *MetricOne) collectInfoChannels(lightning *glightning.Lightning, channels []*glightning.FundingChannel, event string) error {
	cache := make(map[string]bool)
	for _, channel := range channels {

		switch channel.State {
		// state of a channel where there is any type of communication yet
		// we skip this type of state
		case "CHANNELD_AWAITING_LOCKIN", "DUALOPEND_OPEN_INIT",
			"DUALOPEND_AWAITING_LOCKIN":
			continue
		default:
			if err := instance.collectInfoChannel(lightning, channel, event); err != nil {
				// void returning error here? We can continue to make the analysis over the channels
				log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
				return err
			}
			directions, err := instance.getChannelDirections(lightning, channel.ShortChannelId)
			if err != nil {
				log.GetInstance().Errorf("Error: %s", err)
				return nil
			}
			for _, direction := range directions {
				key := strings.Join([]string{channel.ShortChannelId, direction}, "_")
				cache[key] = true
			}
		}
	}

	// make intersection of the channels in the cache and a
	// channels in the metrics plugin
	// this is useful to remove the metrics over closed channels
	// in the metrics one we are not interested to have a story of
	// of the old channels (for the moments).
	for key := range instance.ChannelsInfo {
		_, found := cache[key]
		if !found {
			delete(instance.ChannelsInfo, key)
		}
	}

	return nil
}

func (instance *MetricOne) getChannelDirections(lightning *glightning.Lightning, channelID string) ([]string, error) {
	directions := make([]string, 0)

	channels, err := lightning.GetChannel(channelID)

	if err != nil {
		// This should happen when a channel is no longer inside the gossip
		// map.
		log.GetInstance().Errorf("Error: %s", err)
		directions = append(directions, "UNKNOWN")
		return directions, nil
	}

	for _, channel := range channels {
		direction := ChannelDirections[1]
		if channel.Source == instance.NodeID {
			direction = ChannelDirections[0]
		}
		directions = append(directions, direction)
	}

	return directions, nil
}

func (instance *MetricOne) collectInfoChannel(lightning *glightning.Lightning,
	channel *glightning.FundingChannel, event string) error {

	shortChannelId := channel.ShortChannelId
	var timestamp int64 = 0
	// avoid to store the wrong data related to the gossip delay.
	if instance.pingNode(lightning, channel.Id) {
		timestamp = time.Now().Unix()
	}

	directions, err := instance.getChannelDirections(lightning, shortChannelId)
	if err != nil {
		log.GetInstance().Errorf("Error: %s", err)
		return err
	}

	for _, direction := range directions {
		key := strings.Join([]string{shortChannelId, direction}, "_")
		infoChannel, found := instance.ChannelsInfo[key]

		infoMap, err := instance.getChannelInfo(lightning, channel, infoChannel)
		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error during get the information about the channel: %s", err))
			return err
		}

		info, infoFound := infoMap[direction]
		if !infoFound {
			log.GetInstance().Errorf("Error: channel not exist for direction %s", direction)
			// this should never happen, because we fill the channel with the same
			// method that we derive the directions
			return fmt.Errorf("Error: channel not exist for direction %s", direction)
		}
		// A new channels found
		channelStat := channelStatus{
			Event:     event,
			Timestamp: timestamp,
			Status:    channel.State,
		}

		if !found {
			upTimes := make([]*channelStatus, 1)
			upTimes[0] = &channelStat
			newInfoChannel := statusChannel{
				ChannelId:  shortChannelId,
				NodeId:     info.NodeId,
				NodeAlias:  info.Alias,
				Color:      info.Color,
				Capacity:   channel.ChannelSatoshi,
				LastUpdate: info.LastUpdate,
				Forwards:   info.Forwards,
				UpTimes:    upTimes,
				Online:     channel.Connected,
				Direction:  info.Direction,
				Fee:        info.Fee,
				Limits:     info.Limits,
			}
			instance.ChannelsInfo[key] = &newInfoChannel
		} else {
			infoChannel.Capacity = channel.ChannelSatoshi
			infoChannel.UpTimes = append(infoChannel.UpTimes, &channelStat)
			infoChannel.Color = info.Color
			infoChannel.Online = channel.Connected
			infoChannel.Fee = info.Fee
			infoChannel.Limits = info.Limits
		}
	}
	return nil
}

func (instance *MetricOne) pingNode(lightning *glightning.Lightning, nodeId string) bool {
	if _, err := lightning.Ping(nodeId); err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error during pinging node: %s", err))
		return false
	}
	return true
}

func NewUnknownChannel() *glightning.Channel {
	return &glightning.Channel{
		LastUpdate:               0,
		BaseFeeMillisatoshi:      0,
		FeePerMillionth:          0,
		HtlcMinimumMilliSatoshis: "0msat",
		HtlcMaximumMilliSatoshis: "0msat",
	}
}

// Get the information about the channel that is open with the node id
//
// lightning: Is the Go API for c-lightning
// channel: Contains the channels information of the command listfunds returned by c-lightning
// prevInstance: the previous instance of the channel info stored inside the map, if exist.
//
// as return:
// map[string]*ChannelsInfo: Information on how the channel with a specific short channel id is splitted.
// error: If any error during this operation occurs
func (instance *MetricOne) getChannelInfo(lightning *glightning.Lightning,
	channel *glightning.FundingChannel, prevInstance *statusChannel) (map[string]*ChannelInfo, error) {

	result := make(map[string]*ChannelInfo)

	subChannels, err := lightning.GetChannel(channel.ShortChannelId)

	// This error should never happen
	if err != nil {
		log.GetInstance().Errorf("Error: %s", err)
		subChannels = []*glightning.Channel{NewUnknownChannel()}
	}

	for _, subChannel := range subChannels {
		nodeInfo, err := instance.checkChannelInCache(lightning, channel.Id)
		if err != nil {
			continue
		}
		// Init the default data here
		channelInfo := &ChannelInfo{
			NodeId:     channel.Id,
			Alias:      "unknown",
			Color:      "unknown",
			Direction:  "UNKNOWN",
			LastUpdate: subChannel.LastUpdate,
			Forwards:   make([]*PaymentInfo, 0),
			Fee: &ChannelFee{
				Base:    subChannel.BaseFeeMillisatoshi,
				PerMSat: subChannel.FeePerMillionth,
			},
			Limits: &ChannelLimits{
				Min: getMSatValue(subChannel.HtlcMinimumMilliSatoshis),
				Max: getMSatValue(subChannel.HtlcMaximumMilliSatoshis),
			},
		}

		channelInfo.Direction = ChannelDirections[1]
		if subChannel.Source == instance.NodeID {
			channelInfo.Direction = ChannelDirections[0]
		} else if subChannel.Source == "" {
			channelInfo.Direction = "UNKNOWN"
		}

		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error during the call listNodes: %s", err))
			if prevInstance != nil {
				channelInfo.Alias = prevInstance.NodeAlias
				channelInfo.Color = prevInstance.Color
			}
			// We avoid to return the error because it is correct that the node
			// it is not up and running, this means that it is fine admit an
			// error here.i
			result[channelInfo.Direction] = channelInfo
			continue
		}

		channelInfo.Alias = nodeInfo.Alias
		channelInfo.Color = nodeInfo.Color

		listForwards, err := lightning.ListForwards()

		if err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error during the listForwards call: %s", err))
			return nil, err
		}

		for _, forward := range listForwards {
			receivedTime := utime.FromDecimalUnix(forward.ReceivedTime)
			// The duration of 30 minutes are relative to the plugin uptime event,
			// however, this can change in the future and can be dynamic.
			if !utime.InRangeFromUnix(time.Now().Unix(), receivedTime, 30*time.Minute) {
				// If is an old payments
				continue
			}

			// by default we assume that the payment has a out direction
			paymentInfo := &PaymentInfo{
				Direction: ChannelDirections[1],
				Status:    forward.Status,
				Timestamp: utime.FromDecimalUnix(forward.ReceivedTime),
			}

			// if we have the forward payment is inside the our direction
			// we change the direction from out to in.
			if channel.ShortChannelId == forward.OutChannel {
				paymentInfo.Direction = ChannelDirections[0]
			}

			// TODO: we are assuming that from a in channel we can receive
			// only in forward payment, and from outcoming payment we can forward
			// only if the channel is in outcoming state
			//
			// is correct the intuition?
			if channelInfo.Direction != "UNKNOWN" &&
				paymentInfo.Direction != channelInfo.Direction {
				continue
			}

			channelInfo.Forwards = append(channelInfo.Forwards, paymentInfo)

			switch forward.Status {
			case "settled", "offered", "failed":
				// do nothing
				continue
			case "local_failed":
				// store the information about the failure
				if len(channelInfo.Forwards) == 0 {
					continue
				}
				paymentInfo := channelInfo.Forwards[len(channelInfo.Forwards)-1]
				paymentInfo.FailureReason = forward.FailReason
				paymentInfo.FailureCode = forward.FailCode
			default:
				return nil, fmt.Errorf("Status %s unexpected", forward.Status)
			}
		}
		result[channelInfo.Direction] = channelInfo
	}
	return result, nil
}

//FIXME put inside the utils functions
func getMSatValue(msatStr string) int64 {
	msatTokens := strings.Split(msatStr, "msat")
	if len(msatTokens) == 0 {
		return -1
	}
	msatValue := msatTokens[0]
	value, err := strconv.ParseInt(msatValue, 10, 64)
	if err != nil {
		log.GetInstance().Errorf("Error parsing msat: %s", err)
	}
	return value
}
