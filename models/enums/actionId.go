package enums

type ActionID string

const (
	ActionIDDealDamageToCreature ActionID = "deal_damage_to_creature"
	ActionIDDrawCards            ActionID = "draw_cards"
	ActionShackle                ActionID = "shackle"
	ActionHeal                   ActionID = "heal"
)
