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
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/intel/cri-resource-manager/pkg/cri/resource-manager/sockets"
	logger "github.com/intel/cri-resource-manager/pkg/log"
	"github.com/intel/cri-resource-manager/pkg/utils"
	"google.golang.org/grpc"
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

type fpsDataMsg struct {
	Cpus		string	`json:"cpus"`
	Fps			string	`json:"fps"`
	Game 		string	`json:"game"`	
	IsFpsDrop	string	`json:"isFpsDrop"`
	KeySched	string	`json:"keySched"`
	MicroTimeStamp	string	`json:"microTimeStamp"`
	PodName		string		`json:"pod"`
	TotalSched	string		`json:"totalSched"`
}

type FpsData struct {
	Cpus		[]int
	Fps			float64
	Game 		string
	IsFpsDrop	bool
	KeyRunning	float64
	KeyRunnable	float64
	KeyThread		int
	MicroTimeStamp	time.Time
	PodName	string
	TotalRunning	float64
	TotalRunnable	float64
	TotalThread		int
}

type FPSDropHandler interface{
	RecordFPSData(FpsData) error
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
	options      Options                   // server options
	fpsDropHandler *FPSDropHandler
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
			fpsDatas := s.parseMsg(buf[:nr])
			s.handlefpsDataMsg(fpsDatas)

			s.Info("msg:" + string(buf[0:nr]))
		}
	}
	s.Info("### HandleServerConn end")
}

func(s *server) parseMsg(buf []byte) []FpsData {
	var start int
	var end	int
	var msg fpsDataMsg
	var fpsDatas []FpsData
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
			s.Info(string(buf[start:end+1]))
			if err := json.Unmarshal(buf[start:end+1], &msg); err != nil {
				panic(err)
			}
			s.Info("%#v",msg)
			start = end + 1

			data := msg.convertToFpsData()
			s.Info("%#v",data)
			fpsDatas = append(fpsDatas, data)
		
		}
	}
	return fpsDatas
}

func(rd *fpsDataMsg) convertToFpsData() FpsData {
	data := FpsData{}

	cpus := strings.Split(rd.Cpus, "-")
	startcpu,_ := strconv.Atoi(cpus[0])
	endcpu,_ := strconv.Atoi(cpus[len(cpus)-1])
	for i:= startcpu; i<=endcpu; i = i+1 {
		data.Cpus = append(data.Cpus, i)
	}

	data.Cpus = append(data.Cpus, )
	data.Fps, _ = strconv.ParseFloat(rd.Fps,64)
	data.Game = rd.Game
	data.IsFpsDrop,_ = strconv.ParseBool(rd.IsFpsDrop)

	keySched := strings.Split(rd.KeySched,",")
	data.KeyRunning,_ = strconv.ParseFloat(keySched[0],64)
	data.KeyRunnable,_ = strconv.ParseFloat(keySched[1], 64)
	data.KeyThread,_ = strconv.Atoi(keySched[2])

	data.PodName = rd.PodName
	microsecond,_ := strconv.ParseInt(rd.MicroTimeStamp,10,64)
	sec := microsecond / 1000000
	microsecond = microsecond % 1000000
	data.MicroTimeStamp = time.Unix(sec, microsecond * int64(time.Microsecond))

	totalSched := strings.Split(rd.TotalSched,",")
	data.TotalRunning,_ = strconv.ParseFloat(totalSched[0],64)
	data.TotalRunnable,_ = strconv.ParseFloat(totalSched[1],64)
	data.TotalThread,_ = strconv.Atoi(totalSched[2])

	return data
}

func(s *server) handlefpsDataMsg(fpsDatas []FpsData) {
	for _,data := range fpsDatas {
		(*s.fpsDropHandler).RecordFPSData(data)
	}
}
// Stop serving CRI requests.
func (s *server) Stop() {
	s.Debug("stopping server on socket %s...", s.options.Socket)
	// TODO
}

// createFPSServer creates a FPS server instance on our socket.
func (s *server) createFPSServer() error {
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

// Return a formatter server error.
func serverError(format string, args ...interface{}) error {
	return fmt.Errorf("cri/server: "+format, args...)
}
