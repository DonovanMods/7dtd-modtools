package xmltools

import "regexp"

func RemoveClosingXMLTags(xml string) string {
	return removeClosingXMLTags(xml, false)
}

func RemoveClosingXMLTagsCompact(xml string) string {
	return removeClosingXMLTags(xml, true)
}

func removeClosingXMLTags(xml string, compact bool) string {
	var closingTag = "/>"

	if !compact {
		closingTag = " />"
	}

	re := regexp.MustCompile("></[[:alnum:]]*>")

	return re.ReplaceAllString(string(xml), closingTag)
}
