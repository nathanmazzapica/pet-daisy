package db

import (
	"log"
	"time"
)

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"`
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

func UserToLeaderboardRowData(user User, position int) LeaderboardRowData {
	return LeaderboardRowData{
		DisplayName: user.DisplayName,
		PetCount:    user.PetCount,
		Position:    position,
	}
}

func (s *UserStore) GetTopPlayers() []LeaderboardRowData {

	var topUsers []LeaderboardRowData

	rows, err := s.DB.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT 10")

	if err != nil {
		log.Println("Error getting top players:", err)
		return []LeaderboardRowData{}
	}

	position := 1
	for rows.Next() {
		user := &User{}
		rows.Scan(&user.UserID, &user.DisplayName, &user.PetCount)

		topUsers = append(topUsers, UserToLeaderboardRowData(*user, position))
		position++
	}

	s.LastLeaderboardUpdate = time.Now().UnixMilli()

	return topUsers
}

func (s *UserStore) shouldUpdateLeaderboard() bool {
	return time.Now().UnixMilli() < s.LastLeaderboardUpdate+150
}

// getDelay temporarily uses length of the cache for calculating delay
func (s *UserStore) getDelay() int64 {
	delay := int64(100*len(s.Cache.Rows)) / 2

	if delay > 1000 {
		return 1000
	}

	return delay
}
