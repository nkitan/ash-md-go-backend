package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PersonalInfo struct {
    Fullname    string `gorm:"type:varchar(255);not null" json:"fullname"`
    Email       string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
    Password    string `gorm:"type:varchar(255)" json:"password,omitempty"`
    Username    string `gorm:"type:varchar(255);uniqueIndex" json:"username,omitempty"`
    Bio         string `gorm:"type:varchar(200);default:''" json:"bio,omitempty"`
    ProfileImg string `gorm:"type:varchar(255)" json:"profile_img,omitempty"`
}

type SocialLinks struct {
    Youtube   string `gorm:"type:varchar(255);default:''" json:"youtube,omitempty"`
    Instagram string `gorm:"type:varchar(255);default:''" json:"instagram,omitempty"`
    Facebook  string `gorm:"type:varchar(255);default:''" json:"facebook,omitempty"`
    Twitter   string `gorm:"type:varchar(255);default:''" json:"twitter,omitempty"`
    Github    string `gorm:"type:varchar(255);default:''" json:"github,omitempty"`
    Website   string `gorm:"type:varchar(255);default:''" json:"website,omitempty"`
}

type AccountInfo struct {
    TotalPosts int `gorm:"default:0" json:"total_posts,omitempty"`
    TotalReads int `gorm:"default:0" json:"total_reads,omitempty"`
}

type User struct {
    ID           uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
    PersonalInfo PersonalInfo `gorm:"embedded" json:"personal_info"`
    SocialLinks  SocialLinks  `gorm:"embedded" json:"social_links"`
    AccountInfo  AccountInfo  `gorm:"embedded" json:"account_info"`
    GoogleAuth   bool         `gorm:"default:false" json:"google_auth,omitempty"`
    Blogs        []Blog       `gorm:"foreignKey:AuthorID" json:"blogs,omitempty"`
    JoinedAt     time.Time    `gorm:"autoCreateTime" json:"joinedAt"`
}

func UsersMigrate(db *gorm.DB) {
    db.AutoMigrate(&User{})
}
