package server

import "encoding/json"

type ClientMessage struct {
	Client  *Client `json:"-"`
	Name    string  `json:"name"`
	Message string  `json:"message"`
}

// Name = type; message = data. Don't feel like refactoring too much rn
type ServerMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func prepareRawMessage(rawData []byte, client *Client) (ClientMessage, error) {
	var message ClientMessage

	if err := json.Unmarshal(rawData, &message); err != nil {
		return message, err
	}

	message.Client = client

	return message, nil
}

// I don't think I actually need this; review soon thx
func (message *ClientMessage) toBytes() []byte {
	data, _ := json.Marshal(message)
	return data
}

func (message *ServerMessage) toBytes() []byte {
	data, _ := json.Marshal(message)
	return data
}

func (message *ClientMessage) toServerMessage() ServerMessage {
	var broadcast ServerMessage
	broadcast.Name = message.Name
	broadcast.Message = message.Message

	return broadcast
}
