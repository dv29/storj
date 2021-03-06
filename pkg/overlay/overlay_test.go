// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"storj.io/storj/internal/test"
	"storj.io/storj/pkg/kademlia"
	proto "storj.io/storj/protos/overlay" // naming proto to avoid confusion with this package
)

func TestFindStorageNodes(t *testing.T) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	assert.NoError(t, err)

	id, err := kademlia.NewID()
	assert.NoError(t, err)
	id2, err := kademlia.NewID()
	assert.NoError(t, err)

	srv := NewMockServer(test.KvStore{id.String(): NewNodeAddressValue(t, "127.0.0.1:9090"), id2.String(): NewNodeAddressValue(t, "127.0.0.1:9090")})
	assert.NotNil(t, srv)

	go srv.Serve(lis)
	defer srv.Stop()

	address := lis.Addr().String()
	c, err := NewClient(address, grpc.WithInsecure())
	assert.NoError(t, err)

	r, err := c.FindStorageNodes(context.Background(), &proto.FindStorageNodesRequest{})
	assert.NoError(t, err)
	assert.NotNil(t, r)

	assert.Len(t, r.Nodes, 2)
}

func TestOverlayLookup(t *testing.T) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	assert.NoError(t, err)

	id, err := kademlia.NewID()

	assert.NoError(t, err)

	srv := NewMockServer(test.KvStore{id.String(): NewNodeAddressValue(t, "127.0.0.1:9090")})
	go srv.Serve(lis)
	defer srv.Stop()

	address := lis.Addr().String()
	c, err := NewClient(address, grpc.WithInsecure())
	assert.NoError(t, err)

	r, err := c.Lookup(context.Background(), &proto.LookupRequest{NodeID: id.String()})
	assert.NoError(t, err)
	assert.NotNil(t, r)
}
