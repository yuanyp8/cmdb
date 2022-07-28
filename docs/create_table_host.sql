
CREATE TABLE `host` (
                        `resource_id` char(64) NOT NULL,
                        `cpu` tinyint NOT NULL,
                        `memory` int NOT NULL,
                        `gpu_amount` tinyint DEFAULT NULL,
                        `gpu_spec` varchar(255)  DEFAULT NULL,
                        `os_type` varchar(255)  DEFAULT NULL,
                        `os_name` varchar(255) DEFAULT NULL,
                        `serial_number` varchar(120) DEFAULT NULL,
                        PRIMARY KEY (`resource_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;