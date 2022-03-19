// Copyright 2019 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fpsserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	// "strings"
	"time"

	"google.golang.org/grpc"

	// api "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	"github.com/intel/cri-resource-manager/pkg/cri/resource-manager/sockets"
	// "github.com/intel/cri-resource-manager/pkg/dump"
	logger "github.com/intel/cri-resource-manager/pkg/log"
	"github.com/intel/cri-resource-manager/pkg/utils"

	"github.com/intel/cri-resource-manager/pkg/instrumentation"
	// "go.opencensus.io/trace"
)

// Options contains the configurable options of our CRI server.
type Options struct {
	// Socket is the path of our gRPC servers unix-domain socket.
	Socket string
	// User is the user ID for our gRPC socket.
	User int
	// Group is the group ID for our gRPC socket.
	Group int
	// Mode is the permission mode bits for our gRPC socket.
	Mode os.FileMode
	// QualifyReqFn produces return context for disambiguating a CRI request/reply.
	// QualifyReqFn func(interface{}) string
}

// Handler is a CRI server generic request handler.
type Handler grpc.UnaryHandler

// Interceptor is a hook that intercepts processing a request by a handler.
type Interceptor func(context.Context, string, interface{}) (interface{}, error)

type FPSDropHandler interface{
	HandleFPSDrop(string) error
}

// Server is the interface we expose for controlling our CRI server.
type Server interface {
	// RegisterFPSService registers the server.
	RegisterFPSService() error
	// RegisterInterceptors registers the given interceptors with the server.
	RegisterFPSDropHandler(FPSDropHandler) error
	// Start starts the request processing loop (goroutine) of the server.
	Start() error
	// Stop stops the request processing loop (goroutine) of the server.
	Stop()
	// Chmod changes the permissions of the server's socket.
	Chmod(mode os.FileMode) error
	// Chown changes ownership of the server's socket.
	Chown(uid, gid int) error
}

// server is the implementation of Server.
type server struct {
	logger.Logger
	listener     net.Listener              // socket our gRPC server listens on
	server       *grpc.Server              // our gRPC server
	options      Options                   // server options
	fpsDropHandler *FPSDropHandler
	UnimplementedFPSServiceServer
} 

// NewServer creates a new server instance.
func NewServer(options Options) (Server, error) {
	if !filepath.IsAbs(options.Socket) {
		return nil, serverError("invalid socket '%s', expecting absolute path",
			options.Socket)
	}

	s := &server{
		Logger:  logger.NewLogger("cri/fpsserver"),
		options: options,
	}

	return s, nil
}

func(s *server) FPSDrop(cxt context.Context, request *FPSDropRequest) (*FPSDropReply, error) {
	s.Info("receive FPS drop message from %s", request.Id)
	reply := FPSDropReply{
	}
	(*s.fpsDropHandler).HandleFPSDrop(request.Id)
	return &reply, nil
}


// RegisterImageService registers an image service with the server.
func (s *server) RegisterFPSService() error {
	if err := s.createGrpcServer(); err != nil {
		return err
	}

	RegisterFPSServiceServer(s.server, s)

	return nil
}

func (s *server) RegisterFPSDropHandler(handler FPSDropHandler) error {
	s.fpsDropHandler = &handler

	return nil
}


// Start starts the servers request processing goroutine.
func (s *server) Start() error {

	s.Debug("starting server on socket %s...", s.options.Socket)
	go func() {
		s.server.Serve(s.listener)
	}()

	s.Debug("waiting for server to become ready...")
	if err := utils.WaitForServer(s.options.Socket, time.Second); err != nil {
		return serverError("starting CRI server failed: %v", err)
	}

	return nil
}

// Stop serving CRI requests.
func (s *server) Stop() {
	s.Debug("stopping server on socket %s...", s.options.Socket)
	s.server.Stop()
}

// createGrpcServer creates a gRPC server instance on our socket.
func (s *server) createGrpcServer() error {
	if s.server != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.options.Socket), sockets.DirPermissions); err != nil {
		return serverError("failed to create directory for socket %s: %v",
			s.options.Socket, err)
	}

	l, err := net.Listen("unix", s.options.Socket)
	if err != nil {
		if ls, lsErr := utils.IsListeningSocket(s.options.Socket); ls || lsErr != nil {
			return serverError("failed to create server: socket %q already exists",
				s.options.Socket)
		}
		s.Warn("removing abandoned socket %q...", s.options.Socket)
		os.Remove(s.options.Socket)
		l, err = net.Listen("unix", s.options.Socket)
		if err != nil {
			return serverError("failed to create server on socket %s: %v",
				s.options.Socket, err)
		}
	}

	s.listener = l

	if s.options.User >= 0 {
		if err := s.Chown(s.options.User, s.options.Group); err != nil {
			l.Close()
			s.listener = nil
			return err
		}
	}

	if s.options.Mode != 0 {
		if err := s.Chmod(s.options.Mode); err != nil {
			l.Close()
			s.listener = nil
			return err
		}
	}

	s.server = grpc.NewServer(instrumentation.InjectGrpcServerTrace()...)

	return nil
}

// Chmod changes the permissions of the server's socket.
func (s *server) Chmod(mode os.FileMode) error {
	if s.listener != nil {
		if err := os.Chmod(s.options.Socket, mode); err != nil {
			return serverError("failed to change permissions of socket %q to %v: %v",
				s.options.Socket, mode, err)
		}
		s.Info("changed permissions of socket %q to %v", s.options.Socket, mode)
	}

	s.options.Mode = mode

	return nil
}

// Chown changes ownership of the server's socket.
func (s *server) Chown(uid, gid int) error {
	if s.listener != nil {
		userName := strconv.FormatInt(int64(uid), 10)
		if u, err := user.LookupId(userName); u != nil && err == nil {
			userName = u.Name
		}
		groupName := strconv.FormatInt(int64(gid), 10)
		if g, err := user.LookupGroupId(groupName); g != nil && err == nil {
			groupName = g.Name
		}
		if err := os.Chown(s.options.Socket, uid, gid); err != nil {
			return serverError("failed to change ownership of socket %q to %s/%s: %v",
				s.options.Socket, userName, groupName, err)
		}
		s.Info("changed ownership of socket %q to %s/%s", s.options.Socket, userName, groupName)
	}

	s.options.User = uid
	s.options.Group = gid

	return nil
}

// collectStatistics collects (should collect) request processing statistics.
func (s *server) collectStatistics(kind, name string, start, send, recv, end time.Time) {
	if kind == "passthrough" {
		return
	}

	pre := send.Sub(start)
	server := recv.Sub(send)
	post := end.Sub(recv)

	s.Debug(" * latency for %s: preprocess: %v, CRI server: %v, postprocess: %v, total: %v",
		name, pre, server, post, pre+server+post)
}


// // qualifier pulls a qualifier for disambiguation from a CRI request message.
// func (s server) qualifier(msg interface{}) string {
// 	if fn := s.options.QualifyReqFn; fn != nil {
// 		return fn(msg)
// 	}
// 	return ""
// }

// Return a formatter server error.
func serverError(format string, args ...interface{}) error {
	return fmt.Errorf("cri/server: "+format, args...)
}
