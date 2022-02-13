CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table check_users
(
    id    uuid default uuid_generate_v4()    not null
        constraint check_users_pk
            primary key,
    phone varchar(255)                       not null,
    code  text default ''::character varying not null
);

create unique index check_users_id_uindex
    on check_users (id);

create unique index check_users_phone_uindex
    on check_users (phone);

create table users
(
    id         uuid                     default uuid_generate_v4()        not null
        constraint users_pk
            primary key,
    first_name varchar(255)             default ''::character varying     not null,
    phone      varchar(255)                                               not null,
    email      varchar(255)             default ''::character varying     not null,
    code       text                     default ''::character varying     not null,
    created_at timestamp with time zone default CURRENT_TIMESTAMP         not null,
    updated_at timestamp with time zone default CURRENT_TIMESTAMP         not null,
    balance    double precision         default 0                         not null,
    action     double precision         default 4000                      not null,
    role       varchar(255)             default 'USER'::character varying not null
);

create unique index users_id_uindex
    on users (id);

create unique index users_phone_uindex
    on users (phone);


create table cheques
(
    id           uuid                     default uuid_generate_v4()               not null
        constraint cheques_pk
            primary key,
    created_at   timestamp with time zone default CURRENT_TIMESTAMP                not null,
    updated_at   timestamp with time zone default CURRENT_TIMESTAMP                not null,
    date         varchar(255)             default ''::character varying            not null,
    check_amount varchar(255)             default ''::character varying            not null,
    fn           varchar(255)             default ''::character varying            not null,
    fd           varchar(255)             default ''::character varying            not null,
    winning      double precision         default 0                                not null,
    user_id      uuid                                                              not null
        constraint cheques_users_id_fk
            references users,
    "check"      varchar(255)             default ''::character varying not null,
    fp           varchar(255)             default ''::character varying            not null,
    photo_id     varchar(255)             default ''::character varying            not null
);


create unique index cheques_id_uindex
    on cheques (id);

create table cities
(
    id   serial                                     not null
        constraint cities_pk
            primary key,
    city varchar(255) default ''::character varying not null
);


create unique index cities_id_uindex
    on cities (id);

create table email_problems
(
    id         uuid                     default uuid_generate_v4()    not null
        constraint email_problems_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP     not null,
    data       text                     default ''::character varying not null,
    email      text                     default ''::character varying not null
);


create unique index email_problems_id_uindex
    on email_problems (id);

create table files
(
    id         uuid                     default uuid_generate_v4()    not null
        constraint files_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP     not null,
    url        text                     default ''::character varying not null,
    length     integer                  default 0                     not null,
    mime       text                     default ''::character varying not null,
    bucket     varchar(255)             default ''::character varying not null,
    object     text                     default ''::character varying not null,
    name       text                     default ''::character varying not null,
    user_id    uuid                                                   not null,
    review     boolean                  default false                 not null
);

create unique index files_id_uindex
    on files (id);


create table gifts
(
    id    serial                                         not null
        constraint gifts_pk
            primary key,
    name  text             default ''::character varying not null,
    price double precision default 0                     not null,
    sum   double precision default 0                     not null
);

create table good_resp_check
(
    id                         uuid                     default uuid_generate_v4()    not null
        constraint good_resp_check_pk
            primary key,
    created_at                 timestamp with time zone default CURRENT_TIMESTAMP     not null,
    code                       integer                  default 0                     not null,
    "user"                     varchar(255)             default ''::character varying not null,
    fns_url                    varchar(255)             default ''::character varying not null,
    user_inn                   varchar(255)             default ''::character varying not null,
    data_time                  varchar(255)             default ''::character varying not null,
    kkt_req_id                 varchar(255)             default ''::character varying not null,
    total_sum                  bigint                   default 0                     not null,
    fiscal_sign                bigint                   default 0                     not null,
    retail_place               varchar(255)             default ''::character varying not null,
    shift_number               bigint                   default 0                     not null,
    operation_type             bigint                   default 0                     not null,
    request_number             bigint                   default 0                     not null,
    fiscal_driver_number       text                     default ''::character varying not null,
    retail_place_address       text                     default ''::character varying not null,
    fiscal_document_number     bigint                   default 0                     not null,
    fiscal_document_format_ver bigint                   default 0                     not null,
    url                        text                     default ''::character varying not null,
    user_id                    uuid                                                   not null
);

create unique index good_resp_check_id_uindex
    on good_resp_check (id);

create table log_sms
(
    id         uuid                     default uuid_generate_v4()    not null
        constraint log_sms_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP     not null,
    resp       text                     default ''::character varying not null,
    user_id    uuid                                                   not null
);


create unique index log_sms_id_uindex
    on log_sms (id);

create table logger_req_check
(
    id         uuid                     default uuid_generate_v4()    not null
        constraint logger_req_check_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP     not null,
    logger     text                     default ''::character varying not null,
    user_id    uuid                                                   not null
);

create unique index logger_req_check_id_uindex
    on logger_req_check (id);


create table logger_resp_check
(
    id                  uuid                     default uuid_generate_v4()    not null
        constraint logger_resp_check_pk
            primary key,
    created_at          timestamp with time zone default CURRENT_TIMESTAMP     not null,
    logger              text                     default ''::character varying not null,
    user_id             uuid                                                   not null,
    logger_req_check_id text                                                   not null
);


create unique index logger_resp_check_id_uindex
    on logger_resp_check (id);


create table position_in_check
(
    id                 uuid                     default uuid_generate_v4()    not null
        constraint position_in_check_pk
            primary key,
    created_at         timestamp with time zone default CURRENT_TIMESTAMP     not null,
    name               text                     default ''::character varying not null,
    price              double precision         default 0                     not null,
    count              double precision         default 0                     not null,
    sum                double precision         default 0                     not null,
    marker             varchar(255)             default ''::character varying not null,
    good_resp_check_id uuid                                                   not null,
    user_id            uuid                                                   not null
);

create unique index position_in_check_id_uindex
    on position_in_check (id);


create table products
(
    id         uuid                     default uuid_generate_v4() not null
        constraint products_pk
            primary key,
    gift_id    integer                                             not null,
    count      integer                                             not null,
    user_id    uuid                                                not null
        constraint products_users_id_fk
            references users,
    created_at timestamp with time zone default CURRENT_TIMESTAMP  not null
);

create unique index products_id_uindex
    on products (id);



create table request_gift
(
    id          uuid                     default uuid_generate_v4()    not null
        constraint request_gift_pk
            primary key,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP     not null,
    phone       varchar(255)             default ''::character varying not null,
    gift_id     integer                                                not null,
    prize_name  varchar(255)             default ''::character varying not null,
    certificate text                     default ''::character varying not null,
    product_id  uuid                                                   not null,
    user_id     uuid                                                   not null,
    sent        boolean                  default false                 not null,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP     not null
);

create unique index request_gift_id_uindex
    on request_gift (id);



create table shops
(
    id   serial                                     not null
        constraint shops_pk
            primary key,
    name text         default ''::character varying not null,
    inn  varchar(255) default ''::character varying not null
);

create unique index shops_id_uindex
    on shops (id);


create table whiskey
(
    id           serial not null
        constraint whiskey_pk
            primary key,
    product_name text   not null
);

create unique index whiskey_id_uindex
    on whiskey (id);


create table yoomoney_log
(
    id         uuid                     default uuid_generate_v4()    not null
        constraint yoomoney_log_pk
            primary key,
    created_at timestamp with time zone default CURRENT_TIMESTAMP     not null,
    user_id    uuid                                                   not null,
    type       varchar(255)             default ''::character varying not null,
    error      varchar(255)             default ''::character varying not null,
    request_id varchar(255)             default ''::character varying not null,
    status     varchar(255)             default ''::character varying not null,
    amount     varchar(255)             default ''::character varying not null,
    payment_id varchar(255)             default ''::character varying not null,
    invoice_id varchar(255)             default ''::character varying not null,
    phone      varchar(255)             default ''::character varying not null
);

create unique index yoomoney_log_id_uindex
    on yoomoney_log (id);
