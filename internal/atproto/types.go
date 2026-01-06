package atproto

import "encoding/json"

// Generic wrapper to handle different record types
type RecordValue struct {
	Type string `json:"$type"`
	// We will unmarshal into specific types based on Type
}

// Old WhiteWind Format
type BlogEntry struct {
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	CreatedAt string   `json:"createdAt"`
	Facets    []Facet  `json:"facets,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Embed     *Embed   `json:"embed,omitempty"`
}

// New Leaflet Format
type LeafletDocument struct {
	Type        string   `json:"$type"`
	Title       string   `json:"title"`
	PublishedAt string   `json:"publishedAt"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Publication string   `json:"publication"` // URI of the publication
	Pages       []Page   `json:"pages"`
}

type LeafletPublication struct {
	Type string `json:"$type"`
	Name string `json:"name"`
}

type Page struct {
	Type   string         `json:"$type"` // pub.leaflet.pages.linearDocument
	Blocks []BlockWrapper `json:"blocks"`
}

type BlockWrapper struct {
	Type  string          `json:"$type"` // pub.leaflet.pages.linearDocument#block
	Block json.RawMessage `json:"block"` // Deferred unmarshaling
}

// We'll use a helper to determine block type and unmarshal accordingly
type BaseBlock struct {
	Type string `json:"$type"`
}

type TextBlock struct {
	Type      string  `json:"$type"`
	Plaintext string  `json:"plaintext"`
	Facets    []Facet `json:"facets"`
}

type CodeBlock struct {
	Type      string `json:"$type"`
	Language  string `json:"language"`
	Plaintext string `json:"plaintext"`
}

type UnorderedListBlock struct {
	Type     string     `json:"$type"`
	Children []ListItem `json:"children"`
}

type ListItem struct {
	Type     string          `json:"$type"`    // pub.leaflet.blocks.unorderedList#listItem
	Content  json.RawMessage `json:"content"`  // Usually a TextBlock
	Children []ListItem      `json:"children"` // Nested lists?
}

type ImageBlock struct {
	Type  string `json:"$type"`
	Image Blob   `json:"image"`
	Alt   string `json:"alt"`
}

type BskyPostBlock struct {
	Type    string  `json:"$type"`
	PostRef PostRef `json:"postRef"`
}

type PostRef struct {
	Uri string `json:"uri"`
	Cid string `json:"cid"`
}

// Shared Types

type Facet struct {
	Index    Features  `json:"index"`
	Features []Feature `json:"features"`
}

type Features struct {
	ByteStart int `json:"byteStart"`
	ByteEnd   int `json:"byteEnd"`
}

type Feature struct {
	Type string `json:"$type"`
	URI  string `json:"uri,omitempty"`
	Did  string `json:"did,omitempty"`
}

type Embed struct {
	Type     string         `json:"$type"`
	Images   []ImageEmbed   `json:"images,omitempty"`
	External *ExternalEmbed `json:"external,omitempty"`
}

type ImageEmbed struct {
	Image Blob   `json:"image"`
	Alt   string `json:"alt"`
}

type ExternalEmbed struct {
	Uri         string `json:"uri"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumb       *Blob  `json:"thumb,omitempty"`
}

type Blob struct {
	Ref  BlobRef `json:"ref"`
	Mime string  `json:"mimeType"`
	Size int     `json:"size"`
}

type BlobRef struct {
	Link string `json:"$link"`
}
