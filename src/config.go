package src

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	cursor, err := MongoDatabase.Collection("configs").Find(ctx, filter)
	if err != nil {
		err = fmt.Errorf("%v: %v", getConfigsErrorOnDB, err)
		return
	}

	if err = cursor.Decode(&config); err != nil {
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

	objectID, err := primitive.ObjectIDFromHex(config.ID)
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
	jsonFile, err := os.Open("../configs.json")
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

	return
}
