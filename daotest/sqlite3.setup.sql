
insert into `accessor`
      (`accessor_id`
      ,`accessor_name`
      ,`enabled`
      ,`accessor_pub_key`
      ,`pay_notify_url`
      ,`refund_notify_url`
      ,`pub_key`
      ,`prv_key`
      ,`nonsense`)
values(1001
      ,'1001'
      ,1
      ,'#PublicKeyString#'
      ,null
      ,null
      ,''
      ,''
      ,'test');

insert into `accessor`
      (`accessor_id`
      ,`accessor_name`
      ,`enabled`
      ,`accessor_pub_key`
      ,`pay_notify_url`
      ,`refund_notify_url`
      ,`pub_key`
      ,`prv_key`
      ,`nonsense`)
values(1002
      ,'1002'
      ,1
      ,'#PublicKeyString#'
      ,null
      ,null
      ,''
      ,''
      ,'test');

insert into `merchant`
      (`merchant_id`
      ,`merchant_name`
      ,`merchant_code`
      ,`enabled`
      ,`update_accessor`)
values(1001
      ,'1001'
      ,'m1001'
      ,1
      ,1001);

insert into `accessor_merchant`
      (`accessor_id`
      ,`merchant_id`
      ,`enabled`)
values(1001
      ,1001
      ,1)
