package game

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/db"
)

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"`
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

func GetTopX(count int) []LeaderboardRowData {
	var topUsers []LeaderboardRowData
	fmt.Printf("Getting top %v users\n", count)

	rows, err := db.DB.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT ?", count)

	if err != nil {
		fmt.Printf("Error getting top %v users: %v\n", count, err)
		return []LeaderboardRowData{}
	}

	position := 1
	for rows.Next() {
		user := &db.User{}
		rows.Scan(&user.UserID, &user.DisplayName, &user.PetCount)

		topUsers = append(topUsers, userToLeaderboardRowData(*user, position))
		position++
	}

	return topUsers
}

func userToLeaderboardRowData(user db.User, position int) LeaderboardRowData {
	return LeaderboardRowData{
		DisplayName: user.DisplayName,
		PetCount:    user.PetCount,
		Position:    position,
	}
}
