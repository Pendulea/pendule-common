package pcommon

import "time"

const DAY = 24 * time.Hour
const WEEK = 7 * DAY
const MONTH = 30 * DAY
const QUARTER = 90 * DAY

const MIN_TIME_FRAME = time.Second
const MAX_TIME_FRAME = QUARTER

var ARCHIVES_DIR = "./archives"
var DATABASES_DIR = "./databases"

var MAX_SIMULTANEOUS_PARSING = 3
var MAX_SIMULTANEOUS_INDEXING = 3

const FUTURES_KEY = "_futures"
const SPOT_KEY = "_spot"
