--set timezone = 'Asia/Dushanbe';

create table users (
                       id         bigserial primary key,
                       name       text not null,
                       username   text not null unique,
                       password   text not null,
                       phone      text not null,
                       is_active  bool default true,
                       created_at timestamptz not null default current_timestamp,
                       updated_at timestamptz,
                       deleted_at timestamptz
);

create table tokens (
                        id         bigserial primary key,
                        user_id    int references users (id),
                        token      text not null,
                        created_at timestamptz not null default  current_timestamp + interval '5 hour',
                        expire     timestamptz not null default  current_timestamp + interval '5 hour' + interval '30 minute'
);

create table types (
                       id    bigserial primary key,
                       title text not null
);

create table pets (
                      id      bigserial primary key,
                      title   text not null
);

create table cities (
                        id    bigserial primary key,
                        title text not null
);

create table ads (
                     id            bigserial primary key,
                     user_id       int references users (id),
                     type_id       int references types (id),
                     pet_id        int references pets (id),
                     city_id       int references cities (id),
                     title         text not null,
                     description   text,
                     photo_path    text,
                     reward        int,
                     is_active     bool default true,
                     created_at    timestamptz not null default current_timestamp,
                     updated_at    timestamptz,
                     deleted_at    timestamptz
);

--SET timezone = 'Asia/Dushanbe';

--set timezone to 'Asia/Dushanbe';