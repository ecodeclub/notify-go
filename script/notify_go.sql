CREATE TABLE `delivery` (
                            `id` integer PRIMARY KEY,
                            `template_id` integer,
                            `status` integer,
                            `send_channel` varchar(255) COMMENT '消息发送渠道 10.IM 20.Push 30.短信 40.Email 50.公众号',
                            `msg_type` integer COMMENT '10.通知类消息 20.营销类消息 30.验证码类消息',
                            `proposer` varchar(255) COMMENT '业务方',
                            `creator` varchar(255),
                            `updator` varchar(255),
                            `is_delted` integer,
                            `created` timestamp,
                            `updated` timestamp
);

CREATE TABLE `target` (
                          `id` integer PRIMARY KEY,
                          `target_id_type` varchar(8) COMMENT '接收目标id类型',
                          `target_id` varchar(255) COMMENT '接收目标id',
                          `delivery_id` integer,
                          `status` integer,
                          `msg_content` text
);

CREATE TABLE `template` (
                            `id` integer PRIMARY KEY,
                            `country` varchar(255),
                            `type` integer COMMENT 'sms|email|push',
                            `en_content` text,
                            `chs_content` text,
                            `cht_content` text,
                            `creator` varchar(255),
                            `updator` varchar(255),
                            `is_delted` integer,
                            `created` timestamp,
                            `updated` timestamp
);

ALTER TABLE `target` ADD FOREIGN KEY (`delivery_id`) REFERENCES `delivery` (`id`);

ALTER TABLE `delivery` ADD FOREIGN KEY (`template_id`) REFERENCES `template` (`id`);
