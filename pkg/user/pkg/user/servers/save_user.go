package servers

import (
	"context"

	"github.com/JohnnyS318/RoyalAfgInGo/pkg/protos"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/user/pkg/user/database"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *UserServer) SaveUser(ctx context.Context, m *protos.User) (*protos.User, error) {

	s.l.Infof("Called SaveUser Grpc %v", m)

	user := fromMessageUserExact(m)

	if err := user.Validate(); err != nil {
		s.l.Errorw("Validation error", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Validate unique Username and Email

	err := s.db.CreateUser(user)

	if err != nil {

		if database.IsDup(err) {
			return nil, status.Error(codes.AlreadyExists, "Username or Email already used pleace try again using a different Username or Email")
		}

		s.l.Errorw("Error during parsing", "error", err)
		return nil, status.Error(codes.Internal, "User could not be saved to the database")
	}

	return toMessageUser(user), nil
}