package main

import (
	"fmt"
	"net/http"

	"github.com/RomanKovalev007/mai_news/internal/handlers"
)


func main(){
	r := http.NewServeMux()


	r.HandleFunc("GET /posts/", handlers.GetPostsHandler)
	r.HandleFunc("POST /posts/", handlers.CreatePostHandler)

	r.HandleFunc("GET /posts/{id}/",handlers.GetPostHandler)
	r.HandleFunc("PATCH /posts/{id}/",handlers.PatchPostHandler)
	r.HandleFunc("DELETE /posts/{id}/",handlers.DeletePostHandler)

	fmt.Println("server started: http/localhost:8000/" )
	if err := http.ListenAndServe(":8000", r); err != nil{
		fmt.Println("start server error ", err.Error())
	}
	
}