package interceptors

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/rfruffer/go-musthave-shortener/internal/utils"
)

// TrustedSubnetUnaryInterceptor проверяет IP адрес клиента в доверенной подсети для unary RPC
func TrustedSubnetUnaryInterceptor(trustedSubnet string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := checkTrustedSubnet(ctx, trustedSubnet); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// TrustedSubnetStreamInterceptor проверяет IP адрес клиента в доверенной подсети для stream RPC
func TrustedSubnetStreamInterceptor(trustedSubnet string) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := checkTrustedSubnet(stream.Context(), trustedSubnet); err != nil {
			return err
		}
		return handler(srv, stream)
	}
}

// checkTrustedSubnet выполняет проверку IP адреса клиента
func checkTrustedSubnet(ctx context.Context, trustedSubnet string) error {
	if trustedSubnet == "" {
		return nil // Если доверенная подсеть не настроена, разрешаем доступ
	}

	clientIP := getClientIP(ctx)
	if clientIP == "" {
		return status.Error(codes.PermissionDenied, "Не удалось определить IP адрес клиента")
	}

	allowed, err := utils.IsIPInTrustedSubnet(clientIP, trustedSubnet)
	if err != nil {
		return status.Error(codes.Internal, "Ошибка проверки доверенной подсети")
	}

	if !allowed {
		return status.Error(codes.PermissionDenied, "Доступ запрещен: IP адрес не в доверенной подсети")
	}

	return nil
}

// getClientIP извлекает IP адрес клиента из контекста
func getClientIP(ctx context.Context) string {
	// Сначала проверяем заголовок X-Real-IP из metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if realIPs := md.Get("x-real-ip"); len(realIPs) > 0 {
			return realIPs[0]
		}
		
		// Также проверяем X-Forwarded-For
		if forwardedIPs := md.Get("x-forwarded-for"); len(forwardedIPs) > 0 {
			// Берем первый IP из списка
			if ips := strings.Split(forwardedIPs[0], ","); len(ips) > 0 {
				return strings.TrimSpace(ips[0])
			}
		}
	}

	// Если не найден в metadata, извлекаем из peer
	if peer, ok := peer.FromContext(ctx); ok {
		if tcpAddr, ok := peer.Addr.(*net.TCPAddr); ok {
			return tcpAddr.IP.String()
		}
		// Для других типов адресов пытаемся распарсить строку
		addr := peer.Addr.String()
		if host, _, err := net.SplitHostPort(addr); err == nil {
			return host
		}
		return addr
	}

	return ""
}