package repocache

import (
	"context"
	"encoding/json"
	"fmt"
	"small-mic/user/domain"

	"github.com/go-redis/redis/v7"
)

type redisSvc struct {
	redisClient *redis.Client
}

func NewRedisClient(redisClient *redis.Client) *redisSvc {
	return &redisSvc{
		redisClient: redisClient,
	}
}

func (cacheSvc *redisSvc) AddUser(ctx context.Context, u domain.User) (id string, err error) {
	payload, err := json.Marshal(u)
	if err != nil {
		return
	}
	cacheSvc.redisClient.HSet("SMALLMIC", u.ID, payload)
	id = u.ID
	fmt.Println("ADDED in REDIS")
	return
}
func (cacheSvc *redisSvc) GetUser(ctx context.Context, uId string) (user domain.User, err error) {
	payload, err := cacheSvc.redisClient.HGet("SMALLMIC", uId).Result()
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(payload), &user)
	if err != nil {
		return
	}
	fmt.Println("GET from REDIS")
	return
}
func (cacheSvc *redisSvc) ListUser(ctx context.Context) (users []domain.User, err error) {
	payload, err := cacheSvc.redisClient.HGetAll("SMALLMIC").Result()
	if err != nil {
		return
	}
	for _, item := range payload {
		var user domain.User
		err = json.Unmarshal([]byte(item), &user)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	fmt.Println("GET from REDIS")
	return
}

/*

func (cache redisCache) CreateMovie(movie *Movie) (*Movie, error) {
   c := cache.getClient()
   movie.Id = uuid.New().String()
   json, err := json.Marshal(movie)
   if err != nil {
       return nil, err
   }
   c.HSet("movies", movie.Id, json)
   if err != nil {
       return nil, err
   }
   return movie, nil
}

func (cache redisCache) GetMovie(id string) (*Movie, error) {
   c := cache.getClient()
   val, err := c.HGet("movies", id).Result()

   if err != nil {
       return nil, err
   }
   movie := &Movie{}
   err = json.Unmarshal([]byte(val), movie)

   if err != nil {
       return nil, err
   }
   return movie, nil
}

func (cache redisCache) GetMovies() ([]*Movie, error) {
   c := cache.getClient()
   movies := []*Movie{}
   val, err := c.HGetAll("movies").Result()
   if err != nil {
       return nil, err
   }
   for _, item := range val {
       movie := &Movie{}
       err := json.Unmarshal([]byte(item), movie)
       if err != nil {
           return nil, err
       }
       movies = append(movies, movie)
   }

   return movies, nil
}

func (cache redisCache) UpdateMovie(movie *Movie) (*Movie, error) {
   c := cache.getClient()
   json, err := json.Marshal(&movie)
   if err != nil {
       return nil, err
   }
   c.HSet("movies", movie.Id, json)
   if err != nil {
       return nil, err
   }
   return movie, nil
}
func (cache redisCache) DeleteMovie(id string) error {
   c := cache.getClient()
   numDeleted, err := c.HDel("movies", id).Result()
   if numDeleted == 0 {
       return errors.New("movie to delete not found")
   }
   if err != nil {
       return err
   }
   return nil
}


*/
