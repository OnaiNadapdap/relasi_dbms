package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title    string   `json:"title" form:"title" gorm:"not null"`
	Content  string   `json:"content" form:"content" gorm:"not null"`
	Tags     []Tag    `json:"tags" gorm:"many2many:post_tags;foreignKey:id;joinForeignKey:post_id;references:id;joinReferences:tag_id"`
	TagNames []string `json:"tags_name" form:"tags_name" gorm:"-"`
}

type Tag struct {
	gorm.Model
	Name  string `json:"name" gorm:"not null"`
	Posts []Post `json:"posts" gorm:"many2many:post_tags;foreignKey:id;joinForeignKey:tag_id;references:id;joinReferences:post_id"`
}

type PostTag struct {
	PostID uint `json:"post_id"`
	TagID  uint `json:"tag_id"`
}

type AppendTagToPost struct {
	TagNames []string `json:"tags_name" form:"tags_name" gorm:"-"`
}
