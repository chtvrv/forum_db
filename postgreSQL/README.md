sudo -u postgres psql

create user forum_user with encrypted password 'forum1234';

grant all privileges on database forum_db to forum_user;

ALTER USER forum_user WITH SUPERUSER;

для psql:
\l - список баз
\с database - приконнектиться к базе
\dt - список имеющихся таблиц
\du - список ролей
\s - история команд