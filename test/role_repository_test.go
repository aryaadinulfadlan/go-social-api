package test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aryaadinulfadlan/go-social-api/entity/role"
	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetRoleByName(t *testing.T) {
	gormDB, mock := SetupMockDB(t)
	repo := role.NewRepository(gormDB)
	roleId, _ := uuid.Parse("4b30ed16-06bc-4f7f-8293-6cb8a040267e")
	expectedRole := &db.Role{
		Id:          roleId,
		Name:        "Manager",
		Description: "Has full access to manage users, content, and system settings.",
	}
	roleRows := sqlmock.NewRows([]string{
		"id", "name", "description",
	}).AddRow(
		roleId,
		"Manager",
		"Has full access to manage users, content, and system settings.",
	)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE name = $1 LIMIT $2`)).
		WithArgs("Manager", 1).
		WillReturnRows(roleRows)
	role, err := repo.GetRole(ctx, "Manager")
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, expectedRole.Id, role.Id)
	assert.Equal(t, expectedRole.Name, role.Name)
	assert.Equal(t, expectedRole.Description, role.Description)
}
