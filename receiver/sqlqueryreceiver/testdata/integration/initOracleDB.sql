create table movie
(
    title       VARCHAR2(256),
    genre       VARCHAR(256),
    imdb_rating NUMBER
);

insert into movie (title, genre, imdb_rating)
values ('E.T.', 'SciFi', 7.9);
insert into movie (title, genre, imdb_rating)
values ('Blade Runner', 'SciFi', 8.1);
insert into movie (title, genre, imdb_rating)
values ('Star Wars', 'SciFi', 8.6);
insert into movie (title, genre, imdb_rating)
values ('Die Hard', 'Action', 8.2);
insert into movie (title, genre, imdb_rating)
values ('Mission Impossible', 'Action', 7.1);

/* The alter session command is required to enable user creation in an Oracle docker container
   This command shouldn't be used outside of test environments. */
alter session set "_ORACLE_SCRIPT"=true;
CREATE USER OTEL IDENTIFIED BY password;
GRANT CREATE SESSION TO OTEL;
GRANT ALL ON movie TO OTEL;
