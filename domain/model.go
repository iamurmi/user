package domain

type User struct{ // User Class
    ID string `json:"id" bson:"_id"`
    FirstName string `json:"first_name" bson:"first_name"`
    Roles []string `json:"roles" bson:"roles"`
}
