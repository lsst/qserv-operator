set @repl_ctl_dn := 'repl-ctl-%.qserv.default.svc.cluster.local';

SET @query = CONCAT('CREATE USER `qsreplica`@`', @repl_ctl_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @query = CONCAT('GRANT ALL ON qservReplica.* TO `qsreplica`@`', @repl_ctl_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

set @repl_wrk_dn := 'qserv-%.qserv.default.svc.cluster.local';

SET @query = CONCAT('CREATE USER `qsreplica`@`', @repl_wrk_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @query = CONCAT('GRANT ALL ON qservReplica.* TO `qsreplica`@`', @repl_wrk_dn, '`');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE USER `probe`@`localhost`;

FLUSH PRIVILEGES;
