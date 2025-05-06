package test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aryaadinulfadlan/go-social-api/entity/post"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := post.NewRepository(gormDB)
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	postId, _ := uuid.Parse("023165f2-1893-4eea-b9b8-4b4d2c42ae2e")
	post := &db.Post{
		Id:        postId,
		UserId:    userId,
		Title:     "Title here",
		Content:   "Post content here",
		Tags:      pq.StringArray{"golang", "gorm", "testing"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts" ("user_id","title","content","tags","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(post.UserId, post.Title, post.Content, pq.Array(post.Tags), post.CreatedAt, post.UpdatedAt, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(postId.String()))
	mock.ExpectCommit()
	err := repo.Create(ctx, post)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
func TestGetPostById(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := post.NewRepository(gormDB)
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	postId, _ := uuid.Parse("07112266-83d1-408f-bd19-ec4b80f760a0")
	expectedPost := &db.Post{
		Id:      postId,
		UserId:  userId,
		Title:   "Princess First Post",
		Content: "The first post by Princess Diana",
		Tags:    pq.StringArray{"princess-diana", "first"},
	}
	postRows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "content", "tags",
	}).AddRow(
		postId,
		userId,
		"Princess First Post",
		"The first post by Princess Diana",
		pq.StringArray{"princess-diana", "first"},
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE id = $1 LIMIT $2`)).
		WithArgs(postId, 1).
		WillReturnRows(postRows)
	post, err := repo.GetById(ctx, postId)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedPost.Id, post.Id)
	assert.Equal(t, expectedPost.UserId, post.UserId)
	assert.Equal(t, expectedPost.Title, post.Title)
	assert.Equal(t, expectedPost.Content, post.Content)
	assert.Equal(t, expectedPost.Tags, post.Tags)
}
func TestDeletePostById(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := post.NewRepository(gormDB)
	postId, _ := uuid.Parse("07112266-83d1-408f-bd19-ec4b80f760a0")
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "posts" WHERE id = $1`)).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.Delete(ctx, postId)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
