CREATE USER 'qsreplica'@'%' IDENTIFIED BY '<MYSQL_REPLICA_PASSWORD>';
GRANT ALL ON qservReplica.* TO 'qsreplica'@'%';

CREATE USER 'probe'@'localhost';

FLUSH PRIVILEGES;
