package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/OpenLNMetrics/lnmetrics.utils/log"
)

type Client struct {
	// The graph ql can contains a list of server where
	// make the request.
	BaseUrl []string
	// Token to autenticate to the server
	Token  *string
	Client *http.Client
}

// Builder method to make a new client
func New(baseUrl []string) *Client {
	return &Client{BaseUrl: baseUrl, Client: &http.Client{Timeout: time.Second * 10}}
}

// Make Request is the method to make the http request
func (instance *Client) MakeRequest(query map[string]string) error {
	jsonValue, err := json.Marshal(query)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return err
	}

	failure := 0
	log.GetInstance().Debug(fmt.Sprintf("Push payload on server(s): %s", jsonValue))
	for _, url := range instance.BaseUrl {
		log.GetInstance().Info(fmt.Sprintf("Request to URL %s", url))
		request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
		if err != nil {
			failure++
			log.GetInstance().Error(fmt.Sprintf("Error with the message \"%s\" during the request to endpoint %s", err, url))
			continue
		}
		response, err := instance.Client.Do(request)
		defer func() {
			if err := response.Body.Close(); err != nil {
				log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			}
		}()
		if err != nil {
			failure++
			log.GetInstance().Error(fmt.Sprintf("error with the message \"%s\" during the request to endpoint %s", err, url))
			continue
		}
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
func (instance *Client) UploadMetrics(nodeId string, payloads *string) error {
	//TODO: generalize this method
	payload := fmt.Sprintf("mutation { addNodeMetrics(input: { node_id: \"%s\", payload_metric_one: \"%s\") { node_id } }", nodeId, *payloads)
	log.GetInstance().Info(fmt.Sprintf("Query payload is: %s", payload))
	query := instance.MakeQuery(payload)
	return instance.MakeRequest(query)
}
