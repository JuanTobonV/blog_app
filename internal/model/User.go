package model

type User struct {
    Id           int       `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    Password     string    `json:"-" db:"password"` 
    TokenSession string    `json:"token_session,omitempty" db:"token_session"`
    RefreshToken string    `json:"refresh_token,omitempty" db:"refresh_token"`
    Blogs        []Blog    `json:"blogs,omitempty" db:"-"` 
}