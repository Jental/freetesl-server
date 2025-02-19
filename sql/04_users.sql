INSERT INTO players (login, password, display_name, avatar_name)
VALUES ('player0', 'CF83E1357EEFB8BDF1542850D66D8007D620E4050B5715DC83F4A921D36CE9CE47D0D13C5D85F2B0FF8318D2877EEC2F63B931BD47417A81A538327AF927DA3E', 'Test player 0', 'DBH_NPC_CRDL_02_022_avatar_png'),
       ('player1', 'CF83E1357EEFB8BDF1542850D66D8007D620E4050B5715DC83F4A921D36CE9CE47D0D13C5D85F2B0FF8318D2877EEC2F63B931BD47417A81A538327AF927DA3E', 'Test player 1', 'crdl_04_119_avatar_png');

INSERT INTO decks (name, player_id)
VALUES ('Test deck 0', 1),
       ('Test deck 1', 2);

INSERT INTO deck_cards (deck_id, card_id, count)
VALUES (1, 1, 7),
       (1, 2, 8),
       (1, 3, 9),
       (2, 1, 8),
       (2, 2, 7),
       (2, 3, 9);
