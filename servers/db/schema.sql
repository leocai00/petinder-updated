create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(254) not null unique,
    pass_hash binary(60) not null,
    user_name varchar(255) not null unique,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    photo_url blob not null
);

create table if not exists sign_ins (
    id int not null auto_increment primary key, 
    user_id int not null,
    date_time datetime not null,
    ip_address varchar(39) not null
);