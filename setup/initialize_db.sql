CREATE DATABASE revere;

use revere;

CREATE TABLE configurations (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config longtext,
  emails char(255)
)

CREATE TABLE alert_history (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config_id int(11),
  value int(20),
  state char(8),
  FOREIGN KEY (`config_id`) REFERENCES configurations (id) ON DELETE CASCADE
)
