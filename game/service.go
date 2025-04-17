package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"log"
	"sync/atomic"
)

type Service struct {
	store     *db.UserStore
	PetCount  int64
	UserCache map[string]db.User
}

var Counter int64

func NewController(store *db.UserStore) *Service {
	controller := &Service{
		store,
		0,
		make(map[string]db.User),
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
	s.store.SaveUserScore(user)
}

func CheckPersonalMilestone(count int) bool {
	return count == 10 || count == 25 || count == 50 || count == 100 || count%1000 == 0
}

func (s *Service) CheckMilestone() bool {
	return s.PetCount%25_000 == 0
}

func (s *Service) Autosave() {
	for _, user := range s.UserCache {
		err := s.store.SaveUserScore(&user)
		if err != nil {
			log.Printf("save user score error: %v", err)
			log.Printf("user info dump: %+v", user)
		}
	}
}
