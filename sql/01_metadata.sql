INSERT INTO attributes (id, color, name)
VALUES (1, 'yellow', 'willpower'),
       (2, 'red', 'strength'),
       (3, 'blue', 'intelligence'),
       (4, 'green', 'agility'),
       (5, 'purple', 'endurance'),
       (6, 'gray', 'neutral');

INSERT INTO classes (id, name)
VALUES (1, 'willpower'),
       (2, 'strength'),
       (3, 'intelligence'),
       (4, 'agility'),
       (5, 'endurance'),
       (6, 'neutral'),
       (7, 'archer'),
       (8, 'assassin'),
       (9, 'battlemage'),
       (10, 'crusader'),
       (11, 'mage'),
       (12, 'monk'),
       (13, 'scout'),
       (14, 'sorcerer'),
       (15, 'spellsword'),
       (16, 'warrior'),
       (17, 'house redoran'),
       (18, 'house telvanni'),
       (19, 'house hlaalu'),
       (20, 'tribunal temple'),
       (21, 'house dagoth'),
       (22, 'ebonheart pact'),
       (23, 'daggerfall covenant'),
       (24, 'aldmeri dominion'),
       (25, 'the guildsworn'),
       (26, 'the empire of cyrodiil');

INSERT INTO classes_to_attributes (class_id, attribute_id)
VALUES (1, 1),
       (2, 2),
       (3, 3),
       (4, 4),
       (5, 5),
       (7, 2),
       (7, 4),
       (8, 3),
       (8, 4),
       (9, 2),
       (9, 3),
       (10, 1),
       (10, 2),
       (11, 1),
       (11, 3),
       (12, 2),
       (12, 4),
       (13, 4),
       (13, 5),
       (14, 3),
       (14, 5),
       (15, 1),
       (15, 5),
       (16, 2),
       (16, 5),
       (17, 1),
       (17, 2),
       (17, 5),
       (18, 3),
       (18, 4),
       (18, 5),
       (19, 1),
       (19, 2),
       (19, 4),
       (20, 1),
       (20, 3),
       (20, 5),
       (21, 2),
       (21, 3),
       (21, 4),
       (22, 2),
       (22, 4),
       (22, 5),
       (23, 2),
       (23, 3),
       (23, 5),
       (24, 1),
       (24, 3),
       (24, 4),
       (25, 1),
       (25, 2),
       (25, 3),
       (26, 1),
       (26, 4),
       (26, 5);

INSERT INTO races (id, name)
VALUES (1, 'argonian'),
       (2, 'breton'),
       (3, 'dark elf'),
       (4, 'high elf'),
       (5, 'imperial'),
       (6, 'khajiit'),
       (7, 'nord'),
       (8, 'orc'),
       (9, 'redguard'),
       (10, 'wood elf');

INSERT INTO card_types (id, name)
VALUES (1, 'action'),
       (2, 'creature'),
       (3, 'item'),
       (4, 'support');

INSERT INTO keywords (id, name)
VALUES (1, 'breakthrough'),
       (2, 'charge'),
       (3, 'drain'),
       (4, 'guard'),
       (5, 'lethal'),
       (6, 'mobilize'),
       (7, 'rally'),
       (8, 'regenerate'),
       (9, 'ward'),
       (10, 'prophecy');