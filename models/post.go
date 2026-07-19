package models

import (
	"strings"
	"time"
)

type Post struct {
	Id          int       `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"column:Title;type:varchar(200);not null" json:"title"`
	Content     string    `gorm:"column:Content;type:text;not null" json:"content"`
	Category    string    `gorm:"column:Category;type:varchar(100);not null" json:"category"`
	CreatedDate time.Time `gorm:"column:Created_date;type:timestamp;autoCreateTime" json:"created_date"`
	UpdatedDate time.Time `gorm:"column:Updated_date;type:timestamp;autoUpdateTime" json:"updated_date"`
	Status      string    `gorm:"column:Status;type:varchar(100);not null" json:"status"`
}

// TableName overrides the default pluralized table name of GORM to map exactly to "posts"
func (Post) TableName() string {
	return "posts"
}

// Validate checks the structural constraints for the article
func (p *Post) Validate() map[string]string {
	errs := make(map[string]string)

	p.Title = strings.TrimSpace(p.Title)
	p.Content = strings.TrimSpace(p.Content)
	p.Category = strings.TrimSpace(p.Category)
	p.Status = strings.ToLower(strings.TrimSpace(p.Status))

	if p.Title == "" {
		errs["title"] = "Title is required"
	} else if len(p.Title) < 20 {
		errs["title"] = "Title must be at least 20 characters"
	}

	if p.Content == "" {
		errs["content"] = "Content is required"
	} else if len(p.Content) < 200 {
		errs["content"] = "Content must be at least 200 characters"
	}

	if p.Category == "" {
		errs["category"] = "Category is required"
	} else if len(p.Category) < 3 {
		errs["category"] = "Category must be at least 3 characters"
	}

	if p.Status == "" {
		errs["status"] = "Status is required"
	} else if p.Status != "publish" && p.Status != "draft" && p.Status != "thrash" {
		errs["status"] = "Status must be 'publish', 'draft', or 'thrash'"
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}
