package gapi

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

// extractMetadata extracts metadata from context
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	data := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// for HTTP
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			data.UserAgent = userAgent[0]
		}

		// for gRPC
		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			data.UserAgent = userAgent[0]
		}

		// for HTTP
		if clientIP := md.Get(xForwardedForHeader); len(clientIP) > 0 {
			data.ClientIP = clientIP[0]
		}
	}

	// for gRPC
	if p, ok := peer.FromContext(ctx); ok {
		data.ClientIP = p.Addr.String()
	}

	return data
}
