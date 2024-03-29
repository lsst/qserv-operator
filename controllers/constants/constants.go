package constants

// All constants ending with 'Name' might have their value hard-coded in configmap/ directory
// Do not change their value.
const (
	BaseName = "lsst"

	CmsdName     = "cmsd"
	CmsdPort     = 2131
	CmsdPortName = string(CmsdName)

	CorePathVolumeName = "corepath"

	DatabaseURLFormat = "mysql://%s:%s@%s:%d"

	CzarDatabase = "qservMeta"

	ConfigmapPathEtc   = "/cm-etc"
	ConfigmapPathStart = "/cm-start"
	ConfigmapPathSQL   = "/cm-sql"

	DataVolumeClaimTemplateName = QservName + "-data"

	Localhost = "127.0.0.1"

	MariadbPort                = 3306
	MariadbPortName            = string(MariadbName)
	MariadbQservPassword       = ""
	MariadbQservUser           = "qsmaster"
	MariadbReplicationPassword = ""
	MariadbReplicationUser     = "qsreplica"
	MariadbRootUser            = "root"
	/* #nosec G101 no hard-coded password but environment variable */
	MariadbRootPassword = "${MYSQL_ROOT_PASSWORD}"
	MariadbSocket       = "/qserv/data/mysql/mysql.sock"

	ProbeTimeoutSeconds = 2
	ProxyPort           = 4040
	ProxyPortName       = string(ProxyName)

	QservName = "qserv"
	DotQserv  = "dot-qserv"

	HTTPPort     = 8080
	HTTPPortName = "http"

	XrootdAdminPath           = "/var/run/xrootd"
	XrootdAdminPathVolumeName = "xrootd-adminpath"
	XrootdName                = "xrootd"
	XrootdPort                = 1094
	XrootdPortName            = string(XrootdName)

	GraceTime = 30

	ReplicationDatabase             = "qservReplica"
	ReplicationWorkerDefaultThreads = 16
	ReplicationWorkerThreadFactor   = 2

	WorkerDatabase = "qservw_worker"

	// See https://confluence.lsstcorp.org/display/DM/5.+Automatic+retries+for+the+failed+contribution+requests
	WorkerIngestNumRetries = 2
	WorkerIngestMaxRetries = 10

	// TODO set this port in replication worker configmap script and use 8080 instead of default
	WorkerHTTPLoaderPort     = 25004
	WorkerHTTPLoaderPortName = "httpworker"
)

// QservGID qserv user gid
var QservGID int64 = 1000

// QservUID qserv user uid
var QservUID int64 = 1000

// ProbeInitialDelaySeconds default initial delay for all probes
var ProbeInitialDelaySeconds int32 = 30

// IngestDatabaseReplicas Number of replicas for ingest database
var IngestDatabaseReplicas int32 = 1

// ReplicationControllerReplicas Number of replicas for replication controller
var ReplicationControllerReplicas int32 = 1

// ReplicationDatabaseReplicas Number of replicas for replication database
var ReplicationDatabaseReplicas int32 = 1

// ReplicationRegistryReplicas Number of replicas for replication registry
var ReplicationRegistryReplicas int32 = 1

// ContainerName name all containers
type ContainerName string

const (
	// CmsdRedirectorName name for cmsd containers
	CmsdRedirectorName ContainerName = "cmsd-redirector"
	// CmsdServerName name for cmsd containers
	CmsdServerName ContainerName = "cmsd-server"
	// DebuggerName name for debugger container
	DebuggerName ContainerName = "debugger"
	// IngestDbName name for ingest database container
	IngestDbName ContainerName = "ingest-db"
	// InitDbName name for database initialization containers
	InitDbName ContainerName = "initdb"
	// MariadbName name for mariadb container
	MariadbName ContainerName = "mariadb"
	// ProxyName name for proxy container
	ProxyName ContainerName = "proxy"
	// ReplCtlName name for replication controller container
	ReplCtlName ContainerName = "repl-ctl"
	// ReplDbName name for replication database container
	ReplDbName ContainerName = "repl-db"
	// ReplRegistryName name for replication registry container
	ReplRegistryName ContainerName = "repl-registry"
	// XrootdRedirectorName Name name for xrootd manager container
	XrootdRedirectorName ContainerName = "xrootd-redirector"
	// XrootdServerName name for xrootd containers
	XrootdServerName ContainerName = "xrootd-server"
	// ReplWrkName name for replication worker container
	ReplWrkName ContainerName = "repl-wrk"
)

// PodClass name all classes of pod
// used to generate pod labels
type PodClass string

// GetDbContainerName return name of a database container for a given pod
func GetDbContainerName(pod PodClass) ContainerName {
	var dbName ContainerName
	if pod == ReplDb {
		dbName = ReplDbName
	} else if pod == IngestDb {
		dbName = IngestDbName
	} else {
		dbName = MariadbName
	}
	return dbName
}

const (
	// Czar name pods of class Czar
	Czar PodClass = "czar"
	// IngestDb name pods of class Ingest database
	IngestDb PodClass = PodClass(IngestDbName)
	// ReplCtl name pods of class Replication controller
	ReplCtl PodClass = PodClass(ReplCtlName)
	// ReplDb name pods of class Replication database
	ReplDb PodClass = PodClass(ReplDbName)
	// ReplRegistry name pods of class Replication registry
	ReplRegistry PodClass = PodClass(ReplRegistryName)
	// Worker name pods of class Replication worker
	Worker PodClass = "worker"
	// XrootdRedirector name pods of class Xrootd redirector
	XrootdRedirector PodClass = "xrootd-redirector"
)

// Command contain the default command used to launch a container
var Command = []string{ConfigmapPathStart + "/start.sh"}

// CommandDebug is a prerequisite for interactive debugging
var CommandDebug = []string{"sleep", "infinity"}

// WithEtcStartConfigmaps contain names of all containers which require configmaps both named:
// '<prefix>-<microservice-name>-etc' and '<prefix>-<microservice-name>-start'
var WithEtcStartConfigmaps = []ContainerName{CmsdServerName, IngestDbName, MariadbName, ProxyName, ReplCtlName, ReplDbName, ReplRegistryName, ReplWrkName, XrootdServerName}

// WithStartConfigmap contain names of all containers which require configmaps named:
// '<prefix>-<microservice-name>-start'
var WithStartConfigmap = []ContainerName{CmsdRedirectorName, XrootdRedirectorName}

// WithMariadbImage list container based on Mariadb image
var WithMariadbImage = []ContainerName{InitDbName, IngestDbName, MariadbName, ReplDbName}

// WithQservImage list container based on Qserv image
var WithQservImage = []ContainerName{CmsdRedirectorName, CmsdServerName, XrootdName, ProxyName, ReplCtlName, ReplWrkName, XrootdServerName, XrootdRedirectorName}

// Databases contains names of all Qserv pods which embed a database container
var Databases = []PodClass{Czar, ReplDb, Worker, IngestDb}

// MicroserviceSecrets contains names of all micro-services which require secrets
var MicroserviceSecrets = []ContainerName{MariadbName, ReplDbName}
