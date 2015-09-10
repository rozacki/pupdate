create database if not exists monitoring;

drop dable if exists `events`;
'CREATE TABLE `events` (
  `id` int(11) NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `session` varchar(45) NOT NULL,
  `task_name` varchar(45) DEFAULT NULL COMMENT ''each task has its logical name for example: make, vehicle, test'',
  `task_id` varchar(45) DEFAULT NULL,
  `job` varchar(45) DEFAULT NULL,
  `event` varchar(255) NOT NULL,
  `data` mediumtext,
  PRIMARY KEY (`id`),
  KEY `task_idx` (`task_id`),
  KEY `session_idx` (`session`),
  KEY `job_idx` (`job`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1'