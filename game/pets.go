package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"log"
	"sync/atomic"
)

var Counter int64

func InitCounter(store *db.UserStore) {
	res, err := store.GetTotalPetCount()
	if err != nil {
		log.Fatal(err)
	}

	Counter = int64(res)
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
