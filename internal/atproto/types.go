package atproto

type Entry struct {
	Content   string `json:"content"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
	// RichText facets might be here if it's using standard RT, 
	// but WhiteWind usually stores Markdown in 'content' 
	// and facets in a separate field if applicable.
	// Actually, WhiteWind uses Markdown.
}

type BlogEntry struct {
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	Slug      string   `json:"slug"`
	CreatedAt string   `json:"createdAt"`
	Facets    []Facet  `json:"facets,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Embed     *Embed   `json:"embed,omitempty"`
}

type Embed struct {
	Type   string        `json:"$type"`
	Images []ImageEmbed  `json:"images,omitempty"`
	External *ExternalEmbed `json:"external,omitempty"`
}

type ImageEmbed struct {
	Image Blob   `json:"image"`
	Alt   string `json:"alt"`
}

type ExternalEmbed struct {
	Uri   string `json:"uri"`
	Title string `json:"title"`
	Description string `json:"description"`
	Thumb *Blob `json:"thumb,omitempty"`
}

type Blob struct {
	Ref  BlobRef `json:"ref"`
	Mime string  `json:"mimeType"`
	Size int     `json:"size"`
}

type BlobRef struct {
	Link string `json:"$link"`
}

type Facet struct {
	Index Features `json:"index"`
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
