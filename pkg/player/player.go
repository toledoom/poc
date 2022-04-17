package player

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const tableName = "Players"

type Player struct {
	ID    string
	Name  string
	Score int64
}

func New(id, name string) Player {
	return Player{
		ID:   id,
		Name: name,
	}
}

type PlayerRepository interface {
	Add(p Player) error
	GetByID(id string) (Player, error)
	Update(p Player) error
}

type DynamoPlayerRepository struct {
	d *dynamodb.Client
}

func NewDynamoPlayerRepository(d *dynamodb.Client) *DynamoPlayerRepository {
	return &DynamoPlayerRepository{
		d: d,
	}
}

func (pr *DynamoPlayerRepository) GetByID(id string) (Player, error) {
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
	}

	playerDynamodb, err := pr.d.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return Player{}, err
	}

	p := Player{}
	err = attributevalue.UnmarshalMap(playerDynamodb.Item, &p)
	if err != nil {
		return Player{}, err
	}

	return p, nil
}

func (pr *DynamoPlayerRepository) Add(p Player) error {
	marshaledItem, err := attributevalue.MarshalMap(p)
	if err != nil {
		return err
	}

	addItemInput := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: aws.String(tableName),
	}

	_, err = pr.d.PutItem(context.TODO(), addItemInput)
	return err
}

func (br *DynamoPlayerRepository) Update(p Player) error {
	updateItemInput := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: p.ID},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET Score = :Score"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(p.Score, 10)},
		},
	}

	_, err := br.d.UpdateItem(context.TODO(), updateItemInput)

	return err
}
