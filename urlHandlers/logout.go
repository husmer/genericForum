package urlHandlers

import (
	"forum/dbconnections"
	"net/http"
	"time"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("UserCookie")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	dbconnections.LogoutUser(cookie.Value)
	exp := time.Now().Add(1 * time.Millisecond)
	cookie = &http.Cookie{
		Name:     "UserCookie",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  exp,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
