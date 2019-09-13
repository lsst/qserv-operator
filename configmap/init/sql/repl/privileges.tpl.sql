-- TODO: implement password security, DNS security is not reliable on k8s
set @repl_ctl_dn := '%';

SET @query = CONCAT('CREATE USER `qsreplica`@`', @repl_ctl_dn, '` IDENTIFIED BY `<MYSQL_REPLICA_PASSWORD>`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @query = CONCAT('GRANT ALL ON qservReplica.* TO `qsreplica`@`', @repl_ctl_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

set @repl_wrk_dn := '%.example-qserv-worker.default.svc.cluster.local';

SET @query = CONCAT('CREATE USER `qsreplica`@`', @repl_wrk_dn, '` IDENTIFIED BY `<MYSQL_REPLICA_PASSWORD>`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @query = CONCAT('GRANT ALL ON qservReplica.* TO `qsreplica`@`', @repl_wrk_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE USER `probe`@`localhost`;

FLUSH PRIVILEGES;
