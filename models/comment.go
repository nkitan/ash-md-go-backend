package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type CommentModel struct {
    ID           uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    BlogID       uuid.UUID    `gorm:"type:uuid;not null" json:"blog_id"`
    BlogAuthor   uuid.UUID    `gorm:"type:uuid;not null" json:"blog_author"`
    Comment      string       `gorm:"type:text;not null" json:"comment"`
    Children     []uuid.UUID  `gorm:"type:uuid[]" json:"children,omitempty"`
    CommentedBy  uuid.UUID    `gorm:"type:uuid;not null" json:"commented_by"`
    IsReply      bool         `gorm:"default:false" json:"is_reply,omitempty"`
    Parent       *uuid.UUID   `gorm:"type:uuid" json:"parent,omitempty"`
    CommentedAt  time.Time    `gorm:"autoCreateTime" json:"commented_at"`
}

func CommentsMigrate(db *gorm.DB) {
    db.AutoMigrate(&CommentModel{})
}
