package helpers

import "forum/structs"

func contains(elems []structs.Category, v string) bool {
	for _, s := range elems {
		if v == s.Category {
			return true
		}
	}
	return false
}

func FilterByCat(m structs.MegaData, value string) []structs.Post{
	var newPosts []structs.Post
	for i := 0; i < len(m.AllPosts); i++ {
		if contains(m.AllPosts[i].Categories, value) {
			newPosts = append(newPosts, m.AllPosts[i])
		}
	}
	return newPosts
}