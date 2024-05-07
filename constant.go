package pcommon

import "time"

const DAY = 24 * time.Hour
const WEEK = 7 * DAY
const MONTH = 30 * DAY
const QUARTER = 90 * DAY

const MIN_TIME_FRAME = time.Second
const MAX_TIME_FRAME = QUARTER

const FUTURES_KEY = "_futures"
const SPOT_KEY = "_spot"

type env struct {
	ARCHIVES_DIR              string
	DATABASES_DIR             string
	MAX_SIMULTANEOUS_PARSING  int
	PARSER_SERVER_PORT        int
	INDEXER_SERVER_PORT       int
	MAX_SIMULTANEOUS_INDEXING int
}

var Env = env{
	ARCHIVES_DIR:              "archives",
	DATABASES_DIR:             "databases",
	MAX_SIMULTANEOUS_PARSING:  3,
	MAX_SIMULTANEOUS_INDEXING: 3,
	PARSER_SERVER_PORT:        8889,
	INDEXER_SERVER_PORT:       8890,
}
