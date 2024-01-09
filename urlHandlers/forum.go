package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"forum/helpers"
	"html/template"
	"net/http"
	"strings"
)

func HandleForum(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/forum.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	if r.URL.Path != "/" && r.URL.Path != "/register" && r.URL.Path != "/login" && r.URL.Path != "/logout" && r.URL.Path != "/post" && r.URL.Path != "/postcontent" && !strings.HasPrefix(r.URL.Path, "/postcontent?PostId=") {
		http.Error(w, "Bad Request: 404", http.StatusNotFound)
		return
	}

	m := dbconnections.GetMegaDataValues(r, "Forum")

	if r.Method != http.MethodPost {
		m.CategoryChoice[0].Selected = "true"
		template.Execute(w, m)
		return
	}
	r.ParseForm()
	if r.Form["Category"] != nil {
		m.AllPosts = helpers.FilterByCat(m, r.Form["Category"][0])
		for i := 0; i < len(m.CategoryChoice); i++ {
			if m.CategoryChoice[i].Category == r.Form["Category"][0] {
				m.CategoryChoice[i].Selected = "true"
			}
		}
		executeCat := template.Execute(w, m)
		if executeCat != nil {
			fmt.Println("Template error: ", executeCat)
		}
	}

	userCurrLike := dbconnections.GetPostLike(m.User.Id, r.Form["postId"][0])
	if r.Form["like"] != nil {
		if userCurrLike == "1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "1")
		}
	}
	if r.Form["dislike"] != nil {
		if userCurrLike == "-1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "-1")
		}
	}
	m = dbconnections.GetMegaDataValues(r, "Forum")
	m.CategoryChoice[0].Selected = "true"

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
