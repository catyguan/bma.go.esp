CREATE TABLE `golua_sqlloader` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `module` varchar(100) NOT NULL DEFAULT '',
  `script` varchar(100) NOT NULL DEFAULT '',
  `content` longtext NULL,
  `create_time` int(11) NOT NULL DEFAULT 0,
  `modify_time` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE(`module`, `script`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;