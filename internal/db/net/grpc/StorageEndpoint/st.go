package StorageEndpoint

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"jaeger/internal/db/repo"
	"jaeger/internal/proto/db"
	"log/slog"
)

type Server struct {
	db.UnimplementedStorageEndpointServer
	st repo.Storage
	lg *slog.Logger
}

func NewStorageEndpoint(st repo.Storage, lg *slog.Logger) *Server {
	return &Server{st: st, lg: lg}
}

func (s *Server) Get(ctx context.Context, key *db.Key) (*db.KeyValue, error) {
	lg := s.logger(ctx)
	k := key.GetKey()
	if checkClosed(ctx.Done()) {
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", ctx.Err().Error())
	}
	val, ok, err := s.st.Get(ctx, k)
	lg.Debug("result", "key", k, "val", val, "ok", ok)
	if err != nil {
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", err.Error())
	}
	if !ok {
		lg.Error("key not found", "key", k)
		return nil, status.Errorf(codes.NotFound, "key %s not found", k)
	}
	return &db.KeyValue{Key: k, Val: val}, nil
}

func (s *Server) Set(ctx context.Context, value *db.KeyValue) (*emptypb.Empty, error) {
	lg := s.logger(ctx)
	k, v := value.GetKey(), value.GetVal()
	lg.Debug("result", "key", k, "val", v)
	if checkClosed(ctx.Done()) {
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", ctx.Err().Error())
	}
	return &emptypb.Empty{}, s.st.Set(ctx, k, v)
}

func (s *Server) Del(ctx context.Context, key *db.Key) (*emptypb.Empty, error) {
	lg := s.logger(ctx)
	k := key.GetKey()
	lg.Debug("result", "key", k)
	if checkClosed(ctx.Done()) {
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", ctx.Err().Error())
	}
	return &emptypb.Empty{}, s.st.Del(ctx, k)
}

func (s *Server) logger(ctx context.Context) *slog.Logger {
	uid := uuid.New()
	lg := s.lg.With("ID", uid)
	p, ok := peer.FromContext(ctx)
	fn := func() {
		if p.AuthInfo == nil {
			lg.Info("connection accepted",
				"Addr", p.Addr.String(),
				"LocalAddr", p.LocalAddr.String())
		} else {
			lg.Info("connection accepted",
				"Addr", p.Addr.String(),
				"LocalAddr", p.LocalAddr.String(),
				"AuthType", p.AuthInfo.AuthType())
		}
	}
	if !ok {
		lg.Warn("peer not found in context")
	} else {
		fn()
	}
	return lg
}

func checkClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
