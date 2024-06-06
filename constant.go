package pcommon

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const DAY = 24 * time.Hour
const WEEK = 7 * DAY
const MONTH = 30 * DAY
const QUARTER = 90 * DAY
const TIME_UNIT_DURATION = time.Millisecond

// const MIN_TIME_FRAME = TIME_UNIT_DURATION * 1000
const MAX_TIME_FRAME = QUARTER

type env struct {
	ARCHIVES_DIR                      string
	DATABASES_DIR                     string
	MAX_SIMULTANEOUS_PARSING          int
	PARSER_SERVER_PORT                string
	MIN_TIME_FRAME                    time.Duration
	MAX_DAYS_BACKWARD_FOR_CONSISTENCY int
}

var Env = env{
	ARCHIVES_DIR:                      "archives",
	DATABASES_DIR:                     "databases",
	MAX_SIMULTANEOUS_PARSING:          3,
	PARSER_SERVER_PORT:                "8889",
	MIN_TIME_FRAME:                    1000 * TIME_UNIT_DURATION,
	MAX_DAYS_BACKWARD_FOR_CONSISTENCY: 2,
}

func (e env) Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	minDaysBackwardForConsistency := os.Getenv("MAX_DAYS_BACKWARD_FOR_CONSISTENCY")
	if minDaysBackwardForConsistency != "" {
		min, err := strconv.Atoi(minDaysBackwardForConsistency)
		if err != nil {
			log.Fatal("Error parsing MAX_DAYS_BACKWARD_FOR_CONSISTENCY")
		} else {
			Env.MAX_DAYS_BACKWARD_FOR_CONSISTENCY = min
		}
	}

	// Archives directory
	minTimeFrame := os.Getenv("MIN_TIME_FRAME")
	if minTimeFrame != "" {
		min, err := strconv.Atoi(minTimeFrame)
		if err != nil {
			log.Fatal("Error parsing MIN_TIME_FRAME")
		} else {
			Env.MIN_TIME_FRAME = time.Duration(min) * time.Millisecond
		}
	}

	// Archives directory
	archiveDir := os.Getenv("ARCHIVES_DIR")
	if archiveDir != "" {
		if stat, err := os.Stat(archiveDir); os.IsNotExist(err) || !stat.IsDir() {
			log.Fatalf("archives directory not found or is not a directory")
		} else {
			Env.ARCHIVES_DIR = archiveDir
		}
	}

	// Databases directory
	dbDir := os.Getenv("DATABASES_DIR")
	if dbDir != "" {
		if stat, err := os.Stat(dbDir); os.IsNotExist(err) || !stat.IsDir() {
			log.Fatalf("databases directory not found or is not a directory")
		} else {
			Env.DATABASES_DIR = dbDir
		}
	}

	// Max simultaneous parsing workerss
	maxSimultaneousParsing := os.Getenv("MAX_SIMULTANEOUS_PARSING")
	if maxSimultaneousParsing != "" {
		max, err := strconv.Atoi(maxSimultaneousParsing)
		if err != nil {
			log.Fatal("Error parsing MAX_SIMULTANEOUS_PARSING")
		} else {
			Env.MAX_SIMULTANEOUS_PARSING = max
		}
	}

	// Parser server port
	parserServerPort := os.Getenv("PARSER_SERVER_PORT")
	if parserServerPort != "" {
		serverPortInt, err := strconv.Atoi(parserServerPort)
		if err != nil {
			log.Fatal("Error parsing PARSER_SERVER_PORT")
		} else {
			if serverPortInt < 0 || serverPortInt > 65535 {
				log.Fatal("Invalid port PARSER_SERVER_PORT")
			}
		}
		Env.PARSER_SERVER_PORT = strconv.Itoa(serverPortInt)
	}
}
