package ln

import (
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/model"
	cln4go "github.com/vincenzopalazzo/cln4go/client"
)

func ListFunds(client cln4go.Client) (*model.ListFundsResp, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.ListFundsResp](client, "listfunds", model.EmpityPayload)
}

func ListForwards(client cln4go.Client) ([]*model.Forward, error) {
	resp, err := cln4go.Call[cln4go.Client, map[string]any, model.ListForwardsResp](client, "listforwards", model.EmpityPayload)
	if err != nil {
		return nil, err
	}
	return resp.Forwards, nil
}

func ListConfig(client cln4go.Client) (*model.ListConfigResp, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.ListConfigResp](client, "listconfig", model.EmpityPayload)
}

func GetInfo(client cln4go.Client) (*model.GetInfoResp, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.GetInfoResp](client, "getinfo", model.EmpityPayload)
}

func SignMessage(client cln4go.Client, content *string) (*model.SignMessageResp, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.SignMessageResp](client, "signmessage", model.EmpityPayload)
}

// FIXME: can exist some node with mode channels
func GetChannel(client cln4go.Client, nodeID *string) (*model.ListChannelsChannel, error) {
	res, err := ListChannels(client, nodeID)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func ListChannels(client cln4go.Client, nodeId *string) ([]*model.ListChannelsChannel, error) {
	req := model.ListChannelReq{
		ChannelId: nodeId,
	}
	resp, err := cln4go.Call[cln4go.Client, model.ListChannelReq, model.ListChannelsResp](client, "listchannels", req)
	if err != nil {
		return nil, err
	}
	return resp.Channels, nil
}

func GetNode(client cln4go.Client, channelId *string) (*model.ListNodesNode, error) {
	res, err := ListNodes(client, channelId)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func ListNodes(client cln4go.Client, channelId *string) ([]*model.ListNodesNode, error) {
	req := model.ListNodesReq{
		ChannelId: channelId,
	}
	resp, err := cln4go.Call[cln4go.Client, model.ListNodesReq, model.ListNodesResp](client, "listnodes", req)
	if err != nil {
		return nil, err
	}
	return resp.Nodes, nil
}

func ListPeers(client cln4go.Client, nodeId *string) ([]*model.ListPeersPeer, error) {
	req := model.ListPeersReq{
		PeerId: nodeId,
	}
	resp, err := cln4go.Call[cln4go.Client, model.ListPeersReq, model.ListPeersResp](client, "listpeers", req)
	if err != nil {
		return nil, err
	}
	return resp.Peers, nil
}
