package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"pkg.world.dev/world-engine/cardinal/ecs"
)

func SendRequestWithJsonBody(target_endpoint string, json_value []byte, world *ecs.World) (*http.Response, error) {

	buff := bytes.NewBuffer(json_value)
	world.Logger.Debug().Msg("buff")
	req, req_err := http.NewRequest("POST", target_endpoint, buff)
	if req_err != nil {
		return nil, fmt.Errorf("request to %q failed: %w", req.URL, req_err)
	}

	req.Header.Add("Content-Type", "application/json")

	world.Logger.Debug().Msg("req")
	client := &http.Client{
		Timeout: time.Second,
	}
	world.Logger.Debug().Msg("client")

	resp, do_err := client.Do(req)
	world.Logger.Debug().Msg("resp")
	if do_err != nil {
		return nil, fmt.Errorf("request to %q failed: %w", req.URL, do_err)
	} else if resp.StatusCode != 200 {
		buf, do_err := io.ReadAll(resp.Body)
		world.Logger.Debug().Msg("body")
		return nil, fmt.Errorf("got response of %d: %v, %w", resp.StatusCode, string(buf), do_err)
	}

	return resp, nil
}
