delete from ko_cis_task where id !='';

alter table ko_cis_task
    add result mediumtext null;

alter table ko_cis_task
    add total_pass int default 0 null;

alter table ko_cis_task
    add total_info int default 0 null;

alter table ko_cis_task
    add total_fail int default 0 null;

alter table ko_cis_task
    add total_warn int default 0 null;

alter table ko_cis_task
    add policy varchar(255) default '' null;