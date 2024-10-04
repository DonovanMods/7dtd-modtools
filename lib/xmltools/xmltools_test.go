package xmltools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveClosingXMLTags(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<tag></tag>",
			expected: "<tag />",
		},
		{
			input:    "<tag1>testing</tag1>",
			expected: "<tag1>testing</tag1>",
		},
		{
			input:    "<tag1><tag2></tag2></tag1>",
			expected: "<tag1><tag2 /></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3></tag1>",
			expected: "<tag1><tag2 /><tag3 /></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2 /><tag3 /><tag4 /></tag1>",
		},
		{
			input:    "<tag1 attribute=\"value\"></tag1>",
			expected: "<tag1 attribute=\"value\" />",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2></tag1>",
			expected: "<tag1><tag2 attribute=\"value\" /></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3 attribute=\"value\"></tag3></tag1>",
			expected: "<tag1><tag2 /><tag3 attribute=\"value\" /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3></tag3></tag1>",
			expected: "<tag1><tag2 attribute=\"value\" /><tag3 /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3 attribute=\"value\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute=\"value\" /><tag3 attribute=\"value\" /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"></tag2></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\" /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\"></tag2><tag3 attribute2=\"value2\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" /><tag3 attribute2=\"value2\" /></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3><tag4></tag4><tag5></tag5></tag1>",
			expected: "<tag1><tag2 /><tag3 /><tag4 /><tag5 /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2 attribute=\"value\" /><tag3 /><tag4 /></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3 attribute=\"value\"></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2 /><tag3 attribute=\"value\" /><tag4 /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"></tag2><tag3 attribute3=\"value3\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\" /><tag3 attribute3=\"value3\" /></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\"></tag2><tag3 attribute2=\"value2\"></tag3><tag4 attribute3=\"value3\"></tag4></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" /><tag3 attribute2=\"value2\" /><tag4 attribute3=\"value3\" /></tag1>",
		},
	}

	for _, test := range tests {
		result := RemoveClosingXMLTags(test.input)
		assert.Equal(t, test.expected, result, "RemoveClosingXMLTags(%q)", test.input)
	}
}

func TestRemoveClosingXMLTagsCompact(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "<tag></tag>",
			expected: "<tag/>",
		},
		{
			input:    "<tag1>testing</tag1>",
			expected: "<tag1>testing</tag1>",
		},
		{
			input:    "<tag1><tag2></tag2></tag1>",
			expected: "<tag1><tag2/></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3></tag1>",
			expected: "<tag1><tag2/><tag3/></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2/><tag3/><tag4/></tag1>",
		},
		{
			input:    "<tag1 attribute=\"value\"></tag1>",
			expected: "<tag1 attribute=\"value\"/>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2></tag1>",
			expected: "<tag1><tag2 attribute=\"value\"/></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3 attribute=\"value\"></tag3></tag1>",
			expected: "<tag1><tag2/><tag3 attribute=\"value\"/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3></tag3></tag1>",
			expected: "<tag1><tag2 attribute=\"value\"/><tag3/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3 attribute=\"value\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute=\"value\"/><tag3 attribute=\"value\"/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"></tag2></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\"></tag2><tag3 attribute2=\"value2\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\"/><tag3 attribute2=\"value2\"/></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3></tag3><tag4></tag4><tag5></tag5></tag1>",
			expected: "<tag1><tag2/><tag3/><tag4/><tag5/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute=\"value\"></tag2><tag3></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2 attribute=\"value\"/><tag3/><tag4/></tag1>",
		},
		{
			input:    "<tag1><tag2></tag2><tag3 attribute=\"value\"></tag3><tag4></tag4></tag1>",
			expected: "<tag1><tag2/><tag3 attribute=\"value\"/><tag4/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"></tag2><tag3 attribute3=\"value3\"></tag3></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\" attribute2=\"value2\"/><tag3 attribute3=\"value3\"/></tag1>",
		},
		{
			input:    "<tag1><tag2 attribute1=\"value1\"></tag2><tag3 attribute2=\"value2\"></tag3><tag4 attribute3=\"value3\"></tag4></tag1>",
			expected: "<tag1><tag2 attribute1=\"value1\"/><tag3 attribute2=\"value2\"/><tag4 attribute3=\"value3\"/></tag1>",
		},
	}

	for _, test := range tests {
		result := RemoveClosingXMLTagsCompact(test.input)
		assert.Equal(t, test.expected, result, "RemoveClosingXMLTagsCompact(%q)", test.input)
	}
}
