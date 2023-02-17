package contracts

import (
	"context"
	userPB "small-mic/user/domain/protobuf"
)

type UserSeriveContract interface {
	AddUser(ctx context.Context, u *userPB.AddUserRequest) (res *userPB.AddUserResponse, err error) // declared and defined at service.go
	GetUser(ctx context.Context, req *userPB.GetUserRequest) (user *userPB.GetUserResponse, err error)
	GetUsers(ctx context.Context) (usersData *userPB.ListUsersResponse, err error)
}
