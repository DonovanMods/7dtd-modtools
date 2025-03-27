/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package gamexml

import "errors"

var ErrNoId = errors.New("GameXML must have an id")

// GameXMLs is a map of GameXMLs keyed to the filenname of the game's XML
// e.g. 'blocks'
type GameXMLs map[string]GameXML

// New creates an empty GameXMLs map and returns a pointer to it
func NewGameXMLs() *GameXMLs {
	gameXMLs := make(GameXMLs, 100)

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
	*G = *(NewGameXMLs())
}
