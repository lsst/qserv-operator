
-- Used by xrootd Qserv plugin:
-- to publish LSST databases and chunks

DROP DATABASE IF EXISTS qservw_worker;
CREATE DATABASE qservw_worker;

GRANT ALL ON qservw_worker.* TO 'qsmaster'@'localhost';

CREATE TABLE qservw_worker.Dbs (

  `db` CHAR(200) NOT NULL,

  PRIMARY KEY (`db`)

) ENGINE=InnoDB;

CREATE TABLE qservw_worker.Chunks (

  `db`    CHAR(200)    NOT NULL,
  `chunk` INT UNSIGNED NOT NULL,

  UNIQUE KEY(`db`,`chunk`)

) ENGINE=InnoDB;

CREATE TABLE qservw_worker.Id (

  `id`      VARCHAR(64)  NOT NULL,
  `type`    ENUM('UUID') DEFAULT 'UUID',
  `created` TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,

  UNIQUE KEY (`type`)

) ENGINE=InnoDB;

-- TODO: This needs to be changed to generate a unique identifier by calling MySQL function UUID()
-- insted of the template parameter. The proposed change could be done later after a refined operational
-- model of Qserv will be implementated. In proposed model there will be a special step for installing
-- and initializing Qserv databases. This would also include auto-generating worker identities and
-- populating Replication system's database table "qservReplica.config_worker" with unique identifiers
-- of workers managed by the system.
INSERT INTO qservw_worker.Id (`id`) VALUES ('<HOST>');

CREATE TABLE IF NOT EXISTS qservw_worker.QMetadata (

  `metakey` CHAR(64) NOT NULL COMMENT 'Key string',
  `value`   TEXT         NULL COMMENT 'Value string',

  PRIMARY KEY (`metakey`)

) ENGINE = InnoDB COMMENT = 'Metadata about database as a whole, key-value pairs';

INSERT INTO qservw_worker.QMetadata (`metakey`, `value`) VALUES ('version', '2');

GRANT ALL ON `q\_memoryLockDb`.* TO 'qsmaster'@'localhost';

-- Subchunks databases
GRANT ALL ON `Subchunks\_%`.* TO 'qsmaster'@'localhost';