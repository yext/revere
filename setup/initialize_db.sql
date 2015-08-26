CREATE DATABASE revere;

use revere;

CREATE TABLE configurations (
  id int(11) AUTO_INCREMENT PRIMARY KEY,
  config longtext,
  emails char(255)
)
