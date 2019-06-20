CREATE USER IF NOT EXISTS 'qsmaster'@'localhost';

GRANT ALL ON `q\_memoryLockDb`.* TO 'qsmaster'@'localhost';

-- Subchunks databases
GRANT ALL ON `Subchunks\_%`.* TO 'qsmaster'@'localhost';

-- Create user for external monitoring applications
CREATE USER IF NOT EXISTS 'monitor'@'localhost' IDENTIFIED BY '<MYSQL_MONITOR_PASSWORD>';
GRANT PROCESS ON *.* TO 'monitor'@'localhost';

FLUSH PRIVILEGES;