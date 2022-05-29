//go:build unit
// +build unit

package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/TrevorEdris/hypernova-bot/api/config"
	"github.com/TrevorEdris/hypernova-bot/api/domain"
)

type itemSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockDDB *mockDynamodbClient
	storage *ItemRepo
}

func TestItem(t *testing.T) {
	suite.Run(t, &itemSuite{})
}

func (s *itemSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockDDB = NewmockDynamodbClient(s.ctrl)
	s.storage = NewItemRepo(&config.Config{
		DynamoDB: config.DynamoDB{
			ItemTable: "testTableName",
		},
	}, s.mockDDB)
}

func (s *itemSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *itemSuite) TestItem_GetItem_Error() {
	testID := "1234"
	testErr := errors.New("GetItem error")
	s.mockDDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: testID}},
		TableName: aws.String(s.storage.table),
	}).Return(nil, testErr)

	result, err := s.storage.Get(context.Background(), testID)
	assert.ErrorIs(s.T(), err, testErr)
	assert.Empty(s.T(), result)
}

func (s *itemSuite) TestItem_GetItem_NotFound() {
	testID := "1234"
	s.mockDDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: testID}},
		TableName: aws.String(s.storage.table),
	}).Return(&dynamodb.GetItemOutput{
		Item: nil,
	}, nil)

	_, err := s.storage.Get(context.Background(), testID)
	assert.ErrorIs(s.T(), err, domain.ErrItemNotFound)
}

func (s *itemSuite) TestItem_GetItem() {
	expectedItem := domain.Item{
		ID:          "1234",
		Name:        gofakeit.Noun(),
		Description: gofakeit.Sentence(42), // Random dice roll decided the length
		Price:       gofakeit.Price(0.0, 13.37),
	}

	s.mockDDB.EXPECT().GetItem(gomock.Any(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: expectedItem.ID}},
		TableName: aws.String(s.storage.table),
	}).Return(&dynamodb.GetItemOutput{
		Item: map[string]types.AttributeValue{
			"id":          &types.AttributeValueMemberS{Value: expectedItem.ID},
			"title":       &types.AttributeValueMemberS{Value: expectedItem.Name},
			"description": &types.AttributeValueMemberS{Value: expectedItem.Description},
			"price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", expectedItem.Price)},
		},
	}, nil)

	result, err := s.storage.Get(context.Background(), expectedItem.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectedItem.ID, result.ID)
	assert.Equal(s.T(), expectedItem.Name, result.Name)
	assert.Equal(s.T(), expectedItem.Description, result.Description)
	assert.Equal(s.T(), expectedItem.Price, result.Price)
}
