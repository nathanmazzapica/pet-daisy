package server

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/game"
	"math/rand"
	"strconv"
)

var daisyMessages = []string{
	"arf bark arf arf arf",
	"i love being pet",
	"i am a dog",
	"please never stop patting me",
	"awoooo",
	"yay i love being pet",
	"i love sleeping and laying next to my mom",
	"i cant wait for my program to be on",
	"i love sudoku",
	"did you guys see what the kardashians did.",
	"whats the red heads name again.... miranda?",
	"nathan is so cool.",
	"bark arf ARF ARF RRRRRRARF",
}

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

func daisyMessage() ClientMessage {
	return ClientMessage{"Daisy", daisyMessages[rand.Intn(len(daisyMessages))]}
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
