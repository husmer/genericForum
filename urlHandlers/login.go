package urlHandlers

import (
	"fmt"
	"forum/cleanData"
	"forum/dbconnections"
	"forum/validateData"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Template not found!"+err.Error(), http.StatusInternalServerError)
	}

	m := dbconnections.GetMegaDataValues(r, "Login")

	if r.Method != http.MethodPost {
		template.Execute(w, m)
		return
	}

	formDataUsername := r.FormValue("username")
	formDataPassword := r.FormValue("password")

	errorLog := []string{}
	dataValid := true
	validatePassword, _ := validateData.ValidatePassword(formDataPassword, formDataPassword)
	if !validateData.ValidateName(formDataUsername) {
		errorLog = append(errorLog, "Username should be minimum 3 letters!")
		dataValid = false
	}
	if !validatePassword {
		errorLog = append(errorLog, "Password should be at least 6 letters long!")
		dataValid = false
	}
	if !dataValid {
		m.Errors = errorLog
		executeErr := template.Execute(w, m)
		if executeErr != nil {
			fmt.Println("Template error: ", executeErr)
		}
		return
	}

	username := cleanData.CleanName(formDataUsername)

	if dbconnections.LoginUser(username, formDataPassword) {
		id := uuid.New()
		exp := time.Now().Add(10 * time.Minute)
		cookie := &http.Cookie{
			Name:     "UserCookie",
			Value:    id.String(),
			Path:     "/",
			HttpOnly: true,
			Expires:  exp}
		http.SetCookie(w, cookie)
		dbconnections.ApplyHash(dbconnections.GetID(username), id.String())

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.Errors = append(m.Errors, "Username or Password incorrect")
	executeErr := template.Execute(w, m)
	if executeErr != nil {
		fmt.Println("Template error: ", executeErr)
	}
}
