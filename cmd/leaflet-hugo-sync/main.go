package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/marius/leaflet-hugo-sync/internal/atproto"
	"github.com/marius/leaflet-hugo-sync/internal/config"
	"github.com/marius/leaflet-hugo-sync/internal/generator"
	"github.com/marius/leaflet-hugo-sync/internal/media"
)

func lastPathPart(uri string) string {
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

func main() {
	configPath := flag.String("config", ".leaflet-sync.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	client := atproto.NewClient("")

	did, err := client.ResolveHandle(ctx, cfg.Source.Handle)
	if err != nil {
		log.Fatalf("failed to resolve handle: %v", err)
	}

	fmt.Printf("Resolved %s to %s\n", cfg.Source.Handle, did)

	records, err := client.FetchEntries(ctx, did, cfg.Source.Collection)
	if err != nil {
		log.Fatalf("failed to fetch entries: %v", err)
	}

	fmt.Printf("Found %d entries\n", len(records))

	downloader := media.NewDownloader(cfg.Output.ImagesDir, cfg.Output.ImagePathPrefix, client.XRPC.Host)
	gen := generator.NewGenerator(cfg)

	for _, rec := range records {
		var entry atproto.BlogEntry
		valBytes, _ := json.Marshal(rec.Value)
		if err := json.Unmarshal(valBytes, &entry); err != nil {
			fmt.Printf("Failed to unmarshal record %s: %v\n", rec.Uri, err)
			continue
		}

		if entry.Slug == "" {
			entry.Slug = lastPathPart(rec.Uri)
		}

		fmt.Printf("Processing: %s\n", entry.Title)

		// Handle images
		if entry.Embed != nil && len(entry.Embed.Images) > 0 {
			imageMarkdown := "\n\n"
			for _, img := range entry.Embed.Images {
				localPath, err := downloader.DownloadBlob(ctx, did, img.Image.Ref.Link)
				if err != nil {
					fmt.Printf("  Failed to download image: %v\n", err)
					continue
				}
				imageMarkdown += fmt.Sprintf("![%s](%s)\n", img.Alt, localPath)
			}
			entry.Content += imageMarkdown
		}

		postData := generator.PostData{
			Title:     entry.Title,
			CreatedAt: entry.CreatedAt,
			Slug:      entry.Slug,
			Handle:    cfg.Source.Handle,
			Content:   entry.Content,
		}

		if err := gen.GeneratePost(postData); err != nil {
			fmt.Printf("  Failed to generate post: %v\n", err)
		}
	}

	fmt.Println("Done!")
}
