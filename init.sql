CREATE TABLE `guest` (
  `gid` int(10) NOT NULL AUTO_INCREMENT,
  `code` varchar(64) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `tag` varchar(64) DEFAULT NULL,
  `prize` varchar(64) DEFAULT NULL,
  `step` int(10) DEFAULT NULL,
  `imgUrl` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`gid`),
  UNIQUE KEY `code` (`code`)
);

CREATE TABLE IF NOT EXISTS assignment (
  tag VARCHAR(64) NOT NULL,
  prize VARCHAR(64) NOT NULL,
  maxWin INT(10) NULL DEFAULT 0,
  PRIMARY KEY (tag, prize)
);

CREATE TABLE IF NOT EXISTS prize (
  id INT(10) NOT NULL,
  prize VARCHAR(64) NOT NULL,
  maxStepWin INT(10) NULL DEFAULT 1,
  PRIMARY KEY (prize),
  UNIQUE KEY `prize` (`prize`)
);

CREATE VIEW `assignment-sum` AS
SELECT assignment.tag, assignment.prize, assignment.maxWin
FROM assignment
UNION SELECT COUNT(guest.gid) as sumWin FROM guest JOIN assignment ON guest.prize = assignment.prize AND guest.tag = assignment.tag;