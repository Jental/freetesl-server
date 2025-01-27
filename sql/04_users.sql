INSERT INTO players (login, password, display_name)
VALUES ('player0', null, 'Test player 0'),
       ('player1', null, 'Test player 1');

INSERT INTO decks (name, player_id)
VALUES ('Test deck 0', 1),
       ('Test deck 1', 2);

INSERT INTO deck_cards (deck_id, card_id, count)
VALUES (1, 1, 7),
       (1, 2, 8),
       (2, 1, 8),
       (2, 2, 7);
