package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// GetString retrieves a message from the specified topic, waiting up to the specified timeout using the default connection.
func GetString(ctx context.Context, topic string, timeout time.Duration) (string, error) {
	if defaultConnection == nil {
		return "", ErrNoConnection
	}
	return defaultConnection.GetString(ctx, topic, timeout)
}

// GetString retrieves a message from the specified topic, waiting up to the specified timeout.
func (cn *Conn) GetString(ctx context.Context, topic string, timeout time.Duration) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, errReq := http.NewRequestWithContext(cctx, http.MethodGet, fmt.Sprintf("%s/api/v1/get", cn.address), http.NoBody)
	if errReq != nil {
		if errors.Is(errReq, context.Canceled) {
			return "", ErrNoMessage
		}
		return "", errReq
	}

	q := req.URL.Query()
	q.Add("topic", topic)
	q.Add("timeout", timeout.String())
	req.URL.RawQuery = q.Encode()

	resp, errDo := cn.httpClient.Do(req)
	if errDo != nil {
		return "", errDo
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var om OutputMessage
	errDecode := json.NewDecoder(resp.Body).Decode(&om)
	resp.Body.Close()
	if errDecode != nil {
		return "", errDecode
	}

	return om.Data, nil
}
