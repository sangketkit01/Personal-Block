package models

import "time"

//nolint

type User struct {
	ID    int
	Username  string
	Email     string
	Phone     string
	Name string
	Password  string
	ProfileImage string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Block struct{
	ID int
	UserID int
	Title string
	Content string
	CreatedAt time.Time
	UpdatedAt time.Time
	User User
}

type Comment struct{
	ID int
	UserID int
	BlockID int
	Content string
	CreatedAt time.Time
	UpdatedAt time.Time
	User User
	Block Block
}

type BlockLike struct{
	ID int
	BlockID int
	UserID int
	Block Block
	User User
}

type CommentLike struct{
	ID int
	UserID int
	CommentID int
	BlockID int
	User User
	Comment Comment
}
