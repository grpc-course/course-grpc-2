package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	MDKeyName   = "authorization"
	UserKeyName = "user"
)

func ValidateAuthToken(token string) error {
	// jwt like validation
	return nil
}

func GetUserFromRequest(ctx context.Context) (*User, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	return &User{
		Username: "test_user",
	}, nil
}

func PutUserToContext(ctx context.Context, user *User) context.Context {
	ctx = context.WithValue(ctx, UserKeyName, user)
	return ctx
}

func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserKeyName).(*User)
	return user, ok
}

func CreateClientMD() metadata.MD {
	md := metadata.New(map[string]string{
		MDKeyName: "authorization key",
	})
	return md
}
