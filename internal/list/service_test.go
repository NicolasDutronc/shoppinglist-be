package list_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	mocks "github.com/NicolasDutronc/shoppinglist-be/mocks/pkg/hub"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListServiceTestSuite struct {
	suite.Suite
	srv        *list.ServiceImpl
	mockedRepo *list.MockRepository
	mockedHub  *mocks.Hub
	list       *list.Shoppinglist
}

func (s *ListServiceTestSuite) SetupTest() {
	s.mockedRepo = &list.MockRepository{}
	s.mockedHub = &mocks.Hub{}
	s.srv = list.NewService(s.mockedRepo, s.mockedHub).(*list.ServiceImpl)

	s.list = &list.Shoppinglist{
		BaseModel: common.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: "list",
		Items: []*list.Item{
			{
				Name:     "item1",
				Quantity: "a lot",
			},
			{
				Name:     "item2",
				Quantity: "a little",
			},
		},
	}
}

func (s *ListServiceTestSuite) TestFindListByID() {
	ctx := context.Background()

	s.mockedRepo.On("FindListByID", ctx, s.list.ID.Hex()).Return(s.list, nil)
	list, err := s.srv.FindListByID(ctx, s.list.ID.Hex())
	assert.Equal(s.T(), s.list.ID.Hex(), list.ID.Hex())
	assert.NoError(s.T(), err)

	s.mockedRepo.On("FindListByID", ctx, "wrongID").Return(nil, errors.New("error retrieving the list"))
	list, err = s.srv.FindListByID(ctx, "wrongID")
	assert.Nil(s.T(), list)
	assert.Error(s.T(), err)
}

func (s *ListServiceTestSuite) TestFindAllLists() {
	ctx := context.Background()

	s.mockedRepo.On("FindAllLists", ctx).Return([]*list.Shoppinglist{
		{
			BaseModel: common.BaseModel{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "list",
			Items: []*list.Item{
				{
					Name:     "item1",
					Quantity: "a lot",
				},
				{
					Name:     "item2",
					Quantity: "a little",
				},
			},
		},
		{
			BaseModel: common.BaseModel{
				ID:        primitive.NewObjectID(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name: "list",
			Items: []*list.Item{
				{
					Name:     "item1",
					Quantity: "a lot",
				},
				{
					Name:     "item2",
					Quantity: "a little",
				},
			},
		},
	}, nil).Once()
	lists, err := s.srv.FindAllLists(ctx)
	assert.Len(s.T(), lists, 2)
	assert.NoError(s.T(), err)

	s.mockedRepo.On("FindAllLists", ctx).Return(nil, assert.AnError)
	lists, err = s.srv.FindAllLists(ctx)
	assert.Nil(s.T(), lists)
	assert.Error(s.T(), err)
}

func (s *ListServiceTestSuite) TestStoreList() {
	ctx := context.Background()

	// case 1 : the repo returns an error
	s.mockedRepo.On("StoreList", ctx, "nameThatAlreadyExists").Return(nil, assert.AnError).Once()
	list, err := s.srv.StoreList(ctx, "nameThatAlreadyExists")
	assert.Nil(s.T(), list)
	assert.Error(s.T(), err)

	// for other cases, the repo will return the list
	s.mockedRepo.On("StoreList", ctx, s.list.Name).Return(s.list, nil).Times(3)

	// case 2 : AddTopic returns an error
	s.mockedHub.On("AddTopic", ctx, hub.TopicFromString(s.list.ID.Hex())).Return(assert.AnError).Once()
	list, err = s.srv.StoreList(ctx, s.list.Name)
	assert.Nil(s.T(), list)
	assert.Error(s.T(), err)

	// for other cases, the hub will not return an error while creating a topic
	s.mockedHub.On("AddTopic", ctx, hub.TopicFromString(s.list.ID.Hex())).Return(nil).Times(2)

	// case 3 : Publish returns an error
	s.mockedHub.On("Publish", ctx, mock.Anything).Return(assert.AnError).Once()
	list, err = s.srv.StoreList(ctx, s.list.Name)
	assert.Nil(s.T(), list)
	assert.Error(s.T(), err)

	// for other cases, the hub will not return an error while publishing a message
	s.mockedHub.On("Publish", ctx, mock.Anything).Return(nil).Once()

	// case 4 : everything goes well
	list, err = s.srv.StoreList(ctx, s.list.Name)
	assert.NotNil(s.T(), list)
	assert.NoError(s.T(), err)
}

func (s *ListServiceTestSuite) TestDeleteList() {

}

func (s *ListServiceTestSuite) TestAddItem() {

}

func (s *ListServiceTestSuite) TestUpdateItem() {

}

func (s *ListServiceTestSuite) TestToggleItem() {

}

func (s *ListServiceTestSuite) TestRemoveItem() {

}

func (s *ListServiceTestSuite) TestRemoveAllItems() {

}

func TestListServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ListServiceTestSuite))
}
