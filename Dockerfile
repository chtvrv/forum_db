FROM ubuntu:18.04

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get -y update
RUN apt install -y git wget gcc gnupg

ENV PGVER 11

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN wget https://www.postgresql.org/media/keys/ACCC4CF8.asc
RUN apt-key add ACCC4CF8.asc

RUN apt-get update

RUN apt-get install -y  postgresql-$PGVER

ENV GOVER 1.14

RUN wget https://dl.google.com/go/go$GOVER.linux-amd64.tar.gz
RUN tar -xzf go$GOVER.linux-amd64.tar.gz
RUN cp -r go/ $HOME/

ENV GOROOT $HOME/go
ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

WORKDIR /server
COPY . .

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER forum_user WITH SUPERUSER PASSWORD 'forum1234';" &&\
    createdb -O forum_user forum_db &&\
    psql forum_db -f /server/postgreSQL/create_database.sql &&\
    /etc/init.d/postgresql stop

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root
#RUN go mod vendor
RUN go build /server/cmd/main.go
CMD service postgresql start && ./main -c /server/cmd/config.yaml

EXPOSE 5000