package battle

import (
	"context"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const tableName = "Battles"

type Result int64

const (
	Lose Result = iota
	Win
)

type Battle struct {
	ID           string
	Player1ID    string
	Player2ID    string
	Player1Score int64
	Player2Score int64
}

func New(player1ID, player2ID string) Battle {
	id := uuid.New()
	return Battle{
		ID:        id.String(),
		Player1ID: player1ID,
		Player2ID: player2ID,
	}
}

func (b Battle) Id() string {
	return b.ID
}

type BattleRepository interface {
	GetByID(id string) (Battle, error)
	Add(b Battle) error
	Update(b Battle) error
}

type DynamoBattleRepository struct {
	d *dynamodb.Client
}

func NewDynamoBattleRepository(d *dynamodb.Client) *DynamoBattleRepository {
	return &DynamoBattleRepository{
		d: d,
	}
}

func (br *DynamoBattleRepository) GetByID(id string) (Battle, error) {
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(tableName),
	}

	battleDynamodb, err := br.d.GetItem(context.TODO(), getItemInput)
	if err != nil {
		return Battle{}, err
	}

	b := Battle{}
	err = attributevalue.UnmarshalMap(battleDynamodb.Item, &b)
	if err != nil {
		return Battle{}, err
	}

	return b, nil
}

func (br *DynamoBattleRepository) Add(b Battle) error {
	marshaledItem, err := attributevalue.MarshalMap(b)
	if err != nil {
		return err
	}

	addItemInput := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: aws.String(tableName),
	}

	_, err = br.d.PutItem(context.TODO(), addItemInput)
	return err
}

func (br *DynamoBattleRepository) Update(b Battle) error {

	updateItemInput := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: b.ID},
		},
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String("SET Player1Score = :Player1Score, Player2Score = :Player2Score"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":Player1Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(b.Player1Score, 10)},
			":Player2Score": &types.AttributeValueMemberN{Value: strconv.FormatInt(b.Player2Score, 10)},
		},
	}

	_, err := br.d.UpdateItem(context.TODO(), updateItemInput)

	return err
}

type ScoreCalculator interface {
	Calculate(score1, score2 int64, battleResult Result) int64
}

type EloScoreCalculator struct {
	K, S int64 // K normally varies depending on several factors: battle result (win, lose), player score... Here I am simplifying
}

func NewEloScoreCalculator(k, s int64) *EloScoreCalculator {
	return &EloScoreCalculator{
		K: k,
		S: s,
	}
}

func (e *EloScoreCalculator) Calculate(score1, score2 int64, battleResult Result) int64 {
	diffRatio := float64((score1 - score2) / e.S)
	x := math.Pow(10, diffRatio)
	expectedScore := 1 / (1 + x)
	br := float64(int64(battleResult))
	delta := float64(e.K) * (br - expectedScore)
	finalScore := float64(score1) + delta

	return int64(finalScore)
}
