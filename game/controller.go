package game

import (
	"github.com/nathanmazzapica/pet-daisy/db"
	"log"
	"sync/atomic"
)

type Controller struct {
	store    *db.UserStore
	PetCount int64
}

var Counter int64

func NewController(store *db.UserStore) *Controller {
	controller := &Controller{store, 0}
	controller.InitCounter()

	return controller
}

func (c *Controller) InitCounter() {
	res, err := c.store.GetTotalPetCount()

	if err != nil {
		log.Fatal(err)
	}

	c.PetCount = int64(res)
}

func (c *Controller) PetDaisy() {
	atomic.AddInt64(&c.PetCount, 1)
}

func CheckPersonalMilestone(count int) bool {
	return count == 10 || count == 25 || count == 50 || count == 100 || count%1000 == 0
}

func (c *Controller) CheckMilestone() bool {
	return c.PetCount%25_000 == 0
}
