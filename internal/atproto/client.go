package atproto

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	XRPC *xrpc.Client
}

type Record struct {
	Uri   string          `json:"uri"`
	Cid   string          `json:"cid"`
	Value json.RawMessage `json:"value"`
}

type ListRecordsResponse struct {
	Cursor  *string  `json:"cursor"`
	Records []Record `json:"records"`
}

func NewClient(pdsHost string) *Client {
	if pdsHost == "" {
		pdsHost = "https://bsky.social"
	}
	return &Client{
		XRPC: &xrpc.Client{
			Host: pdsHost,
		},
	}
}

func (c *Client) ResolveHandle(ctx context.Context, handle string) (string, error) {
	out, err := atproto.IdentityResolveHandle(ctx, c.XRPC, handle)
	if err != nil {
		return "", fmt.Errorf("resolving handle: %w", err)
	}
	return out.Did, nil
}

type DIDDocument struct {
	Service []Service `json:"service"`
}

type Service struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// ResolvePDS finds the PDS endpoint for a given DID using plc.directory
func (c *Client) ResolvePDS(ctx context.Context, did string) (string, error) {
	// 1. Try describeRepo first (if we are on the AppView, it works and is faster)
	desc, err := atproto.RepoDescribeRepo(ctx, c.XRPC, did)
	if err == nil && desc.DidDoc != nil {
		// Parse the DidDoc from the response
		// Note: indigo might return it as a map or struct.
		// Let's rely on the public directory if this fails or is complex to parse from the generic type.
		// Actually, let's just use the PLC directory directly for simplicity and reliability.
	}

	// 2. Fallback to PLC Directory
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://plc.directory/%s", did), nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch DID doc: %s", resp.Status)
	}

	var doc DIDDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return "", err
	}

	for _, svc := range doc.Service {
		if svc.ID == "#atproto_pds" || svc.Type == "AtprotoPersonalDataServer" {
			return svc.ServiceEndpoint, nil
		}
	}

	return "", fmt.Errorf("no PDS service found for DID %s", did)
}

func (c *Client) FetchEntries(ctx context.Context, repo string, collection string) ([]Record, error) {
	var records []Record
	cursor := ""

	for {
		params := map[string]interface{}{
			"repo":       repo,
			"collection": collection,
			"limit":      100,
		}
		if cursor != "" {
			params["cursor"] = cursor
		}

		var out ListRecordsResponse
		if err := c.XRPC.Do(ctx, xrpc.Query, "", "com.atproto.repo.listRecords", params, nil, &out); err != nil {
			return nil, fmt.Errorf("listing records: %w", err)
		}

		records = append(records, out.Records...)

		if out.Cursor == nil || *out.Cursor == "" {
			break
		}
		cursor = *out.Cursor
	}

	return records, nil
}
