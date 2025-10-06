package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	TenantID  *uint          `json:"tenant_id" gorm:"index"`
	Role      string         `json:"role" gorm:"not null;default:'user'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Tenant *Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
	Posts  []Post  `json:"posts,omitempty" gorm:"foreignKey:UserID"`
}

type Tenant struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Domain    string         `json:"domain" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Users []User `json:"users,omitempty" gorm:"foreignKey:TenantID"`
	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:TenantID"`
}

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"type:text"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	TenantID  uint           `json:"tenant_id" gorm:"not null;index"`
	IsPublic  bool           `json:"is_public" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Tenant Tenant `json:"tenant" gorm:"foreignKey:TenantID"`
}
