package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
)


type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

var inMemoryDB = make([]Post, 0)
var counter = 1

func main() {
	r := chi.NewRouter()

	r.Get("/posts", GetPosts)
	r.Get("/posts/{id}", GetPost)
	r.Post("/posts", CreatePost)
	r.Put("/posts/{id}", UpdatePost)
	r.Delete("/posts/{id}", DeletePost)

	http.ListenAndServe(":8080", r)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inMemoryDB)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") //url = universal resource locator
	postID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	for _, post := range inMemoryDB {
		if post.ID == postID {
			// w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(post)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newPost.ID = counter
	counter++
	inMemoryDB = append(inMemoryDB, newPost)
	json.NewEncoder(w).Encode("New Post Created !")
	w.WriteHeader(http.StatusCreated)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	for i, post := range inMemoryDB {
		if post.ID == postID {
			var updatedPost Post
			if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			updatedPost.ID = postID
			inMemoryDB[i] = updatedPost

			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	for i, post := range inMemoryDB {
		if post.ID == postID {
			inMemoryDB = append(inMemoryDB[:i], inMemoryDB[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Post not found", http.StatusNotFound)
}
