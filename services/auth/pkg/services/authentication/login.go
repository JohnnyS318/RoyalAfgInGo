package authentication

import (
	"errors"

	"github.com/spf13/viper"

	pAuth "github.com/JohnnyS318/RoyalAfgInGo/pkg/auth"
	"github.com/JohnnyS318/RoyalAfgInGo/pkg/models"
	"github.com/JohnnyS318/RoyalAfgInGo/services/auth/pkg/security"
	"github.com/JohnnyS318/RoyalAfgInGo/services/auth/pkg/serviceconfig"
)

//Login authenticates the user and generates a jwt token to enable the session
func (auth *Authentication) Login(username, password string) (*models.User, string, error) {

	user, err := auth.UserService.GetUserByUsernameOrEmail(username)
	if err != nil {
		return nil, "", err
	}

	if !security.ComparePassword(password, user.Hash, viper.GetString(serviceconfig.Pepper)) {
		return nil, "", errors.New("passwords did not match")
	}
	//TODO: Execute other login schemes (2FA)
	//user := models.NewUser("JohnnyS318", "jonas.max.schneider@gmail.com", "Jonas Schneider", time.Date(2003, 6, 6, 0, 0, 0, 0, time.UTC).Unix())
	//user.ID = primitive.NewObjectID()
	token, err := pAuth.GetJwt(user)
	return user, token, err
}
