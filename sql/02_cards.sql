INSERT INTO cards (id, name, description, power, health, cost, class_id, type_id)
VALUES (1, 'Bruma Profiteer', 'When you summon another creature, you gain 1 health.', 3, 2, 2, 1, 2),
       (2, 'Thieves Guild Recruit', 'Summon: Draw a card. If it costs 7 or more, reduce it''s cost by 1.', 1, 2, 2, 4, 2),
       (3, 'Mournhold Guardian', 'Guard', 2, 1, 1, 4, 2);

INSERT INTO card_races (card_id, race_id)
VALUES (1, 5),
       (2, 1),
       (3, 3);

INSERT INTO card_keywords (card_id, keyword_id)
VALUES (3, 4);        
