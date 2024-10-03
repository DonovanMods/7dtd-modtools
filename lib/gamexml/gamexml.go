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
