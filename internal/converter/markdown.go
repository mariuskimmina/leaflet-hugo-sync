package converter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marius/leaflet-hugo-sync/internal/atproto"
)

type Converter struct {
	// We might need to track image downloads here or just return image references?
	// For now, let's just return the Markdown text and a list of image blobs to download.
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
				// We use a placeholder path that main.go will resolve
				// Actually, main.go needs to know about this.
				// Let's assume standard markdown image syntax: ![alt](cid)
				// The downloader will replace it or we pre-calculate the path?
				// Better: Return the blob CID as the URL, and main.go does a string replace or we handle it here if we pass config.
				// For now, let's use the blob ref link as the URL.
				sb.WriteString(fmt.Sprintf("![%s](%s)\n\n", imgBlock.Alt, imgBlock.Image.Ref.Link))
				images = append(images, ImageRef{Blob: imgBlock.Image, Alt: imgBlock.Alt})

            case "pub.leaflet.blocks.bskyPost":
                var postBlock atproto.BskyPostBlock
                if err := json.Unmarshal(blockWrapper.Block, &postBlock); err != nil {
                    continue
                }
                // Render a link to the post for now, as we can't easily embed it without JS
                postURL := fmt.Sprintf("https://bsky.app/profile/%s/post/%s", "did:...", lastPathPart(postBlock.PostRef.Uri))
                // We don't have the handle here easily to make a pretty URL, but we can try.
                // Actually, let's just make a blockquote link.
                sb.WriteString(fmt.Sprintf("> [View on Bluesky](%s)\n\n", postURL)) // TODO: Improve this
			}
		}
	}

	return &ConversionResult{
		Markdown: sb.String(),
		Images:   images,
	}, nil
}

func (c *Converter) renderText(block *atproto.TextBlock) string {
	// Apply facets
	// Facets are ranges. We need to insert markdown syntax at specific indices.
	// This is tricky because inserting characters shifts indices.
	// Best approach: Slice the string and reconstruct it.

	// Sort facets by start index (descending) to avoid shifting issues?
	// Actually, we should iterate from start to end, keeping track of current index.
	
	// Simplify: Just handle links for now.
    // NOTE: Facets in ATProto are byte-offsets, not rune-offsets. Go strings are UTF-8.
    
    // Convert string to byte slice for easier indexing
    data := []byte(block.Plaintext)
    
    // Map of byte_index -> string_to_insert
    // But we wrap text.
    
    // Let's just do a linear pass if facets are non-overlapping and sorted.
    // They should be.
    
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
            if feat.Type == "pub.leaflet.richtext.facet#link" {
                replacement = fmt.Sprintf("[%s](%s)", text, feat.URI)
            } else if feat.Type == "pub.leaflet.richtext.facet#didMention" {
                 replacement = fmt.Sprintf("[%s](https://bsky.app/profile/%s)", text, feat.Did)
            }
             // code facet? pub.leaflet.richtext.facet#code -> `text`
            if feat.Type == "pub.leaflet.richtext.facet#code" {
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
