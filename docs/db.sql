-- chrome
CREATE TABLE `chromes` (
    `id`               INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `ip`               VARCHAR(20) NOT NULL COMMENT 'ip',
    `port`             INT(10) UNSIGNED NOT NULL COMMENT 'port',
    `debugger_url`     VARCHAR(100) NOT NULL COMMENT 'debugger_url',
    `state`            TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'state',
    `created_at`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`),
    UNIQUE KEY unique_addr (ip, port)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='chromes';

-- 爬虫账号
CREATE TABLE `accounts` (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `category`          VARCHAR(50) NOT NULL COMMENT 'category',
    `name`              VARCHAR(255) NOT NULL COMMENT 'name',
    `last_url`          VARCHAR(255) NOT NULL COMMENT 'last url',
    `cookies`           BLOB COMMENT 'cookies',
    `instance_id`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'instance id',
    `state`             TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'status',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl accounts';

-- 爬虫任务
CREATE TABLE tasks (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `category`          VARCHAR(50) NOT NULL COMMENT 'category',
    `target_id`         VARCHAR(100) NOT NULL COMMENT 'target id',
    `meta_type`         VARCHAR(50) NOT NULL COMMENT 'meta type',
    `meta_id`           INT(10) UNSIGNED NOT NULL COMMENT 'meta id',
    `mark`              VARCHAR(100) NOT NULL COMMENT 'mark',
    `cron_spec`         VARCHAR(20) NOT NULL COMMENT 'cron spec',
    `first_time`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'first time',
    `last_time`         INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'last time',
    `count`             INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'count',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`),
    UNIQUE KEY idx_category__meta (`category`, `meta_type`, `meta_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl tasks';



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
