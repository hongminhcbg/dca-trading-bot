DROP TABLE IF EXISTS `dca_bot`.`order_tracking`;
CREATE TABLE `dca_bot`.`order_tracking`
(
    `id`           bigint    NOT NULL AUTO_INCREMENT,
    `index_num`    INT(10),
    `selected_num` INT(10),
    `status`       VARCHAR(128),
    `error`        longtext,
    `raw_response` longtext,
    `created_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY            `idx_index_num` (`index_num`)
);


-- postgresql version
DROP TABLE IF EXISTS order_tracking;

CREATE TABLE order_tracking
(
    id           serial PRIMARY KEY,
    index_num    INT UNIQUE NOT NULL,
    selected_num INT,
    status       VARCHAR(128),
    error        TEXT,
    raw_response TEXT,
    executed_qty VARCHAR(64),
    usdt_qty     VARCHAR(64),
    created_at   timestamp,
    updated_at   timestamp
);
