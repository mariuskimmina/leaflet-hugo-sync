package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	"mariuskimmina.com/leaflet-hugo-sync/internal/atproto"
)

type Converter struct {
	// No state needed; conversion is stateless
}

type ConversionResult struct {
	Markdown string
	Images   []ImageRef
}

type ImageRef struct {
	Blob atproto.Blob
	Alt  string
}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertLeaflet(doc *atproto.LeafletDocument) (*ConversionResult, error) {
	var sb strings.Builder
	var images []ImageRef

	for _, page := range doc.Pages {
		for _, blockWrapper := range page.Blocks {
			// Unmarshal base block to check type
			var base atproto.BaseBlock
			if err := json.Unmarshal(blockWrapper.Block, &base); err != nil {
				continue
			}

			switch base.Type {
			case "pub.leaflet.blocks.text":
				var textBlock atproto.TextBlock
				if err := json.Unmarshal(blockWrapper.Block, &textBlock); err != nil {
					continue
				}
				sb.WriteString(c.renderText(&textBlock) + "\n\n")

			case "pub.leaflet.blocks.code":
				var codeBlock atproto.CodeBlock
				if err := json.Unmarshal(blockWrapper.Block, &codeBlock); err != nil {
					continue
				}
				// Ensure language is not nil or empty
				lang := codeBlock.Language
				if lang == "" {
					lang = "text"
				}
				sb.WriteString(fmt.Sprintf("\n```%s\n%s\n```\n\n", lang, codeBlock.Plaintext))

			case "pub.leaflet.blocks.unorderedList":
				var listBlock atproto.UnorderedListBlock
				if err := json.Unmarshal(blockWrapper.Block, &listBlock); err != nil {
					continue
				}
				c.renderList(&sb, listBlock.Children, 0)
				sb.WriteString("\n")

			case "pub.leaflet.blocks.image":
				var imgBlock atproto.ImageBlock
				if err := json.Unmarshal(blockWrapper.Block, &imgBlock); err != nil {
					continue
				}
				// Use blob CID as placeholder URL; main.go replaces it with the local path
				sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", imgBlock.Alt, imgBlock.Image.Ref.Link))
				images = append(images, ImageRef{Blob: imgBlock.Image, Alt: imgBlock.Alt})

			case "pub.leaflet.blocks.bskyPost":
				var postBlock atproto.BskyPostBlock
				if err := json.Unmarshal(blockWrapper.Block, &postBlock); err != nil {
					continue
				}
				// Render as a blockquote link to the Bluesky post
				// Parse AT-URI: at://did:plc:abc123/app.bsky.feed.post/postID
			did, postID := parseATUri(postBlock.PostRef.Uri)
			postURL := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", did, postID)
				sb.WriteString(fmt.Sprintf("> [View on Bluesky](%s)\n\n", postURL))
			}
		}
	}

	return &ConversionResult{
		Markdown: sb.String(),
		Images:   images,
	}, nil
}

func (c *Converter) renderText(block *atproto.TextBlock) string {
	// Apply facets to convert rich text to markdown
	// Note: ATProto facets use byte offsets, not rune offsets
	data := []byte(block.Plaintext)

	var sb strings.Builder
	lastPos := 0

	for _, facet := range block.Facets {
		start := facet.Index.ByteStart
		end := facet.Index.ByteEnd

		if start < lastPos {
			continue // Overlap or out of order
		}

		// Append text before facet
		sb.Write(data[lastPos:start])

		// Handle feature
		text := string(data[start:end])
		replacement := text

		for _, feat := range facet.Features {
			switch feat.Type {
			case "pub.leaflet.richtext.facet#link":
				replacement = fmt.Sprintf("[%s](%s)", text, feat.URI)
			case "pub.leaflet.richtext.facet#didMention":
				replacement = fmt.Sprintf("[%s](https://bsky.app/profile/%s)", text, feat.Did)
			case "pub.leaflet.richtext.facet#code":
				replacement = fmt.Sprintf("`%s`", text)
			}
		}

		sb.WriteString(replacement)
		lastPos = end
	}

	sb.Write(data[lastPos:])

	return sb.String()
}

func (c *Converter) renderList(sb *strings.Builder, items []atproto.ListItem, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, item := range items {
		// Unmarshal content (TextBlock)
		var textBlock atproto.TextBlock
		if err := json.Unmarshal(item.Content, &textBlock); err == nil {
			sb.WriteString(fmt.Sprintf("%s- %s\n", indent, c.renderText(&textBlock)))
		}
		if len(item.Children) > 0 {
			c.renderList(sb, item.Children, depth+1)
		}
	}
}

func lastPathPart(uri string) string {
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

// parseATUri extracts DID and record key from an AT-URI
// Example: at://did:plc:abc123/app.bsky.feed.post/3mbrxzvw36c22
// Returns: (did:plc:abc123, 3mbrxzvw36c22)
func parseATUri(uri string) (did string, recordKey string) {
	// Remove "at://" prefix
	uri = strings.TrimPrefix(uri, "at://")

	// Split into parts
	parts := strings.Split(uri, "/")
	if len(parts) >= 3 {
		did = parts[0]           // did:plc:abc123
		recordKey = parts[len(parts)-1] // 3mbrxzvw36c22
	}

	return did, recordKey
}
