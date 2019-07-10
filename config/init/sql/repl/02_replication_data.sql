SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 ;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 ;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL' ;

USE qservReplica;

-----------------------------------------------------------
-- Preload configuration parameters for testing purposes --
-----------------------------------------------------------

-- Common parameters of all types of servers

INSERT INTO `config` VALUES ('common', 'request_buf_size_bytes',     '131072');
INSERT INTO `config` VALUES ('common', 'request_retry_interval_sec', '5');

-- Controller-specific parameters

INSERT INTO `config` VALUES ('controller', 'num_threads',            '16');
INSERT INTO `config` VALUES ('controller', 'http_server_port',       '8080');
INSERT INTO `config` VALUES ('controller', 'http_server_threads',    '16');
INSERT INTO `config` VALUES ('controller', 'request_timeout_sec', '57600');   -- 16 hours
INSERT INTO `config` VALUES ('controller', 'job_timeout_sec',     '57600');   -- 16 hours
INSERT INTO `config` VALUES ('controller', 'job_heartbeat_sec',       '0');   -- temporarily disabled

-- Database service-specific parameters

INSERT INTO `config` VALUES ('database', 'services_pool_size', '32');

-- Connection parameters for the Qserv Management Services

INSERT INTO `config` VALUES ('xrootd', 'auto_notify',         '1');
INSERT INTO `config` VALUES ('xrootd', 'host',                'xrootd-mgr');
INSERT INTO `config` VALUES ('xrootd', 'port',                '1094');
INSERT INTO `config` VALUES ('xrootd', 'request_timeout_sec', '600');

-- Default parameters for all workers unless overwritten in worker-specific
-- tables

INSERT INTO `config` VALUES ('worker', 'technology',                 'FS');
INSERT INTO `config` VALUES ('worker', 'svc_port',                   '25000');
INSERT INTO `config` VALUES ('worker', 'fs_port',                    '25001');
INSERT INTO `config` VALUES ('worker', 'num_svc_processing_threads', '16');
INSERT INTO `config` VALUES ('worker', 'num_fs_processing_threads',  '32');       -- double compared to the previous one to allow more elasticity
INSERT INTO `config` VALUES ('worker', 'fs_buf_size_bytes',          '4194304');  -- 4 MB
INSERT INTO `config` VALUES ('worker', 'data_dir',                   '/qserv/data/mysql');

SET SQL_MODE=@OLD_SQL_MODE ;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS ;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS ;
