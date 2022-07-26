CREATE TABLE `resource` (
                            `id` char(64) NOT NULL COMMENT '资源的实例Id',
                            `vendor` tinyint(1) NOT NULL,
                            `region` varchar(64)  NOT NULL,
                            `create_at` bigint NOT NULL,
                            `expire_at` bigint DEFAULT NULL,
                            `type` varchar(120) NOT NULL,
                            `name` varchar(255)  NOT NULL,
                            `description` varchar(255)  DEFAULT NULL,
                            `status` varchar(255)  NOT NULL,
                            `update_at` bigint DEFAULT NULL,
                            `sync_at` bigint DEFAULT NULL,
                            `account` varchar(255)  DEFAULT NULL,
                            `public_ip` varchar(64)  DEFAULT NULL,
                            `private_ip` varchar(64)  DEFAULT NULL,
                            PRIMARY KEY (`id`),
                            KEY `name` (`name`) USING BTREE,
                            KEY `status` (`status`),
                            KEY `private_ip` (`public_ip`) USING BTREE,
                            KEY `public_ip` (`public_ip`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;