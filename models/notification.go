package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type NotificationType string

const (
    Like    NotificationType = "like"
    Comment NotificationType = "comment"
    Reply   NotificationType = "reply"
)

type Notification struct {
    ID                 uuid.UUID         `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    Type               NotificationType  `gorm:"type:varchar(20);not null" json:"type"`
    BlogID             uuid.UUID         `gorm:"type:uuid;not null" json:"blog"`
    NotificationFor    uuid.UUID         `gorm:"type:uuid;not null" json:"notification_for"`
    UserID             uuid.UUID         `gorm:"type:uuid;not null" json:"user"`
    CommentID          *uuid.UUID        `gorm:"type:uuid" json:"comment,omitempty"`
    ReplyID            *uuid.UUID        `gorm:"type:uuid" json:"reply,omitempty"`
    RepliedOnCommentID *uuid.UUID        `gorm:"type:uuid" json:"replied_on_comment,omitempty"`
    Seen               bool              `gorm:"default:false" json:"seen,omitempty"`
    CreatedAt          time.Time         `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt          time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

func NotificationsMigrate(db *gorm.DB) {
    db.AutoMigrate(&Notification{})
}
