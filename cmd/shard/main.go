package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
)

var (
	shard *Shard
)

func main() {
	e := echo.New()
	shard = NewShard()

    e.GET("/keys", handleGetAllKeys)

	e.GET("/:key", handleGetKey)

	e.POST("/key", handleSetKey)

	e.DELETE("/:key", handleDeleteKey)

	port := os.Getenv("PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func handleGetKey(c echo.Context) error {
	key := c.Param("key")
	return c.String(200, shard.Get(key))
}

type SetKeyParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func handleSetKey(c echo.Context) error {
	var params SetKeyParams
	if err := c.Bind(&params); err != nil {
		return c.String(400, "Invalid JSON")
	}
	key := params.Key
	value := params.Value

	shard.Set(key, value)
	fmt.Printf("SET %s=%s\n", key, value)
	return c.String(200, fmt.Sprintf("SET %s=%s", key, value))
}

func handleDeleteKey(c echo.Context) error {
	key := c.Param("key")
	shard.Delete(key)
	return c.String(200, fmt.Sprintf("DELETE %s", key))
}

func handleGetAllKeys(c echo.Context) error {
    return c.JSON(200, shard.GetAll())
}
