package model

var EmpityPayload = make(map[string]any)

type ListFundsResponse struct {
	Channels []*ListFundsChannelResponse
}

type ListFundsChannelResponse struct{}

type ListForwardsResponse struct{}

type ListConfigResp struct{}

type GetInfoResp struct {
	Id      string `json:"id"`
	Color   string `json:"color"`
	Alias   string `json:"alias"`
	Network string `json:"network"`
	Version string `json:"version"`
}

type SignMessageResp struct {
	ZBase string `json:"zbase"`
}

type ListChannelReq struct {
	ChannelId *string
}

type ListChannelsResp struct {
	Channels []*ListChannelsChannel
}

type ListChannelsChannel struct{}

type ListNodeReq struct {
	ChannelId *string
}

type ListNodesResp struct {
	nodes []*ListNodesNode
}

type ListNodesNode struct{}
