package converter

import (
	"encoding/json"
	"strings"
	"testing"

	"mariuskimmina.com/leaflet-hugo-sync/internal/atproto"
)

func TestConvertLeaflet_TextBlock(t *testing.T) {
	doc := &atproto.LeafletDocument{
		Title: "Test Post",
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.TextBlock{
							Type:      "pub.leaflet.blocks.text",
							Plaintext: "Hello world",
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	expected := "Hello world\n\n"
	if result.Markdown != expected {
		t.Errorf("expected %q, got %q", expected, result.Markdown)
	}
}

func TestConvertLeaflet_CodeBlock(t *testing.T) {
	doc := &atproto.LeafletDocument{
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.CodeBlock{
							Type:      "pub.leaflet.blocks.code",
							Language:  "go",
							Plaintext: "fmt.Println(\"hello\")",
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	if !strings.Contains(result.Markdown, "```go") {
		t.Errorf("expected code block with go language, got %q", result.Markdown)
	}
	if !strings.Contains(result.Markdown, "fmt.Println") {
		t.Errorf("expected code content, got %q", result.Markdown)
	}
}

func TestConvertLeaflet_ImageBlock(t *testing.T) {
	doc := &atproto.LeafletDocument{
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.ImageBlock{
							Type: "pub.leaflet.blocks.image",
							Alt:  "Test Image",
							Image: atproto.Blob{
								Ref: atproto.BlobRef{
									Link: "bafytest123",
								},
							},
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	if !strings.Contains(result.Markdown, "![Test Image](bafytest123)") {
		t.Errorf("expected image markdown, got %q", result.Markdown)
	}
	if len(result.Images) != 1 {
		t.Errorf("expected 1 image reference, got %d", len(result.Images))
	}
	if result.Images[0].Alt != "Test Image" {
		t.Errorf("expected alt 'Test Image', got %q", result.Images[0].Alt)
	}
}

func TestRenderText_WithLinkFacet(t *testing.T) {
	block := &atproto.TextBlock{
		Plaintext: "Check out example.com for more info",
		Facets: []atproto.Facet{
			{
				Index: atproto.Features{
					ByteStart: 10,
					ByteEnd:   21,
				},
				Features: []atproto.Feature{
					{
						Type: "pub.leaflet.richtext.facet#link",
						URI:  "https://example.com",
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result := conv.renderText(block)

	expected := "Check out [example.com](https://example.com) for more info"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderText_WithInlineCode(t *testing.T) {
	block := &atproto.TextBlock{
		Plaintext: "Use the fmt.Println function",
		Facets: []atproto.Facet{
			{
				Index: atproto.Features{
					ByteStart: 8,
					ByteEnd:   19,
				},
				Features: []atproto.Feature{
					{
						Type: "pub.leaflet.richtext.facet#code",
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result := conv.renderText(block)

	expected := "Use the `fmt.Println` function"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestRenderText_WithMention(t *testing.T) {
	block := &atproto.TextBlock{
		Plaintext: "Thanks @alice for the help",
		Facets: []atproto.Facet{
			{
				Index: atproto.Features{
					ByteStart: 7,
					ByteEnd:   13,
				},
				Features: []atproto.Feature{
					{
						Type: "pub.leaflet.richtext.facet#didMention",
						Did:  "did:plc:alice123",
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result := conv.renderText(block)

	expected := "Thanks [@alice](https://bsky.app/profile/did:plc:alice123) for the help"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestConvertLeaflet_UnorderedList(t *testing.T) {
	listItem := atproto.ListItem{
		Content: mustMarshal(atproto.TextBlock{
			Type:      "pub.leaflet.blocks.text",
			Plaintext: "First item",
		}),
	}

	doc := &atproto.LeafletDocument{
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.UnorderedListBlock{
							Type:     "pub.leaflet.blocks.unorderedList",
							Children: []atproto.ListItem{listItem},
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("")
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	if !strings.Contains(result.Markdown, "- First item") {
		t.Errorf("expected list item, got %q", result.Markdown)
	}
}

func TestConvertLeaflet_BskyPost_Link(t *testing.T) {
	doc := &atproto.LeafletDocument{
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.BskyPostBlock{
							Type: "pub.leaflet.blocks.bskyPost",
							PostRef: atproto.PostRef{
								Uri: "at://did:plc:abc123/app.bsky.feed.post/3mbrxzvw36c22",
								Cid: "test-cid",
							},
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("link") // Default link mode
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	expected := "[View on Bluesky](https://bsky.app/profile/did:plc:abc123/post/3mbrxzvw36c22)"
	if !strings.Contains(result.Markdown, expected) {
		t.Errorf("expected link format, got %q", result.Markdown)
	}
}

func TestConvertLeaflet_BskyPost_Shortcode(t *testing.T) {
	doc := &atproto.LeafletDocument{
		Pages: []atproto.Page{
			{
				Blocks: []atproto.BlockWrapper{
					{
						Block: mustMarshal(atproto.BskyPostBlock{
							Type: "pub.leaflet.blocks.bskyPost",
							PostRef: atproto.PostRef{
								Uri: "at://did:plc:abc123/app.bsky.feed.post/3mbrxzvw36c22",
								Cid: "test-cid",
							},
						}),
					},
				},
			},
		},
	}

	conv := NewConverter("shortcode")
	result, err := conv.ConvertLeaflet(doc)
	if err != nil {
		t.Fatalf("ConvertLeaflet failed: %v", err)
	}

	expected := `{{< bsky did="did:plc:abc123" postid="3mbrxzvw36c22" >}}`
	if !strings.Contains(result.Markdown, expected) {
		t.Errorf("expected shortcode format, got %q", result.Markdown)
	}
}

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
