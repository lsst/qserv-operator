package constants

const (
	BaseName = "lsst"
	AppLabel = "qserv-operator"

	XrootdConfigName                     = "xrootd"
	RedisStorageVolumeName               = "redis-data"
	RedisConfigurationVolumeName         = "redis-config"
	RedisShutdownConfigurationVolumeName = "redis-shutdown-config"
	RedisRoleName                        = "redis"
	HostnameTopologyKey                  = "kubernetes.io/hostname"
	XrootdRoleName                       = "xrootd"

	GraceTime = 30
)

var WorkerServiceConfigmaps = []string{"mariadb", "xrootd", "wmgr"}
