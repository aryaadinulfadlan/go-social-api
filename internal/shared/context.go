package shared

import (
	"net/http"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
)

type userKey string

var UserCtx userKey = "user"

func GetUserFromContext(r *http.Request) *db.User {
	user, ok := r.Context().Value(UserCtx).(*db.User)
	if !ok {
		return nil
	}
	return user
}
