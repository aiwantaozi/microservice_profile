package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func AuthStreamInterceptor(
	srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) error {
	md, _ := metadata.FromContext(ss.Context())
	if len(md["uid"]) == 0 {
		return grpc.Errorf(codes.Unauthenticated, "authentication required")
	}
	if _, err := ValidateUser(md["uid"]); err != nil {
		return grpc.Errorf(codes.Unauthenticated, "authentication required")
	}

	return handler(srv, ss)
}

func ValidateUser(uidSlice []string) (string, error) {
	// TODO
	uid := uidSlice[0]
	return uid, nil
}
