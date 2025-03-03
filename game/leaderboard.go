package game

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/db"
)

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"` // Change to user_id?
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

var userPets = make(map[string]int)
var topPlayers = make([]LeaderboardRowData, 10)

func populateLeaderboardFromDB() {
	var topUsers []LeaderboardRowData
	fmt.Printf("Getting top 10 users\n")

	rows, err := db.DB.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT 10")

	if err != nil {
		fmt.Printf("Error getting top 10 users: %v\n", err)
		return
	}

	position := 1
	for rows.Next() {
		user := &db.User{}
		rows.Scan(&user.UserID, &user.DisplayName, &user.PetCount)

		userPets[user.UserID] = position

		topUsers = append(topUsers, userToLeaderboardRowData(*user, position))
		position++
	}

	topPlayers = topUsers
}

func updateLeaderboard(user *db.User) {
	fmt.Printf("Updating leaderboard\n")
	pets := user.PetCount

	// checks if player already in top 10
	if _, ok := userPets[user.UserID]; !ok {
		// compare user and 10th best
		if pets > topPlayers[9].PetCount {
			topPlayers[9] = userToLeaderboardRowData(*user, 10)
		}
		return
	}

	position := userPets[user.UserID]

	fmt.Println("Next highest: ", topPlayers[position-2], "pets:", topPlayers[position-2].PetCount)
	fmt.Println("current user:", topPlayers[position-1], "pets", topPlayers[position-1].PetCount)

	if topPlayers[position-2].PetCount > pets {
		fmt.Println("idk")
		a := topPlayers[position-2]
		b := topPlayers[position-1]

		topPlayers[position-1] = a
		topPlayers[position-2] = b
	}
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
