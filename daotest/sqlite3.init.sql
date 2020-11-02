
drop table if exists `accessor` ;

create table `accessor` (
  `accessor_id` bigint unsigned not null,
  `accessor_name` varchar(20) default null,
  `enabled` tinyint unsigned not null default '1',
  `accessor_pub_key` text not null,
  `pay_notify_url` text,
  `refund_notify_url` text,
  `pub_key` text,
  `prv_key` text,
  `nonsense` varchar(32) not null,
  `update_time` timestamp not null default current_timestamp,
  primary key (`accessor_id`)
);

drop table if exists `merchant` ;

create table `merchant` (
  `merchant_id` bigint unsigned not null,
  `merchant_name` varchar(100) not null,
  `merchant_code` varchar(50) default null,
  `enabled` tinyint unsigned not null default '1',
  `update_time` timestamp not null default current_timestamp,
  `update_accessor` bigint unsigned not null,
  primary key (`merchant_id`),
  unique (`merchant_code`)
);

drop table if exists `accessor_merchant` ;

create table `accessor_merchant` (
  `accessor_id` bigint unsigned not null,
  `merchant_id` bigint unsigned not null,
  `enabled` tinyint unsigned not null default '1',
  `update_time` timestamp not null default current_timestamp,
  primary key (`accessor_id`,`merchant_id`)
);
