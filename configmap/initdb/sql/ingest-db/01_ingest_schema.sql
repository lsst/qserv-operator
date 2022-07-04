
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 ;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 ;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL' ;


CREATE DATABASE qservIngest;
USE qservIngest;

-- --------------------------------------------------------------
-- Table `contribfile_queue`
-- --------------------------------------------------------------
--
-- The list of contributions file to load inside a Qserv database
-- Used as a queue by ingest jobs
CREATE TABLE `contribfile_queue` (

  `id`                    INTEGER UNSIGNED    NOT NULL AUTO_INCREMENT,
  `chunk_id`              INTEGER UNSIGNED    DEFAULT NULL ,              -- the id of the chunk to load
  `filepath`              VARCHAR(255)        NOT NULL ,                  -- the path of the chunk file to load
  `database`              VARCHAR(255)        NOT NULL ,                  -- the name of the target database
  `is_overlap`            BOOLEAN             DEFAULT NULL ,              -- is this file an overlap
  `locking_pod`           VARCHAR(255)        DEFAULT NULL,               -- the id of the latest pod which has locked the chunk
  `succeed`               BOOLEAN             NULL ,                      -- the status of the file:
                                                                          --   - NULL (pending),
                                                                          --   - 0 (error during latest ingest attempt),
                                                                          --   - 1 (success during latest ingest attempt)
  `table`                 VARCHAR(255)        NOT NULL ,                  -- the name of the target table

  PRIMARY KEY (`id`),
  UNIQUE KEY (`filepath`, `database`, `table`)
)
ENGINE = InnoDB;

-- --------------------------------------------------------------
-- Table `mutex`
-- --------------------------------------------------------------
--
-- Used to store the name of the pod which is currently locking contribfiles in the queue
-- Only store one row
CREATE TABLE `mutex` (
  `pod`           VARCHAR(255) DEFAULT NULL,  -- the id of the pod which is currently using the mutex to lock a set of contribution file
  `latest_move`   DATETIME NOT NULL           -- the latest time when the mutex was acquire/released
)
ENGINE = InnoDB;

INSERT `mutex`(`latest_move`) VALUES (NOW());
