package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
)

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"`
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

func UserToLeaderboardRowData(user db.User, position int) LeaderboardRowData {
	return LeaderboardRowData{
		DisplayName: user.DisplayName,
		PetCount:    user.PetCount,
		Position:    position,
	}
}
