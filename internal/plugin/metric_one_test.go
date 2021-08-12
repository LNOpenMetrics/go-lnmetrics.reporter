package plugin

import (
	"encoding/json"
	"testing"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/kinbiko/jsonassert"
)

func TestJSONSerializzation(t *testing.T) {
	sys, err := sysinfo.Host()
	if err != nil {
		t.Errorf("Test Failure caused by: %s", err)
	}
	metric := NewMetricOne("1234", sys.Info())
	jsonString, err := json.Marshal(metric)

	jsonTest := jsonassert.New(t)
	jsonTest.Assertf(string(jsonString), `{
   "channels_info": [],
   "color": "<<PRESENCE>>",
   "metric_name": "metric_one",
   "node_id": "<<PRESENCE>>",
   "os_info": {
      "architecture": "<<PRESENCE>>",
      "os": "<<PRESENCE>>",
      "version": "<<PRESENCE>>"
   },
   "timezone": "<<PRESENCE>>",
   "up_time": []
}`)
}

func TestJSONDeserializzation(t *testing.T) {
	jsonString := `{
   "channels_info": [
{
         "capacity": 450000,
         "color": "fe903f",
         "direction": "",
         "forwards": [
            {
               "direction": "INCOOMING",
               "failure_code": 4103,
               "failure_reason": "WIRE_TEMPORARY_CHANNEL_FAILURE",
               "status": "local_failed"
            }
         ],
         "last_update": 0,
         "node_alias": "carrot",
         "node_id": "036d2ac71176151db04fdac839a0ddea9f3a584f6c23bb0b4ac72c323124ec506b",
         "online": true,
         "public": false,
"channel_id": "fake",
"up_times": []
}
],
   "color": "02bf81",
   "metric_name": "metric_one",
   "node_id": "033904095f082d5fe8ff8d7ee96172e69f166f1b498ccfd3a1e4e5d139d1fad597",
   "os_info": {
      "architecture": "x86_64",
      "os": "Linux Mint",
      "version": "20.1 (Ulyssa)"
   },
   "timezone": "CEST",
   "up_time": [
      {
         "channels": {
            "summary": [],
            "tot_channels": 0
         },
         "event": "on_start",
         "forwards": {
            "completed": 0,
            "failed": 0
         },
         "timestamp": 1627742938
      }
   ]
}`
	var metric MetricOne
	err := json.Unmarshal([]byte(jsonString), &metric)
	if err != nil {
		t.Errorf("Test failure cause from the following error %s", err)
	}
	_, found := metric.ChannelsInfo["fake"]
	if found == false {
		t.Errorf("Test failure cause from the missing key in the channels info map. Key \"fake\" missed")
	}
}