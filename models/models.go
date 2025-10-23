package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 1. User Model
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	Email        string    `gorm:"size:255;not null;unique" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"` // Exclude from JSON responses
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 2. Tag Model
type Tag struct {
	ID    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name  string    `gorm:"size:100;not null;unique" json:"name"`
	Posts []*Post   `gorm:"many2many:post_tags;" json:"posts,omitempty"` // Many-to-Many relationship with Posts
}

// 3. Post Model
type Post struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Title            string    `gorm:"size:200;not null" json:"title"`
	Content          string    `gorm:"type:text;not null" json:"content"`
	Category         string    `gorm:"size:100;not null" json:"category"`
	Status           string    `gorm:"size:50;not null;default:'draft'" json:"status"`
	FeaturedImageURL string    `gorm:"type:text" json:"featured_image_url"`

	// Author Relationship (Many-to-One)
	AuthorID uuid.UUID `gorm:"not null" json:"author_id"`
	Author   User      `gorm:"foreignKey:AuthorID" json:"author"`

	// Tags Relationship (Many-to-Many)
	Tags []*Tag `gorm:"many2many:post_tags;" json:"tags"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // For soft deletes
}

// We don't need to create a struct for 'post_tags'.
// GORM will handle it automatically based on the tag `gorm:"many2many:post_tags;"`.
