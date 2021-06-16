CREATE TABLE IF NOT EXISTS `fabric`
(
    id BINARY(16) NOT NULL,
    name VARCHAR (255) NOT NULL,
    cost FLOAT NOT NULL,
    amount FLOAT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;