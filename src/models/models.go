package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"unique"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	Posts     []Post   // ความสัมพันธ์แบบ HasMany
	Followers []Follow `gorm:"foreignKey:FollowingID;"`
	Following []Follow `gorm:"foreignKey:FollowerID;"`
}

type Post struct {
	gorm.Model
	UserID   uint   // นี่คือ foreign key ที่เชื่อมโยงกับ User
	User     User   // ความสัมพันธ์ระหว่าง Post และ User
	Content  string `gorm:"type:text;not null"`
	ImageURL string `gorm:"type:varchar(255)"`
	Comments []Comment
	Likes    []Like
}

type Comment struct {
	gorm.Model
	PostID      uint
	UserID      uint
	CommentText string `gorm:"type:text;not null"`
	User        User
}

type Like struct {
	gorm.Model
	PostID uint
	UserID uint
}

type Follow struct {
	gorm.Model
	FollowerID  uint
	FollowingID uint
}
