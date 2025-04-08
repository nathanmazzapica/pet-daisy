package game

import (
	"fmt"
	"github.com/nathanmazzapica/pet-daisy/db"
	"sync/atomic"
)

var Counter int64

func InitCounter() {
	result := db.DB.QueryRow("SELECT SUM(pets) FROM users")
	result.Scan(&Counter)
	fmt.Println("Init Counter:", Counter)
}

func PetDaisy(user *db.User) {
	atomic.AddInt64(&Counter, 1)
	user.PetCount++
	user.SaveToDB()

	//fmt.Printf("%s pet Daisy! Total pets: %d\n", user.DisplayName, Counter)
}

func CheckPersonalMilestone(count int) bool {
	return count == 10 || count == 25 || count == 50 || count == 100 || count%1000 == 0
}

func CheckMilestone() bool {
	return Counter%25_000 == 0
}
