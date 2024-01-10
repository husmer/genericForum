package helpers

import (
	"genericforum/structs"

	"golang.org/x/crypto/bcrypt"
)

func contains(elems []structs.Category, v string) bool {
	for _, s := range elems {
		if v == s.Category {
			return true
		}
	}
	return false
}

func FilterByCat(m structs.MegaData, value string) []structs.Post {
	var newPosts []structs.Post
	for i := 0; i < len(m.AllPosts); i++ {
		if contains(m.AllPosts[i].Categories, value) {
			newPosts = append(newPosts, m.AllPosts[i])
		}
	}
	return newPosts
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
