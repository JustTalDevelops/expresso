package text

import "encoding/json"

// Text represents the custom JSON text format in Minecraft.
// TODO: Add options to make creating this type of text simple.
type Text struct {
	Text string `json:"text,omitempty"`

	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	UnderLined    bool   `json:"underlined,omitempty"`
	StrikeThrough bool   `json:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty"`
	Color         string `json:"color,omitempty"`

	Translate string            `json:"translate,omitempty"`
	With      []json.RawMessage `json:"with,omitempty"`
	Extra     []Text            `json:"extra,omitempty"`
}
