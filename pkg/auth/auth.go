package auth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// For safety use in context
type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func AuthInterceptor(jwtKey []byte) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// For use register and login without token
		if info.FullMethod == "/user.UserService/RegisterUser" || info.FullMethod == "/user.UserService/Login" {
			return handler(ctx, req)
		}

		// Metadata from context request
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		// header - authorization
		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		// Get token without "Bearer" and " "
		tokenString := strings.TrimPrefix(authHeaders[0], "Bearer ")
		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			return nil, status.Error(codes.Unauthenticated, "auth token is empty")
		}

		// Token parse with check(algoritm - HMAC)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, status.Error(codes.Unauthenticated, "unexpected signing method")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Parse data from token (how map)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		}

		// Get role from claims
		role, ok := claims["role"].(string)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "role claim missing")
		}

		// Get user id from claims
		userID, _ := claims["user_id"].(string)

		// Add user id and role under unique keys
		newCtx := context.WithValue(ctx, UserIDKey, userID)
		newCtx = context.WithValue(newCtx, RoleKey, role)

		return handler(newCtx, req)
	}
}
