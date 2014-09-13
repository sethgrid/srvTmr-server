CREATE TABLE `timer` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `place_id` varchar(100) NOT NULL,
  `timer_ms` int(11) NOT NULL,
  `ip` int(10) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)