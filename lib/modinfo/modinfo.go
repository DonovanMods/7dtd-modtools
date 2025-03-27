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
package modinfo

import (
	XML "encoding/xml"
	"fmt"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/xmltools"
)

const XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`

type xmlValue struct {
	Value  string `xml:"value,attr"`
	Compat string `xml:"compat,attr,omitempty"`
}

type modInfoMeta struct {
	path string
}

type ModInfo struct {
	XMLName     XML.Name `xml:"xml"`
	Name        xmlValue `xml:"Name"`
	DisplayName xmlValue `xml:"DisplayName"`
	Description xmlValue `xml:"Description"`
	Author      xmlValue `xml:"Author"`
	Version     xmlValue `xml:"Version"`
	Website     xmlValue `xml:"Website"`
	meta        modInfoMeta
}

func (M *ModInfo) SetPath(path string) {
	(*M).meta.path = path
}

func (M *ModInfo) Path() string {
	return M.meta.path
}

func (M *ModInfo) Dir() string {
	return filepath.Dir(M.meta.path)
}

func (M *ModInfo) Filename() string {
	return filepath.Base(M.meta.path)
}

func (M *ModInfo) XML() (string, error) {
	xml, err := XML.MarshalIndent(M, "", "  ")
	if err != nil {
		return "", err
	}

	xml = []byte(xmltools.RemoveClosingXMLTags(string(xml)))

	return fmt.Sprintf("%s\n%s", XMLHeader, xml), nil
}

func NewModInfo(modletName string) *ModInfo {
	return &ModInfo{
		Name:        xmlValue{Value: modletName},
		DisplayName: xmlValue{Value: "My New Modlet"},
		Description: xmlValue{Value: fmt.Sprintf("This is the description for %s -- please change it", modletName)},
		Author:      xmlValue{Value: "Unknown"},
		Version:     xmlValue{Value: "0.1.0", Compat: "V1"},
		Website:     xmlValue{Value: "https://example.org"},
	}
}
