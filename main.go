package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type AuthHeader struct {
	Username string
	Password string
}

func startHLSStream(rtspURL string, _ AuthHeader) error {
	cmd := exec.Command("ffmpeg",
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-c:v", "libx264",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_allow_cache", "1",
		"/var/www/html/stream/index.m3u8")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func hlsServer(c *gin.Context) {

	filePath := filepath.Join("/var/www/html/stream", c.Param("filename"))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.File(filePath)
}

func main() {
	rtspURL := "rtsp://10.10.10.10/profile2/media.smp"

	auth := AuthHeader{
		Username: "admin",
		Password: "pass",
	}

	go func() {
		if err := startHLSStream(rtspURL, auth); err != nil {
			log.Fatalf("Failed to start HLS stream: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	r.GET("/stream/*filename", hlsServer)
}
