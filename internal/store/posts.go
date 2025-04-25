package store

import (
	"context"
	"fmt"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Post struct {
	Id        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid" json:"user_id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      pq.StringArray `gorm:"type:varchar[]" json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	User      *User          `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
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
func (post_store *PostStore) GetPostFeed(ctx context.Context, userId uuid.UUID, params *model.PostParams) ([]*PostWithMetadata, int64, error) {
	var post_feed []*PostWithMetadata
	var total int64
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN user_followers f ON f.following_id = p.user_id AND f.follower_id = $1
		WHERE (p.user_id = $1 OR f.follower_id IS NOT NULL)
	`
	countArgs := []any{userId}
	countIndex := 2
	if params.Search != "" {
		countQuery += fmt.Sprintf("\nAND (p.title ILIKE $%d OR p.content ILIKE $%d)", countIndex, countIndex)
		countArgs = append(countArgs, "%"+params.Search+"%")
		countIndex++
	}
	if len(params.Tags) > 0 {
		countQuery += fmt.Sprintf("\nAND p.tags && $%d::varchar[]", countIndex)
		countArgs = append(countArgs, pq.Array(params.Tags))
		countIndex++
	}
	if params.Since != "" && params.Until != "" {
		countQuery += fmt.Sprintf("\nAND (p.created_at::date BETWEEN $%d AND $%d)", countIndex, countIndex+1)
		countArgs = append(countArgs, params.Since, params.Until)
		countIndex += 2
	}
	err := post_store.db.WithContext(ctx).
		Raw(countQuery, countArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	feedQuery := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.tags, p.created_at, p.updated_at,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN user_followers f ON f.following_id = p.user_id AND f.follower_id = $1
		WHERE (p.user_id = $1 OR f.follower_id IS NOT NULL)
	`
	feedArgs := []any{userId}
	feedIndex := 2
	if params.Search != "" {
		feedQuery += fmt.Sprintf("\nAND (p.title ILIKE $%d OR p.content ILIKE $%d)", feedIndex, feedIndex)
		feedArgs = append(feedArgs, "%"+params.Search+"%")
		feedIndex++
	}
	if len(params.Tags) > 0 {
		feedQuery += fmt.Sprintf("\nAND p.tags && $%d::varchar[]", feedIndex)
		feedArgs = append(feedArgs, pq.Array(params.Tags))
		feedIndex++
	}
	if params.Since != "" && params.Until != "" {
		feedQuery += fmt.Sprintf("\nAND (p.created_at::date BETWEEN $%d AND $%d)", feedIndex, feedIndex+1)
		feedArgs = append(feedArgs, params.Since, params.Until)
		feedIndex += 2
	}
	feedQuery += `
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + params.Sort + `
		LIMIT $` + fmt.Sprint(feedIndex) + ` OFFSET $` + fmt.Sprint(feedIndex+1)
	feedArgs = append(feedArgs, params.PerPage, (params.Page-1)*params.PerPage)
	err = post_store.db.WithContext(ctx).
		Raw(feedQuery, feedArgs...).Scan(&post_feed).Error
	if err != nil {
		return nil, 0, err
	}
	return post_feed, total, nil
}
