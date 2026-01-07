# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

leaflet-hugo-sync is a Go CLI tool that syncs Leaflet/Bluesky blog posts to Hugo-compatible markdown files. It fetches blog entries from the AT Protocol (ATProto) network and converts them to markdown with downloaded media.

## Build & Run Commands

### Using Nix (recommended)

```bash
# Build with Nix
nix build

# Run directly
nix run . -- -config .leaflet-sync.yaml

# Enter development shell
nix develop

# Within nix develop shell:
go build -o leaflet-hugo-sync ./cmd/leaflet-hugo-sync
go test ./...
```

### Using Go directly

```bash
# Build the binary
go build -o leaflet-hugo-sync ./cmd/leaflet-hugo-sync

# Run directly
go run ./cmd/leaflet-hugo-sync/main.go -config .leaflet-sync.yaml

# Run tests
go test ./...

# Run specific package tests
go test ./internal/config
go test ./internal/generator
```

## Architecture

### Core Data Flow

The application follows this pipeline:

1. **ATProto Resolution** (`internal/atproto/client.go`)
   - Resolves Bluesky handle → DID using public resolver (bsky.social)
   - Resolves DID → PDS endpoint via plc.directory
   - Creates client pointing to user's personal data server

2. **Record Fetching** (`internal/atproto/client.go:FetchEntries`)
   - Fetches all records from specified collection (e.g., `pub.leaflet.document`)
   - Handles pagination using cursors
   - Optionally filters by publication URI

3. **Content Conversion** (`internal/converter/markdown.go`)
   - Converts Leaflet's block-based format to markdown
   - Handles multiple block types: text, code, unorderedList, image, bskyPost
   - Processes rich text facets (links, mentions, inline code)
   - Returns markdown string + list of image references

4. **Media Download** (`internal/media/downloader.go`)
   - Downloads blob images from PDS via `com.atproto.sync.getBlob` XRPC endpoint
   - Determines file extension from Content-Type header
   - Caches downloaded images (skips if file exists)
   - Returns Hugo-compatible image paths

5. **Post Generation** (`internal/generator/hugo.go`)
   - Uses Go templates to generate frontmatter from config
   - Combines frontmatter + markdown content
   - Writes final `.md` files to configured output directory

### Key Types

**ATProto Types** (`internal/atproto/types.go`):
- `LeafletDocument`: Top-level document with title, publishedAt, pages
- `Page`: Contains array of BlockWrappers
- `BlockWrapper`: Wraps a block with deferred JSON unmarshaling
- Block types: `TextBlock`, `CodeBlock`, `UnorderedListBlock`, `ImageBlock`, `BskyPostBlock`
- `Facet`: Rich text annotation with byte-based ranges (links, mentions, inline code)

**Configuration** (`internal/config/config.go`):
- `Source`: Specifies handle, collection, optional publication_name
- `Output`: Defines posts_dir, images_dir, image_path_prefix, bsky_embed_style
- `Template`: Go template strings for frontmatter and content

### Important Implementation Details

**Facet Processing**: Facets use byte offsets, not rune offsets. The converter handles this by working with `[]byte` instead of string indices when applying rich text formatting.

**Collection Migration**: The code defaults `com.whtwnd.blog.entry` to `pub.leaflet.document` as the old collection is deprecated (see main.go:104-110).

**Blob Download**: Uses the pattern `{PDS_HOST}/xrpc/com.atproto.sync.getBlob?did={DID}&cid={CID}` to fetch images.

**Publication Filtering**: When `publication_name` is configured, the tool first fetches `pub.leaflet.publication` records, finds the matching publication URI, then filters documents by that URI.

## Configuration File

The tool requires a YAML config file (default: `.leaflet-sync.yaml`):

```yaml
source:
  handle: "username.bsky.social"
  collection: "pub.leaflet.document"
  publication_name: "optional-publication-name"  # Filter by specific publication

output:
  posts_dir: "content/posts/leaflet"
  images_dir: "static/images/leaflet"
  image_path_prefix: "/images/leaflet"
  bsky_embed_style: "link"  # "link" (default) or "shortcode" for Hugo embeds

template:
  frontmatter: |
    ---
    title: "{{ .Title }}"
    date: {{ .CreatedAt }}
    ---
```

Template variables available: `.Title`, `.CreatedAt`, `.Slug`, `.Handle`, `.OriginalURL`, `.Content`

## Dependencies

- `github.com/bluesky-social/indigo`: Official ATProto/Bluesky Go library for XRPC client and ATProto methods
- `gopkg.in/yaml.v3`: YAML config parsing
