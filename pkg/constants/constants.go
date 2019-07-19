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

	CZAR           = "czar-0"
	REPL_CTL       = "repl-ctl"
	REPL_DB        = "repl-db-0"
	QSERV_DOMAIN   = "qserv"
	XROOTD_MANAGER = "xrootd-0"
)

var WorkerServiceConfigmaps = []string{"mariadb", "xrootd", "wmgr"}
var WorkerServiceSecrets = []string{"mariadb", "wmgr"}
var Databases = []string{"czar", "repl", "worker"}
