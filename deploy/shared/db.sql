USE `umkm`;
CREATE TABLE `goods` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `stocks` int(11) DEFAULT 0,
    `price` double DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

INSERT INTO `goods` (`id`, `name`, `stocks`, `price`) VALUES
(1, 'Kopi', 100, 3000),
(2, 'Pisang Goreng', 45, 1500),
(3, 'Bakwan', 50, 1500),
(4, 'Teh manis', 100, 2000),
(5, 'Teh tawar', 100, 1000),
(6, 'Pisang Keju', 25, 2500),
(7, 'Lumpia Udang', 30, 2500);

CREATE TABLE `transaction_details` (
    `id_transaction` bigint(20) NOT NULL,
    `id_goods` int(11) NOT NULL,
    `total_goods` int(11) NOT NULL DEFAULT '1',
    `created_at` bigint(20) NOT NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

CREATE TABLE `transactions` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `id_user` int(11) DEFAULT NULL,
    `total_amount` double DEFAULT NULL,
    `status` tinyint(4) DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;