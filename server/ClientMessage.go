package server

import "encoding/json"

type ClientMessage struct {
	Client *Client `json:"-"`
	Name   string  `json:"name"`
	Data   string  `json:"data"`
}

// Name = type; message = data. Don't feel like refactoring too much rn
type ServerMessage struct {
	Name string `json:"name"`
	Data string `json:"message"`
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

// Deprecated: toBytes converts a ClientMessage to a slice of bytes.
func (message *ClientMessage) toBytes() []byte {
	data, _ := json.Marshal(message)
	return data
}

// Deprecated: toBytes converts a ServerMessage to a slice of bytes.
func (message *ServerMessage) toBytes() []byte {
	data, _ := json.Marshal(message)
	return data
}

// toServerMessage converts a ClientMessage into a ServerMessage.
func (message *ClientMessage) toServerMessage() ServerMessage {
	var broadcast ServerMessage
	broadcast.Name = message.Name
	broadcast.Data = message.Data

	return broadcast
}
