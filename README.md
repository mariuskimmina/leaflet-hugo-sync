# leaflet-hugo-sync

A CLI tool to hydrate Hugo blogs with content hosted on Leaflet (via AT Protocol).

## Features

- **Configurable Mapping**: Use a YAML config file to define source and output paths.
- **Rich Text to Markdown**: Converts Leaflet blog entries to Hugo-compatible Markdown.
- **Image Mirroring**: Automatically downloads image blobs and updates paths.
- **Customizable Frontmatter**: Use Go templates to define your Hugo frontmatter.

## Installation

```bash
go install github.com/marius/leaflet-hugo-sync/cmd/leaflet-hugo-sync@latest
```

## Configuration

Create a `.leaflet-sync.yaml` in your Hugo project root:

```yaml
source:
  handle: "yourname.bsky.social"
  collection: "com.whtwnd.blog.entry"

output:
  posts_dir: "content/posts/leaflet"
  images_dir: "static/images/leaflet"
  image_path_prefix: "/images/leaflet"

template:
  frontmatter: |
    ---
    title: "{{ .Title }}"
    date: {{ .CreatedAt }}
    original_url: "https://leaflet.pub/{{ .Handle }}/{{ .Slug }}"
    ---
```

## Usage

Run the tool from your Hugo project root:

```bash
leaflet-hugo-sync
```

You can specify a custom config path:

```bash
leaflet-hugo-sync -config my-config.yaml
```

## License

MIT
