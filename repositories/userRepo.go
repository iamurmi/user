package repositories

// Performe DB opperation
import (
	"context"
	"fmt"
	"small-mic/user/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "User"

type userRepoStruct struct {
	client *mongo.Client
}

func NewUserRepoConstructor(client *mongo.Client) *userRepoStruct {
	return &userRepoStruct{
		client: client,
	}
}

func (repo *userRepoStruct) AddUser(ctx context.Context, u domain.User) (id string, err error) {
	resp, err := repo.client.Database("SMALLMIC").Collection(userCollection).InsertOne(ctx, u)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(resp.InsertedID), nil
}

func (repo *userRepoStruct) GetUser(ctx context.Context, uId string) (user domain.User, err error) {
	err = repo.client.Database("SMALLMIC").Collection(userCollection).FindOne(ctx, bson.M{"_id": uId}).Decode(&user)
	if err != nil {
		return
	}
	return
}
func (repo *userRepoStruct) ListUser(ctx context.Context) (users []domain.User, err error) {
	resp, err := repo.client.Database("SMALLMIC").Collection(userCollection).Find(ctx, bson.M{})
	if err != nil {
		return
	}
	defer resp.Close(ctx)
	for resp.Next(ctx) {
		var user domain.User
		resp.Decode(&user)
		users = append(users, user)
	}
	return
}
