package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RomanKovalev007/mai_news/internal/services"
)



func GetPostsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.GetPosts())
}


func CreatePostHandler(w http.ResponseWriter, r *http.Request){
	var post services.InputPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	createdPost := services.CreatePost(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(createdPost)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	post, err := services.GetPost(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request){
	id := r.PathValue("id")
	err := services.DeletePost(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
}

func PatchPostHandler(w http.ResponseWriter, r *http.Request){
	var inputPost services.InputPost
	if err := json.NewDecoder(r.Body).Decode(&inputPost); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	post, err := services.PatchPost(id, inputPost)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(post)
}