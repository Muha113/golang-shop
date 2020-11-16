package parser

import strip "github.com/grokify/html-strip-tags-go"

func RemoveTagsHTML(value string) string {
	return strip.StripTags(value)
}
