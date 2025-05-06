package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/entity/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Repository interface {
	CreateAndInvite(ctx context.Context, user *db.User, user_invitation *db.UserInvitation) error
	GetDetail(ctx context.Context, userId uuid.UUID) (*db.User, error)
	GetByInvitation(ctx context.Context, token string) (*db.User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*db.User, error)
	GetByUsernameEmail(ctx context.Context, username string, email string) (*db.User, error)
	FollowUnfollow(ctx context.Context, targetId uuid.UUID, senderId uuid.UUID) (string, error)
	GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*db.User, error)
	Activate(ctx context.Context, user *db.User) (*db.User, error)
	Delete(ctx context.Context, userId uuid.UUID) error
	GetFeeds(ctx context.Context, userId uuid.UUID, params *post.PostParams) ([]*db.PostWithMetadata, int64, error)
}

type RepositoryImplementation struct {
	gorm *gorm.DB
}

func NewRepository(gorm *gorm.DB) Repository {
	return &RepositoryImplementation{gorm: gorm}
}

func (repository *RepositoryImplementation) CreateAndInvite(ctx context.Context, user *db.User, user_invitation *db.UserInvitation) error {
	err := repository.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&user).Error
		if err != nil {
			return err
		}
		err = tx.Create(&user_invitation).Error
		if err != nil {
			return err
		}
		user.UserInvitation = user_invitation
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryImplementation) GetDetail(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	var user db.User
	err := repository.gorm.WithContext(ctx).
		Preload("Following").Preload("Followers").Preload("Role").
		Take(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (repository *RepositoryImplementation) GetByInvitation(ctx context.Context, token string) (*db.User, error) {
	var user db.User
	err := repository.gorm.WithContext(ctx).
		Joins("JOIN user_invitations ui ON ui.user_id = users.id").
		Where("ui.token = ? AND ui.expired_at > ?", token, time.Now().UTC()).
		// Preload("UserInvitation").
		Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (repository *RepositoryImplementation) GetById(ctx context.Context, userId uuid.UUID) (*db.User, error) {
	var user db.User
	err := repository.gorm.WithContext(ctx).Where("id = ?", userId).
		Preload("Role").Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (repository *RepositoryImplementation) GetByUsernameEmail(ctx context.Context, username string, email string) (*db.User, error) {
	var user db.User
	err := repository.gorm.WithContext(ctx).Where("username = ? OR email = ?", username, email).
		Preload("Role").Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (repository *RepositoryImplementation) FollowUnfollow(ctx context.Context, targetId uuid.UUID, senderId uuid.UUID) (string, error) {
	// err := user_store.db.WithContext(ctx).Select("id").Model(&User{Id: senderId}).Association("Following").Append(&User{Id: targetId})
	// relation := make(map[string]interface{})
	var message string
	var count int64
	err := repository.gorm.WithContext(ctx).
		Table("user_followers").
		Where("follower_id = ? AND following_id = ?", senderId, targetId).
		Count(&count).Error
	if err != nil {
		return "", err
	}
	if count > 0 {
		err_delete := repository.gorm.WithContext(ctx).
			Table("user_followers").
			Where("follower_id = ? AND following_id = ?", senderId, targetId).
			Delete(nil).Error
		if err_delete != nil {
			return "", err_delete
		}
		message = "Successfully Unfollowed"
	} else {
		err_insert := repository.gorm.WithContext(ctx).
			Table("user_followers").
			Create(map[string]any{
				"follower_id":  senderId,
				"following_id": targetId,
			}).Error
		if err_insert != nil {
			return "", err_insert
		}
		message = "Successfully Followed"
	}
	return message, nil
}
func (repository *RepositoryImplementation) GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*db.User, error) {
	var users []*db.User
	err := repository.gorm.Model(&db.User{Id: userId}).Association(actionType).Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (repository *RepositoryImplementation) Activate(ctx context.Context, user *db.User) (*db.User, error) {
	err := repository.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&db.User{}).
			Where("id = ?", user.Id).
			Updates(db.User{
				IsActivated: user.IsActivated,
			}).Error
		if err != nil {
			return err
		}
		err = tx.Where("user_id = ?", user.Id).Delete(&db.UserInvitation{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repository *RepositoryImplementation) Delete(ctx context.Context, userId uuid.UUID) error {
	var user *db.User
	err := repository.gorm.WithContext(ctx).Where("id = ?", userId).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (repository *RepositoryImplementation) GetFeeds(ctx context.Context, userId uuid.UUID, params *post.PostParams) ([]*db.PostWithMetadata, int64, error) {
	var post_feed []*db.PostWithMetadata
	var total int64
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN user_followers f ON f.following_id = p.user_id AND f.follower_id = $1
		WHERE (p.user_id = $1 OR f.follower_id IS NOT NULL)
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
		WHERE (p.user_id = $1 OR f.follower_id IS NOT NULL)
	`
	countArgs := []any{userId}
	countIndex := 2
	feedArgs := []any{userId}
	feedIndex := 2

	if params.Search != "" {
		countQuery += fmt.Sprintf("\nAND (p.title ILIKE $%d OR p.content ILIKE $%d)", countIndex, countIndex)
		countArgs = append(countArgs, "%"+params.Search+"%")
		countIndex++
		feedQuery += fmt.Sprintf("\nAND (p.title ILIKE $%d OR p.content ILIKE $%d)", feedIndex, feedIndex)
		feedArgs = append(feedArgs, "%"+params.Search+"%")
		feedIndex++
	}
	if len(params.Tags) > 0 {
		countQuery += fmt.Sprintf("\nAND p.tags && $%d::varchar[]", countIndex)
		countArgs = append(countArgs, pq.Array(params.Tags))
		countIndex++
		feedQuery += fmt.Sprintf("\nAND p.tags && $%d::varchar[]", feedIndex)
		feedArgs = append(feedArgs, pq.Array(params.Tags))
		feedIndex++
	}
	if params.Since != "" && params.Until != "" {
		countQuery += fmt.Sprintf("\nAND (p.created_at::date BETWEEN $%d AND $%d)", countIndex, countIndex+1)
		countArgs = append(countArgs, params.Since, params.Until)
		countIndex += 2
		feedQuery += fmt.Sprintf("\nAND (p.created_at::date BETWEEN $%d AND $%d)", feedIndex, feedIndex+1)
		feedArgs = append(feedArgs, params.Since, params.Until)
		feedIndex += 2
	}
	err := repository.gorm.WithContext(ctx).
		Raw(countQuery, countArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}
	feedQuery += `
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + params.Sort + `
		LIMIT $` + fmt.Sprint(feedIndex) + ` OFFSET $` + fmt.Sprint(feedIndex+1)
	feedArgs = append(feedArgs, params.PerPage, (params.Page-1)*params.PerPage)
	err = repository.gorm.WithContext(ctx).
		Raw(feedQuery, feedArgs...).Scan(&post_feed).Error
	if err != nil {
		return nil, 0, err
	}
	return post_feed, total, nil
}
