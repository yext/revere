CREATE DATABASE revere DEFAULT CHARACTER SET = utf8mb4;

use revere;

CREATE TABLE pfx_monitors (
  monitorid INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  owner VARCHAR(60) NOT NULL DEFAULT '',
  description TEXT NOT NULL,
  response TEXT NOT NULL,
  probetype SMALLINT NOT NULL,
  probe TEXT NOT NULL,
  changed DATETIME NOT NULL,
  version INTEGER NOT NULL DEFAULT 0,
  archived DATETIME DEFAULT NULL
) ENGINE = InnoDB;

CREATE TABLE pfx_triggers (
  triggerid INTEGER AUTO_INCREMENT PRIMARY KEY,
  level TINYINT NOT NULL DEFAULT 2,
  triggeronexit BOOLEAN NOT NULL DEFAULT TRUE,
  periodmilli INTEGER NOT NULL DEFAULT 0,
  targettype SMALLINT NOT NULL,
  target TEXT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE pfx_monitor_triggers (
  monitorid INTEGER NOT NULL,
  subprobes TEXT NOT NULL,
  triggerid INTEGER NOT NULL,
  PRIMARY KEY (`monitorid`, `triggerid`),
  FOREIGN KEY (`monitorid`) REFERENCES pfx_monitors(`monitorid`) ON DELETE CASCADE,
  FOREIGN KEY (`triggerid`) REFERENCES pfx_triggers(`triggerid`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE pfx_subprobes (
  subprobeid INTEGER AUTO_INCREMENT PRIMARY KEY,
  monitorid INTEGER NOT NULL,
  name VARCHAR(150) NOT NULL,
  archived DATETIME DEFAULT NULL,
  CONSTRAINT UNIQUE (monitorid,name),
  FOREIGN KEY (`monitorid`) REFERENCES pfx_monitors(`monitorid`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE pfx_subprobe_statuses (
  subprobeid INTEGER NOT NULL,
  recorded DATETIME NOT NULL,
  state TINYINT NOT NULL,
  silenced BOOLEAN NOT NULL DEFAULT FALSE,
  enteredstate DATETIME NOT NULL,
  lastnormal DATETIME NOT NULL,
  FOREIGN KEY (`subprobeid`) REFERENCES pfx_subprobes(`subprobeid`) ON DELETE CASCADE,
  INDEX (`state`, `silenced`, `enteredstate`, `recorded`, `subprobeid`)
) ENGINE = InnoDB;

CREATE TABLE pfx_readings (
  readingid BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  subprobeid INTEGER NOT NULL,
  recorded DATETIME NOT NULL,
  state TINYINT NOT NULL,
  FOREIGN KEY (`subprobeid`) REFERENCES pfx_subprobes(`subprobeid`) ON DELETE CASCADE,
  INDEX (subprobeid, recorded, readingid)
) ENGINE = InnoDB;

CREATE TABLE pfx_silences (
  silenceid BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  monitorid INTEGER NOT NULL,
  subprobes TEXT NOT NULL,
  start DATETIME NOT NULL,
  end DATETIME NOT NULL,
  FOREIGN KEY (`monitorid`) REFERENCES pfx_monitors(`monitorid`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE pfx_labels (
  labelid INTEGER AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(30) NOT NULL,
  description TEXT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE pfx_labels_monitors (
  labelid INTEGER NOT NULL,
  monitorid INTEGER NOT NULL,
  subprobes TEXT NOT NULL,
  PRIMARY KEY (`monitorid`, `labelid`),
  FOREIGN KEY (`monitorid`) REFERENCES pfx_monitors(`monitorid`) ON DELETE CASCADE,
  FOREIGN KEY (`labelid`) REFERENCES pfx_labels(`labelid`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE pfx_label_triggers (
  labelid INTEGER NOT NULL,
  triggerid INTEGER NOT NULL,
  PRIMARY KEY (`triggerid`, `labelid`),
  FOREIGN KEY (`triggerid`) REFERENCES pfx_triggers(`triggerid`) ON DELETE CASCADE,
  FOREIGN KEY (`labelid`) REFERENCES pfx_labels(`labelid`) ON DELETE CASCADE
) ENGINE = InnoDB;

CREATE TABLE pfx_data_sources (
  sourceid INTEGER AUTO_INCREMENT PRIMARY KEY,
  sourcetype SMALLINT NOT NULL,
  source TEXT NOT NULL
) ENGINE = InnoDB;

CREATE TABLE pfx_settings (
  settingid INTEGER AUTO_INCREMENT PRIMARY KEY,
  settingtype SMALLINT NOT NULL,
  setting TEXT NOT NULL
) ENGINE = InnoDB;
