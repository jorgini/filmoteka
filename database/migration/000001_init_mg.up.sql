CREATE TABLE users
(
    id        serial PRIMARY KEY,
    login     varchar(255) not null unique,
    password  varchar(255) not null,
    user_role varchar(255) not null,
    CHECK (user_role = 'regular' or user_role = 'admin')
);

CREATE TABLE actors
(
    id       serial PRIMARY KEY,
    name     varchar(255) not null,
    surname  varchar(255) not null,
    sex      varchar(255) not null,
    CHECK (sex = 'male' or sex = 'female'),
    birthday date         not null
);

ALTER TABLE actors
    ADD CONSTRAINT uqc_actor_name
        UNIQUE (name, surname);

CREATE TABLE films
(
    id          serial PRIMARY KEY,
    title       varchar(255) not null unique,
    CHECK (LENGTH(title) > 0 and LENGTH(title) <= 150),
    description text,
    CHECK (LENGTH(description) <= 1000),
    issue_date  date         not null,
    rating      integer      not null,
    CHECK (rating >= 0 and rating <= 10)
);

CREATE TABLE starred_in
(
    actor_id integer not null,
    film_id  integer not null,
    FOREIGN KEY (actor_id) REFERENCES actors (id) ON DELETE CASCADE,
    FOREIGN KEY (film_id) REFERENCES films (id) ON DELETE CASCADE,
    PRIMARY KEY (actor_id, film_id)
);