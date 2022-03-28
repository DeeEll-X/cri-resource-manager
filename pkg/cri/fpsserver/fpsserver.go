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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"github.com/intel/cri-resource-manager/pkg/cri/resource-manager/sockets"
	logger "github.com/intel/cri-resource-manager/pkg/log"
	"github.com/intel/cri-resource-manager/pkg/utils"

	// "github.com/intel/cri-resource-manager/pkg/instrumentation"
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

type fpsData struct {
	cpus		[]int
	fps			float64
	game 		string
	isFpsDrop	bool
	keyRunning	float64
	keyRunnable	float64
	keyThread		int
	microTimeStamp	time.Time
	podSandboxId	string
	totalRunning	float64
	totalRunnable	float64
	totalThread		int
}

type FPSDropHandler interface{
	HandleFPSDrop(string) error
	RecordFPSData(string, float64, float64) error
}

// Server is the interface we expose for controlling our CRI server.
type Server interface {
	// RegisterFPSDropHandler registers the given FPSDropHandler with the server.
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

func(s *server) RecordFPSData(ctx context.Context, request *FPSStatistic) (*FPSStatisticReply, error) {
	s.Info("receive FPS statistics message, fps %f, schedule time %f", request.Fps, request.SchedTime)
	reply := FPSStatisticReply{
	}
	(*s.fpsDropHandler).RecordFPSData(request.PodSandboxId, request.Fps, request.SchedTime)
	if request.IsFpsDrop {
		s.Info("podSandbox %s fps drop", request.PodSandboxId)
		(*s.fpsDropHandler).HandleFPSDrop(request.PodSandboxId)
	}
	return &reply, nil
}


func (s *server) RegisterFPSDropHandler(handler FPSDropHandler) error {
	s.fpsDropHandler = &handler

	return nil
}


// Start starts the servers request processing goroutine.
func (s *server) Start() error {
	if err := s.createFPSServer(); err != nil {
		return serverError("starting CRI server failed: %v", err)
	}
	
	s.Debug("starting server on socket %s...", s.options.Socket)
	go func() {
		for {
			c, err := s.listener.Accept()
			if err != nil {
				s.Info("Accept: " + err.Error())
			} else {
				go s.HandleServerConn(c)
			}
		}
	}()

	return nil
}

func(s *server) HandleServerConn(c net.Conn) {
	s.Info("### HandleServerConn start", c)
	defer c.Close()
	buf := make([]byte, 10480)
	for {
		nr, err := c.Read(buf)
		if err != nil {
			s.Error("Read: " + err.Error())
			break
		} else {
			s.parseMsg(buf[:nr])
			s.Error("msg:" + string(buf[0:nr]))
		}
	}
	s.Info("### HandleServerConn end")
}

func(s *server) parseMsg(buf []byte){
	var start int
	var end	int
	var data map[string]string
	length := len(buf)
	
	for start = 0;start <length;{
		if buf[start] != '{' {
			start = start + 1
			continue
		} else {
			for end = start; end<length; end = end + 1{
				if buf[end] == '}'{
					break
				}
			}
			if err := json.Unmarshal(buf[start:end+1], &data); err != nil {
				panic(err)
			}
			// handle data
			start = end + 1
		}
	}

}

// Stop serving CRI requests.
func (s *server) Stop() {
	s.Debug("stopping server on socket %s...", s.options.Socket)
	s.server.Stop()
}

// createFPSServer creates a FPS server instance on our socket.
func (s *server) createFPSServer() error {
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
