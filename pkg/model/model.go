package model

var EmpityPayload = make(map[string]any)

type ListFundsResp struct {
	Channels []*ListFundsChannel
}

type ListFundsChannel struct{}

type ListForwardsResp struct {
	Forwards []*Forward `json:"forwards"`
}

type Forward struct{}

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

type ListNodesReq struct {
	ChannelId *string
}

type ListNodesResp struct {
	Nodes []*ListNodesNode
}

type ListNodesNode struct{}

type ListPeersResp struct {
	Peers []*ListPeersPeer
}

type ListPeersReq struct {
	PeerId *string `json:"id,omitempty"`
}

type ListPeersPeer struct {
	Connected bool `json:"connected"`
}
