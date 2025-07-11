package db

import (
	"log"
	"sort"
	"sync"
	"time"
)

type Leaderboard struct {
	Top10      []LeaderboardRowData
	LastTop10  []LeaderboardRowData
	LastUpdate time.Time
	mu         sync.Mutex
}

type LeaderboardRowData struct {
	DisplayName string `json:"display_name"`
	PetCount    int    `json:"pet_count"`
	Position    int    `json:"position"`
}

func UserToLeaderboardRowData(user *User, position int) LeaderboardRowData {
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

		topUsers = append(topUsers, UserToLeaderboardRowData(user, position))
		position++
	}

	s.LastLeaderboardUpdate = time.Now().UnixMilli()

	return topUsers
}

func (s *UserStore) shouldUpdateLeaderboard() bool {
	return time.Now().UnixMilli() < s.LastLeaderboardUpdate+150
}

// getDelay temporarily uses length of the cache for calculating delay
// deprecated
func (s *UserStore) getDelay() int64 {
	delay := int64(100*len(s.Cache.Rows)) / 2

	if delay > 1000 {
		return 1000
	}

	return delay
}

type LeaderboardDelta struct {
	position     int
	lastPosition int
	petCount     int64
	lastPetCount int64
}

func NewLeaderboard(initial []LeaderboardRowData) *Leaderboard {
	lb := &Leaderboard{
		Top10:      make([]LeaderboardRowData, len(initial)),
		LastTop10:  make([]LeaderboardRowData, len(initial)),
		LastUpdate: time.Now(),
	}
	copy(lb.Top10, initial)
	copy(lb.LastTop10, initial)
	return lb
}

func (lb *Leaderboard) GetAll() []LeaderboardRowData {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	data := make([]LeaderboardRowData, len(lb.Top10))
	copy(data, lb.Top10)
	return data
}

func (lb *Leaderboard) UpdateUser(user *User) []LeaderboardRowData {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.LastTop10 = append([]LeaderboardRowData(nil), lb.Top10...)

	found := false
	for i := range lb.Top10 {
		if lb.Top10[i].DisplayName == user.DisplayName {
			lb.Top10[i].PetCount = user.PetCount
			found = true
			break
		}
	}

	if !found {
		if len(lb.Top10) < 10 || user.PetCount > lb.Top10[len(lb.Top10)-1].PetCount {
			lb.Top10 = append(lb.Top10, LeaderboardRowData{
				DisplayName: user.DisplayName,
				PetCount:    user.PetCount,
			})
		}
	}

	sort.Slice(lb.Top10, func(i, j int) bool { return lb.Top10[i].PetCount > lb.Top10[j].PetCount })

	if len(lb.Top10) > 10 {
		lb.Top10 = lb.Top10[:10]
	}

	for i := range lb.Top10 {
		lb.Top10[i].Position = i + 1
	}

	oldMap := make(map[string]LeaderboardRowData)
	for _, r := range lb.LastTop10 {
		oldMap[r.DisplayName] = r
	}
	newMap := make(map[string]LeaderboardRowData)
	for _, r := range lb.Top10 {
		newMap[r.DisplayName] = r
	}

	var changed []LeaderboardRowData
	for name, row := range newMap {
		if old, ok := oldMap[name]; !ok || old.PetCount != row.PetCount || old.Position != row.Position {
			changed = append(changed, row)
		}
	}

	for name, old := range oldMap {
		if _, ok := newMap[name]; !ok {
			changed = append(changed, LeaderboardRowData{DisplayName: name, PetCount: old.PetCount, Position: 0})
		}
	}

	lb.LastUpdate = time.Now()
	return changed
}
