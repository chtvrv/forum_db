drop table if exists users cascade;

/* нечувствительный к регистру text */
create extension if not exists citext;

create table users (
    nickname citext    primary key not null collate "C",
    fullname citext    not null,
    email    citext    not null unique collate "C",
    about    text
);