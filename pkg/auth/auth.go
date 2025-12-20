package auth

import (
	"google.golang.org/grpc/metadata"
)

const (
	MDKeyName = "authorization"
)

func ValidateAuthToken(token string) error {
	// jwt like validation
	return nil
}

func CreateClientMD() metadata.MD {
	md := metadata.New(map[string]string{
		MDKeyName: "authorization key",
	})
	return md
}
