package constants

const (
	BaseName                             = "lsst"
	AppLabel                             = "qserv-operator"
	XrootdConfigName                     = "xrootd-etc"
	RedisStorageVolumeName               = "redis-data"
	RedisConfigurationVolumeName         = "redis-config"
	RedisShutdownConfigurationVolumeName = "redis-shutdown-config"
	RedisRoleName                        = "redis"
	HostnameTopologyKey                  = "kubernetes.io/hostname"

	GraceTime = 30
)
