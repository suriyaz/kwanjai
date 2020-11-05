package models

import (
	"kwanjai/libraries"
	"net/http"
	"time"
)

// Post model.
type Post struct {
	ID           string     `json:"id"`
	Board        string     `json:"board" binding:"required"`
	Project      string     `json:"project"` // It's an advantage for checking project membership.
	User         string     `json:"username"`
	Title        string     `json:"title" binding:"required"`
	Body         string     `json:"body" binding:"required"`
	Comments     []*Comment `json:"comments"`
	Completed    bool       `json:"is_completed"`
	Urgent       bool       `json:"is_urgent"`
	People       []string   `json:"people"`
	AddedDate    time.Time  `json:"added_date"`
	LastModified time.Time  `json:"last_modified"`
}

// Comment model.
type Comment struct {
	UUID         string    `json:"uuid"`
	User         string    `json:"username"`
	Body         string    `json:"body"`
	AddedDate    time.Time `json:"added_date"`
	LastModified time.Time `json:"last_modified"`
}

func (post *Post) initialize() {
	post.People = []string{}
	post.Comments = []*Comment{}
	now := time.Now().Truncate(time.Millisecond)
	post.AddedDate = now
	post.LastModified = now
}

func (post *Post) CreatePost() (int, string, *Post) {
	post.initialize()
	reference, _, err := libraries.FirestoreAdd("posts", post)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	go libraries.FirestoreDeleteField("posts", reference.ID, "ID")
	post.ID = reference.ID
	return http.StatusCreated, "Created post.", post
}

func (post *Post) FindPost() (int, string, *Post) {
	if post.ID == "" {
		return http.StatusNotFound, "Post not found.", nil
	}
	getPost, _ := libraries.FirestoreFind("posts", post.ID)
	if getPost.Exists() {
		getPost.DataTo(post)
		post.ID = getPost.Ref.ID
		return http.StatusOK, "Get post successfully.", post
	}
	return http.StatusNotFound, "Post not found.", nil
}

func (post *Post) UpdatePost(field string, value interface{}) (int, string, *Post) {
	_, err := libraries.FirestoreUpdateField("posts", post.ID, field, value)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", post
}

func (post *Post) DeletePost() (int, string, *Post) {
	_, err := libraries.FirestoreDelete("posts", post.ID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
