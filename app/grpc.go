package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"microservice/app/core"
	"net"
	"strconv"
)

var (
	grpcServer *grpc.Server
	grpcMux    *runtime.ServeMux
)

// GRPC
func InitGRPCServer() (*grpc.Server, *runtime.ServeMux, error) {

	var options []grpc.ServerOption

	// TSL

	tslEnable := viper.GetBool("app.grpc.tsl")
	if tslEnable {
		crt := "./cert/service.pem"
		key := "./cert/service.key"
		caN := "./cert/ca.cert"

		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caN)
		if err != nil {
			return nil, nil, errors.New("cannot initialize GRPC Server")
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, nil, errors.New("failed to append client certs")
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})
		options = append(options, grpc.Creds(creds))
	}

	// Middleware
	options = append(options, []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(errorLogging),
		grpc.ChainUnaryInterceptor(anyLogging),
		grpc.ChainUnaryInterceptor(fromGWOnly),
	}...)

	// Create server
	grpcServer = grpc.NewServer(options...)
	if grpcServer == nil {
		return nil, nil, errors.New("cannot initialize GRPC Server")
	}

	grpcMux = runtime.NewServeMux()

	return grpcServer, grpcMux, nil
}

func RunGRPCServer() {

	gRPCPort := viper.GetString("app.grpc.port")

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":"+gRPCPort)
	if err != nil {
		log.Fatal("%v", err)
	}
	log.Info("GRPC server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("%v", err)
	}
}

// Register

type DeliveryService interface {
	Init() error
}

func InitDelivery(fn interface{}) error {

	di := core.GetDI()
	scope := di.Scope("in")

	err := scope.Provide(fn, dig.As(new(DeliveryService)))
	if err != nil {
		return errors.Wrap(err, "cannot init rest delivery")
	}

	return scope.Invoke(func(d DeliveryService) error {
		err := d.Init()
		if err != nil {
			return err
		}
		return nil
	})
}

func InitGRPCService[T any](s func(grpc.ServiceRegistrar, T), src T) {
	s(grpcServer, src)
}

// Logging interceptor

func fromGWOnly(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	isAuth := false
	m, ok := metadata.FromIncomingContext(ctx)
	if ok {
		tokens := m.Get("Authorization")
		if len(tokens) > 0 && viper.GetString("app.secret") == tokens[0] {
			isAuth = true
		}
	}

	if isAuth {
		return handler(ctx, req)
	} else {
		return nil, errors.Errorf("DENIED access without AUTH! %s", info.FullMethod)
	}
}

func errorLogging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// Calls the handler
	h, err := handler(ctx, req)

	// Log if error
	if err != nil {
		log.Error("%v", err)
		return h, status.Error(codes.Internal, err.Error())
	}

	return h, nil
}

func anyLogging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Info("New call %s", info.FullMethod)
	return handler(ctx, req)
}

// Tools

func ExtractRequestUserId(ctx context.Context) (int32, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userIds := m.Get("user_id")
		if len(userIds) > 0 {
			userId, err := strconv.ParseInt(userIds[0], 10, 32)
			if err != nil {
				return -1, errors.Wrap(err, "cannot parse user_id")
			}
			return int32(userId), nil
		}
	}
	return -1, errors.New("user_id was not found into context")
}
