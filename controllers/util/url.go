package util

import (
	"fmt"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
)

// WorkerDatabaseURL return mariadb url for a worker database
func WorkerDatabaseURL(workerFqdn string) string {
	url := fmt.Sprintf(constants.DatabaseURLFormat+"/%s", constants.MariadbQservUser, constants.MariadbQservPassword, workerFqdn, constants.MariadbPort, constants.WorkerDatabase)
	return url
}

// databaseSocketURL return mariadb socket url for a database
func databaseSocketURL(user string, password string) string {
	url := fmt.Sprintf(constants.DatabaseURLFormat+"?socket=%s", user, password, constants.Localhost, constants.MariadbPort, constants.MariadbSocket)
	return url
}

// databaseURL return mariadb url for a database
func databaseURL(host string, user string, password string, database string) string {
	url := fmt.Sprintf(constants.DatabaseURLFormat+"/%s", user, password, host, constants.MariadbPort, database)
	return url
}

// GetCzarDatabaseRootURL returns the url of the replication database
func GetCzarDatabaseRootURL(cr *qservv1beta1.Qserv) string {
	czarDatabaseDn := GetName(cr, string(constants.Czar))
	return databaseURL(czarDatabaseDn, constants.MariadbRootUser, constants.MariadbRootPassword, constants.CzarDatabase)
}

// GetReplicationDatabaseRootURL returns the url of the replication database
func GetReplicationDatabaseRootURL(cr *qservv1beta1.Qserv) string {
	replicationDatabaseDn := GetName(cr, string(constants.ReplDbName))
	return databaseURL(replicationDatabaseDn, constants.MariadbRootUser, constants.MariadbRootPassword, constants.ReplicationDatabase)
}

// GetReplicationDatabaseURL returns the url of the replication database
func GetReplicationDatabaseURL(cr *qservv1beta1.Qserv) string {
	replicationDatabaseDn := GetName(cr, string(constants.ReplDbName))
	return databaseURL(replicationDatabaseDn, constants.MariadbReplicationUser, constants.MariadbReplicationPassword, constants.ReplicationDatabase)
}

// CzarDatabaseLocalRootURL URL to connect as root@127.0.0.1 to czar database
var CzarDatabaseLocalRootURL = databaseURL(constants.Localhost, constants.MariadbRootUser, constants.MariadbRootPassword, constants.CzarDatabase)

// ReplicationDatabaseLocalRootURL URL to connect as root@127.0.0.1 to replication database
var ReplicationDatabaseLocalRootURL = databaseURL(constants.Localhost, constants.MariadbRootUser, constants.MariadbRootPassword, constants.ReplicationDatabase)

// SocketQservUser URL to connect as qserv user to mariadb using socket file
var SocketQservUser = databaseSocketURL(constants.MariadbQservUser, constants.MariadbQservPassword)

// SocketRootUser URL to connect as root user to mariadb using socket file
var SocketRootUser = databaseSocketURL(constants.MariadbRootUser, constants.MariadbRootPassword)

// WorkerDatabaseLocalURL URL to connect as 'qserv user'@127.0.0.1 to worker database
var WorkerDatabaseLocalURL = databaseURL(constants.Localhost, constants.MariadbQservUser, constants.MariadbQservPassword, constants.WorkerDatabase)

// WorkerDatabaseLocalRootURL URL to connect as root@127.0.0.1 to worker database
var WorkerDatabaseLocalRootURL = databaseURL(constants.Localhost, constants.MariadbRootUser, constants.MariadbRootPassword, constants.WorkerDatabase)
