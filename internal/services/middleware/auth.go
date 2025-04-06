package middleware

import (
    "context"
    "errors"
)

var (
    ErrUnauthorized = errors.New("unauthorized")
)

type ContextKey string

const UserIDKey ContextKey = "user_id"
const CurrentUserKey ContextKey = "current_user"

func Auth(next func(context.Context) (interface{}, error)) func(context.Context) (interface{}, error) {
    return func(ctx context.Context) (interface{}, error) {
        userID := ctx.Value(ContextKey("user_id"))
        if userID == nil {
            return nil, ErrUnauthorized
        }
        return next(ctx)
    }
}
