package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id             uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name           string          `json:"name"`
	Username       string          `gorm:"type:citext;unique" json:"username"`
	Email          string          `gorm:"type:citext;unique" json:"email"`
	Password       string          `json:"-"`
	IsActivated    bool            `json:"is_activated"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Posts          []Post          `json:"posts,omitempty"`
	UserInvitation *UserInvitation `json:"user_invitation,omitempty"`
	Comments       []Comment       `json:"comments,omitempty"`
	Following      []*User         `gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:following_id" json:"following,omitempty"`
	Followers      []*User         `gorm:"many2many:user_followers;joinForeignKey:following_id;joinReferences:follower_id" json:"followers,omitempty"`
}

type UserStore struct {
	db *gorm.DB
}

func (user_store *UserStore) CreateUserAndInvite(ctx context.Context, user *User, user_invitation *UserInvitation) error {
	err := user_store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (user_store *UserStore) GetUser(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	err := user_store.db.WithContext(ctx).Preload("Following").Preload("Followers").Take(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user_store *UserStore) GetUserByInvitation(ctx context.Context, token string) (*User, error) {
	var user User
	err := user_store.db.WithContext(ctx).
		Joins("JOIN user_invitations ui ON ui.user_id = users.id").
		Where("ui.token = ? AND ui.expired_at > ?", token, time.Now().UTC()).
		// Preload("UserInvitation").
		Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user_store *UserStore) GetExistingUser(ctx context.Context, field string, value any) (*User, error) {
	var user User
	validFields := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
	}
	if !validFields[field] {
		return nil, errors.New("invalid field name")
	}
	query := fmt.Sprintf("%s = ?", field)
	err := user_store.db.WithContext(ctx).Where(query, value).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (user_store *UserStore) IsUserExists(ctx context.Context, field string, value any) (int64, error) {
	var count int64
	validFields := map[string]bool{
		"id":       true,
		"username": true,
		"email":    true,
	}
	if !validFields[field] {
		return 0, errors.New("invalid field name")
	}
	query := fmt.Sprintf("%s = ?", field)
	err := user_store.db.WithContext(ctx).Model(&User{}).Where(query, value).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (user_store *UserStore) FollowUnfollowUser(ctx context.Context, targetId uuid.UUID, senderId uuid.UUID) error {
	// err := user_store.db.WithContext(ctx).Select("id").Model(&User{Id: senderId}).Association("Following").Append(&User{Id: targetId})
	// relation := make(map[string]interface{})
	var count int64
	err := user_store.db.WithContext(ctx).
		Table("user_followers").
		Where("follower_id = ? AND following_id = ?", senderId, targetId).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		err_delete := user_store.db.WithContext(ctx).
			Table("user_followers").
			Where("follower_id = ? AND following_id = ?", senderId, targetId).
			Delete(nil).Error
		if err_delete != nil {
			return err_delete
		}
	} else {
		err_insert := user_store.db.WithContext(ctx).
			Table("user_followers").
			Create(map[string]any{
				"follower_id":  senderId,
				"following_id": targetId,
			}).Error
		if err_insert != nil {
			return err_insert
		}
	}
	return nil
}

func (user_store *UserStore) GetConnections(ctx context.Context, userId uuid.UUID, actionType string) ([]*User, error) {
	var users []*User
	err := user_store.db.Model(&User{Id: userId}).Association(actionType).Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (user_store *UserStore) ActivateUser(ctx context.Context, user *User) (*User, error) {
	err := user_store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&User{}).
			Where("id = ?", user.Id).
			Updates(User{
				IsActivated: user.IsActivated,
			}).Error
		if err != nil {
			return err
		}
		err = tx.Where("user_id = ?", user.Id).Delete(&UserInvitation{}).Error
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

func (user_store *UserStore) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	var user User
	err := user_store.db.WithContext(ctx).Where("id = ?", userId).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}
