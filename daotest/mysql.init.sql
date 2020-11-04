
drop table if exists `accessor` ;

create table `accessor` (
  `accessor_id` bigint unsigned not null comment '访问者标识',
  `accessor_name` varchar(20) default null comment '访问者名称',
  `enabled` tinyint unsigned not null default '1' comment '有效状态 0-无效 1-有效',
  `accessor_pub_key` text not null comment '访问者公钥',
  `pay_notify_url` text comment '支付回调地址',
  `refund_notify_url` text comment '退款回调地址',
  `pub_key` text comment '平台公钥',
  `prv_key` text comment '平台私钥',
  `nonsense` varchar(32) not null comment '随机字符串',
  `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  primary key (`accessor_id`)
) comment='访问者信息';

drop table if exists `merchant` ;

create table `merchant` (
  `merchant_id` bigint unsigned not null comment '商户标识',
  `merchant_name` varchar(100) not null comment '商户名称',
  `merchant_code` varchar(50) default null comment '商户编码',
  `enabled` tinyint unsigned not null default '1' comment '有效状态 0-无效 1-有效',
  `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `update_accessor` bigint unsigned not null comment '更新访问者',
  primary key (`merchant_id`),
  unique key `merchant_code_unique` (`merchant_code`)
) comment='商户信息';

drop table if exists `accessor_merchant` ;

create table `accessor_merchant` (
  `accessor_id` bigint unsigned not null comment '访问者标识',
  `merchant_id` bigint unsigned not null comment '商户标识',
  `enabled` tinyint unsigned not null default '1' comment '有效状态 0-无效 1-有效',
  `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  primary key (`accessor_id`,`merchant_id`)
) comment='访问者与商户关联信息';

drop table if exists `merchant_api_params` ;

create table `merchant_api_params` (
  `merchant_id` bigint unsigned not null comment '商户标识',
  `api_name` varchar(20) not null comment '支付接口名称',
  `param_name` varchar(50) not null comment '配置参数名',
  `param_value` text not null comment '配置参数值',
  `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `update_accessor` bigint unsigned not null comment '更新访问者',
  primary key (`merchant_id`,`api_name`,`param_name`)
) comment='商户接口配置表';
