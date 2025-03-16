package fly

type Aircrafts []*Aircraft

func (a *Aircrafts) Add(aircraft *Aircraft) {
	if a == nil {
		*a = []*Aircraft{}
	}
	if aircraft == nil {
		return
	}
	*a = append(*a, aircraft)
}

func (a *Aircrafts) Remove(other *Aircraft) {
	if a == nil || other == nil {
		return
	}
	var aircrafts []*Aircraft
	for _, aircraft := range *a {
		if aircraft.ID != other.ID {
			aircrafts = append(aircrafts, aircraft)
		}
	}
	*a = aircrafts
}

func (a *Aircrafts) Contains(other *Aircraft) bool {
	if a == nil || other == nil {
		return false
	}
	for _, aircraft := range *a {
		if aircraft.ID == other.ID {
			return true
		}
	}
	return false
}

func (a *Aircrafts) IsEmpty() bool {
	return a == nil || len(*a) == 0
}
