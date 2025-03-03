package api

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/metadata_api"

	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	swgui "github.com/swaggest/swgui/v5emb"
	openapi "github.com/xmtp/xmtpd/pkg/proto/openapi"

	"github.com/xmtp/xmtpd/pkg/api"
	"github.com/xmtp/xmtpd/pkg/config"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/message_api"
	"github.com/xmtp/xmtpd/pkg/proto/xmtpv4/payer_api"
)

// HTTP Server Gateway for xmtpd using grpc-gateway
// https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/
type HttpServerGateway struct {
	apiServer *api.ApiServer
	gwmux     *runtime.ServeMux

	ctx     context.Context
	cancel  context.CancelFunc
	log     *zap.Logger
	options config.ServerOptions
}

func NewHTTPGateway(ctx context.Context, apiServer *api.ApiServer, log *zap.Logger, options config.ServerOptions) (*HttpServerGateway, error) {
	var err error

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/x-protobuf", &runtime.ProtoMarshaller{}),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithStreamErrorHandler(runtime.DefaultStreamErrorHandler),
		// runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
	)

	s := &HttpServerGateway{
		options:   options,
		log:       log,
		apiServer: apiServer,
		gwmux:     gwmux,
	}
	s.ctx, s.cancel = context.WithCancel(ctx)

	client, err := apiServer.DialGRPC(ctx)
	if err != nil {
		log.Fatal("dialing GRPC from HTTP Gateway")
		return nil, err
	}

	if options.Replication.Enable {
		err = metadata_api.RegisterMetadataApiHandler(ctx, gwmux, client)
		log.Info("Metadata http gateway enabled")

		if err != nil {
			return nil, errors.Wrap(err, "creating metadata api handler")
		}

		err = message_api.RegisterReplicationApiHandler(ctx, gwmux, client)
		log.Info("Replication http gateway enabled")
		if err != nil {
			return nil, errors.Wrap(err, "creating replication api handler")
		}
	}

	if options.Payer.Enable {
		err = payer_api.RegisterPayerApiHandler(ctx, gwmux, client)
		log.Info("Payer http gateway enabled")
		if err != nil {
			return nil, errors.Wrap(err, "creating payer api handler")
		}
	}

	err = s.createSwaggerUI()
	if err != nil {
		return nil, err
	}

	return s, nil
}

type APISpec struct {
	Name        string
	FilePath    string
	Description string
}

// Define our API specifications
var apiSpecs = []APISpec{
	{
		Name:        "Envelopes API",
		FilePath:    "xmtpv4/envelopes/envelopes.swagger.json",
		Description: "Types for handling XMTP message envelopes",
	},
	{
		Name:        "Message API",
		FilePath:    "xmtpv4/message_api/message_api.swagger.json",
		Description: "messaging functionality",
	},
	{
		Name:        "Misbehavior API",
		FilePath:    "xmtpv4/message_api/misbehavior_api.swagger.json",
		Description: "reporting and handling misbehavior",
	},
	{
		Name:        "Metadata API",
		FilePath:    "xmtpv4/metadata_api/metadata_api.swagger.json",
		Description: "managing message metadata",
	},
	{
		Name:        "Payer API",
		FilePath:    "xmtpv4/payer_api/payer_api.swagger.json",
		Description: "payment-related functionality",
	},
}

func (s *HttpServerGateway) createSwaggerUI() error {

	for _, apiSpec := range apiSpecs {
		// ex: /docs/xmtpv4/EnvelopesApi
		urlPath := fmt.Sprintf("/docs/xmtpv4/%s", strings.ToLower(strings.ReplaceAll(apiSpec.Name, " ", "")))
		filePath := apiSpec.FilePath
		s.gwmux.HandlePath("GET", urlPath, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			data, err := openapi.SwaggerFS.ReadFile(filePath)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error reading Swagger file: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().set("Content-Type", "application/json")
			_, err = w.Write(data)
		})
	}
	// swaggerUI := swgui.NewHandler("API", "/swagger.json", "/")
	s.gwmux.HandlePath("GET", "/docs/xmtpv4/{something}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

		for _, apiSpec := range apiSpecs {
		}
		// _, err := w.Write(xmtpv4openapi.xmtpv4EnvelopesJSON)
		//	if err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//	}
	})

	return nil
}

/*
// Create a mapping between URL paths and the embedded file paths

	var pathMapping = map[string]string{
		"/api/v4/envelopes":   "xmtpv4/envelopes/envelopes.swagger.json",
		"/api/v4/message":     "xmtpv4/message_api/message_api.swagger.json",
		"/api/v4/misbehavior": "xmtpv4/message_api/misbehavior_api.swagger.json",
		"/api/v4/metadata":    "xmtpv4/metadata_api/metadata_api.swagger.json",
		"/api/v4/payer":       "xmtpv4/payer_api/payer_api.swagger.json",
	}
*/

//func startHTTPServer(ctx context.Context, log *zap.Logger, options config.ServerOptions) error {
//
//	// swaggerUI := swgui.NewHandler("API", "/swagger.json", "/")
//	// gwmux.HandlePath("GET", "/swagger-ui/xmtpv4/{something}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
//	// _, err := w.Write(xmtpv4openapi.xmtpv4EnvelopesJSON)
//	//		if err != nil {
//	//			w.WriteHeader(http.StatusInternalServerError)
//	//		}
//	// })
//
//	/*
//
//		addr := addrString(s.HTTPAddress, s.HTTPPort)
//		s.httpListener, err = net.Listen("tcp", addr)
//		if err != nil {
//			return errors.Wrap(err, "creating grpc-gateway listener")
//		}
//
//		// Add two handler wrappers to mux: gzipWrapper and allowCORS
//		server := http.Server{
//			Addr:    addr,
//			Handler: allowCORS(gzipWrapper(mux)),
//		}
//
//		tracing.GoPanicWrap(s.ctx, &s.wg, "http", func(ctx context.Context) {
//			s.Log.Info("serving http", zap.String("address", s.httpListener.Addr().String()))
//			if s.httpListener == nil {
//				s.Log.Fatal("no http listener")
//			}
//			err = server.Serve(s.httpListener)
//			if err != nil && err != http.ErrServerClosed && !isErrUseOfClosedConnection(err) {
//				s.Log.Error("serving http", zap.Error(err))
//			}
//		})
//	*/
//	return nil
//}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Client-Version",
		"X-App-Version",
		"Baggage",
		"DNT",
		"Sec-CH-UA",
		"Sec-CH-UA-Mobile",
		"Sec-CH-UA-Platform",
		"Sentry-Trace",
		"User-Agent",
		"x-libxmtp-version",
		"x-app-version",
	}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
			preflightHandler(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// func incomingHeaderMatcher(key string) (string, bool) {
// 	switch strings.ToLower(key) {
// 	case apicontext.ClientVersionMetadataKey:
// 		return key, true
// 	case apicontext.AppVersionMetadataKey:
// 		return key, true
// 	default:
// 		return key, false
// 	}
// }
