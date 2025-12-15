package model

type Blog struct {
    Id        int    `json:"id" db:"id"`
    Title     string `json:"title" db:"title"`
    Content   string `json:"content" db:"content"`
    CreatedAt string `json:"created_at" db:"created_at"`
    UserId    int    `json:"user_id" db:"user_id"`
    Author    *User  `json:"author,omitempty" db:"-"` 
}