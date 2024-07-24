package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Activity struct {
    TotalLikes          int `gorm:"default:0" json:"total_likes,omitempty"`
    TotalComments       int `gorm:"default:0" json:"total_comments,omitempty"`
    TotalReads          int `gorm:"default:0" json:"total_reads,omitempty"`
    TotalParentComments int `gorm:"default:0" json:"total_parent_comments,omitempty"`
}

type Blog struct {
    BlogID      string      `gorm:"type:varchar(255);primary_key" json:"blog_id"`
    Title       string      `gorm:"type:varchar(255);not null" json:"title"`
    Banner      string      `gorm:"type:varchar(255)" json:"banner,omitempty"`
    Des         string      `gorm:"type:varchar(200)" json:"des,omitempty"`
    Content     []byte      `gorm:"type:bytea" json:"content,omitempty"`
    Tags        []byte      `gorm:"type:jsonb" json:"tags,omitempty"`  // Changed to jsonb
    AuthorID    uuid.UUID   `gorm:"type:uuid;not null" json:"author"`
    Author      User        `gorm:"foreignKey:AuthorID" json:"author"`
    Activity    Activity    `gorm:"embedded" json:"activity,omitempty"`
    Comments    []uuid.UUID `gorm:"type:uuid[]" json:"comments,omitempty"`
    Draft       bool        `gorm:"default:false" json:"draft,omitempty"`
    PublishedAt time.Time   `gorm:"autoCreateTime" json:"published_at"`
}

func BlogsMigrate(db *gorm.DB) {
    db.AutoMigrate(&Blog{})
}
