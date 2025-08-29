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
	om, err := cn.GetMessage(ctx, topic, timeout)
	if err != nil {
		return "", err
	}
	if om == nil {
		return "", ErrNoMessage
	}
	return om.Data, nil
}

// GetMessage retrieves a message from the specified topic using the default connection, waiting up to the specified timeout.
func GetMessage(ctx context.Context, topic string, timeout time.Duration) (*OutputMessage, error) {
	if defaultConnection == nil {
		return nil, ErrNoConnection
	}
	return defaultConnection.GetMessage(ctx, topic, timeout)
}

// GetMessage retrieves a message from the specified topic, waiting up to the specified timeout.
func (cn *Conn) GetMessage(ctx context.Context, topic string, timeout time.Duration) (*OutputMessage, error) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, errReq := http.NewRequestWithContext(reqCtx, http.MethodGet, fmt.Sprintf("%s/api/v1/get", cn.address), http.NoBody)
	if errReq != nil {
		if errors.Is(errReq, context.Canceled) {
			return nil, ErrNoMessage
		}
		return nil, errReq
	}

	q := req.URL.Query()
	q.Add("name", cn.name)
	q.Add("topic", topic)
	q.Add("timeout", timeout.String())
	req.URL.RawQuery = q.Encode()

	resp, errDo := cn.httpClient.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var om OutputMessage
	errDecode := json.NewDecoder(resp.Body).Decode(&om)
	_ = resp.Body.Close()
	if errDecode != nil {
		return nil, errDecode
	}

	return &om, nil
}
