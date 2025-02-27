INSERT INTO players (login, password, display_name, avatar_name)
VALUES ('player0', 'E51FA927933824E9A2F265270E9C973DC83DBBC996BE39F043817FFADA63DC598C37A35ACB006850D93D48C28298B609F31C3031EAFDE476A1005BD3F7FFDAD8', 'Test player 0', 'DBH_NPC_CRDL_02_022_avatar_png'),
       ('player1', 'E51FA927933824E9A2F265270E9C973DC83DBBC996BE39F043817FFADA63DC598C37A35ACB006850D93D48C28298B609F31C3031EAFDE476A1005BD3F7FFDAD8', 'Test player 1', 'crdl_04_119_avatar_png');

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
