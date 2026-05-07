package profile

import "fmt"

var registry = []Profile{
	Berryer,
}

func All() []Profile {
	return registry
}

func ByID(id string) (Profile, error) {
	for _, p := range registry {
		if p.ID == id {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("profil inconnu: %q", id)
}
