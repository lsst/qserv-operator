package constants

// All constants ending with 'Name' might have their value hard-coded in configmap/ directory
// Do not change their value.
const (
	BaseName = "lsst"
	AppLabel = "qserv"

	CmsdPort     = 2131
	CmsdPortName = string(CmsdName)

	MariadbPort     = 3306
	MariadbPortName = string(MariadbName)

	ProxyPort     = 4040
	ProxyPortName = string(ProxyName)

	QservName = "qserv"

	RedisName = "redis"

	ReplicationControllerPort     = 25080
	ReplicationControllerPortName = string(ReplCtlName)

	WmgrPort     = 5012
	WmgrPortName = string(WmgrName)

	XrootdAdminPathVolumeName = "xrootd-adminpath"
	XrootdPort                = 1094
	XrootdPortName            = string(XrootdName)

	GraceTime = 30
)

type ContainerName string

const (
	CmsdName    ContainerName = "cmsd"
	InitDbName  ContainerName = "initdb"
	MariadbName ContainerName = "mariadb"
	ProxyName   ContainerName = "proxy"
	XrootdName  ContainerName = "xrootd"
	ReplCtlName ContainerName = "repl-ctl"
	ReplDbName  ContainerName = "repl-db"
	WmgrName    ContainerName = "wmgr"
	ReplWrkName ContainerName = "repl-wrk"
)

type ComponentName string

const (
	CzarName             ComponentName = "czar"
	ReplName             ComponentName = "repl"
	WorkerName           ComponentName = "worker"
	XrootdRedirectorName ComponentName = "xrootd-redirector"
)

// ContainerConfigmaps contains names of all micro-services which require configmaps named:
// '<prefix>-<microservice-name>-etc' and '<prefix>-<microservice-name>-start'
var ContainerConfigmaps = []ContainerName{MariadbName, XrootdName, ProxyName, WmgrName, ReplCtlName, ReplDbName, ReplWrkName}

// MicroserviceSecrets contains names of all micro-services which require secrets
var MicroserviceSecrets = []ContainerName{MariadbName, WmgrName, ReplDbName}

// Databases contains names of all Qserv components which have a database
var Databases = []ComponentName{CzarName, ReplName, WorkerName}

// Command contains the default command used to launch a container
var Command = []string{"/config-start/start.sh"}
