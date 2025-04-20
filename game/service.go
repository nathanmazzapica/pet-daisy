package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	store    *db.UserStore
	PetCount int64

	UserPetHistory map[string][10]time.Time
	mu             sync.RWMutex
}

var Counter int64

func NewController(store *db.UserStore) *Service {
	controller := &Service{
		store,
		0,
		make(map[string][10]time.Time),
		sync.RWMutex{},
	}
	controller.InitCounter()

	return controller
}

func (s *Service) InitCounter() {
	res, err := s.store.GetTotalPetCount()

	if err != nil {
		log.Fatal(err)
	}

	s.PetCount = int64(res)
}

func (s *Service) PetDaisy(user *db.User) {
	atomic.AddInt64(&s.PetCount, 1)
	user.SafeIncrementPet()

	if user.PetCount < 10_000 && user.PetCount%100 != 0 {
		return
	}

	s.store.SaveUserScore(user)
}

func CheckPersonalMilestone(count int) bool {
	return count == 10 || count == 25 || count == 50 || count == 100 || count%1000 == 0
}

func (s *Service) CheckMilestone() bool {
	return s.PetCount%25_000 == 0
}

// Autosave and UserCache belong in UserStore.
