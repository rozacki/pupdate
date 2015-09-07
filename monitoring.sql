create database if not exists monitoring;

drop dable if exists `events`;
'CREATE TABLE `events` (
  `id` int(11) NOT NULL,
  `ts` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `session` varchar(45) NOT NULL,
  `task` int(11) NOT NULL,
  `job` int(11) DEFAULT NULL,
  `event` varchar(255) NOT NULL,
  `data` mediumtext,
  PRIMARY KEY (`id`),
  KEY `task_idx` (`task`),
  KEY `session_idx` (`session`),
  KEY `job_idx` (`job`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1'