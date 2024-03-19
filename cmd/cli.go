package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getKeysCommand)
	rootCmd.AddCommand(getStatsCommand)
	rootCmd.AddCommand(seedShardsCommand)
}

var getKeysCommand = &cobra.Command{
	Use:   "list",
	Short: "Get all keys from shards",
	Run: func(cmd *cobra.Command, args []string) {
		data := getData()
		for k, v := range data {
			fmt.Printf("Node: %s\n", k)
			for key, value := range v {
				fmt.Printf("\t%s: %s\n", key, value)
			}
		}
	},
}

var getStatsCommand = &cobra.Command{
	Use:   "stats",
	Short: "Get stats from shards",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("http://localhost:8080/api/stats")
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(body))
	},
}

var seedShardsCommand = &cobra.Command{
	Use:   "seed",
	Short: "Seed shards with random data",
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < 100; i++ {
			key := genRandomString()
			value := genRandomString()

			nodeAddr := "http://localhost:8080/api"
			payload := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, key, value)

			req, err := http.Post(
				nodeAddr,
				"application/json",
				bytes.NewBuffer([]byte(payload)),
			)

			if err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			defer req.Body.Close()
		}
	},
}

type KeyValue map[string]string

func getData() map[string]KeyValue {
	results := make(map[string]KeyValue)

	resp, err := http.Get("http://localhost:8080/api/keys")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Println(err)
	}

	return results
}

func genRandomString() string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand_str := make([]byte, rand.Intn(10)+5)

	for i := range rand_str {
		rand_str[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(rand_str)
}
