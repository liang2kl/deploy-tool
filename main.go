package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	lastTime time.Time = time.Now().Add(time.Hour * -1)
)

func main() {
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/update/:service/:branch", updateHandler)

	hostname := viper.GetString("hostname")
	port := viper.GetInt("port")
	router.Run(fmt.Sprintf("%s:%v", hostname, port))

}

func updateHandler(c *gin.Context) {
	latestTime := lastTime

	threshold := viper.GetDuration("interval")

	if time.Since(latestTime) < threshold {
		sendResponse(c, http.StatusBadRequest, "access too frequent")
		return
	}

	lastTime = time.Now()

	service := c.Param("service")
	branch := c.Param("branch")

	projectDir := getInfo(service)

	if projectDir == "" || branch == "" {
		sendResponse(c, http.StatusBadRequest, "invalid arguments")
		return
	}

	script := viper.GetString("script")
	dockerComposeDir := viper.GetString("docker-compose-file")
	// projectDir := viper.GetString("project-dir")

	cmd := exec.Command("sh", script, projectDir, dockerComposeDir, branch, service)
	log.Println(cmd.String())

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		msg := fmt.Sprintf("error: %s; output: %s", err.Error(), stderr.String())
		log.Println(msg)
		sendResponse(c, http.StatusBadRequest, msg)
		return
	}

	sendResponse(c, http.StatusOK, "success")
}

func getInfo(serviceName string) string {
	projects := viper.GetStringMapString("projects")
	dir, ok := projects[serviceName]
	if !ok {
		return ""
	}
	return dir
}

func sendResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, map[string]interface{}{
		"message": msg,
	})
}
