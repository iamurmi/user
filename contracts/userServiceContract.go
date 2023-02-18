package contracts

import (
	"context"

	userPB "github.com/iamurmi/user/domain/protobuf"
)

type UserSeriveContract interface {
	AddUser(ctx context.Context, u *userPB.AddUserRequest) (res *userPB.AddUserResponse, err error) // declared and defined at service.go
	GetUser(ctx context.Context, req *userPB.GetUserRequest) (user *userPB.GetUserResponse, err error)
	GetUsers(ctx context.Context, _ *userPB.ListUsersRequest) (usersData *userPB.ListUsersResponse, err error)
}
