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

type PostWithMetadata struct {
	Post
	Username      string `json:"username"`
	CommentsCount int64  `json:"comments_count"`
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	post.Id = uuid.New()
	return
}

type PostStore struct {
	db *gorm.DB
}

func (post_store *PostStore) CreatePost(ctx context.Context, post *Post) error {
	err := post_store.db.WithContext(ctx).Create(&post).Error
	if err != nil {
		return err
	}
	return nil
}
func (post_store *PostStore) GetPost(ctx context.Context, postId uuid.UUID) (*Post, error) {
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
func (post_store *PostStore) UpdatePost(ctx context.Context, post *Post) (*Post, error) {
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
func (post_store *PostStore) DeletePost(ctx context.Context, postId uuid.UUID) error {
	var post Post
	err := post_store.db.WithContext(ctx).Where("id = ?", postId).Delete(&post).Error
	if err != nil {
		return err
	}
	return nil
}
func (post_store *PostStore) GetPostFeed(ctx context.Context, userId uuid.UUID, paginatedQuery *PaginatedFeedQuery) ([]*PostWithMetadata, int64, error) {
	var post_feed []*PostWithMetadata
	var total int64
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN user_followers f ON f.following_id = p.user_id AND f.follower_id = $1
		WHERE p.user_id = $1 OR f.follower_id IS NOT NULL
	`
	feedQuery := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN user_followers f ON f.following_id = p.user_id AND f.follower_id = $1
		WHERE p.user_id = $1 OR f.follower_id IS NOT NULL
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + paginatedQuery.Sort + `
		LIMIT $2 OFFSET $3;
	`
	err := post_store.db.WithContext(ctx).
		Raw(countQuery, userId).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = post_store.db.WithContext(ctx).
		Raw(feedQuery, userId, paginatedQuery.PerPage, (paginatedQuery.Page-1)*paginatedQuery.PerPage).
		Scan(&post_feed).Error
	if err != nil {
		return nil, 0, err
	}
	return post_feed, total, nil
}
