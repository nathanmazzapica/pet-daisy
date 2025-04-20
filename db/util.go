package db

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func (s *UserStore) Autosave() {
	for {
		time.Sleep(3 * time.Minute)
		s.mu.RLock()
		for _, user := range s.Cache {
			err := s.SaveUserScore(user)
			if err != nil {
				log.Printf("save user score error: %v", err)
				log.Printf("user info dump: %+v", user)
				continue
			}
		}
		s.mu.RUnlock()
	}
}

// generateSyncCode generates a random 6 digit 'syncCode' used for account recovery/syncing
func generateSyncCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

func getRandomDisplayName() string {
	adjectives := []string{"big", "long", "small", "golden", "yellow", "black",
		"red", "short", "cunning", "silly", "radical", "sluggish",
		"speedy", "humorous", "shy", "scared", "brave", "intelligent", "stupid",
		"orange", "medium", "austere", "gaudy", "ugly", "beautiful", "sexy",
		"intellectual", "philosophical", "charged", "empty", "full",
		"serious", "vengeful", "malignant", "generous", "complacent",
		"ambitious", "lazy", "dull", "sharp", "splendid", "sexy", "cute",
		"loving", "hateful", "spiteful", "rude", "polite", "dastardly", "depressed"}

	nouns := []string{"Dog", "Watermelon", "Crusader", "Lancer", "Envisage", "Frog",
		"Beetle", "Cellphone", "Python", "Lizard", "Butterfly", "Dragon",
		"Automobile", "Cow", "Henry", "Levi", "Array", "Buzzer", "Balloon", "Book",
		"Calendar", "Burrito", "Corgi", "Pencil", "Pen", "Marker", "Bookshelf",
		"Sharpener", "Can", "Lightbulb", "Flower", "Daisy", "Eraser", "Battery",
		"Butter", "Cantaloupe", "Fridge", "Computer", "Programmer", "Kitty", "Barbell", "Bottle", "Toad", "Beryllium", "Consumer", "President", "Orange", "Entity"}

	fmt.Printf("%d\n", len(adjectives)*len(nouns)*1_000)

	adjI := rand.Intn(len(adjectives))
	nounI := rand.Intn(len(nouns))
	num := fmt.Sprintf("%04d", rand.Intn(1_000))

	return adjectives[adjI] + nouns[nounI] + num

}
