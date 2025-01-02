package main

import (
	"fmt"
)

func GetTopX(count int) []User {
	var topUsers []User
	fmt.Printf("Getting top %v users\n", count)

	rows, err := db.Query("SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT ?", count)

	if err != nil {
		fmt.Printf("Error getting top %v users: %v\n", count, err)
		return []User{}
	}

	for rows.Next() {
		user := &User{}
		rows.Scan(&user.userID, &user.displayName, &user.petCount)
		topUsers = append(topUsers, *user)
	}

	return topUsers
}
