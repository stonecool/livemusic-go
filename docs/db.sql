-- 爬虫账号
CREATE TABLE `crawl_account` (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `account_type`      VARCHAR(100) NOT NULL COMMENT 'account type',
    `account_id`        VARCHAR(100) NOT NULL COMMENT 'account id',
    `account_name`      VARCHAR(100) NOT NULL COMMENT 'account name',
    `cookies`           BLOB COMMENT 'cookies',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl account';

-- 爬虫消息生产者
CREATE TABLE crawl_msg (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `data_type`         VARCHAR(100) NOT NULL COMMENT 'data type',
    `data_id`           INT(10) UNSIGNED NOT NULL COMMENT 'data id',
    `account_type`      VARCHAR(100) NOT NULL COMMENT 'account type',
    `target_account_id` VARCHAR(100) NOT NULL COMMENT 'target account id',
    `mark`              VARCHAR(100) NOT NULL COMMENT 'mark',
    `count`             INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'count',
    `first_time`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'first time',
    `last_time`         INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'last time',
    `created_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'created time',
    `updated_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'updated time',
    `deleted_at`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT 'deleted time',

    PRIMARY KEY (`id`)
    UNIQUE KEY data_type__data_id__crawl_type (`data_type`, `data_id`, `account_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='crawl msg';

-- CREATE TABLE `musician` (
--     `id`            INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
--     `name`          VARCHAR(100) DEFAULT '' COMMENT '姓名',
-- --     `intro`,
-- --     `gender`        TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '性别，0男，1女',
-- --     `avatar_url`,
-- --     `background_url`,
-- --     `location`      VARCHAR(100) DEFAULT '' COMMENT '',
-- --     `email`,
-- --     `telephone`,
-- --     `wx`,
-- --     `weibo`,
-- --     `wx_public_no`,
-- --     `163_music_home_page`,
-- --     `qq_music_home_page`,
-- --     `website`,
-- --     `type` COMMENT '0单人，1乐队',
-- --     `member_id_list`,
-- --     `status` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0未确认。1已确认'
--     `created_on`    INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '创建时间',
--     `modified_on`   INT(10) UNSIGEND NOT NULL DEFAULT 0 COMMENT '更新时间',
--     `deleted_on`    INT(10) UNSIGEND NOT NULL DEFAULT 0 COMMENT '删除时间',
--     PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='音乐人';

CREATE TABLE `livehouse` (
    `id`            INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name`          VARCHAR(100) DEFAULT '' COMMENT 'name',
    `location`      VARCHAR(100) DEFAULT '' COMMENT 'location',
    `telephone`     VARCHAR(30) DEFAULT '' COMMENT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='livehouse';

-- CREATE TABLE `music_festival` (
--     `id`            INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
--     `name`          VARCHAR(100) DEFAULT '' COMMENT '姓名',
--     PRIMARY KEY (`id`)
-- )
--
-- CREATE TABLE `show` (
--     `id`            INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
--     `musician_id`   INT(10),
-- --     `theme`         VARCHAR(100) DEFAULT '' COMMENT '主题',
-- --     `style`,
-- --     `space`,
-- --     `time`,
-- --     `ticket_price`        tinyint DEFAULT 0 COMMENT '性别',
-- --     `buy_ticket_qr_code`,
-- --     `buy_ticket_link`,
--     PRIMARY KEY (`id`)
-- ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='演出信息';





-- 爬取记录
CREATE TABLE `crawl_log` (
    `id`                INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `instance_id`       INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '源id',
    `crawl_time`        INT(10) UNSIGNED NOT NULL DEFAULT 0 COMMENT '爬取时间',
    `metadata`          BLOB NOT NULL COMMENT '原始数据',
    `http_status`       SMALLINT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'http status',
    `status`            TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '记录状态，0爬取未解析、1已解析',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='server log';

