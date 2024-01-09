package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandlePostContent(w http.ResponseWriter, r *http.Request) {

	template, err := template.ParseFiles("./templates/postContent.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "PostContent")
	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	r.ParseForm()
	for key := range r.Form {
		if key == "createPostComment" {
			if len(r.Form[key][0]) < 1 {
				m.Errors = append(m.Errors, "Comment can not be empty")
				template.Execute(w, m)
				return
			} else {
				dbconnections.InsertComment(r.Form["PostId"][0], m.User.Id, r.Form["createPostComment"][0])
			}
		}

		if key == "like" {
			userCurrLike := dbconnections.GetCommentLike(m.User.Id, r.Form["CommentId"][0])
			if userCurrLike == "1" {
				dbconnections.SetCommentLikes(m.User.Id, r.Form["CommentId"][0], "0")
			} else {
				dbconnections.SetCommentLikes(m.User.Id, r.Form["CommentId"][0], "1")
			}
		}
		if key == "dislike" {
			userCurrLike := dbconnections.GetCommentLike(m.User.Id, r.Form["CommentId"][0])
			if userCurrLike == "-1" {
				dbconnections.SetCommentLikes(m.User.Id, r.Form["CommentId"][0], "0")
			} else {
				dbconnections.SetCommentLikes(m.User.Id, r.Form["CommentId"][0], "-1")
			}
		}
	}

	m = dbconnections.GetMegaDataValues(r, "PostContent")
	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
