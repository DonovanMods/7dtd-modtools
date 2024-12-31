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

/*
 * GameXML Struct
 */

// GameXML is a struct that represents the data associated with a game XML file
type GameXML struct {
	Id       string   // the filename of the game's XML
	Modified bool     // whether the game's XML has been modified
	Elements []string // the elements in the game's XML that were modified
}

func (g *GameXML) AddElement(e string) {
	g.Elements = append(g.Elements, e)
}
