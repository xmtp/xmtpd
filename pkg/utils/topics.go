package utils

import "strings"

var topicCategoryByPrefix = map[string]string{
	"test":         "test",
	"contact":      "contact",
	"intro":        "v1-intro",
	"dm":           "v1-conversation",
	"invite":       "v2-invite",
	"m":            "v2-conversation",
	"privatestore": "private",
}

func CategoryFromTopic(contentTopic string) string {
	if strings.HasPrefix(contentTopic, "test-") {
		return "test"
	}
	topic := strings.TrimPrefix(contentTopic, "/xmtp/0/")
	if len(topic) == len(contentTopic) {
		return "invalid"
	}
	prefix, _, hasPrefix := strings.Cut(topic, "-")
	if hasPrefix {
		if category, found := topicCategoryByPrefix[prefix]; found {
			return category
		}
	}
	return "invalid"
}
