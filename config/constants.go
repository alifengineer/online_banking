package config

import (
	"time"
)

const (
	DatabaseQueryTimeLayout = `'YYYY-MM-DD"T"HH24:MI:SS"."MS"Z"TZ'`
	// DatabaseTimeLayout
	DatabaseTimeLayout string = time.RFC3339
	// AccessTokenExpiresInTime ...
	AccessTokenExpiresInTime time.Duration = 1 * 24 * 60 * time.Minute
	// RefreshTokenExpiresInTime ...
	RefreshTokenExpiresInTime time.Duration = 30 * 24 * 60 * time.Minute
	// SigningKey ...
	SigningKey = "amsdklma345345345lsdmvjrbvuidj345345vyuvhsndsbdvnjshd"
)

// Environment
const (
	DEBUG_ENVIRONMENT      = "debug"
	TEST_ENVIRONMENT       = "test"
	PRODUCTION_ENVIRONMENT = "release"
)
const TimestampFormat = "2006-01-02 15:04:05.000000"

const (
	RECORD_NOT_FOUND    = "record not found"
	USER_ALREADY_EXISTS = "user already exists"

	SYSTEM_ERROR    = "system error"
	PASSWORD_WRONG  = "password is wrong, please check and try again"
	DUBLICATE_PHONE = "pq: duplicate key value violates unique constraint \"users_phone_key\""

	INVALID_PASSWORD_LENGTH = "password must not be less than 6 characters"
)
