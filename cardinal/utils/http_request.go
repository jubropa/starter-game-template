package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
