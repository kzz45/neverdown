package redis

const (
	RedisDefaultPort = 6379
)

const (
	EnvRedisMasterHost = "ENV_REDIS_MASTER"
	EnvRedisMasterPort = "ENV_REDIS_MASTER_PORT"
	EnvRedisDir        = "ENV_REDIS_DIR"
	EnvRedisDbFileName = "ENV_REDIS_DBFILENAME"
	EnvRedisConf       = "ENV_REDIS_CONF"
	EnvRedisPort       = "ENV_REDIS_PORT"

	EnvRedisConfTemplate       = "redis-%s.conf"
	EnvRedisDbFileNameTemplate = "redis-%s.rdb"
)
