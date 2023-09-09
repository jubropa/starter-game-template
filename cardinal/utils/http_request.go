package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"pkg.world.dev/world-engine/cardinal/ecs"
	"pkg.world.dev/world-engine/cardinal/ecs/inmem"
	"pkg.world.dev/world-engine/cardinal/ecs/storage"
)

func SendRequestWIthJsonBody(target_endpoint string, json_string string) (*http.Response, error) {

	json_value, _ := json.Marshal(json_string)
	buff := bytes.NewBuffer(json_value)
	req, err := http.NewRequest("POST", target_endpoint, buff)

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %q failed: %w", req.URL, err)
	} else if resp.StatusCode != 200 {
		buf, err := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("got response of %d: %v, %w", resp.StatusCode, string(buf), err)
	}

	return resp, nil
}

// NewWorld is the recommended way to run the game
func NewWorld(addr string, password string) *ecs.World {
	log.Log().Msg("Running in normal mode, using external Redis")
	if addr == "" {
		log.Log().Msg("Redis address is not set, using fallback - localhost:6379")
		addr = "localhost:6379"
	}
	if password == "" {
		log.Log().Msg("Redis password is not set, make sure to set up redis with password in prod")
		password = ""
	}

	rs := storage.NewRedisStorage(storage.Options{
		Addr:     addr,
		Password: password, // make sure to set this in prod
		DB:       0,        // use default DB
	}, "world")
	worldStorage := storage.NewWorldStorage(&rs)
	world, err := ecs.NewWorld(worldStorage)
	if err != nil {
		panic(err)
	}

	return world
}

// NewEmbeddedWorld is the most convenient way to run the game
// because it doesn't require spinning up Redis in a container.
// It runs a Redis server as a part of the Go process.
// NOTE: worlds with embedded redis are incompatible with Cardinal Editor.
func NewEmbeddedWorld() *ecs.World {
	log.Log().Msg("Running in embedded mode, using embedded miniredis")
	return inmem.NewECSWorld()
}
