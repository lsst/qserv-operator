package constants

const (
	BaseName = "lsst"
	AppLabel = "qserv"

	CzarName = "czar"

	CmsdName     = "cmsd"
	CmsdPort     = 2131
	CmsdPortName = "cmsd"

	MariadbName     = "mariadb"
	MariadbPort     = 3306
	MariadbPortName = MariadbName

	MysqlProxyName     = "proxy"
	MysqlProxyPort     = 4040
	MysqlProxyPortName = MysqlProxyName

	WmgrName     = "wmgr"
	WmgrPort     = 5012
	WmgrPortName = WmgrName

	WorkerName = "worker"

	XrootdAdminPathVolumeName = "xrootd-adminpath"
	XrootdName                = "xrootd"
	XrootdPort                = 1094
	XrootdPortName            = XrootdName
	XrootdRedirectorName      = "xrootd-redirector"

	GraceTime = 30

	CZAR         = "czar-0"
	REPL_CTL     = "repl-ctl"
	REPL_DB      = "repl-db-0"
	QSERV_DOMAIN = "qserv"
)

var MicroserviceConfigmaps = []string{MariadbName, XrootdName, MysqlProxyName, WmgrName}
var MicroserviceSecrets = []string{MariadbName, WmgrName}
var Databases = []string{"czar", "repl", "worker"}

var Command = []string{"/config-start/start.sh"}
