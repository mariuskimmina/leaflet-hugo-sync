# BlueSky Post Embed Setup

The sync tool can generate BlueSky post embeds in two ways:
- **Simple links** (default) - Plain markdown links that work everywhere
- **Hugo shortcodes** (optional) - Rich, styled embeds with custom styling

## Quick Setup for Hugo Shortcodes

By default, BlueSky posts are rendered as simple markdown links. To enable rich Hugo shortcode embeds:

1. **Enable shortcode mode in your config file** (`.leaflet-sync.yaml`):
   ```yaml
   output:
     posts_dir: "content/posts/leaflet"
     images_dir: "static/images/leaflet"
     image_path_prefix: "/images/leaflet"
     bsky_embed_style: "shortcode"  # Add this line
   ```

3. **Copy the shortcode to your Hugo site:**
   ```bash
   cp bsky.html /path/to/your/hugo/site/layouts/shortcodes/bsky.html
   ```

4. **Run the sync:**
   ```bash
   leaflet-hugo-sync -config .leaflet-sync.yaml
   ```
