package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load("default.env", ".env"); err != nil {
		err = fmt.Errorf("envs file not found: %s", err.Error())
		return err
	}
	return nil
}

func GetTimeoutDB() int {
	timeout := 5

	strTimeout, ok := os.LookupEnv("DB_TIMEOUT")
	if ok {
		parsedTimeout, err := strconv.Atoi(strTimeout)
		if err == nil {
			timeout = parsedTimeout
		}
	}
	return timeout
}
