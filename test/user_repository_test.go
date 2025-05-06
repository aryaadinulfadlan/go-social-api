package test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aryaadinulfadlan/go-social-api/entity/user"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserById(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := user.NewRepository(gormDB)
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	roleId, _ := uuid.Parse("4b30ed16-06bc-4f7f-8293-6cb8a040267e")
	expectedUser := &db.User{
		Id:       userId,
		RoleId:   roleId,
		Name:     "Princess Diana",
		Username: "princess_diana",
		Email:    "princess_diana@gmail.com",
	}
	expectedRole := &db.Role{
		Id:   roleId,
		Name: "Admin",
	}
	userRows := sqlmock.NewRows([]string{
		"id", "role_id", "name", "username", "email",
	}).AddRow(
		userId,
		roleId,
		"Princess Diana",
		"princess_diana",
		"princess_diana@gmail.com",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 LIMIT $2`)).
		WithArgs(userId, 1).
		WillReturnRows(userRows)
	roleRows := sqlmock.NewRows([]string{
		"id", "name",
	}).AddRow(
		roleId,
		"Admin",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE "roles"."id" = $1`)).
		WithArgs(roleId).
		WillReturnRows(roleRows)

	user, err := repo.GetById(ctx, userId)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.RoleId, user.RoleId)
	assert.Equal(t, expectedUser.Name, user.Name)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedRole.Id, user.RoleId)
	assert.Equal(t, expectedRole.Name, user.Role.Name)
}
func TestDeleteUserById(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := user.NewRepository(gormDB)
	userId, _ := uuid.Parse("50b466de-2de4-4e40-bdec-08270f23a8c8")
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
		WithArgs(userId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	err := repo.Delete(ctx, userId)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
