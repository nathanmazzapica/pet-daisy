package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"log"
	"sync/atomic"
)

type Service struct {
	store    *db.UserStore
	PetCount int64
}

var Counter int64

func NewController(store *db.UserStore) *Service {
	controller := &Service{store, 0}
	controller.InitCounter()

	return controller
}

func (c *Service) InitCounter() {
	res, err := c.store.GetTotalPetCount()

	if err != nil {
		log.Fatal(err)
	}

	c.PetCount = int64(res)
}

func (c *Service) PetDaisy() {
	atomic.AddInt64(&c.PetCount, 1)
}

func CheckPersonalMilestone(count int) bool {
	return count == 10 || count == 25 || count == 50 || count == 100 || count%1000 == 0
}

func (c *Service) CheckMilestone() bool {
	return c.PetCount%25_000 == 0
}
