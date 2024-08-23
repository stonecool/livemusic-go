-- 爬虫账号
CREATE TABLE `crawl_account` (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `account_type`      VARCHAR(100) NOT NULL COMMENT 'account type',
    `account_id`        VARCHAR(100) NOT NULL COMMENT 'account id',
    `account_name`      VARCHAR(100) NOT NULL COMMENT 'account name',
    `last_login_url`    VARCHAR(100) NOT NULL COMMENT 'last login url',
    `cookies`           BLOB COMMENT 'cookies',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl account';

-- 爬虫例程
CREATE TABLE crawl_routine (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `account_type`      VARCHAR(100) NOT NULL COMMENT 'account type',
    `data_type`         VARCHAR(100) NOT NULL COMMENT 'data type, livehouse',
    `data_id`           INT(10) UNSIGNED NOT NULL COMMENT 'data id',
    `target_account_id` VARCHAR(100) NOT NULL COMMENT 'target account id',
    `mark`              VARCHAR(100) NOT NULL COMMENT 'mark',
    `count`             INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'count',
    `first_time`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'first time',
    `last_time`         INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'last time',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`),
    UNIQUE KEY idx_data_type_data_id_crawl_type (`data_type`, `data_id`, `account_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl msg';

CREATE TABLE `crawl_data_wechat` (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `rid`               INT(10) UNSIGNED NOT NULL,
    `uid`               VARCHAR(100) NOT NULL COMMENT 'title_datetime',
    `title`             VARCHAR(100) NOT NULL COMMENT 'title',
    `cover`             VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'cover link',
    `link`              VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'cover link',
    `original_link`     VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'original link',
    `datetime`          INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'date time',
    `state`             TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'state',
    `type`              TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'type',
    `sub_type`          TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'sub type',
    `raw_id`            INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'raw id',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl wechat data';

CREATE TABLE `crawl_data_wechat_raw` (
     `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
     `rid`               INT(10) UNSIGNED NOT NULL,
     `uid`               VARCHAR(100) NOT NULL COMMENT 'title_datetime',
     `title`             VARCHAR(100) NOT NULL COMMENT 'title',
     `cover`             VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'cover link',
     `link`              VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'cover link',
     `original_link`     VARCHAR(1000) NOT NULL DEFAULT '' COMMENT 'original link',
     `datetime`          INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'date time',
     `state`             TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'state',
     `type`              TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'type',
     `sub_type`          TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'sub type',
     `raw_id`            INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'raw id',
     `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
     `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
     `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl wechat raw data';

-- livehouse
CREATE TABLE `livehouse` (
     `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
     `name`              VARCHAR(100) DEFAULT '' COMMENT 'name',
     `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
     `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
     `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

     PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='livehouse';
