package queue

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SendString sends a non-persistent message to the specified topic using the default connection.
func SendString(ctx context.Context, topic string, msg string) error {
	if defaultConnection == nil {
		return ErrNoConnection
	}
	return defaultConnection.SendString(ctx, topic, msg)
}

// SendString sends a non-persistent message to the specified topic.
func (cn *Conn) SendString(ctx context.Context, topic string, msg string) error {
	im := InputMessage{
		Data:       msg,
		Persistent: false,
	}

	return cn.SendMessage(ctx, topic, im)
}

func SendMessage(ctx context.Context, topic string, msg InputMessage) (err error) {
	if defaultConnection == nil {
		return ErrNoConnection
	}
	return defaultConnection.SendMessage(ctx, topic, msg)
}

func (cn *Conn) SendMessage(ctx context.Context, topic string, msg InputMessage) (err error) {
	type request struct {
		Name  string `json:"name"`
		Topic string `json:"topic,omitempty"`
		InputMessage
	}

	reqMsg := request{
		Name:         cn.name,
		Topic:        topic,
		InputMessage: msg,
	}

	data, errEncode := json.Marshal(reqMsg)
	if errEncode != nil {
		return errEncode
	}

	u := fmt.Sprintf("%s/api/v1/send", cn.address)

	req, errReq := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(data))
	if errReq != nil {
		return errReq
	}

	resp, errDo := cn.httpClient.Do(req)
	if errDo != nil {
		return errDo
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusGone {
		return ErrNoConsumers
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("bad status code: %d, expect 201", resp.StatusCode)
	}

	return nil
}
