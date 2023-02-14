package main

// docs for this weviate api https://weaviate.io/developers/weaviate/api/rest/nodes
// status: Status of the node (one of: HEALTHY, UNHEALTHY, UNAVAILABLE).
// http://localhost:8080/v1/nodes
// {
//   "nodes": [
//     {
//       "name": "weaviate-7",
//       "status": "HEALTHY",
//       "version": "1.16-alpha.0",
//       "gitHash": "8cd2efa",
//       "stats": {
//         "shardCount":2,
//         "objectCount": 23328
//       },
//       "shards": [
//         {
//           "name":"azuawSAd9312F",
//           "class": "Class_7",
//           "objectCount": 13328
//         }, {
//           "name":"cupazAaASdfPP",
//           "class": "Foo",
//           "objectCount": 10000
//         }
//       ]
//     }
//       ]
//     }
//    ]
// }

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// struct for results data, including name and health status of each node in the cluster
type Nodes struct {
	Nodes []struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		version string `json:"version"`
		gitHash string `json:"gitHash"`
		Stats   struct {
			ShardCount  int `json:"shardCount"`
			ObjectCount int `json:"objectCount"`
		} `json:"stats"`
		Shards []struct {
			Name        string `json:"name"`
			Class       string `json:"class"`
			ObjectCount int    `json:"objectCount"`
		} `json:"shards"`
	} `json:"nodes"`
}

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Url string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-weaviate-health",
			Short:    "Basic Check for Weaviate using REST",
			Keyspace: "sensu.io/plugins/https://github.com/DoctorOgg/sensu-weaviate-health/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "url",
			Env:       "WEAVIATE_URL",
			Argument:  "url",
			Shorthand: "u",
			Default:   "",
			Usage:     "URL of your Weaviate instance, e.g. http://localhost:8080",
			Value:     &plugin.Url,
		},
	}
)

func main() {
	useStdin := false
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Error check stdin: %v\n", err)
		panic(err)
	}
	//Check the Mode bitmask for Named Pipe to indicate stdin is connected
	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Println("using stdin")
		useStdin = true
	}

	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, useStdin)
	check.Execute()
}

func checkArgs(event *corev2.Event) (int, error) {
	if len(plugin.Url) == 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--url or WEAVIATE_URL environment variable is required")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *corev2.Event) (int, error) {
	log.Println("executing check with url:", plugin.Url)

	// get the data from the weaviate api
	var results Nodes
	resp, err := http.Get(plugin.Url + "/v1/nodes")
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("Error getting data from weaviate: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("Error reading response body: %v", err)
	}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return sensu.CheckStateCritical, fmt.Errorf("Error unmarshalling json: %v", err)
	}

	found_unhealthy := false
	found_unavailable := false

	for _, node := range results.Nodes {

		fmt.Printf("Node %s is %s	", node.Name, node.Status)
		if node.Status == "UNHEALTHY" {
			found_unhealthy = true
		}
		if node.Status == "UNAVAILABLE" {
			found_unavailable = true
		}
	}

	if found_unhealthy {
		return sensu.CheckStateCritical, fmt.Errorf("One or more nodes are unhealthy")
	} else if found_unavailable {
		return sensu.CheckStateWarning, fmt.Errorf("One or more nodes are unavailable")
	} else {
		return sensu.CheckStateOK, nil
	}
}
