package main



type User struct {
	ID       string `json:"id" readOnly:"true" doc:"User ID"`
	Username string `json:"username" minLength:"1" maxLength:"255" example:"bob" doc:"User name"`
	Email    string `json:"email" minLength:"1" maxLength:"255" example:"bob@gmail.com" doc:"User email"`
}

type AuthMixin struct {
	Auth string `header:"Authorization`
}

type PaginationMixin struct {
	Cursor int `query:"cursor" minimum:"0" default:"0"`
	Limit  int `query:"limit" minimum:"1" maximum:"100" default:"10"`
}

type CreateUserInput struct {
	AuthMixin
	Body *User
}

type CreateUserOutput struct {
	Body *User
}

type ListUsersInput struct {
	AuthMixin
	PaginationMixin
}

type ListUsersOutput struct {
	Body []*User
}

type GetUserInput struct {
	AuthMixin
	// ID string `path:"id" minLength:"1" maxLength:"36" example:"ed7d12e1-9a0e-4706-9848-233566f3d6b3" doc:"User ID"`
}

type GetUserOutput struct {
	Body *User
}

type DeleteUserInput struct {
	AuthMixin
	ID string `path:"id" minLength:"1" maxLength:"36" example:"ed7d12e1-9a0e-4706-9848-233566f3d6b3" doc:"User ID"`
}

type DeleteUserOutput struct {
	Body *User
}


