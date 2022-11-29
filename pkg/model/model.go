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
	PeerId         string  `json:"peer_id"`
	State          string  `json:"state"`
	ShortChannelId *string `json:"short_channel_id"`
	Connected      bool    `json:"connected"`
	/// FIXME: where cln remove the deprecation period
	/// for the string please put this as `uint64`.
	InternalTotAmountMsat any `json:"amount_msat"`
}

// TotAmountMsat return a string version of the total capacity
// of the channel.
func (self *ListFundsChannel) TotAmountMsat() uint64 {
	return ParseDeprecatedMsat(self.InternalTotAmountMsat)
}

// ParseMsatStrToInt  Parse a string that is the amount msat with the following form
// "XXXmsat" the result from this function will be "xxx" as int
func ParseMsatStrToInt(obj *string) (uint64, error) {
	intPart := strings.Split(*obj, "msat")[0]
	value, err := strconv.Atoi(intPart)
	if err != nil {
		log.GetInstance().Errorf("value %s parsing invalid %s", *obj, err)
		return 0, err
	}
	return uint64(value), nil
}

// / ParseDeprecatedMsat wrap the deprecated cast method
// / inside a simple function.
// FIXME: remove the function when core lightning will remove
// the type string from the API.
func ParseDeprecatedMsat(msat any) uint64 {
	if msat != nil {
		switch val := msat.(type) {
		case string:
			res, _ := ParseMsatStrToInt(&val)
			return res
		default:
			res, _ := val.(float64)
			return uint64(res)
		}
	}
	return 0
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
	HtlcMinimumMsat     uint64 `json:"minimum_htlc_out_msat"`
	HtlcMaximumMsat     uint64 `json:"maximum_htlc_out_msat"`
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
