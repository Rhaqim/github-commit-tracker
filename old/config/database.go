package config

import (
	"strconv"

	ut "github.com/Rhaqim/savannahtech/old/utils"
)

// Postgres
var (
	Database   = ut.Env("DATABASE", "savannah")
	PgHost     = ut.Env("DB_HOST", "localhost")
	PgPort     = ut.Env("DB_PORT", "5432")
	PgUser     = ut.Env("DB_USER", "savannah")
	PgPassword = ut.Env("DB_PASSWORD", "postgres")
	SSLMode    = ut.Env("SSL_MODE", "disable")
)

// Redis
var (
	RedisHost = ut.Env("REDIS_HOST", "localhost")
	RedisPort = ut.Env("REDIS_PORT", "6379")

	RedisAddress  = ut.Env("REDIS_ADDRESS", "localhost:6379")
	RedisPassword = ut.Env("REDIS_PASSWORD", "")
	RedisDB, _    = strconv.Atoi(ut.Env("REDIS_DB", "0"))
)
