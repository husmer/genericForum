package urlHandlers

import (
	"fmt"
	"forum/dbconnections"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

	r.Body = http.MaxBytesReader(w, r.Body, 20971520)
	// 32 << 20 peaks olema 32 MB
	if err := r.ParseMultipartForm(20971520); err != nil {
		m.Errors = append(m.Errors, "The uploaded file is too big. Please choose a file that's less than 20MB in size")
		template.Execute(w, m)
		return
	}

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

	if len(r.MultipartForm.File) > 0 {
		file, fileHeader, err := r.FormFile("img")
		if err != nil {
			fmt.Println("File upload error", err)
		}
		defer file.Close()

		fileExtension := filepath.Ext(fileHeader.Filename)
		if fileExtension != ".jpeg" && fileExtension != ".jpg" && fileExtension != ".png" && fileExtension != ".gif" {
			m.Errors = append(m.Errors, "Only jpeg, png and gif extensions allowed")
			template.Execute(w, m)
			return
		}

		imageName := time.Now().UnixNano()
		dst, err := os.Create(fmt.Sprintf("./static/images/%d%s", imageName, filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dbconnections.InsertMessage(userForm, m.User.Id, strconv.Itoa(int(imageName))+filepath.Ext(fileHeader.Filename))
	} else {
		dbconnections.InsertMessage(userForm, m.User.Id, "")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
