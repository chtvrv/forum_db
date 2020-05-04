/* нечувствительный к регистру text */
create extension if not exists citext;

/* ПОЛЬЗОВАТЕЛЬ */
drop table if exists users cascade;
create table users (
    nickname citext    primary key not null collate "C",
    fullname citext    not null,
    email    citext    not null unique collate "C",
    about    text
);

/* ФОРУМ */
drop table if exists forums cascade;
drop table if exists forum_user cascade;
create table forums (
    slug    citext    not null primary key collate "C",
    title   text      not null,
    usr     citext    not null references users (nickname) collate "C"
);
create table forum_user (
    slug      citext not null references forums (slug) collate "C",
    nickname  citext not null references users (nickname) collate "C",
    PRIMARY KEY (slug, nickname)
);

/* ВЕТКА ОБСУЖДЕНИЯ */
drop table if exists threads;
create table threads (
    id         serial      not null primary key,
    title      text        not null,
    author     citext      not null references users (nickname) collate "C",
    forum      citext      not null references forums (slug) collate "C",
    message    text,
    votes      integer     default 0,
    slug       citext      default null unique collate "C",
    created    timestamptz default current_timestamp
);