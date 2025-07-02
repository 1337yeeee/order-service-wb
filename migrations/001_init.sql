create table if not exists orders (
    order_uid uuid primary key,
    track_number varchar,
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shardkey varchar,
    sm_id integer,
    date_created timestamp,
    oof_shard varchar
);

create table if not exists deliveries (
    order_uid uuid references orders(order_uid),
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar,
    primary key (order_uid)
);

create table if not exists payments (
    order_uid uuid references orders(order_uid),
    transaction varchar,
    request_id varchar,
    currency varchar,
    provider varchar,
    amount float,
    payment_dt integer,
    bank varchar,
    delivery_cost float,
    goods_total float,
    custom_fee float,
    primary key (order_uid)
);

create table if not exists items (
    order_uid uuid references orders(order_uid),
    chrt_id integer,
    track_number varchar,
    price float,
    rid varchar,
    name varchar,
    sale integer,
    size varchar,
    total_price float,
    nm_id integer,
    brand varchar,
    status integer
);
