package src

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/klever-io/inject-liqduidity-job/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func findByID(configID string) (config LiquidityPoolConfigs, err error) {
	objectID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.GetTimeoutDB())*time.Second)
	defer cancel()

	filter := primitive.M{"_id": objectID}

	rs := MongoDatabase.Collection("configs").FindOne(ctx, filter)
	if err != nil {
		err = fmt.Errorf("%v: %v", getConfigsErrorOnDB, err)
		return
	}

	if err = rs.Decode(&config); err != nil {
		if err == io.EOF {
			config.ID = primitive.NilObjectID
			return config, nil
		}
		err = fmt.Errorf("%v: %v", getConfigsErrorOnCursor, err)
		return
	}

	log.Default().Printf("current base: %f / quote: %f", config.BaseAmount, config.QuoteAmount)
	return
}

func LoadConfig(configID string) (config LiquidityPoolConfigs, err error) {
	if MongoClient == nil {
		err = fmt.Errorf("mongo client not initialized")
		return
	}

	return findByID(configID)
}

func UpdateConfig(config LiquidityPoolConfigs, randomNumber int) (err error) {
	if MongoClient == nil {
		err = fmt.Errorf("mongo client not initialized")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(config.ConfigID)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.GetTimeoutDB())*time.Second)
	defer cancel()

	// 10% of current value
	percentToMoveBaseAmount := 0.0
	percentToMoveQuoteAmount := 0.0

	// Generate a random number between 0 and 1

	if randomNumber == 0 {
		// add liquidity
		percentToMoveBaseAmount = config.BaseAmount + (config.BaseAmount * 0.1)
		percentToMoveQuoteAmount = config.QuoteAmount + (config.QuoteAmount * 0.1)
		log.Default().Println("adding liquidity...")
	} else {
		// remove liquidity
		percentToMoveBaseAmount = config.BaseAmount - (config.BaseAmount * 0.1)
		percentToMoveQuoteAmount = config.QuoteAmount - (config.QuoteAmount * 0.1)
		log.Default().Println("removing liquidity...")
	}

	filter := bson.M{
		"_id": objectID,
	}

	configToUpdate := primitive.D{
		{Key: "baseAmount", Value: percentToMoveBaseAmount},
		{Key: "quoteAmount", Value: percentToMoveQuoteAmount},
	}

	rs, err := MongoDatabase.Collection("configs").UpdateByID(ctx, filter, configToUpdate)
	if err != nil {
		return
	}

	if rs.UpsertedID == nil {
		err = fmt.Errorf("unexpected resultset after update: %v", rs)
		return
	}

	log.Default().Printf("updated base: %f / quote: %f", percentToMoveBaseAmount, percentToMoveQuoteAmount)
	return
}

func AddConfig() (err error) {
	path := os.Getenv("SEED_PATH")
	path = filepath.Join(path, "seed_configs.json")

	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}

	var config LiquidityPoolConfigs
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return
	}

	objectID, err := primitive.ObjectIDFromHex(config.ConfigID)
	if err != nil {
		return
	}
	config.ID = objectID
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(utils.GetTimeoutDB())*time.Second)
	defer cancel()

	rs, err := MongoDatabase.Collection("configs").InsertOne(ctx, config)
	if err != nil {
		return
	}

	log.Default().Printf("updated config with ID: %v", rs.InsertedID)
	return
}
