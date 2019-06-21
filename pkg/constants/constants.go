package constants

const (
	BaseName                             = "bc"
	AppLabel                             = "redis-operator"
	RedisName                            = "redis"
	RedisStorageVolumeName               = "redis-data"
	RedisConfigurationVolumeName         = "redis-config"
	RedisShutdownConfigurationVolumeName = "redis-shutdown-config"
	RedisRoleName                        = "redis"
	SentinelRoleName                     = "sentinel"
	RedisConfigFileName                  = "redis.conf"
	SentinelConfigFileName               = "sentinel.conf"
	HostnameTopologyKey                  = "kubernetes.io/hostname"

	GraceTime = 30
)
