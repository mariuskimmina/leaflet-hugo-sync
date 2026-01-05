package atproto

import (
	"context"
	"fmt"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type Client struct {
	XRPC *xrpc.Client
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

func (c *Client) FetchEntries(ctx context.Context, repo string, collection string) ([]*atproto.RepoListRecords_Record, error) {
	var records []*atproto.RepoListRecords_Record
	cursor := ""

	for {
		out, err := atproto.RepoListRecords(ctx, c.XRPC, collection, cursor, 100, repo, false)
		if err != nil {
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
