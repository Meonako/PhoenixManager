package party

type Party struct {
	Players   []string
	Target    TARGET
	MaxPlayer int
}

// Create a Party with specify target.
func NewParty(Target TARGET, Members ...string) (newParty *Party) {
	newParty = &Party{Players: Members, Target: Target, MaxPlayer: Target.GetMaxPlayer()}
	ActiveParties.Add(newParty)
	return
}

// Add a player to party
func (p *Party) Join(ID string) {
	p.Players = append(p.Players, ID)
}

// Remove a player from party
func (p *Party) Leave(ID string) {
	for i, v := range p.Players {
		if v == ID {
			p.Players = append(p.Players[:i], p.Players[i+1:]...)
		}
	}
}

// Return number of players in the party
func (p *Party) PlayersCount() int {
	return len(p.Players)
}

// Return true if party is empty. (i.e. numbers of players is 0)
func (p *Party) IsEmpty() bool {
	return p.Players == nil || p.PlayersCount() <= 0
}

// ----------------------------------------
// ------------- ACTIVE PARTY -------------
// ----------------------------------------

type ActiveParty struct {
	AllParty []*Party
}

var ActiveParties = ActiveParty{
	AllParty: []*Party{},
}

// Add a party to list of active parties
func (ap *ActiveParty) Add(pt *Party) {
	ap.AllParty = append(ap.AllParty, pt)
}

// Return number of active parties
func (ap *ActiveParty) Count() int {
	return len(ap.AllParty)
}

// Find party that has this ID
func (ap *ActiveParty) FindByPlayer(ID string) *Party {
	if ap.IsEmpty() {
		pt := Party{}
		return &pt
	}

	for _, party := range ap.AllParty {
		for _, player := range party.Players {
			if player == ID {
				return party
			}
		}
	}

	pt := Party{}
	return &pt
}

// Find party that has the specify Target
func (ap *ActiveParty) FindByTarget(Target TARGET) *Party {
	if ap.IsEmpty() {
		pt := Party{}
		return &pt
	}

	for _, party := range ap.AllParty {
		if party.Target == Target {
			return party
		}
	}

	pt := Party{}
	return &pt
}

func (ap *ActiveParty) IsEmpty() bool {
	return ap.Count() <= 0
}
