SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 ;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 ;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL' ;

CREATE DATABASE qservIngest;
USE qservIngest;

-- --------------------------------------------------------------
-- Table `task`
-- --------------------------------------------------------------
--
-- The list of chunks to load inside a Qserv database
-- Used as a queue by loader jobs
CREATE TABLE `task` (

  `chunk_id`              INTEGER             NOT NULL ,                  -- the id of the chunk to load
  `chunk_file_url`        VARCHAR(255)        NOT NULL ,                  -- the url of the chunk file to load
  `database`              VARCHAR(255)        NOT NULL ,                  -- the name of the target database
  `pod`                   VARCHAR(255)        DEFAULT NULL ,              -- the name of the pod which launch the ingest
  `table`                 VARCHAR(255)        NOT NULL ,                  -- the name of the target table

  PRIMARY KEY (`chunk_id`, `chunk_file_url`, `database`, `table`)
)
ENGINE = InnoDB;
