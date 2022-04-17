package server

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-redis/redis/v8"
	"github.com/toledoom/poc/pkg/battle"
	"github.com/toledoom/poc/pkg/framework"
	"github.com/toledoom/poc/pkg/leaderboard"
	"github.com/toledoom/poc/pkg/player"
)

func NewGrpc(redisAddr, dynamoAddr string) *GameServer {
	dynamodbClient := createDynamoDBLocalClient(dynamoAddr)
	redisClient := createRedisLocalClient(redisAddr)
	ranking := leaderboard.NewRedisRanking(redisClient, "tournament")
	evtPublisher := registerEventHandlers(ranking)
	return &GameServer{
		br:             battle.NewDynamoBattleRepository(dynamodbClient),
		pr:             player.NewDynamoPlayerRepository(dynamodbClient),
		calculator:     battle.NewEloScoreCalculator(20, 400),
		ranking:        ranking,
		eventPublisher: evtPublisher,
	}
}

func createDynamoDBLocalClient(dynamoAddr string) *dynamodb.Client {
	if dynamoAddr == "" {
		dynamoAddr = "http://127.0.0.1:8000"
	}
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: dynamoAddr}, nil
		})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func createRedisLocalClient(redisAddr string) *redis.Client {
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func registerEventHandlers(r leaderboard.Ranking) *framework.EventPublisher {
	deb := framework.NewEventPublisher()
	deb.Subscribe(leaderboard.NewPlayerScoreUpdateEventHandler(r), &player.PlayerScoreDispatched{})
	return deb
}
