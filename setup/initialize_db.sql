CREATE DATABASE revere;

use revere;

CREATE TABLE configurations (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config longtext,
  emails varchar(255)
);

CREATE TABLE readings (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config_id int(11),
  subprobe varchar(255),
  state int(8),
  time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`config_id`) REFERENCES configurations (id) ON DELETE CASCADE
);

CREATE TABLE alert_histories (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config_id int(11),
  subprobe varchar(255),
  time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
