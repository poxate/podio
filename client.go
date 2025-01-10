package podio

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

type fetchContext struct {
	durations map[string]*time.Duration
}

type Format string

const (
	MP3  Format = "mp3"
	WAV  Format = "wav"
	OPUS Format = "opus"
	// PCM signed 16-bit little-endian
	PCM_16_LE Format = "s16le"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) Compile(ctx context.Context, format Format, ab AudioBuilder, dst io.Writer) error {
	fCtx := &fetchContext{durations: map[string]*time.Duration{}}

	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, "wss://podio.poxate.com/compile?apiKey="+c.apiKey, nil)
	if err != nil {
		if resp != nil {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to dial: %w: %s", err, string(b))
		} else {
			return fmt.Errorf("failed to dial: %w", err)
		}
	}
	defer conn.Close()

	if err := conn.WriteJSON(map[string]any{
		"format": string(format),
		"body":   ab.toJson(fCtx),
	}); err != nil {
		return fmt.Errorf("conn.WriteJSON audio builder configuration: %w", err)
	}

	for {
		messageType, body, err := conn.ReadMessage()
		if err != nil {
			if ce, ok := err.(*websocket.CloseError); ok {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					break
				} else {
					return fmt.Errorf("conn.ReadMessage...unexpected close: %w", ce)
				}
			} else {
				return fmt.Errorf("conn.ReadMessage: %w", err)
			}
		}
		if messageType == websocket.BinaryMessage {
			if _, err := dst.Write(body); err != nil {
				return fmt.Errorf("dst.Write: %w", err)
			}
		} else if messageType == websocket.TextMessage {
			var status struct {
				Type         string `json:"type"`
				SaveDuration struct {
					Tag string `json:"tag"`
					Dur int64  `json:"duration"`
				} `json:"_saveDuration"`
			}
			if err := json.Unmarshal(body, &status); err != nil {
				return fmt.Errorf("json.Unmarshal duration tag: %w", err)
			}

			if status.Type == "saveDuration" {
				if _, ok := fCtx.durations[status.SaveDuration.Tag]; !ok {
					return fmt.Errorf("unknown duration tag: %s", status.SaveDuration.Tag)
				}
				*fCtx.durations[status.SaveDuration.Tag] = time.Duration(status.SaveDuration.Dur)
			}
		} else if messageType == websocket.CloseMessage {
			break
		}
	}

	return nil
}
