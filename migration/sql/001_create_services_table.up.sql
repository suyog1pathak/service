CREATE TABLE `services`
(
    `id`                bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at`        datetime(3) DEFAULT NULL,
    `updated_at`        datetime(3) DEFAULT NULL,
    `deleted_at`        datetime(3) DEFAULT NULL,
    `name`              varchar(50) DEFAULT NULL,
    `description`       LONGTEXT DEFAULT NULL,
    `version`           INT DEFAULT 1,
    `is_active`         BOOLEAN DEFAULT TRUE,
    `tags`              VARCHAR(255)  DEFAULT  NULL, -- Assuming each tag is up to 20 characters and you have 10 tags
    PRIMARY KEY (`id`)
);
