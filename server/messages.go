package server

import "encoding/json"

type ClientMessage struct {
	Client *Client `json:"-"`
	Name   string  `json:"name"`
	Data   string  `json:"data"`
}

type ServerMessage struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// buildClientMessage converts an incoming slice of raw byte data into a usable ClientMessage
func buildClientMessage(rawData []byte, client *Client) (ClientMessage, error) {
	var message ClientMessage

	if err := json.Unmarshal(rawData, &message); err != nil {
		return message, err
	}

	message.Client = client

	return message, nil
}

// toServerMessage converts a ClientMessage into a ServerMessage.
func (message *ClientMessage) toServerMessage() ServerMessage {
	var broadcast ServerMessage
	broadcast.Name = message.Name
	broadcast.Data = message.Data

	return broadcast
}
