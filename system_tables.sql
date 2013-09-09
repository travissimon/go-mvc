CREATE SCHEMA `gomvc` ;

use gomvc;

CREATE TABLE User (
  Id bigint(20) NOT NULL AUTO_INCREMENT,
  IpAddress varchar(256) NOT NULL,
  Username varchar(255) NOT NULL,
  Password varchar(128) NOT NULL,
  RecoveryEmailAddress varchar(512) NOT NULL,
  PRIMARY KEY (Id),
  UNIQUE KEY Id_UNIQUE (Id),
  UNIQUE KEY Username_UNIQUE (Username)
);

CREATE TABLE Authentication (
  SessionId int(11) NOT NULL,
  UserId int(11) NOT NULL,
  IPAddress varchar(45) NOT NULL,
  PRIMARY KEY (SessionId)
);

