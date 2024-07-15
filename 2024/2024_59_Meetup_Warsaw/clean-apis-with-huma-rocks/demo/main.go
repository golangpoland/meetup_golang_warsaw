package main

import (
	"context"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/google/uuid"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

var users []*User

func listUsersHandler(_ context.Context, in *ListUsersInput) (*ListUsersOutput, error) {
	out := make([]*User, 0, in.Limit)
	for i := in.Cursor; i < len(users) && i < in.Cursor+in.Limit; i++ {
		out = append(out, users[i])
	}

	return &ListUsersOutput{
		Body: out,
	}, nil
}

func createUserHandler(_ context.Context, in *CreateUserInput) (*CreateUserOutput, error) {
	exists := slices.ContainsFunc(users, func(u *User) bool { return u.Email == in.Body.Email })
	if exists {
		return nil, huma.Error409Conflict("User already exists")
	}

	in.Body.ID = uuid.New().String()

	users = append(users, in.Body)

	return &CreateUserOutput{
		Body: in.Body,
	}, nil
}

func getUserHandler(_ context.Context, in *GetUserInput) (*GetUserOutput, error) {
	idx := slices.IndexFunc(users, func(u *User) bool { return u.ID == in.ID })
	if idx == -1 {
		return nil, huma.Error404NotFound("User not found")
	}

	return &GetUserOutput{
		Body: users[idx],
	}, nil
}

func deleteUserHandler(_ context.Context, in *DeleteUserInput) (*struct{}, error) {
	idx := slices.IndexFunc(users, func(u *User) bool { return u.ID == in.ID })
	if idx == -1 {
		return nil, huma.Error404NotFound("User not found")
	}

	users[idx] = users[len(users)-1]
	users = users[:len(users)-1]

	return nil, nil
}

func main() {
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("My API", "1.0.0"))

	huma.Post(api, "/users", createUserHandler)
	huma.Get(api, "/users", listUsersHandler)
	huma.Get(api, "/users/{id}", getUserHandler)
	huma.Delete(api, "/users/{id}", deleteUserHandler)

	http.ListenAndServe("127.0.0.1:8888", router)
}
