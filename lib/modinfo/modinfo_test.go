package modinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	modletName := "TestModlet"
	modInfo := New(modletName)

	assert.Equal(t, modletName, modInfo.Name.Value, "Expected Name.Value to be %s", modletName)
	assert.Equal(t, "My New Modlet", modInfo.DisplayName.Value, "Expected DisplayName.Value to be 'My New Modlet'")
	expectedDescription := "This is the description for TestModlet -- please change it"
	assert.Equal(t, expectedDescription, modInfo.Description.Value, "Expected Description.Value to be '%s'", expectedDescription)
	assert.Equal(t, "Unknown", modInfo.Author.Value, "Expected Author.Value to be 'Unknown'")
	assert.Equal(t, "0.1.0", modInfo.Version.Value, "Expected Version.Value to be '0.1.0'")
	assert.Equal(t, "V1", modInfo.Version.Compat, "Expected Version.Compat to be 'V1'")
	assert.Equal(t, "https://example.org", modInfo.Website.Value, "Expected Website.Value to be 'https://example.org'")
}

func TestXML(t *testing.T) {
	modletName := "TestModlet"
	modInfo := New(modletName)

	xmlOutput, err := modInfo.XML()
	assert.NoError(t, err, "Expected no error while generating XML")

	expectedXML := `<modinfo>
  <Name value="TestModlet" />
  <DisplayName value="My New Modlet" />
  <Description value="This is the description for TestModlet -- please change it" />
  <Author value="Unknown" />
  <Version value="0.1.0" compat="V1" />
  <Website value="https://example.org" />
</modinfo>`
	assert.Equal(t, expectedXML, xmlOutput, "Expected XML output to match")
}
