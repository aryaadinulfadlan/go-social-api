package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      *User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Comments  []Comment      `json:"comments,omitempty"`
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	post.Id = uuid.New()
	return
}

type PostStore struct {
	db *gorm.DB
}

func (post_store *PostStore) Create(ctx context.Context, post *Post) error {
	err := post_store.db.WithContext(ctx).Create(&post).Error
	if err != nil {
		return err
	}
	return nil
}
func (post_store *PostStore) GetById(ctx context.Context, postId uuid.UUID) (*Post, error) {
	var post Post
	err := post_store.db.WithContext(ctx).Preload("User").Preload("Comments").Take(&post, "id = ?", postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (post_store *PostStore) CheckPostExists(ctx context.Context, postId uuid.UUID) (*Post, error) {
	var post Post
	err := post_store.db.WithContext(ctx).Take(&post, "id = ?", postId).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (post_store *PostStore) UpdateById(ctx context.Context, post *Post) (*Post, error) {
	err := post_store.db.WithContext(ctx).Model(&Post{}).
		Where("id = ?", post.Id).
		Updates(Post{
			Title:   post.Title,
			Content: post.Content,
			Tags:    post.Tags,
		}).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}
