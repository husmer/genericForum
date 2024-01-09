package dbconnections

import (
	"database/sql"
	"fmt"
	"forum/helpers"
	"forum/structs"
	"forum/validateData"
	"net/http"
	"net/url"
	"time"
)

// Open and Close database connection. Returns database connection.
func DbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	validateData.CheckErr(err)
	return db
}

// Register new user, add user values to users database. Returns true if user/email is database.
func RegisterUser(username, email, password string) (bool, bool) {
	db := DbConnection()
	usernameCheck := CheckValueFromDB("username", username)
	emailCheck := CheckValueFromDB("email", email)
	hashpsw, err := helpers.HashPassword(password)
	if err != nil {
		fmt.Println("Error on password hashing", err)
	}
	if !usernameCheck && !emailCheck {
		_, err := db.Exec("INSERT INTO users(username, password, email) VALUES(?, ?, ?)", username, hashpsw, email)
		validateData.CheckErr(err)
		fmt.Println("New user added to the DB")
		if username != "Guest" && username != "Admin" {
			SetAccessRight(GetID(username), "2")
		}
		fmt.Println("Access granted to user", GetID(username))
	}
	defer db.Close()
	return usernameCheck, emailCheck
}

// Returns True if user inserted credentials are in database.
func LoginUser(username, password string) bool {
	getPassword := CheckPassword(username)
	return helpers.CheckPasswordHash(password, getPassword)
}

// Deletes Cookie from session database
func LogoutUser(hash string) {
	db := DbConnection()
	_, err := db.Exec("DELETE FROM session WHERE hash=?", hash)
	if err != nil {
		fmt.Println("LogoutUser")
		fmt.Println("Error code:", err)
	}
	defer db.Close()
}

// Applies Cookie in session database
func ApplyHash(user, hash string) {
	db := DbConnection()
	_, err := db.Exec("INSERT OR REPLACE INTO session(user, hash) VALUES(?, ?)", user, hash)
	defer db.Close()
	validateData.CheckErr(err)
}

// Returns ID if username exists in the users database
func GetID(username string) string {
	db := DbConnection()
	query := db.QueryRow("SELECT id FROM users WHERE username=?", username).Scan(&username)
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find username with that name to return ID")
		fmt.Println("Error code: ", query)
	}
	return username
}

// Returns All user info by user hash
func GetUserInfo(id string) structs.User {
	db := DbConnection()
	var userInfo structs.User
	var dump string
	query := db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&userInfo.Id, &userInfo.Username, &dump, &userInfo.Email)
	dump = ""
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find userId with that id")
		fmt.Println("Error code: ", query)
	}
	return userInfo
}

// Returns the users access rights
func GetAccessRight(id string) structs.AccessRights {
	db := DbConnection()
	var userAccess structs.AccessRights
	query := db.QueryRow("SELECT user_access FROM user_access WHERE user=?", id).Scan(&userAccess.AccessRight)
	defer db.Close()
	if query != nil {
		fmt.Println("Didn't find user with that id")
		fmt.Println("Error code: ", query)
	}
	return userAccess
}

// Sets the users access rights
func SetAccessRight(user string, access string) {
	db := DbConnection()
	_, err := db.Exec("INSERT INTO user_access (user, user_access) VALUES(?, ?)", user, access)
	if err != nil {
		fmt.Println("SetAccessRight")
		fmt.Println("Error code: ", err)
	}
	defer db.Close()
}

// Returns UserID from session database
func CheckHash(hash string) string {
	db := DbConnection()
	var user string
	var date string
	query := db.QueryRow("SELECT user, date FROM session WHERE hash=?", hash).Scan(&user, &date)
	defer db.Close()
	if query == nil {
		layout := "2006-01-02 15:04:05"
		hashDate, _ := time.Parse(layout, date)
		currTime := time.Now().Add(time.Minute * -10)
		if currTime.Sub(hashDate) > 0 {
			fmt.Println("Hash expired, Executing delete")
			LogoutUser(hash)
			return "1"
		}
	}
	if query != nil {
		return "1"
	}
	return user
}

// Returns True if hash in database
func HashInDatabase(hash string) bool {
	db := DbConnection()
	var user string
	query := db.QueryRow("SELECT user FROM session WHERE hash=?", hash).Scan(&user)
	defer db.Close()
	if query != nil {
		fmt.Println("HashInDatabase: didn't find user with that hash!")
		fmt.Println("Error code: ", query)
		return false
	}
	return true
}

// Return True if value is in users database column
func CheckValueFromDB(column string, valueToCheck string) bool {
	db := DbConnection()
	newUsername := db.QueryRow("SELECT "+column+" FROM users WHERE "+column+"=?", valueToCheck).Scan(&valueToCheck)
	defer db.Close()
	trigger := false
	if newUsername == nil {
		fmt.Println("Username already exists!")
		trigger = true
	}
	return trigger
}

// Returns password from users based on username. Hopefully encrypted.
func CheckPassword(username string) string {
	db := DbConnection()
	var returnPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&returnPassword)
	defer db.Close()
	if err != nil {
		fmt.Println("User does not exist")
		return ""
	}
	return returnPassword
}

// Returns all rows in an array of structs from posts database
func GetAllPosts(postId, userId string) []structs.Post {
	var allPosts []structs.Post
	if len(postId) > 0 {
		allPosts = append(allPosts, GetOnePost(postId))
		return allPosts
	}
	db := DbConnection()
	rows, _ := db.Query("SELECT posts.id, title, posts.user, posts.post, created, username, SUM(post_likes.post_like) FROM posts INNER JOIN users ON posts.user = users.id INNER JOIN post_likes ON posts.id = post_likes.post GROUP BY posts.id, title, posts.user, posts.post, created, username")
	defer db.Close()
	for rows.Next() {
		var post structs.Post
		if err := rows.Scan(&post.Id, &post.Title, &post.User, &post.Post, &post.Created, &post.Username, &post.LikeRating); err != nil {
			fmt.Println("Error: GetAllPosts SQL Query", err)
			rows.Close()
			return allPosts
		}
		// layout := "2006-01-02 15:04:05"
		// postDate, _ := time.Parse(layout, post.Created.String())
		// post.Created = time.Since(postDate).Truncate(time.Second).String()
		post.Categories = GetAllCategoriesForPost(post.Id)
		allPosts = append(allPosts, post)
	}
	defer rows.Close()
	allPostLikes := GetAllPostLikes(userId)
	for i := 0; i < len(allPosts); i++ {
		for j := 0; j < len(allPostLikes); j++ {
			if allPosts[i].Id == allPostLikes[j].Post {
				allPosts[i].CurrUserRate = allPostLikes[j].PostLike
			}
		}
	}
	return allPosts
}

// Returns a struct that contains data from one row in post database
func GetOnePost(postId string) structs.Post {
	db := DbConnection()
	posts := db.QueryRow("SELECT posts.id, posts.title, users.username, posts.post, posts.created FROM posts INNER JOIN users ON posts.user = users.id WHERE posts.id=?", postId)
	defer db.Close()
	var post structs.Post
	if err := posts.Scan(&post.Id, &post.Title, &post.User, &post.Post, &post.Created); err != nil {
		fmt.Println(err)
	}
	post.Media = GetMedia(postId)
	post.Categories = GetAllCategoriesForPost(postId)
	return post
}

func GetAllCategoriesForPost(postId string) []structs.Category {
	db := DbConnection()
	categoryList, _ := db.Query("SELECT category.id, category.category FROM post_category_list INNER JOIN category ON category.id = post_category_list.post_category WHERE post_id=?", postId)
	var allCats []structs.Category
	for categoryList.Next() {
		var cats structs.Category
		categoryList.Scan(&cats.Id, &cats.Category)
		allCats = append(allCats, cats)
	}
	categoryList.Close()
	return allCats
}

// Inserts into posts and post_category_list user inserted data
func InsertMessage(userForm url.Values, userId, mediaName string) {
	db := DbConnection()
	var inputTitle string
	var inputMessage string
	catArray := []string{"1"}

	for key, value := range userForm {
		if key == "title" {
			inputTitle = value[0]
		} else if key == "message" {
			inputMessage = value[0]
		} else {
			catArray = append(catArray, key)
		}
	}

	var data string
	err := db.QueryRow("INSERT INTO posts (title, user, post) VALUES (?, ?, ?) RETURNING id", inputTitle, userId, inputMessage).Scan(&data)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("INSERT INTO post_likes (post, user, post_like) VALUES (?, ?, ?)", data, "1", "0")
	if err != nil {
		fmt.Println(err)
	}

	if len(mediaName) > 0 {
		SetMedia(data, mediaName)
	}

	for _, v := range catArray {
		_, err = db.Exec("INSERT INTO post_category_list (post_category, post_id) VALUES (?, ?)", v, data)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer db.Close()
}

// Inserts into comments database user inserted comment and commentator userID
func InsertComment(postId string, commentatorId string, comment string) {
	db := DbConnection()

	var commentId string
	err := db.QueryRow("INSERT INTO comments (post_id, user, comment) VALUES (?, ?, ?) RETURNING id", postId, commentatorId, comment).Scan(&commentId)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("INSERT INTO comment_likes (comment, user, comment_like) VALUES (?, ?, ?)", commentId, "1", "0")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

}

// Returns all comments by post_id
func GetAllComments(postId, userId string) []structs.Comment {
	db := DbConnection()
	var allComments []structs.Comment
	allCommentsFromData, _ := db.Query("SELECT comments.id, comments.post_id, users.username, comments.comment, comments.created, SUM(comment_likes.comment_like) FROM comments INNER JOIN users ON users.id = comments.user INNER JOIN comment_likes ON comment_likes.comment = comments.id WHERE post_id=? GROUP BY comments.id", postId)
	defer db.Close()
	for allCommentsFromData.Next() {
		var comments structs.Comment
		if err := allCommentsFromData.Scan(&comments.Id, &comments.PostId, &comments.UserId, &comments.Comment, &comments.Created, &comments.Likes); err != nil {
			fmt.Println(err)
		}
		allComments = append(allComments, comments)
	}
	defer allCommentsFromData.Close()
	allCommentLikes := GetAllCommentLikes(userId)
	for i := 0; i < len(allComments); i++ {
		for j := 0; j < len(allCommentLikes); j++ {
			if allComments[i].Id == allCommentLikes[j].CommentId {
				allComments[i].CurrUserRate = allCommentLikes[j].CommentLike
			}
		}
	}
	return allComments
}

func GetAllCategories() []structs.Category {
	db := DbConnection()
	var allCategories []structs.Category
	data, _ := db.Query("SELECT * FROM category")
	defer db.Close()
	for data.Next() {
		var category structs.Category
		if err := data.Scan(&category.Id, &category.Category); err != nil {
			fmt.Println(err)
			return allCategories
		}
		allCategories = append(allCategories, category)
	}
	defer data.Close()
	return allCategories
}

func SetPostLikes(userId, postId, like string) {
	db := DbConnection()
	var postLike structs.PostLikes
	db.QueryRow("SELECT * FROM post_likes WHERE post=? AND user=?", postId, userId).Scan(&postLike.Id, &postLike.Post, &postLike.UserId, &postLike.PostLike)
	if postLike.Id == "" {
		_, err := db.Exec("INSERT INTO post_likes (post, user, post_like) VALUES (?, ?, ?)", postId, userId, like)
		if err != nil {
			fmt.Println("New: Error inserting like to table")
		}
	} else {
		_, err := db.Exec("REPLACE INTO post_likes (post, user, post_like) VALUES (?, ?, ?)", postId, userId, like)
		if err != nil {
			fmt.Println("Exists: Error inserting like to table")
		}
	}
	defer db.Close()
}

func SetCommentLikes(userId, commentId, like string) {
	db := DbConnection()
	var commentLike structs.CommentLikes
	db.QueryRow("SELECT * FROM comment_likes WHERE comment=? AND user=?", commentId, userId).Scan(&commentLike.Id, &commentLike.CommentId, &commentLike.UserId, &commentLike.CommentLike)
	if commentLike.Id == "" {
		_, err := db.Exec("INSERT INTO comment_likes (comment, user, comment_like) VALUES (?, ?, ?)", commentId, userId, like)
		if err != nil {
			fmt.Println("New: Error inserting comment like to table")
		}
	} else {
		_, err := db.Exec("REPLACE INTO comment_likes (comment, user, comment_like) VALUES (?, ?, ?)", commentId, userId, like)
		if err != nil {
			fmt.Println("Exists: Error inserting comment like to table")
		}
	}
	defer db.Close()
}

func GetAllPostLikes(userId string) []structs.PostLikes {
	db := DbConnection()
	rows, _ := db.Query("SELECT * FROM post_likes where user=?", userId)
	defer db.Close()
	var allUserPostLikes []structs.PostLikes
	for rows.Next() {
		var postLike structs.PostLikes
		if err := rows.Scan(&postLike.Id, &postLike.Post, &postLike.UserId, &postLike.PostLike); err != nil {
			fmt.Println(err)
			return allUserPostLikes
		}
		allUserPostLikes = append(allUserPostLikes, postLike)
	}
	defer rows.Close()
	return allUserPostLikes
}

func GetPostLike(userId, postId string) string {
	db := DbConnection()
	var rating string
	db.QueryRow("SELECT post_like FROM post_likes where user=? and post=? ", userId, postId).Scan(&rating)
	defer db.Close()
	return rating
}

func GetAllCommentLikes(userId string) []structs.CommentLikes {
	db := DbConnection()
	rows, _ := db.Query("SELECT * FROM comment_likes where user=?", userId)
	defer db.Close()
	var allUserCommentLikes []structs.CommentLikes
	for rows.Next() {
		var commentLike structs.CommentLikes
		if err := rows.Scan(&commentLike.Id, &commentLike.CommentId, &commentLike.UserId, &commentLike.CommentLike); err != nil {
			fmt.Println(err)
			return allUserCommentLikes
		}
		allUserCommentLikes = append(allUserCommentLikes, commentLike)
	}
	defer rows.Close()
	return allUserCommentLikes
}

func GetCommentLike(userId, commentId string) string {
	db := DbConnection()
	var rating string
	db.QueryRow("SELECT comment_like FROM comment_likes where user=? and comment=? ", userId, commentId).Scan(&rating)
	defer db.Close()
	return rating
}

func SetMedia(postId, imageName string) {
	db := DbConnection()
	_, err := db.Exec("INSERT INTO media (post_id, image_name) VALUES (?, ?)", postId, imageName)
	if err != nil {
		fmt.Println(err)
	}
}

func GetMedia(postId string) string {
	db := DbConnection()
	var imageName string
	db.QueryRow("SELECT image_name FROM media where post_id=? ", postId).Scan(&imageName)
	defer db.Close()
	return imageName
}

func GetMegaDataValues(r *http.Request, handler string) structs.MegaData {
	var userId string
	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		userId = "1"
	} else {
		userId = CheckHash(cookie.Value)
	}

	postId := r.URL.Query().Get("PostId")
	m := structs.MegaData{
		User:   GetUserInfo(userId),
		Access: GetAccessRight(userId),
	}
	if handler == "Forum" {
		m.AllPosts = GetAllPosts(postId, userId)
		m.CategoryChoice = GetAllCategories()
	}
	if handler == "Post" {
		m.CategoryChoice = GetAllCategories()[1:]
	}
	if handler == "PostContent" {
		m.AllPosts = GetAllPosts(postId, userId)
		m.AllPosts[0].Categories = m.AllPosts[0].Categories[1:]
		m.AllComments = GetAllComments(postId, userId)
	}
	return m
}
