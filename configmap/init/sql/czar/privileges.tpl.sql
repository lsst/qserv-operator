CREATE USER IF NOT EXISTS 'qsmaster'@'localhost';

GRANT ALL ON qservResult.* TO 'qsmaster'@'localhost';

-- Secondary index database (i.e. objectId/chunkId relation)
-- created by integration test script/loader for now
GRANT ALL ON qservMeta.* TO 'qsmaster'@'localhost';

-- CSS database
GRANT ALL ON qservCssData.* TO 'qsmaster'@'localhost';

-- Create user for external monitoring applications
CREATE USER IF NOT EXISTS 'monitor'@'localhost' IDENTIFIED BY '<MYSQL_MONITOR_PASSWORD>';
GRANT PROCESS ON *.* TO 'monitor'@'localhost';

FLUSH PRIVILEGES;
