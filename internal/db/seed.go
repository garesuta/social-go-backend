package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"math/rand"

	"github.com/backend-production-go-1/internal/store"
)

var usernames = []string{"bob", "dave", "alice", "charlie", "eve", "frank", "grace", "hank", "irene", "jack",
	"karen", "leo", "mia", "nathan", "olivia", "paul", "quincy", "rachel", "steve", "tina",
	"uma", "victor", "wendy", "xander", "yara", "zane", "abby", "ben", "chloe", "daniel",
	"elena", "felix", "gia", "harry", "isla", "james", "katy", "liam", "molly", "noah",
	"ophelia", "peter", "quinn", "ruby", "simon", "terry", "ursula", "violet", "will", "zoe"}

var titles = []string{
	"Top Travel Apps for 2025",
	"Tech Gadgets for Digital Nomads",
	"AI in Travel Planning",
	"The Rise of Smart Luggage",
	"Best Travel Drones to Buy",
	"Wi-Fi on the Go: Staying Connected",
	"Eco-Friendly Travel Tech",
	"VR Travel Experiences",
	"The Future of Travel Booking",
	"Tech Tips for Remote Work Abroad",
	"Using GPS for Adventure Trips",
	"Travel Security Gadgets You Need",
	"Packing Light with Smart Devices",
	"Wearable Tech for Travelers",
	"Exploring Smart Cities",
	"5G and Its Impact on Travel",
	"Top Travel Gadgets for 2025",
	"How IoT is Changing Tourism",
	"Travel Photography Tech",
	"Blockchain in the Travel Industry",
}

var contents = []string{
	"Exploring the top 10 destinations for digital nomads in 2025.",
	"How AI is revolutionizing the travel booking experience.",
	"A beginner's guide to building REST APIs with Go.",
	"Tips for packing light with smart travel gadgets.",
	"An introduction to generics in Go and their practical use cases.",
	"The future of tourism: Virtual reality travel experiences.",
	"10 must-have tech gadgets for frequent travelers.",
	"How to implement concurrency in Go with Goroutines and Channels.",
	"The rise of blockchain technology in the travel industry.",
	"Tips for creating scalable applications with Go.",
	"Wi-Fi hotspots and portable solutions for staying connected abroad.",
	"Understanding error handling patterns in Go applications.",
	"How IoT devices are reshaping the hospitality industry.",
	"Eco-friendly travel tech to reduce your carbon footprint.",
	"An in-depth look at Go's built-in testing framework.",
	"Top drones for capturing breathtaking travel photos and videos.",
	"A comprehensive guide to working with databases in Go.",
	"How 5G networks are enhancing the travel experience.",
	"The essential toolkit for building CLI tools in Go.",
	"Cybersecurity tips for travelers in a tech-driven world.",
}

var tags = []string{
	"travel", "technology", "golang", "programming", "AI",
	"remote-work", "digital-nomad", "gadgets", "web-development", "tutorial",
	"blockchain", "smart-tech", "eco-friendly", "vr", "coding",
	"databases", "cloud", "innovation", "mobile-tech", "cybersecurity",
}

var comments = []string{
	"Great post! This was really helpful, thanks for sharing.",
	"I never thought about this before—very insightful!",
	"Can you provide more details on this topic?",
	"This is exactly what I was looking for, thanks!",
	"Awesome content! Keep up the great work.",
	"I tried this, and it worked perfectly. Thanks for the tips!",
	"Do you have recommendations for beginners?",
	"Interesting perspective! I learned something new today.",
	"Thanks for breaking this down so clearly.",
	"I disagree with some points, but overall, great post!",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(100)
	//create a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err) // Stop execution
	}
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback() // rollback if it's get error
			log.Println("Error creating user", err)
			return
		}
	}
	tx.Commit()

	posts := generatePosts(100, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}
	log.Println("Seeding has been completed")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			// Role: store.Role{
			// 	Name: "user",
			// },
			// Password: "123456",
		}
	}
	return users
}
func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]
		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts
}
func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}
	return cms
}
