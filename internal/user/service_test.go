package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/NicolasDutronc/autokey"
	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"github.com/NicolasDutronc/shoppinglist-be/internal/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceTestSuite struct {
	suite.Suite
	srv      *user.ServiceImpl
	repo     *user.MockRepository
	consumer *autokey.MockConsumer
	u        *user.User
}

func (s *UserServiceTestSuite) SetupTest() {
	s.u = &user.User{
		BaseModel: common.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     "dummy",
		Password: "1234",
	}
	s.repo = &user.MockRepository{}
	s.consumer = &autokey.MockConsumer{}
	s.srv = user.NewService(s.repo, s.consumer).(*user.ServiceImpl)
}

func (s *UserServiceTestSuite) TestFindByID() {
	ctx := context.Background()
	s.repo.On("FindByID", ctx, s.u.ID.Hex()).Return(s.u, nil)

	user, err := s.srv.FindByID(ctx, s.u.ID.Hex())

	assert.Equal(s.T(), s.u.ID, user.ID)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestFindByName() {
	ctx := context.Background()
	s.repo.On("FindByName", ctx, s.u.Name).Return(s.u, nil)

	user, err := s.srv.FindByName(ctx, s.u.Name)

	assert.Equal(s.T(), s.u.ID, user.ID)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestStore() {
	ctx := context.Background()

	s.repo.On("Store", ctx, s.u.Name, mock.AnythingOfType("string")).Return(
		func(_ context.Context, _ string, arg string) *user.User {
			return &user.User{
				BaseModel: s.u.BaseModel,
				Name:      s.u.Name,
				Password:  arg,
			}
		},
		nil,
	)

	user, err := s.srv.Store(ctx, s.u.Name, s.u.Password)
	assert.NoError(s.T(), err)

	assert.NoError(s.T(), bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.u.Password)))
}

func (s *UserServiceTestSuite) TestUpdateName() {
	ctx := context.Background()

	newName := "newName"
	s.repo.On("FindByID", ctx, s.u.ID.Hex()).Return(s.u, nil)
	s.repo.On("UpdateName", ctx, s.u.ID.Hex(), newName).Return(int64(1), nil)

	n, err := s.srv.UpdateName(ctx, s.u.ID.Hex(), newName)

	assert.Equal(s.T(), int64(1), n)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestUpdatePassword() {
	ctx := context.Background()

	newPassword := "newPassword"
	s.repo.On("FindByID", ctx, mock.AnythingOfType("string")).Return(
		func(_ context.Context, id string) *user.User {
			if id == "notFoundID" {
				return nil
			}

			encrypted, err := bcrypt.GenerateFromPassword([]byte(s.u.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil
			}

			return &user.User{
				BaseModel: s.u.BaseModel,
				Name:      s.u.Name,
				Password:  string(encrypted),
			}

		},
		func(_ context.Context, id string) error {
			if id == "notFoundID" {
				return assert.AnError
			}

			return nil
		},
	).Times(3)
	s.repo.On("UpdatePassword", ctx, s.u.ID.Hex(), s.u.Password, mock.AnythingOfType("string")).Return(
		func(_ context.Context, _ string, _ string, newPasswordHashed string) int64 {
			if err := bcrypt.CompareHashAndPassword([]byte(newPasswordHashed), []byte(newPassword)); err != nil {
				return -1
			}

			return 1
		},
		func(_ context.Context, _ string, _ string, newPasswordHashed string) error {
			if err := bcrypt.CompareHashAndPassword([]byte(newPasswordHashed), []byte(newPassword)); err != nil {
				return err
			}

			return nil
		},
	).Times(3)

	n, err := s.srv.UpdatePassword(ctx, s.u.ID.Hex(), s.u.Password, newPassword)
	assert.Equal(s.T(), int64(1), n)
	assert.NoError(s.T(), err)

	n, err = s.srv.UpdatePassword(ctx, "notFoundID", s.u.Password, newPassword)
	assert.Equal(s.T(), int64(-1), n)
	assert.Error(s.T(), err)

	n, err = s.srv.UpdatePassword(ctx, s.u.ID.Hex(), "wrongPassword", newPassword)
	assert.Equal(s.T(), int64(-1), n)
	assert.Error(s.T(), err)

}

func (s *UserServiceTestSuite) TestDelete() {
	ctx := context.Background()

	s.repo.On("Delete", ctx, s.u.ID.Hex()).Return(int64(1), nil)

	n, err := s.srv.Delete(ctx, s.u.ID.Hex())

	assert.Equal(s.T(), int64(1), n)
	assert.NoError(s.T(), err)
}

func (s *UserServiceTestSuite) TestLogin() {
	ctx := context.Background()

	s.repo.On("FindByName", ctx, s.u.Name).Return(
		func(_ context.Context, name string) *user.User {
			if name != s.u.Name {
				return nil
			}

			encrypted, err := bcrypt.GenerateFromPassword([]byte(s.u.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil
			}

			return &user.User{
				BaseModel: s.u.BaseModel,
				Name:      name,
				Password:  string(encrypted),
			}
		},
		func(_ context.Context, name string) error {
			if name != s.u.Name {
				return assert.AnError
			}

			return nil
		},
	)

	key := "superSecretKey"
	s.consumer.On("Get").Return(key, nil)

	user, token, err := s.srv.Login(ctx, s.u.Name, s.u.Password)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Greater(s.T(), len(token), 0)

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})

	assert.NoError(s.T(), parsedToken.Claims.Valid())
	assert.IsType(s.T(), &jwt.StandardClaims{}, parsedToken.Claims)

}

func (s *UserServiceTestSuite) TestAuthenticate() {
	ctx := context.Background()

	// create token
	key := "superSecretKey"
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   s.u.ID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		s.FailNow(err.Error())
	}

	s.consumer.On("Get").Return(key, nil)

	s.repo.On("FindByID", ctx, s.u.ID.Hex()).Return(s.u, nil)

	user, err := s.srv.Authenticate(ctx, ss)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.u.Name, user.Name)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
