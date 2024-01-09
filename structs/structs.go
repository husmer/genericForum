package structs

type Post struct {
	Id           string
	Title        string
	User         string
	Post         string
	Created      string
	Username     string
	LikeRating   string
	CurrUserRate string
	Categories   []Category
	Media        string
}

type Category struct {
	Id       string
	Category string
	Selected string
}

type Comment struct {
	Id           string
	PostId       string
	UserId       string
	Comment      string
	Created      string
	Likes        string
	CurrUserRate string
}

type User struct {
	Id       string
	Username string
	Email    string
}

type AccessRights struct {
	AccessRight string
}

type PostLikes struct {
	Id       string
	Post     string
	UserId   string
	PostLike string
}

type CommentLikes struct {
	Id          string
	CommentId   string
	UserId      string
	CommentLike string
}

type MegaData struct {
	User           User
	Post           Post
	Categories     []Category
	CategoryChoice []Category
	AllPosts       []Post
	AllPostLikes   []PostLikes
	AllComments    []Comment
	Access         AccessRights
	Errors         []string
}
