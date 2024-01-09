package dbconnections

import (
	"database/sql"
	"fmt"
	"os"

	"forum/validateData"
)

func CreateDB() {
	if err := os.MkdirAll("./database", os.ModeSticky|os.ModePerm); err != nil {
		fmt.Println(err)
	}
	os.Create("./database/forum.db")
}

func CreateUsers() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `users` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `username` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL)")
	db.Close()
}

func CreateAccessRights() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `access_rights` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `access_rights` VARCHAR(255) NOT NULL)")
	for _, v := range []string{"Guest", "User", "Moderator", "Admin"} {
		_, err := db.Exec("INSERT INTO access_rights (access_rights) VALUES (?)", v)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func CreateUserAccess() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `user_access` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `user` INTEGER NOT NULL REFERENCES users(id), `user_access` INTEGER NOT NULL REFERENCES access_rights(id))")
	db.Close()
}

func CreateSessions() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `session` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `user` INTEGER UNIQUE REFERENCES users(id), `hash` VARCHAR(255) NOT NULL, `date` NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	database.Close()
}

func CreateCategory() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `category` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `category` VARCHAR(255) NOT NULL)")

	for _, v := range []string{"All", "Potato", "Carrot", "Tomatoe", "Apple", "Orange"} {
		_, err := db.Exec("INSERT INTO category (category) VALUES (?)", v)
		if err != nil {
			fmt.Println(err)
		}
	}
	db.Close()
}

func CreatePostCategoryList() {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	db.Exec("CREATE TABLE `post_category_list` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post_category` INTEGER NOT NULL REFERENCES category(id), `post_id` INTEGER NOT NULL REFERENCES posts(id))")
}

func CreatePostLikes() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `post_likes` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post` INTEGER NOT NULL REFERENCES posts(id), `user` INTEGER NOT NULL REFERENCES users(id), `post_like` INTEGER, UNIQUE(post, user) ON CONFLICT REPLACE)")
	database.Close()
}

func CreateCommentLikes() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `comment_likes` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `comment` INTEGER NOT NULL REFERENCES comments(id), `user` INTEGER NOT NULL REFERENCES users(id), `comment_like` INTEGER, UNIQUE(comment, user) ON CONFLICT REPLACE)")
	database.Close()
}

func CreatePosts() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `posts` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `title` VARCHAR(255) NOT NULL, `user` INTEGER NOT NULL REFERENCES users(id), `post` VARCHAR(255), `created` NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	database.Close()
}

func CreateComments() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `comments` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post_id` INTEGER NOT NULL REFERENCES posts(id), `user` INTEGER NOT NULL REFERENCES users(id), `comment` VARCHAR(255), `created` NOT NULL DEFAULT CURRENT_TIMESTAMP)")
	database.Close()
}

func CreateMedia() {
	database, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	database.Exec("CREATE TABLE `media` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `post_id` INTEGER NOT NULL REFERENCES posts(id), `image_name` VARCHAR(255) NOT NULL)")
	database.Close()
}
