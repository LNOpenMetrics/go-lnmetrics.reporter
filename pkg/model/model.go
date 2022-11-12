package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

var EmpityPayload = make(map[string]any)

type ListFundsResp struct {
	Channels []*ListFundsChannel `json:"channels"`
}

type ListFundsChannel struct {
	PeerId                string  `json:"peer_id"`
	State                 string  `json:"state"`
	ShortChannelId        *string `json:"short_channel_id"`
	Connected             bool    `json:"connected"`
	InternalTotAmountMsat *uint64 `json:"amount_msat"`
	// FIXME: check the format of the deprecated filed
	InternalChannelTotalSat *string `json:"channel_total_sat"`
}

// TotAmountMsat return a string version of the total capacity
// of the channel.
func (self *ListFundsChannel) TotAmountMsat() uint64 {
	if self.InternalTotAmountMsat != nil {
		return *self.InternalTotAmountMsat
	}
	intPart := strings.Split(*self.InternalChannelTotalSat, "msat")[0]
	value, err := strconv.Atoi(intPart)
	if err != nil {
		log.GetInstance().Errorf("value %s parsing invalid %s", *self.InternalChannelTotalSat, err)
	}
	return uint64(value)
}

type ListForwardsResp struct {
	Forwards []*Forward `json:"forwards"`
}

type Forward struct {
	Status       string  `json:"status"`
	ReceivedTime float64 `json:"received_time"`
	OutChannel   string  `json:"out_channel"`
	InChannel    string  `json:"in_channel"`
	FailCode     *uint32 `json:"failcode"`
	FailReason   string  `json:"failreason"`
}

type GetInfoResp struct {
	Id        string         `json:"id"`
	Color     string         `json:"color"`
	Alias     string         `json:"alias"`
	Network   string         `json:"network"`
	Version   string         `json:"version"`
	Addresses []*NodeAddress `json:"address"`
}

type NodeAddress struct {
	Type string `json:"type"`
	Port uint16 `json:"port"`
	Addr string `json:"address"`
}

type SignMessageResp struct {
	ZBase string `json:"zbase"`
}

type ListChannelReq struct {
	ChannelId *string `json:"short_channel_id"`
}

type ListChannelsResp struct {
	Channels []*ListChannelsChannel `json:"channels"`
}

type ListChannelsChannel struct {
	Public              bool   `json:"public"`
	Source              string `json:"source"`
	LastUpdate          uint64 `json:"last_update"`
	BaseFeeMillisatoshi uint64 `json:"base_fee_millisatoshi"`
	FeePerMillionth     uint64 `json:"fee_per_millionth"`
	// FIXME: this was a string, so this wit the deprecate API should fails?
	HtlcMinimumMsat uint64 `json:"htlc_minimum_msat"`
	HtlcMaximumMsat uint64 `json:"htlc_maximum_msat"`
}

func (self *ListChannelsChannel) HtlcMinMsat() *string {
	result := fmt.Sprintf("%dmsat", self.HtlcMinimumMsat)
	return &result
}

func (self *ListChannelsChannel) HtlcMaxMsat() *string {
	result := fmt.Sprintf("%dmsat", self.HtlcMaximumMsat)
	return &result
}

type ListNodesReq struct {
	ChannelId *string `json:"id"`
}

type ListNodesResp struct {
	Nodes []*ListNodesNode `json:"nodes"`
}

type ListNodesNode struct {
	Id       string  `json:"id"`
	Alias    string  `json:"alias"`
	Color    string  `json:"color"`
	Features *string `json:"features"`
}

type ListPeersResp struct {
	Peers []*ListPeersPeer `json:"peers"`
}

type ListPeersReq struct {
	PeerId *string `json:"id,omitempty"`
}

type ListPeersPeer struct {
	Id        string `json:"id"`
	Connected bool   `json:"connected"`
}
