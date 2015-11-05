CREATE TABLE IF NOT EXISTS guest  (
  gid INT(10) NOT NULL AUTO_INCREMENT,
  code VARCHAR(64) NULL DEFAULT NULL,
  name VARCHAR(64) NULL DEFAULT NULL,
  tag VARCHAR(64) NULL DEFAULT NULL,
  imgUrl VARCHAR(512) NULL DEFAULT NULL,
  PRIMARY KEY (gid)
);

CREATE TABLE IF NOT EXISTS win (
  wid INT(10) NOT NULL AUTO_INCREMENT,
  step INT(10) NOT NULL,
  gid INT(10) NOT NULL,
  prize VARCHAR(64) NULL DEFAULT NULL,
  PRIMARY KEY (wid),
  FOREIGN KEY (gid)
    REFERENCES guest(gid)
    ON DELETE CASCADE
);
