package plugin

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/OpenLNMetrics/go-metrics-reported/pkg/db"

	sysinfo "github.com/elastic/go-sysinfo"
	"github.com/kinbiko/jsonassert"
	//	"github.com/stretchr/testify/assert"
)

func init() {
	// TODO: The database is null in the test method, why?
	rootDir, _ := os.Executable()
	_ = db.GetInstance().InitDB(rootDir)
}

func TestJSONSerializzation(t *testing.T) {
	sys, err := sysinfo.Host()
	if err != nil {
		t.Errorf("Test Failure caused by: %s", err)
	}
	metric := NewMetricOne("1234", sys.Info())
	jsonString, _ := json.Marshal(metric)

	jsonTest := jsonassert.New(t)
	jsonTest.Assertf(string(jsonString), `{
   "channels_info": [],
   "color": "<<PRESENCE>>",
   "metric_name": "metric_one",
   "node_id": "<<PRESENCE>>",
   "node_alias": "<<PRESENCE>>",
   "os_info": {
      "architecture": "<<PRESENCE>>",
      "os": "<<PRESENCE>>",
      "version": "<<PRESENCE>>"
   },
   "timezone": "<<PRESENCE>>",
   "up_time": [],
   "version": "<<PRESENCE>>"
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
   "version": 1,
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

func TestJSONMigrationFrom0to1One(t *testing.T) {
	jsonString := `{
   "channels_info": {"fake":
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
},
   "color": "02bf81",
   "metric_name": "metric_one",
   "node_id": "033904095f082d5fe8ff8d7ee96172e69f166f1b498ccfd3a1e4e5d139d1fad597",
   "os_info": {
      "architecture": "x86_64",
      "os": "Linux Mint",
      "version": "20.1 (Ulyssa)"
   },
   "timezone": "CEST",
   "version": 0,
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

func TestJSONMigrationFrom0to1Two(t *testing.T) {
	jsonString := `{
   "channels_info": {"fake":
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
},
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
	if metric.Name != "metric_one" {
		t.Errorf("Test failure cause with metric name different")
		t.Errorf("We have %s but we expected %s", metric.Name, "metric_one")
	}
	_, found := metric.ChannelsInfo["fake"]
	if found == false {
		t.Errorf("Test failure cause from the missing key in the channels info map. Key \"fake\" missed")
	}
}

func TestJSONMigrationFrom0to1DevPrefix(t *testing.T) {
	jsonString := `{
   "dev_channels_info": {},
   "channels_info": {"fake":
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
},
   "color": "02bf81",
   "metric_name": "metric_one",
   "node_id": "033904095f082d5fe8ff8d7ee96172e69f166f1b498ccfd3a1e4e5d139d1fad597",
   "os_info": {
      "architecture": "x86_64",
      "os": "Linux Mint",
      "version": "20.1 (Ulyssa)"
   },
   "timezone": "CEST",
   "version": 0,
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
