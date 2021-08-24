CREATE TABLE KillSwitch (
  kswitch_id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  feat_id SMALLINT UNSIGNED NOT NULL,
  min_version VARCHAR(50),
  max_version VARCHAR(50),
  active BOOLEAN NOT NULL,

  PRIMARY KEY (kswitch_id),
  KEY (feat_id)
);

CREATE TABLE KillSwitch2Browser (
  kswitch_id SMALLINT UNSIGNED NOT NULL,
  browser SMALLINT UNSIGNED NOT NULL,

  PRIMARY KEY (kswitch_id, browser),
  KEY (kswitch_id)
);

CREATE TABLE KillSwitchAuthorizedUser (
  user_id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  google_uid VARCHAR(64),
  email VARCHAR(255),
  access_level SMALLINT NOT NULL,

  PRIMARY KEY (user_id),
  KEY (google_uid),
  KEY (email)
);

CREATE TABLE KillSwitchAuditLog (
  log_id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
  data BLOB,

  PRIMARY KEY (log_id)
);
