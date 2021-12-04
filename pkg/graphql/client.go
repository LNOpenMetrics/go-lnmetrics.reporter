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

// Partial Wrapper around graphql error response
// FIXME: adding support for the path filed
type GraphQLError struct {
	Message string `json:"message"`
}

// GraphQL Response wrapper
type GraphQLResponse struct {
	Data   *map[string]interface{} `json:"data"`
	Errors []*GraphQLError         `json:"errors"`
}

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
		Client:    &http.Client{Timeout: time.Second * 90},
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

	httpTransport := &http.Transport{DialContext: dialContext, DisableKeepAlives: true}

	return &Client{
		BaseUrl: baseUrl,
		Client: &http.Client{
			Timeout:   time.Second * 90,
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
func (instance *Client) MakeRequest(query map[string]string) ([]*GraphQLResponse, error) {
	jsonValue, err := json.Marshal(query)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, err
	}

	failure := 0
	responses := make([]*GraphQLResponse, 0)
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

		if response.StatusCode != http.StatusOK {
			failure++
			log.GetInstance().Errorf("Non-OK HTTP status: %d", response.StatusCode)
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
		var respModel GraphQLResponse
		if err := json.Unmarshal([]byte(result), &respModel); err != nil {
			failure++
			log.GetInstance().Infof("Raw server response: %s", result)
			log.GetInstance().Error(fmt.Sprintf("Error during graphql response: %s", err))
			continue
		}
		responses = append(responses, &respModel)
		log.GetInstance().Debug(fmt.Sprintf("Result from server %s", result))
	}

	if failure == len(instance.BaseUrl) {
		return nil, fmt.Errorf("All the request to push the data into request are failed. %d Failure over %d request", failure, len(instance.BaseUrl))
	}

	return responses, nil
}

// Private function to clean the payload to migrate strings with " to \"
func (instance *Client) cleanBody(payload *string) *string {
	replace := strings.ReplaceAll(*payload, `"`, `\"`)
	return &replace
}

// TODO: adding variables to give more flexibility
func (instance *Client) MakeQuery(payload string) map[string]string {
	return map[string]string{"query": payload}
}

func (instance *Client) InitMetric(nodeID string, body *string, signature string) error {
	body = instance.cleanBody(body)
	payload := fmt.Sprintf(`mutation {
                                  initMetricOne(node_id: "%s", payload: "%s", signature: "%s") {
                                    node_id
                                  }
                               }`, nodeID, *body, signature)
	query := instance.MakeQuery(payload)
	_, err := instance.MakeRequest(query)
	return err
}

// Utils Function to update the with the last data the metrics on server..
func (instance *Client) UploadMetric(nodeID string, body *string, signature string) error {
	cleanBody := instance.cleanBody(body)
	payload := fmt.Sprintf(`mutation {
                                   updateMetricOne(node_id: "%s", payload: "%s", signature: "%s")
                               }`, nodeID, *cleanBody, signature)
	query := instance.MakeQuery(payload)
	_, err := instance.MakeRequest(query)
	return err
}

// Utils function that call the GraphQL server to get the metrics about the channel
func (instance *Client) GetMetricOneByNodeID(nodeID string, startPeriod int, endPeriod int) error {
	payload := fmt.Sprintf(`query {
                                  getMetricOne(node_id: "%s", start_period: %d, end_period: %d) {
                                     node_id
                                     metric_name
                                  }
                              }`, nodeID, startPeriod, endPeriod)
	query := instance.MakeQuery(payload)
	responses, err := instance.MakeRequest(query)
	for _, resp := range responses {
		if len(resp.Errors) != 0 {
			// Get only the first error.
			// FIXME: It is enough only the first one?
			errorQL := resp.Errors[0]
			return fmt.Errorf(errorQL.Message)
		}
	}
	return err
}
