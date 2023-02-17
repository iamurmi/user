package userservice

import (
	"errors"

	"github.com/iamurmi/user/contracts"
	"github.com/iamurmi/user/domain"
	userPB "github.com/iamurmi/user/domain/protobuf"

	"github.com/gin-gonic/gin"
)

// encode decode and call service method

type ResponseWrap struct { // class
	Sucess bool        `json:"sucess"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error,omitempty"`
}

func SuccessResponse(data interface{}, err error) (resp ResponseWrap) {
	resp = ResponseWrap{
		Sucess: true,
		Data:   data,
		Error:  "",
	}
	return resp
}
func FailedResponse(data interface{}, err error) (resp ResponseWrap) {
	errStr := err.Error()
	resp = ResponseWrap{
		Sucess: false,
		Data:   nil,
		Error:  errStr,
	}
	return resp
}

// End of RESP Wrapper

type Routes struct {
	userClient userPB.UserClient
}

func NewUserRoutes(userClient userPB.UserClient) *Routes {
	return &Routes{
		userClient: userClient,
	}
}

// Handler Apis call
func (r *Routes) NewRoutes(e *gin.Engine, svcObj contracts.UserSeriveContract) { // Initializer of Transport layer
	e.POST("/tusi/api/smallmic/adduser", r.AddUserHandler(svcObj))
	e.GET("/tusi/api/smallmic/getuser/:id", r.GetUserHandler(svcObj))
	e.GET("/tusi/api/smallmic/getusers", r.GetUsersHandler(svcObj))
}

// Handler Functions
// When we want to call this function we doesnt have to pass context, but we need the context for decoding json to struct, so here i create an anonymous function who take context as arguments which doest not required to pass

func (r *Routes) AddUserHandler(svc contracts.UserSeriveContract) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user domain.User
		err := c.ShouldBindJSON(&user) // postman request stored inside of c (gin context).
		if err != nil {
			err = errors.New("Bad Request")
			failedStructWrapper := FailedResponse(nil, err)
			c.JSON(400, failedStructWrapper)
			return
		}
		// Conversion of classes
		userPb := userPB.AddUserRequest{
			User: &userPB.UserData{
				Id:        user.ID,
				FirstName: user.FirstName,
				Roles:     user.Roles,
			},
		}
		//Svc call
		// id, err := svc.AddUser(c, &userPb)
		id, err := r.userClient.AddUser(c, &userPb) // GRPC Client
		if err != nil {
			failedStructWrapper := FailedResponse(nil, err)
			c.JSON(400, failedStructWrapper)
			return
		}
		sucessStructWrapper := SuccessResponse(id, nil)
		c.JSON(200, sucessStructWrapper)
	}
}

func (r *Routes) GetUserHandler(svc contracts.UserSeriveContract) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		//Svc call
		req := userPB.GetUserRequest{
			Id: id,
		}
		// resp, err := svc.GetUser(c, &req)
		resp, err := r.userClient.GetUser(c, &req) // GRPC Client
		if err != nil {
			failedStructWrapper := FailedResponse(nil, err)
			c.JSON(400, failedStructWrapper)
			return
		}
		u := domain.User{
			ID:        resp.User.Id,
			FirstName: resp.User.FirstName,
			Roles:     resp.User.Roles,
		}
		sucessStructWrapper := SuccessResponse(u, nil)
		c.JSON(200, sucessStructWrapper)
	}
}

func (r *Routes) GetUsersHandler(svc contracts.UserSeriveContract) func(c *gin.Context) {
	return func(c *gin.Context) {
		//Svc call
		// resp, err := svc.GetUsers(c)

		var emptyObj *userPB.ListUsersRequest
		resp, err := r.userClient.GetUsers(c, emptyObj)
		if err != nil {
			failedStructWrapper := FailedResponse(nil, err)
			c.JSON(400, failedStructWrapper)
			return
		}
		var users []domain.User
		for i := range resp.Users {
			u := domain.User{
				ID:        resp.Users[i].Id,
				FirstName: resp.Users[i].FirstName,
				Roles:     resp.Users[i].Roles,
			}
			users = append(users, u)
		}
		sucessStructWrapper := SuccessResponse(users, nil)
		c.JSON(200, sucessStructWrapper)
	}
}
