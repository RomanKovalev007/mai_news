package services

import (
	"fmt"
	"time"

	"github.com/RomanKovalev007/mai_news/internal/models"
	"github.com/RomanKovalev007/mai_news/internal/repository"
	"github.com/google/uuid"
)

type InputPost struct{
	Title string `json:"title"` 
	Content string `json:"content"`
}

func GetPosts() []*models.Post{
	postsSlice := make([]*models.Post, 0, len(repository.Posts))
	for _, post := range repository.Posts{
		postsSlice = append(postsSlice, post)
	}
	return postsSlice
}

func CreatePost(inputPost InputPost) models.Post{
	post := models.Post{
		ID: uuid.New().String(),
		Title: inputPost.Title,
		Content: inputPost.Content,
		CreatedAt: time.Now(),
	}
	repository.Posts[post.ID] = &post

	return post
}

func GetPost(id string) (models.Post, error){
	post, ok := repository.Posts[id]
	if !ok{
		return models.Post{}, fmt.Errorf("post not found")
	}
	return *post, nil
}

func DeletePost(id string) error{
	_, ok := repository.Posts[id]
	if !ok{
		return fmt.Errorf("post not found")
	}
	repository.Posts[id] = nil
	return nil
}  

func PatchPost(id string, inputPost InputPost) (models.Post, error){
	post, ok := repository.Posts[id]
	if !ok{
		return models.Post{}, fmt.Errorf("post not found")
	}
	post.Title = inputPost.Title
	post.Content = inputPost.Content

	return *post, nil

} 