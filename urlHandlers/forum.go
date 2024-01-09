package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"forum/helpers"
	"html/template"
	"net/http"
	"strings"
	"time"
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
		return
	}

	if r.Form["like"] != nil {
		userCurrLike := dbconnections.GetPostLike(m.User.Id, r.Form["postId"][0])
		if userCurrLike == "1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "1")
		}
	}
	if r.Form["dislike"] != nil {
		userCurrLike := dbconnections.GetPostLike(m.User.Id, r.Form["postId"][0])
		if userCurrLike == "-1" {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "0")
		} else {
			dbconnections.SetPostLikes(m.User.Id, r.Form["postId"][0], "-1")
		}
	}

	if r.Form["filter"] != nil {
		if r.Form["filter"][0] == "Likes" {
			for i := 0; i < len(m.AllPosts); i++ {
				for j := 0; j < len(m.AllPosts); j++ {
					if m.AllPosts[i].LikeRating > m.AllPosts[j].LikeRating {
						m.AllPosts[i], m.AllPosts[j] = m.AllPosts[j], m.AllPosts[i]
					}
				}
			}
			m.CategoryChoice[0].Selected = "true"
			executeErr := template.Execute(w, m)
			if executeErr != nil {
				fmt.Println("Template error: ", executeErr)
			}
			return
		} else if r.Form["filter"][0] == "Dates" {
			for i := 0; i < len(m.AllPosts); i++ {
				for j := 0; j < len(m.AllPosts); j++ {
					layout := "2006-01-02 15:04:05"
					postDatei, _ := time.Parse(layout, m.AllPosts[i].Created)
					postDatej, _ := time.Parse(layout, m.AllPosts[j].Created)
					if time.Since(postDatei) > time.Since(postDatej) {
						m.AllPosts[i], m.AllPosts[j] = m.AllPosts[j], m.AllPosts[i]
					}
				}
			}
			m.CategoryChoice[0].Selected = "true"
			executeErr := template.Execute(w, m)
			if executeErr != nil {
				fmt.Println("Template error: ", executeErr)
			}
			return
		} else {
			fmt.Println("Nothing yet")
		}
	}

	m = dbconnections.GetMegaDataValues(r, "Forum")
	m.CategoryChoice[0].Selected = "true"

	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
