update records set bought_value = 0 where bought_value is null;
update records set current_value = 0 where current_value is null;
update records set realized_value = 0 where realized_value is null;
update assets set default_increment = 0 where default_increment is null;

create table records_dg_tmp
(
    id             INTEGER
        primary key autoincrement,
    asset_id       INTEGER                  not null
        references assets,
    date           DATE                     not null,
    bought_value   DECIMAL(14, 2) default 0 not null,
    current_value  DECIMAL(14, 2) default 0 not null,
    realized_value DECIMAL(14, 2) default 0 not null,
    note           TEXT           default NULL
);

insert into records_dg_tmp(id, asset_id, date, bought_value, current_value, realized_value, note)
select id, asset_id, date, bought_value, current_value, realized_value, note
from records;

drop table records;

alter table records_dg_tmp
    rename to records;

create unique index idx_record
    on records (asset_id, date);

create index idx_record_date
    on records (date);

create table assets_dg_tmp
(
    id                INTEGER
        primary key autoincrement,
    name              TEXT                     not null,
    broker            TEXT                     not null,
    type_id           INTEGER                  not null
        references asset_types,
    default_increment DECIMAL(14, 2) default 0 not null,
    sequence          INTEGER                  not null,
    is_active         BOOLEAN                  not null
);

insert into assets_dg_tmp(id, name, broker, type_id, default_increment, sequence, is_active)
select id, name, broker, type_id, default_increment, sequence, is_active
from assets;

drop table assets;

alter table assets_dg_tmp
    rename to assets;

create unique index idx_asset
    on assets (name, broker);