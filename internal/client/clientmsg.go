package client

import "github.com/stonecool/livemusic-go/internal"

type ClientMessage struct {
	message *internal.Message
	client  *Client
}
