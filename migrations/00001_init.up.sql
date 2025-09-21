begin;

create table provinces (
  id bigint unsigned not null primary key auto_increment,
  title varchar(255) not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

create table cities (
  id bigint unsigned not null primary key auto_increment,
  province_id bigint unsigned not null,
  title varchar(255) not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  foreign key (province_id) references provinces (id)
);

create table merchants (
  id bigint unsigned not null primary key auto_increment,
  name varchar(255) not null,
  title varchar(255) not null,
  ozone_card_merchant_id bigint unsigned not null,
  website varchar(255) not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

create table branches (
  id bigint unsigned not null primary key auto_increment,
  title varchar(255) not null,
  merchant_id bigint unsigned not null,
  city_id bigint unsigned null,
  province_id bigint unsigned not null,
  oks varchar(255) not null,
  code varchar(255) not null,
  lat decimal(9, 6) not null,
  lon decimal(9, 6) not null,
  address text not null,
  status enum("ACTIVE","DEACTIVE") not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  constraint branches_merchant_id_foreign
    foreign key (merchant_id) references merchants (id),
  constraint branches_city_id_foreign
    foreign key (city_id) references cities (id),
  unique key unique_oks (oks)
);

create table users (
  id bigint unsigned not null primary key,
  cellphone varchar(255) not null,
  kyc_status enum ('PENDING', 'APPROVED', 'DUPLICATED') default 'PENDING' not null,
  email varchar(255) null,
  postal_code varchar(10) null,
  address text null,
  name varchar(255) null,
  last_name varchar(255) null,
  status enum("ACTIVE","SUSPEND") not null,
  city_id bigint unsigned null,
  birth_date date null,
  branch_id bigint unsigned null,
  created_at timestamp not null,
  updated_at timestamp not null,
  constraint users_city_id_foreign
    foreign key (city_id) references cities (id),
  constraint users_branch_id_foreign
    foreign key (branch_id) references branches (id)
);

create table categories (
  id bigint unsigned not null primary key auto_increment,
  title varchar(255) not null,
  image_file_path varchar(255) null,
  parent_id bigint unsigned default null,
  level bigint unsigned not null default 0,
  number_of_products bigint unsigned not null default 0,
  status varchar(20) not null default 'ACTIVE',
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null default null,
  constraint fk_categories_parent foreign key (parent_id) references categories(id)
);

create table brands (
  id bigint unsigned not null primary key,
  title varchar(255) not null,
  created_at timestamp not null,
  updated_at timestamp not null
);

create table products (
  id bigint unsigned not null primary key auto_increment,
  merchant_id bigint unsigned not null,
  barcode varchar(255) not null,
  brand_id bigint unsigned,
  status enum('ACTIVE', 'INACTIVE', 'DELETED') not null default 'ACTIVE',
  title varchar(255) not null,
  description text not null,
  category_id bigint unsigned not null,
  price bigint unsigned not null,
  discounted_price bigint unsigned null,
  discount_percent decimal(5,2) unsigned null,
  discount_amount bigint unsigned null,
  image_file_path text null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null,
  constraint products_category_id_foreign
    foreign key (category_id) references categories (id),
  constraint fk_products_merchants foreign key (merchant_id) references merchants(id),
  constraint fk_products_brand_id foreign key (brand_id) references brands(id),
  constraint unique_merchant_id_barcode unique (merchant_id, barcode)
);

create table branch_products (
  id bigint unsigned not null primary key auto_increment,
  branch_id bigint unsigned not null,
  product_id bigint unsigned not null,
  price bigint unsigned not null,
  discount_percent decimal(5,2) unsigned null,
  discount_amount bigint unsigned null,
  discounted_price bigint unsigned null,
  status enum('ACTIVE', 'INACTIVE', 'DELETED') not null default 'ACTIVE',
  stock bigint unsigned null,
  created_at timestamp not null,
  updated_at timestamp not null,
  constraint branch_products_branch_id_foreign
    foreign key (branch_id) references branches (id),
  constraint branch_products_product_id_foreign
    foreign key (product_id) references products (id),
  unique key unique_branch_id_product_id (branch_id, product_id)
);

create table campaigns (
  id bigint unsigned not null primary key auto_increment,
  title varchar(255) not null,
  subtitle text not null,
  description text not null,
  image_file_path varchar(255) not null,
  sort int not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  start_date timestamp not null,
  end_date timestamp not null
);

create table deals (
  id bigint unsigned not null primary key auto_increment,
  title varchar(255) not null,
  description varchar(255) not null,
  product_id bigint unsigned not null,
  campaign_id bigint unsigned not null,
  discounted_price bigint unsigned not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  constraint deals_product_id_foreign
    foreign key (product_id) references products (id),
  constraint deals_campaign_id_foreign
    foreign key (campaign_id) references campaigns (id)
);

create table user_shopping_lists (
  id bigint unsigned not null primary key auto_increment,
  user_id bigint unsigned not null,
  branch_id bigint unsigned not null,
  title varchar(255) not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null,
  constraint user_shopping_lists_user_id_foreign
    foreign key (user_id) references users (id),
  constraint user_shopping_lists_branch_id_foreign
    foreign key (branch_id) references branches (id)
);

create table user_shopping_list_products (
  id bigint unsigned not null primary key auto_increment,
  user_shopping_list_id bigint unsigned not null,
  product_id bigint unsigned null,
  product_name varchar(255) null,
  count bigint unsigned not null default 1,
  created_at timestamp not null,
  deleted_at timestamp null,
  constraint user_shopping_list_products_user_shopping_list_id_foreign
    foreign key (user_shopping_list_id) references user_shopping_lists (id),
  constraint user_shopping_list_products_product_id_foreign
    foreign key (product_id) references products (id)
);

create table outbox_event_types
(
    id         bigint unsigned auto_increment
        primary key,
    code       varchar(150)                                  not null,
    type       enum ('ASYNQ', 'RABBITMQ') default 'RABBITMQ' not null,
    target     varchar(255)                                  null,
    created_at timestamp                                     not null,
    constraint outbox_event_types_code_unique
        unique (code)
);

create table outbox_events
(
    id           bigint unsigned auto_increment
        primary key,
    type_id      bigint unsigned                                 not null,
    body         text                                            not null,
    status       enum ('PENDING', 'PROCESSED') default 'PENDING' not null,
    created_at   timestamp                                       not null,
    processed_at timestamp                                       null,
    constraint outbox_events_type_id_foreign
        foreign key (type_id) references outbox_event_types (id)
);

create table user_cards
(
    id          bigint unsigned auto_increment primary key,
    user_id     bigint unsigned                          not null,
    card_number varchar(16)  not null,
    status      enum ('PENDING', 'APPROVED', 'REJECTED') not null,
    created_at  timestamp    not null,
    updated_at  timestamp    not null,
    deleted_at  timestamp    null,
    constraint user_cards_user_fk
        foreign key (user_id) references users (id)
);

create table media
(
    id         bigint unsigned auto_increment
        primary key,
    type       enum ('IMAGE', 'VIDEO') not null,
    file_path  text                    null,
    bucket     varchar(255)            not null,
    is_draft   tinyint(1) default 1    not null,
    created_at timestamp               not null,
    updated_at timestamp               not null,
    deleted_at timestamp               null
);

create table media_metadata
(
    id           bigint unsigned auto_increment
        primary key,
    user_id      bigint unsigned                                        not null,
    media_id     bigint unsigned                                        not null,
    object_id    bigint unsigned                                        null,
    object_type  enum ('AD') not null,
    file_size_kb bigint unsigned                                        null,
    mime_type    varchar(255)                                           null,
    created_at   timestamp                                              not null,
    updated_at   timestamp                                              not null,
    constraint media_metadata_pk
        unique (media_id)
);

-- Create ad_groups table
CREATE TABLE ad_groups (
  id bigint unsigned not null primary key auto_increment,
  ad_group_type enum('SINGLE', 'TRIPLE', 'QUADRUPLE') not null,
  title varchar(255) not null,
  created_by bigint unsigned null,
  is_active boolean not null,
  start_at timestamp null,
  end_at timestamp null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null,
  constraint ad_groups_created_by_foreign
    foreign key (created_by) references users(id)
);

-- Create ads table
create table ads (
  id bigint unsigned not null primary key auto_increment,
  ad_group_id bigint unsigned not null,
  redirect_url varchar(255) not null,
  media_id bigint unsigned not null,
  sort int unsigned not null,
  is_active boolean not null,
  aspect_ratio enum("1:1", "2:3", "5:1") not null,
  aspect_ratio_factor float not null,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null,
  constraint ads_ad_group_id_foreign
    foreign key (ad_group_id) references ad_groups (id),
  constraint ads_media_id_foreign
    foreign key (media_id) references media (id)
);

CREATE TABLE branch_ad_groups (
  id          bigint unsigned auto_increment primary key,
  branch_id   BIGINT UNSIGNED NOT NULL,
  ad_group_id BIGINT UNSIGNED NOT NULL,
  created_at  timestamp    not null,
  updated_at  timestamp    not null,
  deleted_at  timestamp    null,

  constraint fk_branch_ad_groups_branch
    foreign key (branch_id) references branches(id),
  constraint fk_branch_ad_groups_ad_group
    foreign key (ad_group_id) references ad_groups(id),
  constraint unique_branch_id_ad_group_id
    unique (branch_id, ad_group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

