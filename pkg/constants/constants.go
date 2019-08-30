package constants

// All constants ending with 'Name' might have their value hard-coded in configmap/ directory
// Do not change their value.
const (
	BaseName = "lsst"
	AppLabel = "qserv"

	CzarName = "czar"

	CmsdName     = "cmsd"
	CmsdPort     = 2131
	CmsdPortName = "cmsd"

	InitDbName = "initdb"

	MariadbName = "mariadb"
	MariadbPort = 3306

	ProxyName = "proxy"
	ProxyPort = 4040

	QservName = "qserv"

	ReplName    = "repl"
	ReplCtlName = "repl-ctl"
	ReplDbName  = "repl-db"
	ReplWrkName = "repl-wrk"

	WmgrName = "wmgr"
	WmgrPort = 5012

	WorkerName = "worker"

	XrootdAdminPathVolumeName = "xrootd-adminpath"
	XrootdName                = "xrootd"
	XrootdPort                = 1094
	XrootdPortName            = XrootdName
	XrootdRedirectorName      = "xrootd-redirector"

	GraceTime = 30
)

// MicroserviceConfigmaps contains names of all micro-services which require configmaps named:
// 'config-<microservice-name>-etc' and 'config-<microservice-name>-start'
var MicroserviceConfigmaps = []string{MariadbName, XrootdName, ProxyName, WmgrName, ReplCtlName, ReplDbName, ReplWrkName}

// MicroserviceSecrets contains names of all micro-services which require secrets
var MicroserviceSecrets = []string{MariadbName, WmgrName}

// Databases contains names of all Qserv components which have a database
var Databases = []string{CzarName, ReplName, WorkerName}

// Command contains the default command used to launch a container
var Command = []string{"/config-start/start.sh"}
