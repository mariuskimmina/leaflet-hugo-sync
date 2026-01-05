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
	"github.com/marius/leaflet-hugo-sync/internal/converter"
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
	
	// 1. Resolve Handle to DID using public resolver (bsky.social)
	baseClient := atproto.NewClient("https://bsky.social")
	did, err := baseClient.ResolveHandle(ctx, cfg.Source.Handle)
	if err != nil {
		log.Fatalf("failed to resolve handle: %v", err)
	}
	fmt.Printf("Resolved %s to %s\n", cfg.Source.Handle, did)

	// 2. Resolve PDS Endpoint for the DID
	pdsEndpoint, err := baseClient.ResolvePDS(ctx, did)
	if err != nil {
		log.Fatalf("failed to resolve PDS: %v", err)
	}
	fmt.Printf("PDS Endpoint: %s\n", pdsEndpoint)

	// 3. Connect to User's PDS
	pdsClient := atproto.NewClient(pdsEndpoint)

		// 4. Resolve Publication (if configured)

		var publicationURI string

		if cfg.Source.PublicationName != "" {

			fmt.Printf("Resolving publication '%s'...\n", cfg.Source.PublicationName)

			pubRecords, err := pdsClient.FetchEntries(ctx, did, "pub.leaflet.publication")

			if err != nil {

				log.Fatalf("failed to fetch publications: %v", err)

			}

			for _, rec := range pubRecords {

				var pub atproto.LeafletPublication

				if err := json.Unmarshal(rec.Value, &pub); err == nil {

					if pub.Name == cfg.Source.PublicationName {

						publicationURI = rec.Uri

						fmt.Printf("Found publication URI: %s\n", publicationURI)

						break

					}

				}

			}

			if publicationURI == "" {

				log.Fatalf("Publication '%s' not found", cfg.Source.PublicationName)

			}

		}

	

		// 5. Fetch Entries

		// Update collection to Leaflet Document if user hasn't specified it

		collection := cfg.Source.Collection

		if collection == "com.whtwnd.blog.entry" {

			fmt.Println("Warning: Defaulting to 'pub.leaflet.document' as 'com.whtwnd.blog.entry' seems deprecated/unused for Leaflet.")

			collection = "pub.leaflet.document"

		}

	

		records, err := pdsClient.FetchEntries(ctx, did, collection)

		if err != nil {

			log.Fatalf("failed to fetch entries: %v", err)

		}

	

		fmt.Printf("Found %d entries\n", len(records))

	

		downloader := media.NewDownloader(cfg.Output.ImagesDir, cfg.Output.ImagePathPrefix, pdsClient.XRPC.Host)

		gen := generator.NewGenerator(cfg)

		conv := converter.NewConverter()

	

		for _, rec := range records {

			// Try to unmarshal as LeafletDocument

			var doc atproto.LeafletDocument

			

			// Check type first

			var typeCheck struct {

				Type string `json:"$type"`

			}

			if err := json.Unmarshal(rec.Value, &typeCheck); err != nil {

				fmt.Printf("Failed to check type for record %s: %v\n", rec.Uri, err)

				continue

			}

			

			if typeCheck.Type != "pub.leaflet.document" {

				// Skip or try legacy

				continue

			}

	

			if err := json.Unmarshal(rec.Value, &doc); err != nil {

				fmt.Printf("Failed to unmarshal record %s: %v\n", rec.Uri, err)

				continue

			}

	

			// Filter by Publication

			if publicationURI != "" && doc.Publication != publicationURI {

				// Skip entries not matching the publication

				continue

			}

			

			fmt.Printf("Processing: %s\n", doc.Title)

	

			// Convert to Markdown

			result, err := conv.ConvertLeaflet(&doc)

			if err != nil {

				fmt.Printf("  Failed to convert document: %v\n", err)

				continue

			}

	

			// Download Images

			finalContent := result.Markdown

			for _, imgRef := range result.Images {

				localPath, err := downloader.DownloadBlob(ctx, did, imgRef.Blob.Ref.Link)

				if err != nil {

					fmt.Printf("  Failed to download image: %v\n", err)

					continue

				}

				finalContent = strings.ReplaceAll(finalContent, imgRef.Blob.Ref.Link, localPath)

			}

	

			// Slug generation

			slug := lastPathPart(rec.Uri)

	

			// Construct Original URL (Assuming toolbox.leaflet.pub structure or generic)

			// User config usually handles the "base" part in template. 

			// But here we construct a full URL for the OriginalURL field.

			// If we don't know the base, we can't fully construct it.

			// However, user provided "https://toolbox.leaflet.pub/{{ .Slug }}" in template.

			// Let's pass the slug and handle.

			// Ideally, we should construct the full URL if we can. 

			// For now, let's just pass what we have.

			

			originalURL := fmt.Sprintf("https://leaflet.pub/%s", slug) // Fallback generic

	

			postData := generator.PostData{

				Title:       doc.Title,

				CreatedAt:   doc.PublishedAt,

				Slug:        slug,

				Handle:      cfg.Source.Handle,

				OriginalURL: originalURL,

				Content:     finalContent,

			}

	

			if err := gen.GeneratePost(postData); err != nil {

				fmt.Printf("  Failed to generate post: %v\n", err)

			}

		}

	

		fmt.Println("Done!")

	}

	