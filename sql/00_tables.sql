CREATE DATABASE tesl;

CREATE TABLE attributes (
    id int NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    color varchar(7) NOT NULL
);

CREATE TABLE classes (
    id int NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE classes_to_attributes (
    class_id int NOT NULL REFERENCES classes(id),
    attribute_id int NOT NULL REFERENCES attributes(id),
    PRIMARY KEY (class_id, attribute_id)
);

CREATE TABLE races (
    id int NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE cards (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    description varchar(1023) NOT NULL,
    power int NOT NULL,
    defence int NOT NULL,
    cost int NOT NULL,
    class_id int NOT NULL REFERENCES classes(id)
);

CREATE TABLE card_races (
    card_id int NOT NULL REFERENCES cards(id),
    race_id int NOT NULL REFERENCES races(id),
    PRIMARY KEY (card_id, race_id)
);

CREATE TABLE players (
    id serial PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    password varchar(255),
    display_name varchar(255),
    avatar_name varchar(255)
);

CREATE TABLE decks (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    player_id int NOT NULL REFERENCES players(id)
);

CREATE TABLE deck_cards (
    deck_id int NOT NULL REFERENCES decks(id),
    card_id int NOT NULL REFERENCES cards(id),
    count int NOT NULL
);
