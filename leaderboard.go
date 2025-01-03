package main

import (
	"fmt"
)

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"`
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

func GetTopX(count int) []LeaderboardRowData {
	var topUsers []LeaderboardRowData
	fmt.Printf("Getting top %v users\n", count)

	rows, err := db.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT ?", count)

	if err != nil {
		fmt.Printf("Error getting top %v users: %v\n", count, err)
		return []LeaderboardRowData{}
	}

	position := 1
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.userID, &user.displayName, &user.petCount)

		topUsers = append(topUsers, userToLeaderboardRowData(*user, position))
		position++
	}

	return topUsers
}

func userToLeaderboardRowData(user User, position int) LeaderboardRowData {
	return LeaderboardRowData{
		DisplayName: user.displayName,
		PetCount:    user.petCount,
		Position:    position,
	}
}
