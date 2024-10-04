package gamexml

import "errors"

var ErrNoId = errors.New("GameXML must have an id")

// GameXMLs is a map of GameXMLs keyed to the filenname of the game's XML
// e.g. 'blocks'
type GameXMLs map[string]GameXML

// New creates an empty GameXMLs map and returns a pointer to it
func New() *GameXMLs {
	var gameXMLs = make(GameXMLs, 100)

	return &gameXMLs
}

// Add adds a GameXML to the GameXMLs map
// Note: this will overwrite any existing GameXML entry with the same id
func (G *GameXMLs) Add(g GameXML) error {
	if g.Id == "" {
		return ErrNoId
	}

	(*G)[g.Id] = g

	return nil
}

// Get returns a GameXML from the GameXMLs map
func (G *GameXMLs) Get(name string) (GameXML, bool) {
	if name == "" {
		return GameXML{}, false
	}

	gameXML, ok := (*G)[name]

	return gameXML, ok
}

// Reset resets the GameXMLs map to an empty map
func (G *GameXMLs) Reset() {
	*G = *(New())
}
