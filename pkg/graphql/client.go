package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

type Client struct {
	// The graph ql can contains a list of server where
	// make the request.
	BaseUrl []string
	// Token to autenticate to the server
	Token     *string
	Client    *http.Client
	WithProxy bool
}

// Builder method to make a new client
func New(baseUrl []string) *Client {
	return &Client{
		BaseUrl:   baseUrl,
		Client:    &http.Client{Timeout: time.Second * 10},
		WithProxy: false,
	}
}

func NewWithProxy(baseUrl []string, hostProxy string, portProxy uint64) (*Client, error) {
	// From: https://www.reddit.com/r/golang/comments/3qbdbf/how_can_i_create_an_http_request_with_socks5_proxy/
	proxyAddr := strings.Join([]string{hostProxy, fmt.Sprint(portProxy)}, ":")
	log.GetInstance().Info(fmt.Sprintf("Proxy url: %s", proxyAddr))
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.GetInstance().Error(
			fmt.Sprintf("Error during connection with proxy: %s", err),
		)
		return nil, err
	}

	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		// do anything with ctx
		return dialer.Dial(network, address)
	}

	httpTransport := &http.Transport{DialContext: dialContext}

	return &Client{
		BaseUrl: baseUrl,
		Client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: httpTransport,
		},
		WithProxy: true,
	}, nil
}

// TODO: move in a utils module
func isOnionUrl(url string) bool {
	return strings.HasPrefix(url, ".onion")
}

// Make Request is the method to make the http request
func (instance *Client) MakeRequest(query map[string]string) error {
	jsonValue, err := json.Marshal(query)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}

	failure := 0
	for _, url := range instance.BaseUrl {
		log.GetInstance().Info(fmt.Sprintf("Request to URL %s", url))
		if !instance.WithProxy && isOnionUrl(url) {
			log.GetInstance().Debug(fmt.Sprintf("Skipped request to url %s because the proxy it is not configured in the plugin", url))
			continue
		}
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
		if err != nil {
			failure++
			log.GetInstance().Error(fmt.Sprintf("Error with the message \"%s\" during the request to endpoint %s", err, url))
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := instance.Client.Do(request)
		if err != nil {
			failure++
			log.GetInstance().Error(fmt.Sprintf("error with the message \"%s\" during the request to endpoint %s", err, url))
			continue
		}
		defer func() {
			if err := response.Body.Close(); err != nil {
				log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			}
		}()
		result, err := ioutil.ReadAll(response.Body)
		if err != nil {
			failure++
			log.GetInstance().Error(fmt.Sprintf("error with the message \"%s\" during the request to endpoint %s", err, url))
			continue
		}
		log.GetInstance().Debug(fmt.Sprintf("Result from server %s", result))
	}

	if failure == len(instance.BaseUrl) {
		return fmt.Errorf("All the request to push the data into request are failed. %d Failure over %d request", failure, len(instance.BaseUrl))
	}

	return nil
}

func (instance *Client) MakeQuery(payload string) map[string]string {
	return map[string]string{"query": payload}
}

// This method is a util function to help the node to push the mertics over the servers.
// the payload is a JSON string of the payloads.
func (instance *Client) UploadMetrics(nodeId string, body *string) error {
	//TODO: generalize this method
	// mutation {
	//    addNodeMetrics(input: { node_id: "%s", payload_metric_one: "{}"} ){
	//	node_id
	//    }
	// }
	cleanBody := strings.ReplaceAll(*body, `"`, `\"`)
	payload := fmt.Sprintf(`mutation {
                                   addNodeMetrics( input: { node_id: "%s", payload_metric_one: "%s" } ) {
                                    node_id
                                   }
                                }`, nodeId, cleanBody)
	query := instance.MakeQuery(payload)
	return instance.MakeRequest(query)
}
