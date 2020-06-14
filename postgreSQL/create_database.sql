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
create index if not exists id_users_nickname on users(nickname);
create index if not exists id_users_email on users(email);

/* ФОРУМ */
drop table if exists forums cascade;
drop table if exists forum_user cascade;
create table forums (
    slug    citext    not null primary key collate "C",
    title   text      not null,
    usr     citext    not null references users (nickname) collate "C",
    posts   integer   not null default 0,
    threads integer   not null default 0
);
create index if not exists id_forum_user on forums(usr);

create table forum_user (
    slug      citext not null references forums (slug) collate "C",
    nickname  citext not null references users (nickname) collate "C",
    PRIMARY KEY (slug, nickname)
);

/* ВЕТКА ОБСУЖДЕНИЯ */
drop table if exists posts cascade;
drop table if exists threads cascade;
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

/* ПОСТЫ */
create table posts (
    id         serial         not null primary key,
    parent     integer        not null default 0,
    author     citext         not null references users (nickname) collate "C",
    message    text,
    is_edited  bool           not null default false,
    forum      citext         not null references forums (slug) collate "C",
    thread     integer        not null references threads (id),
    created    timestamptz    default current_timestamp,
    path       integer[]      default array[]::integer[]
);
create index if not exists id_posts_path on posts(path);
create index if not exists id_posts_thread on posts(thread);

/* ГОЛОС ПОЛЬЗОВАТЕЛЯ */
drop table if exists votes cascade;
create table votes (
    nickname citext   not null references users (nickname) collate "C",
    voice    smallint check (voice in (-1, 1)),
    thread   integer  not null references threads (id),
    --unique (nickname, thread)
    PRIMARY KEY (nickname, thread)
);