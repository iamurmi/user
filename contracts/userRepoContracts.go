package contracts

import (
	"context"

	"github.com/iamurmi/user/domain"
)

type UserRepoContracts interface {
	AddUser(ctx context.Context, u domain.User) (id string, err error) // Function Declare, not contain body
	GetUser(ctx context.Context, uId string) (user domain.User, err error)
	ListUser(ctx context.Context) (users []domain.User, err error)
}
