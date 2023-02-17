package userservice

// Performe Logics and Call repositories methods
import (
	"context"
	"errors"
	"fmt"
	"small-mic/user/contracts"
	"small-mic/user/domain"
	userPB "small-mic/user/domain/protobuf"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repoContract      contracts.UserRepoContracts
	redisRepoContract contracts.UserRepoContracts
}

func NewServiceConstructor(repoContract contracts.UserRepoContracts, redisRepoContract contracts.UserRepoContracts) *service {
	return &service{
		repoContract:      repoContract,
		redisRepoContract: redisRepoContract,
	}
}

func (svc *service) AddUser(ctx context.Context, u *userPB.AddUserRequest) (res *userPB.AddUserResponse, err error) {
	//Logic
	u.User.Id = string(primitive.NewObjectID().Hex())
	// Conversion of classes
	user := domain.User{
		ID:        u.User.Id,
		FirstName: u.User.FirstName,
		Roles:     u.User.Roles,
	}
	res.Id, err = svc.repoContract.AddUser(ctx, user)
	if err != nil {
		return
	}
	// cache write
	_, err = svc.redisRepoContract.AddUser(ctx, user)
	return
}
func (svc *service) GetUser(ctx context.Context, req *userPB.GetUserRequest) (user *userPB.GetUserResponse, err error) {
	//Logic
	if req.Id == "" {
		err = errors.New("Bad Request")
		return
	}
	// cache get
	var u domain.User
	u, err = svc.redisRepoContract.GetUser(ctx, req.Id)
	if err != nil {
		return
	}
	// Conversion of classes
	user = &userPB.GetUserResponse{
		User: &userPB.UserData{
			Id:        u.ID,
			FirstName: u.FirstName,
			Roles:     u.Roles,
		},
	}
	if u.ID == req.Id { // cache HIT
		fmt.Println("CACHE HIT")
		return
	}
	// cache miss
	fmt.Println("CACHE MISS")
	u, err = svc.repoContract.GetUser(ctx, req.Id)
	if err != nil {
		return
	}
	// cache write
	_, err = svc.redisRepoContract.AddUser(ctx, u)
	return
}
func (svc *service) GetUsers(ctx context.Context, _ *userPB.ListUsersRequest) (usersData *userPB.ListUsersResponse, err error) {
	//Logic
	// cache get
	users, err := svc.redisRepoContract.ListUser(ctx)
	if err != nil {
		return
	}
	if len(users) > 0 { // cache HIT
		fmt.Println("CACHE HIT")
		for i := range users {
			usersData.Users = append(usersData.Users, &userPB.UserData{
				Id:        users[i].ID,
				FirstName: users[i].FirstName,
				Roles:     users[i].Roles,
			})
		}
		return
	}
	// cache miss
	fmt.Println("CACHE MISS")
	users, err = svc.repoContract.ListUser(ctx)
	if err != nil {
		return
	}
	// cache write
	for _, item := range users {
		_, err = svc.redisRepoContract.AddUser(ctx, item)
	}
	for i := range users {
		usersData.Users = append(usersData.Users, &userPB.UserData{
			Id:        users[i].ID,
			FirstName: users[i].FirstName,
			Roles:     users[i].Roles,
		})
	}
	return
}
