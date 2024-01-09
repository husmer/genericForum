package urlHandlers

import (
	"forum/dbconnections"
	"html/template"
	"net/http"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/post.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "Post")
	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	r.ParseForm()
	userForm := r.Form

	if len(userForm["title"][0]) < 1 || len(userForm["message"][0]) < 1 {
		if len(userForm["title"][0]) < 1 {
			m.Errors = append(m.Errors, "Title can not be empty!")
		}
		if len(userForm["message"][0]) < 1 {
			m.Errors = append(m.Errors, "Message can not be empty!")
		}
		template.Execute(w, m)
		return
	}

	dbconnections.InsertMessage(userForm, m.User.Id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
