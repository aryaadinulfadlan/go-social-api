package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string    `json:"name"`
	Username  string    `gorm:"type:citext;unique" json:"username"`
	Email     string    `gorm:"type:citext;unique" json:"email"`
	Password  string    `gorm:"type:bytea" json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `json:"posts,omitempty"`
	Comments  []Comment `json:"comments,omitempty"`
	Following []*User   `gorm:"many2many:user_followers;joinForeignKey:follower_id;joinReferences:following_id" json:"following"`
	Followers []*User   `gorm:"many2many:user_followers;joinForeignKey:following_id;joinReferences:follower_id" json:"followers"`
}

func (user *User) BeforeCreate(db *gorm.DB) (err error) {
	user.Id = uuid.New()
	return
}

type UserStore struct {
	db *gorm.DB
}

func (user_store *UserStore) CreateUser(ctx context.Context, user *User) error {
	err := user_store.db.WithContext(ctx).Create(&user).Error
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
func (user_store *UserStore) CheckUserExists(ctx context.Context, userId uuid.UUID) (*User, error) {
	var user User
	err := user_store.db.WithContext(ctx).Take(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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
