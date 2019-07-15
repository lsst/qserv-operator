package constants

const (
	XrootdRoleName             = "xrootd"
	XrootdConfigFileName       = "xrootd.conf"
	XrootdStartupFileName      = "start.sh"
	XrootdFinalStartupFileName = "xrd.sh"
	XrdssiConfigFileName       = "xrdssi.conf"
)

const XrootdConfigFileContent string = `# Unified configuration for xrootd/cmsd for both manager and server instances
# "if"-block separates manager-only and server-only configuration.


############################
# if: manager node
############################
if named manager

	# Use manager mode
	all.role manager

############################
# else: server nodes
############################
else

	# Use server mode
	all.role server

	# Get virtual network id, issued from mysql data
	set vnidfile = $VNID_FILE
	cms.vnid <${vnidfile}

	# Use XrdSsi plugin
	xrootd.fslib libXrdSsi.so
	ssi.svclib libxrdsvc.so
	oss.statlib -2 -arevents libXrdSsi.so

	# Force disable asyncronous access
	# because of XrdSsi
	xrootd.async off

	ssi.trace all debug

fi

########################################
# Shared directives (manager and server)
########################################

# Path to write logging and other information
all.adminpath /qserv/run/tmp/xrd

# Do not change. This specifies valid virtual paths that can be accessed.
# "nolock" directive prevents write-locking and is important for qserv
# qserv is hardcoded for these paths.
all.export / nolock

# Specify that no significant free space is required on servers
# Indeed current configuration doesn't expect to be dynamically
# written to, but export the space in R/W mode
cms.space 1k 2k

ssi.loglib libxrdlog.so

# Optional: Prevent dns resolution in logs.
# This may speed up request processing.
xrd.network nodnr

# This causes hostname resolution to occur at run-time not configuration time
# This is required by k8s
# Andy H. still have to modify the local IP-to-Name cache to account
# for dynamic DNS (it doesn't now). Unfortunately, it's a non-ABI compatible
# change so it will go into Release 5 branch not git master. The caching
# shouldn't really be a problem but if causes you grief simply turn it off by
# also specifying "xrd.network cache 0". Once Andy H. fixes the cache it will work
# correctly with a dynamic DNS with no side-effects (though it's unlikely any of
# them are observed as it is).
xrd.network dyndns
xrd.network cache 0

all.manager xrootd-mgr-0.xrootd-mgr:2131
all.manager xrootd-mgr-1.xrootd-mgr:2131

# - cmsd redirector runs on port 2131
# - cmsd server does not open server socket
#   but only client connection to cmsd redirector
# - xrootd default port is 1094
if exec cmsd
	xrd.port 2131
fi

# Uncomment the following line for detailed xrootd debugging
# xrootd.trace all debug
	`

const XrdssiConfigFileContent string = `# Qserv xrdssi plugin configuration file
# Default values for parameters are commented

[mysql]

# hostname =
# port =

# Username for mysql connections
username = qsmaster
password =

# MySQL socket file path for db connections
socket = /qserv/data/mysql/mysql.sock 

[memman]

# MemMan class to use for managing memory for tables
# can be "MemManReal" or "MemManNone"
# class = MemManReal

# Memory available for locking tables, in MB
# memory = 1000
memory = 7900

# Path to database tables
location = /qserv/data/mysql

[scheduler]

# Thread pool size
# thread_pool_size = 10
thread_pool_size = 20

# Required number of completed tasks for table in a chunk for the average time to be valid
# required_tasks_completed = 25
required_tasks_completed = 1

# Maximum group size for GroupScheduler
# group_size = 1
group_size = 10

# Scheduler priority - higher numbers mean higher priority.
# Running the fast scheduler at high priority tends to make it use significant
# resources on a small number of queries.
# priority_snail = -20
# priority_slow = 2
priority_slow = 4
# priority_med = 3
# priority_fast = 4
priority_fast = 2



# Maximum number of threads to reserve per scan scheduler
# reserve_snail = 2
# reserve_slow = 2
# reserve_med = 2
# reserve_fast = 2

# Maximum number of active chunks per scan scheduler
# maxActiveChunks_snail = 1
# maxActiveChunks_slow = 4
# maxActiveChunks_med = 4
# maxActiveChunks_fast = 4

# Maximum time for all tasks in a user query to complete.
# scanmaxminutes_fast = 60
# scanmaxminutes_med = 480
# scanmaxminutes_slow = 720
# scanmaxminutes_snail = 1440

# Maximum number of Tasks that can take too long before moving a query to the snail scan.
# maxtasksbootedperuserquery = 5
`

const XrootdStartupFileContent string = `#!/bin/sh

# Start cmsd or
# setup ulimit and start xrootd

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
set -x

usage() {
    cat << EOD

Usage: 'basename $0' [options] [cmd]

  Available options:
    -S <service> Service to start, default to xrootd

  Prepare cmsd and xrootd (ulimit setup) startup and
  launch associated startup script using qserv user.
EOD
}

service=xrootd

# get the options
while getopts S: c ; do
    case $c in
        S) service="$OPTARG" ;;
        \?) usage ; exit 2 ;;
    esac
done
shift $(($OPTIND - 1))

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

export XROOTD_DOMAIN="xrootd-mgr"
export XROOTD_DN="${XROOTD_DOMAIN}"

if hostname | egrep "^xrootd-mgr-[0-9]+"
then
    INSTANCE_NAME='manager'
else
    INSTANCE_NAME='worker'
fi
export INSTANCE_NAME

if [ "$service" = "xrootd" -a "$INSTANCE_NAME" = 'worker' ]; then

    # Increase limit for locked-in-memory size
    MLOCK_AMOUNT=$(grep MemTotal /proc/meminfo | awk '{printf("%.0f\n", $2 - 1000000)}')
    ulimit -l "$MLOCK_AMOUNT"

fi

su qserv -c "/config/xrd.sh -S $service"
`
const XrootdFinalStartupFileContent string = `#!/bin/sh

# Start cmsd and xrootd inside pod
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
# set -x

usage() {
    cat << EOD

Usage: "basename $0" [options] [cmd]

Available options:
-S <service> Service to start, default to xrootd

Start cmsd or xrootd.
EOD
}

service=xrootd

# get the options
while getopts S: c ; do
case $c in
	S) service="$OPTARG" ;;
	\?) usage ; exit 2 ;;
esac
done
shift $(($OPTIND - 1))

if [ $# -ne 0 ] ; then
usage
exit 2
fi

# Source pathes to eups packages
. /qserv/run/etc/sysconfig/qserv

CONFIG_DIR="/config"
XROOTD_CONFIG="$CONFIG_DIR/xrootd.cf"

# INSTANCE_NAME is required by xrdssi plugin to
# choose which type of queries to launch against metadata
if [ "$INSTANCE_NAME" = 'worker' ]; then

MYSQLD_SOCKET="/qserv/data/mysql/mysql.sock"
XRDSSI_CONFIG="$CONFIG_DIR/xrdssi.cf"

# Wait for local mysql to be configured and started
while true; do
	if mysql --socket "$MYSQLD_SOCKET" --user="$MYSQLD_USER_QSERV"  --skip-column-names \
		-e "SELECT CONCAT('Mariadb is up: ', version())"
	then
		break
	else
		echo "Wait for MySQL startup"
	fi
	sleep 2
done

# TODO move to /qserv/run/tmp when it is managed as a shared volume
export VNID_FILE="/qserv/data/mysql/cms_vnid.txt"
if [ ! -e "$VNID_FILE" ]
then
	WORKER=$(mysql --socket "$MYSQLD_SOCKET" --batch \
		--skip-column-names --user="$MYSQLD_USER_QSERV" -e "SELECT id FROM qservw_worker.Id;")
	if [ -z "$WORKER" ]; then
		>&2 echo "ERROR: unable to extract vnid from database"
		exit 2
	fi
	echo "$WORKER" > "$VNID_FILE"
fi

# Wait for at least one xrootd redirector readiness
until timeout 1 bash -c "cat < /dev/null > /dev/tcp/${XROOTD_DN}/2131"
do
	echo "Wait for xrootd manager (${XROOTD_DN})..."
	sleep 2
done

OPT_XRD_SSI="-l @libXrdSsiLog.so -+xrdssi $XRDSSI_CONFIG"
fi

# Start service
#
echo "Start $service"
"$service" -c "$XROOTD_CONFIG" -n "$INSTANCE_NAME" -I v4 $OPT_XRD_SSI
`
