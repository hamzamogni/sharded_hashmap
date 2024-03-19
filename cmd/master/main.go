package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

var (
	master *Master
)

func main() {
	master = NewMaster()

	e := echo.New()

	e.GET("/api/stats", handleGetStats)
	e.GET("/api/keys", handleGetAllKeys)
	e.GET("/api/:key", handleGetKey)
	e.POST("/api", handlePostKey)
	e.DELETE("/api/:key", handleDeleteKey)

	e.Logger.Fatal(e.Start(":8080"))

}

func handleGetKey(c echo.Context) error {
	key := c.Param("key")
	return c.String(200, key)
}

type SetKeyParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func handlePostKey(c echo.Context) error {
	var params SetKeyParams
	if err := c.Bind(&params); err != nil {
		return c.String(400, "Invalid JSON")
	}

	node := master.GetNode(params.Key)
	nodeAddr := fmt.Sprintf("http://%s/key", node.String())
	payload := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, params.Key, params.Value)

	req, err := http.Post(
		nodeAddr,
		"application/json",
		bytes.NewBuffer([]byte(payload)),
	)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return c.String(500, "Internal Server Error")
	}
	defer req.Body.Close()

	return c.String(200, fmt.Sprintf("SET %s=%s in node %s", params.Key, params.Value, node))
}

func handleDeleteKey(c echo.Context) error {
	key := c.Param("key")
	node := master.GetNode(key)
	nodeAddr := fmt.Sprintf("http://%s/key", node.String())

	req, err := http.NewRequest(http.MethodDelete, nodeAddr, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return c.String(500, "Internal Server Error")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return c.String(500, "Internal Server Error")
	}
	defer resp.Body.Close()

	return c.String(200, fmt.Sprintf("DELETE %s from node %s", key, node))
}

func handleGetAllKeys(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(len(master.Nodes))

	var mu sync.Mutex
	results := make(map[string]map[string]string)

    // Fetch all keys from all nodes asynchronously
	for _, node := range master.Nodes {
		go func(endpoint string) {
			defer wg.Done()

			resp, err := http.Get(fmt.Sprintf("http://%s/keys", node.String()))
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return 
			}

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return 
			}

            var response map[string]string
            err = json.Unmarshal(body, &response)
            if err != nil {
                fmt.Printf("Error: %s\n", err)
                return 
            }

            mu.Lock()
            results[endpoint] = response
            mu.Unlock()

		}(node.String())
	}

    wg.Wait()
    return c.JSON(200, results)

}


type Stats struct {
    NodeCount int `json:"node_count"`
}

func handleGetStats(c echo.Context) error {
    stats := Stats{
        NodeCount: len(master.Nodes),
    }

    return c.JSON(200, stats)
}

