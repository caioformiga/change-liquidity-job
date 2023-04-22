package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/klever-io/inject-liqduidity-job/src"
	"github.com/klever-io/inject-liqduidity-job/utils"
)

func main() {

	// Run the loop until the timeout has expired
	go func() {
		processTime := time.Tick(5 * time.Minute)
		for {
			<-processTime
			log.Default().Println("running...")
		}
	}()

	if err := utils.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	URI := os.Getenv("WORKSPACE_MONGO_URI")
	database := os.Getenv("WORKSPACE_MONGO_DB")
	dbTimeout := utils.GetTimeoutDB()

	var err error
	src.MongoClient, err = src.InitDB(URI, database, dbTimeout)
	if err != nil {
		log.Fatal(err)
	}

	src.MongoDatabase = src.MongoDatabase.Client().Database(database)

	configID := os.Getenv("CONFIG_ID")
	if len(configID) == 0 {
		return
	}

	go func() {
		iteration := 1

		// Define the function to execute
		doSomething := func() {
			config, err := src.LoadConfig(configID)
			if err != nil {
				return
			}

			// Seed the random number generator
			rand.Seed(time.Now().UnixNano())

			randomNumber := rand.Intn(100) % 2
			src.UpdateConfig(config, randomNumber)
		}

		for {
			// Seed the random number generator
			rand.Seed(time.Now().UnixNano())

			// set random interval of some seconds between 5 - 300 = 5 minutes
			randonInterval := rand.Intn(300) + 5

			// Set the interval to run function
			duration := time.Duration(randonInterval) * time.Second
			timer := time.Tick(duration)

			logMsg := fmt.Sprintf("wait %s to execute %d iteration", duration.String(), iteration)
			log.Default().Println(logMsg)

			<-timer
			doSomething()
			iteration++
		}
	}()

	// Parse the timeout value from the command line arguments, with a default value of 60 minute
	timeout := flag.Duration("timeout", 60*time.Minute, "timeout duration")
	flag.Parse()

	logMsg := fmt.Sprintf("job will run for %s\n", timeout.String())
	log.Default().Println(logMsg)

	// Set a timer to stop the loop after the timeout has expired
	stopTimer := time.AfterFunc(*timeout, func() {
		log.Default().Println("stopping bot...")

		os.Exit(0)
	})
	<-stopTimer.C
}
