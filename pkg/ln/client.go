package ln

import (
	cln4go "github.com/vincenzopalazzo/cln4go/client"
	"github.com/vincenzopalazzo/pkg/model"
)

func ListFunds(client cln4go.Client) (*model.ListFundsResponse, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.ListFundsResponse](client, "listfunds", model.EmpityPayload)
}

func ListForwards(client cln4go.Client) (*model.ListForwards, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.ListForwardsResponse](client, "listforwards", model.EmpityPayload)
}

func ListConfig(client cln4go.Client) (*model.ListConfig, error) {
	return cln4go.Call[cln4go.Client, map[string]any, model.ListConfigResp](client, "listconfig", model.EmpityPayload)
}

func GetInfo(client cln4go.Client) (*model.GetInfo, error) {
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
	return res[0]
}

func ListChannels(client cln4go.Client, nodeId *string) (*[]model.ListChannelsResp, error) {
	req := model.ListChannelReq{
		ChannelId: nodeId,
	}
	return cln4go.Call[cln4go.Client, model.ListChannelReq, model.ListChannelsResp](client, "listchannels", req)
}

func GetNode(client cln4go.Client, channelId *string) (*model.ListNodesNode, error) {
	res, err := ListNodes(client, channelId)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

func ListNodes(client cln4go.Client, channelId *string) (*model.ListNodeResp, error) {
	req := model.ListNodesReq{
		ChannelId: channelId,
	}
	return cln4go.Call(client, "listnodes", req)
}
