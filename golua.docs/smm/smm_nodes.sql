CREATE TABLE `smm_nodes` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL DEFAULT '',
  `host_name` varchar(40) NOT NULL DEFAULT '',
  `api_url` varchar(200) NOT NULL DEFAULT '',
  `code` varchar(100) NOT NULL DEFAULT '',
  `status` int(4) NOT NULL DEFAULT '1',
  `remark` varchar(1000) NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL DEFAULT 0,
  `modify_time` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;