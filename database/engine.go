package database

import (
	"fmt"
	"os"

	"genericforum/dbconnections"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Engine() {
	if !fileExists("./database/forum.db") {
		fmt.Println("Did not find the Database! Starting regeneration!")

		dbconnections.CreateDB()
		fmt.Println("Database Created!")

		dbconnections.CreateUsers()
		fmt.Println("User Table Created!")

		dbconnections.CreateAccessRights()
		fmt.Println("Access Rights Table Created!")

		dbconnections.CreateUserAccess()
		fmt.Println("User Access Table Created!")

		dbconnections.CreateSessions()
		fmt.Println("Session hash Table Created!")

		dbconnections.CreateCategory()
		fmt.Println("Category Table Created!")

		dbconnections.CreatePostCategoryList()
		fmt.Println("Post category list Table Created!")

		dbconnections.CreatePostLikes()
		fmt.Println("Post Likes Table Created!")

		dbconnections.CreateCommentLikes()
		fmt.Println("Comment Likes Table Created!")

		dbconnections.CreatePosts()
		fmt.Println("Posts Table Created!")

		dbconnections.CreateComments()
		fmt.Println("Comments Table Created!")

		dbconnections.CreateMedia()
		fmt.Println("Media Table Created!")

		fmt.Println("Full Database Regeneration Successfull!")
		fmt.Println("Creating dummy users")

		dbconnections.RegisterUser("Guest", "guest@admin.com", "55guest55")
		dbconnections.SetAccessRight(dbconnections.GetID("Guest"), "1")
		dbconnections.RegisterUser("Admin", "admin@admin.com", "admin")
		dbconnections.SetAccessRight(dbconnections.GetID("Admin"), "4")
	}
}
