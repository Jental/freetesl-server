INSERT INTO cards (id, name, description, power, health, cost, class_id, type_id)
VALUES (1, 'Bruma Profiteer', 'When you summon another creature, you gain 1 health.', 3, 2, 2, 1, 2),
       (2, 'Thieves Guild Recruit', 'Summon: Draw a card. If it costs 7 or more, reduce it''s cost by 1.', 1, 2, 2, 4, 2),
       (3, 'Mournhold Guardian', 'Guard', 2, 1, 1, 4, 2),
       (4, 'Arrow in the Knee', 'Shackle a creature and deal 1 damage to it.', 0, 0, 1, 4, 1),
       (5, 'Knight of the Hour', 'Prophecy, Guard\nSummon: You gain 3 health.', 3, 3, 5, 1, 2),
       (6, 'Improvised Weapon', 'Breakthrough\n+1/+1', 0, 0, 0, 2, 3);

INSERT INTO card_races (card_id, race_id)
VALUES (1, 5),
       (2, 1),
       (3, 3),
       (5, 5);

INSERT INTO card_keywords (card_id, keyword_id)
VALUES (3, 4),
       (5, 4),
       (5, 10),
       (6, 1);

INSERT INTO card_actions (card_id, action_id, interceptor_point_id, actions_parameters_values)
VALUES (4, 'deal_damage_to_creature', 'operations.cardPlay', '1'),
       (4, 'shackle', 'operations.cardPlay', NULL),
       (2, 'draw_cards', 'operations.cardPlay', '1'),
       (5, 'heal', 'operations.cardPlay', '3');

INSERT INTO card_effects (card_id, effect_id, name, parameter0, parameter1)
VALUES (6, 3, 'ModifiedPowerHealth: +1/+1', '1', '1');
