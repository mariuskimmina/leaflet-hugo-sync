# leaflet-hugo-sync

Syncs Leaflet blog posts from the AT Protocol network to Hugo-compatible markdown files.

## Configuration

Create a `.leaflet-sync.yaml` file in your `hugo` project:

```yaml
source:
  handle: "username.bsky.social"
  collection: "pub.leaflet.document"
  publication_name: "optional-publication-name"

output:
  posts_dir: "content/posts/leaflet"
  images_dir: "static/images/leaflet"
  image_path_prefix: "/images/leaflet"
  bsky_embed_style: "link"  # Optional: "link" (default) or "shortcode"

template:
  frontmatter: |
    ---
    title: "{{ .Title }}"
    date: {{ .CreatedAt }}
    original_url: "{{ .OriginalURL }}"
    ---
```

## BlueSky Post Embeds

When your Leaflet posts reference BlueSky posts, they can be rendered in two ways:

- **`link`** (default): Simple markdown links that work everywhere
- **`shortcode`**: Rich Hugo shortcodes for custom styling

For shortcode setup instructions, see [SHORTCODE_SETUP.md](SHORTCODE_SETUP.md).

## How it works

The tool resolves your Bluesky handle to find your personal data server, fetches your Leaflet documents, converts them to markdown, downloads embedded images, and writes Hugo-compatible markdown files to your specified output directory.
