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

func newPetNotification() ServerMessage {
	return ServerMessage{Name: "petCounter", Data: strconv.Itoa(int(game.Counter))}
}

func serverNotification(content string) ServerMessage {
	return ServerMessage{"server", content}
}

func newMilestoneNotification() ServerMessage {
	return ServerMessage{"Daisy", fmt.Sprintf("Yay! I have been pet %v times!", game.Counter)}
}

func daisyMessage() ServerMessage {
	return ServerMessage{"Daisy", daisyMessages[rand.Intn(len(daisyMessages))]}
}

func newAchievmentNotification(user string, count int) ServerMessage {
	return ServerMessage{"server", fmt.Sprintf("%v has pet daisy %v times!", user, count)}
}

func playerJoinNotification(user string) ServerMessage {
	return ServerMessage{"server", fmt.Sprintf("%v has joined! say hi!", user)}
}

func playerLeftNotification(user string) ServerMessage {
	return ServerMessage{"server", fmt.Sprintf("%v has disconnected :(", user)}
}

// TODO: Implement
func leaderboardUpdateNotification() ServerMessage {

	// optimistic about errors :D
	return ServerMessage{"leaderboard", "not implemented"}
}
