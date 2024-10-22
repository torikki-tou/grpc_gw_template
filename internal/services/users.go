package services

import (
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"grpc_gw_template/protogen/golang/users"
)

type UsersService struct {
	users.UnimplementedUsersServer
}

func NewUsersService() *UsersService {
	return &UsersService{}
}

func (s *UsersService) Get(_ context.Context, r *users.GetUserRequest) (*users.User, error) {
	return &users.User{
		Id:     r.GetId(),
		Name:   "Test",
		Labels: []string{"Label1", "Label2"},
		Limit:  &wrapperspb.Int64Value{Value: 1321312},
	}, nil
}

func (s *UsersService) Search(_ context.Context, r *users.SearchUserRequest) (*users.SearchUserResponse, error) {
	usrs := []*users.User{
		{
			Id:     "ID",
			Name:   "Test",
			Labels: []string{"Label1", "Label2"},
			Limit:  &wrapperspb.Int64Value{Value: 1321312},
		},
		{
			Id:     "ID",
			Name:   "Test",
			Labels: []string{"Label1", "Label2"},
			Limit:  &wrapperspb.Int64Value{Value: 1321312},
		},
	}

	if r.GetLimit() != nil && r.Limit.Value <= uint64(len(usrs)) {
		usrs = usrs[:r.Limit.Value]
	}

	return &users.SearchUserResponse{Users: usrs}, nil
}

func (s *UsersService) Create(_ context.Context, r *users.CreateUserRequest) (*users.User, error) {
	return nil, nil
}

func (s *UsersService) Patch(_ context.Context, r *users.PatchUserRequest) (*users.Empty, error) {
	return nil, nil
}

func (s *UsersService) Delete(_ context.Context, r *users.DeleteUserRequest) (*users.Empty, error) {
	return nil, nil
}
