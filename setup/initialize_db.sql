CREATE DATABASE revere DEFAULT CHARACTER SET = utf8mb4;

use revere;

CREATE TABLE monitors (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  owner VARCHAR(60) NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  response TEXT NOT NULL DEFAULT '',
  probeType SMALLINT NOT NULL,
  probe TEXT NOT NULL,
  changed DATETIME NOT NULL,
  version INTEGER NOT NULL DEFAULT 0,
  archived DATETIME DEFAULT NULL
) ENGINE = InnoDB;

CREATE TABLE readings (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config_id int(11),
  subprobe varchar(255),
  state int(8),
  time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`config_id`) REFERENCES configurations (id) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE alerts (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config_id int(11),
  subprobe varchar(255),
  time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`config_id`) REFERENCES configurations (id) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE silenced_alerts (
  config_id int(11),
  subprobe  varchar(255),
  silence_time timestamp NOT NULL,
  FOREIGN KEY (`config_id`) REFERENCES configurations (id) ON DELETE CASCADE,
  PRIMARY KEY (`config_id`, `subprobe`)
) ENGINE = InnoDB;
