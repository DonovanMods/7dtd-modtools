package modinfo

import (
	XML "encoding/xml"
	"fmt"

	"github.com/donovanmods/7dmt/lib/xmltools"
)

type xmlValue struct {
	Value  string `xml:"value,attr"`
	Compat string `xml:"compat,attr,omitempty"`
}

type ModInfo struct {
	XMLName     XML.Name `xml:"modinfo"`
	Name        xmlValue `xml:"Name"`
	DisplayName xmlValue `xml:"DisplayName"`
	Description xmlValue `xml:"Description"`
	Author      xmlValue `xml:"Author"`
	Version     xmlValue `xml:"Version"`
	Website     xmlValue `xml:"Website"`
}

func New(modletName string) ModInfo {
	return ModInfo{
		Name:        xmlValue{Value: modletName},
		DisplayName: xmlValue{Value: "My New Modlet"},
		Description: xmlValue{Value: fmt.Sprintf("This is the description for %s -- please change it", modletName)},
		Author:      xmlValue{Value: "Unknown"},
		Version:     xmlValue{Value: "0.1.0", Compat: "V1"},
		Website:     xmlValue{Value: "https://example.org"},
	}
}

func (M ModInfo) XML() (string, error) {
	xml, err := XML.MarshalIndent(M, "", "  ")
	if err != nil {
		return "", err
	}

	xml = []byte(xmltools.RemoveClosingXMLTags(string(xml)))

	return string(xml), nil
}
