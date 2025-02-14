package server

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/game"
	"strconv"
)

func newPetNotification() ClientMessage {
	// CHECK BACK LATER
	return ClientMessage{Name: "petCounter", Message: strconv.Itoa(int(game.Counter))}
}

func serverNotification(content string) ClientMessage {
	return ClientMessage{"server", content}
}

func newMilestoneNotification() ClientMessage {
	return ClientMessage{"Daisy", fmt.Sprintf("Yay! I have been pet %v times!", game.Counter)}
}

func newAchievmentNotification(user string, count int) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has pet daisy %v times!", user, count)}
}

func playerJoinNotification(user string) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has joined! say hi!", user)}
}

func playerLeftNotification(user string) ClientMessage {
	return ClientMessage{"server", fmt.Sprintf("%v has disconnected :(", user)}
}
