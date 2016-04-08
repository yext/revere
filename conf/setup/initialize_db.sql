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

CREATE TABLE triggers (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  level TINYINT NOT NULL DEFAULT 2,
  triggerOnExit BOOLEAN NOT NULL DEFAULT TRUE,
  periodMs INTEGER NOT NULL DEFAULT 0,
  targetType SMALLINT NOT NULL,
  target TEXT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE monitor_triggers (
  monitor_id INTEGER NOT NULL,
  subprobe TEXT NOT NULL DEFAULT '',
  trigger_id INTEGER NOT NULL,
  PRIMARY KEY (`monitor_id`, `trigger_id`),
  FOREIGN KEY (`monitor_id`) REFERENCES monitors(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`trigger_id`) REFERENCES triggers(`id`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE subprobes (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  monitor_id INTEGER NOT NULL,
  name VARCHAR(150) NOT NULL,
  archived DATETIME DEFAULT NULL,
  CONSTRAINT UNIQUE (monitor_id,name),
  FOREIGN KEY (`monitor_id`) REFERENCES monitors(`id`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE subprobe_statuses (
  subprobe_id INTEGER NOT NULL,
  recorded DATETIME NOT NULL,
  state TINYINT NOT NULL,
  silenced BOOLEAN NOT NULL DEFAULT FALSE,
  enteredState DATETIME NOT NULL,
  FOREIGN KEY (`subprobe_id`) REFERENCES subprobes(`id`) ON DELETE CASCADE,
  INDEX (`state`, `silenced`, `enteredState`, `recorded`, `subprobe_id`)
) ENGINE = InnoDB;

CREATE TABLE readings2 (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  subprobe_id INTEGER NOT NULL,
  recorded DATETIME NOT NULL,
  state TINYINT NOT NULL,
  FOREIGN KEY (`subprobe_id`) REFERENCES subprobes(`id`) ON DELETE CASCADE,
  INDEX (subprobe_id, recorded, id)
) ENGINE = InnoDB;

CREATE TABLE silences (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  monitor_id INTEGER NOT NULL,
  subprobes TEXT NOT NULL DEFAULT '',
  start DATETIME NOT NULL,
  end DATETIME NOT NULL,
  FOREIGN KEY (`monitor_id`) REFERENCES monitors(`id`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE labels (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  description TEXT NOT NULL DEFAULT ''
) ENGINE = InnoDB;

CREATE TABLE labels_monitors (
  label_id INTEGER NOT NULL,
  monitor_id INTEGER NOT NULL,
  subprobes TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (`monitor_id`, `label_id`),
  FOREIGN KEY (`monitor_id`) REFERENCES monitors(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`label_id`) REFERENCES labels(`id`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE label_triggers (
  label_id INTEGER NOT NULL,
  trigger_id INTEGER NOT NULL,
  PRIMARY KEY (`trigger_id`, `label_id`),
  FOREIGN KEY (`trigger_id`) REFERENCES triggers(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`label_id`) REFERENCES labels(`id`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE data_sources (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  sourceType SMALLINT NOT NULL,
  source TEXT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE settings (
  id INTEGER AUTO_INCREMENT PRIMARY KEY,
  settingType SMALLINT NOT NULL,
  setting TEXT NOT NULL
) ENGINE = InnoDB;
