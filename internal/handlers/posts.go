package handlers

import (
	"encoding/json"

	"log/slog"
	"net/http"
	"strconv"

	"github.com/RomanKovalev007/mai_news/internal/models"
)


type Poster interface{
	GetAllPosts() ([]models.OutputPost, error)
	GetPost(id int) (models.OutputPost, error)
	SavePost(post models.InputPost) (models.OutputPost, error)
	PatchPost(id int, inputPost models.InputPost) (models.OutputPost, error)
	DeletePost(id int) error
}

func GetAllPostsHandler(poster Poster, log *slog.Logger) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		posts, err := poster.GetAllPosts()
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			log.Error("failed to get all posts", err)
			return
		}
		json.NewEncoder(w).Encode(posts)
	}
}


func CreatePostHandler(poster Poster, log *slog.Logger) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		var post models.InputPost

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		createdPost, err := poster.SavePost(post)
		if err != nil {
			http.Error(w, "failed to save post", http.StatusInternalServerError)
			log.Error("failed to save post", err)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) 
		json.NewEncoder(w).Encode(createdPost)
	}
}

func GetPostHandler(poster Poster, log *slog.Logger) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		post, err := poster.GetPost(id)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			log.Error("failed to get post", err)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}

func PatchPostHandler(poster Poster, log *slog.Logger) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		var inputPost models.InputPost
		if err := json.NewDecoder(r.Body).Decode(&inputPost); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}


		post, err := poster.PatchPost(id, inputPost)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			log.Error("failed to patch post", err)
			return
		}
		json.NewEncoder(w).Encode(post)
	}
}

func DeletePostHandler(poster Poster, log *slog.Logger) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request){
		
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		err = poster.DeletePost(id)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			log.Error("failed to delete post", err)
			return
		}
	}
}	