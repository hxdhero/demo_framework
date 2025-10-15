package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"lls_api/pkg/log"
	"lls_api/pkg/rerr"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	name string
	*gin.Engine
	httpSrv *http.Server
	host    string
	port    int
}

// startVerifyUrl 用于检查服务是否正常启动
func (s *Server) startVerifyUrl() string {
	host := s.host
	if strings.TrimSpace(s.host) == "" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("http://%s:%d/startup/", host, s.port)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/startup/" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			fmt.Println("startup error :", err)
		}
		return
	}
	s.Engine.ServeHTTP(w, r)
}

func NewServer(name string, engine *gin.Engine, host string, port int) *Server {
	return &Server{
		name:   name,
		Engine: engine,
		host:   host,
		port:   port,
	}
}

func (s *Server) Start() error {
	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s,
	}

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil {
			log.DefaultContext().Fatalf("listen: %s\n", err)
		}
	}()

	// 验证服务启动成功
	time.Sleep(500 * time.Millisecond)
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(s.startVerifyUrl())
	if err != nil {
		return rerr.Wrap(fmt.Errorf("服务%s host:%s port:%d 启动验证失败: %v", s.name, s.host, s.port, err))
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.DefaultContext().ErrorErr(rerr.Wrap(err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return rerr.Wrap(fmt.Errorf("服务%s host:%s port:%d 启动验证失败", s.name, s.host, s.port))
	}
	log.DefaultContext().Infof("启动成功, 服务:%s 地址:%s", s.name, s.httpSrv.Addr)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.DefaultContext().Info("Shutting down server...")

	// 5秒后强制退出
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		log.DefaultContext().FatalErr(rerr.WrapS(err, "Server forced to shutdown"))
	}

	log.DefaultContext().Info("Server exiting")
	return nil
}

func (s *Server) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.Engine.Group("/1/g").Group(relativePath, handlers...)
}
