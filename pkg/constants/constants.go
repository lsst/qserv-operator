package constants

const (
	BaseName = "lsst"
	AppLabel = "qserv"

	CzarName   = "czar"
	WorkerName = "worker"

	MariadbName     = "mariadb"
	MariadbPort     = 3306
	MariadbPortName = "mariadb"

	MysqlProxyName     = "proxy"
	MysqlProxyPort     = 4040
	MysqlProxyPortName = "proxy"

	XrootdAdminPathVolumeName = "xrootd-adminpath"
	XrootdName                = "xrootd"
	XrootdPort                = 1094
	XrootdPortName            = "xrootd"
	XrootdRedirectorName      = "xrootd-redirector"

	CmsdName     = "cmsd"
	CmsdPort     = 2131
	CmsdPortName = "cmsd"

	GraceTime = 30

	CZAR           = "czar-0"
	REPL_CTL       = "repl-ctl"
	REPL_DB        = "repl-db-0"
	QSERV_DOMAIN   = "qserv"
	XROOTD_MANAGER = "xrootd-0"
)

var WorkerServiceConfigmaps = []string{"mariadb", XrootdName, "wmgr"}
var WorkerServiceSecrets = []string{"mariadb", "wmgr"}
var Databases = []string{"czar", "repl", "worker"}

var Command = []string{"/config-start/start.sh"}
