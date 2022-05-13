package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var lastTime time.Time = time.Now().Add(time.Hour * -1)

func main() {
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/update/:service/*branch", updateHandler)
	router.GET("/log/:service", logHandler)

	hostname := viper.GetString("hostname")
	port := viper.GetInt("port")
	router.Run(fmt.Sprintf("%s:%v", hostname, port))

}

func updateHandler(c *gin.Context) {
	threshold := viper.GetDuration("interval")

	if time.Since(lastTime) < threshold {
		c.String(http.StatusBadRequest, "access too frequent")
		return
	}

	service := c.Param("service")
	branch := c.Param("branch")

	projectDir := getDirectory(service)

	if projectDir == "" {
		c.String(http.StatusBadRequest, "the service does not exist")
		return
	}

	if branch == "/" || branch == "" {
		c.String(http.StatusBadRequest, "invalid branch name")
		return
	}

	lastTime = time.Now()

	script := viper.GetString("script")
	dockerComposeFile := viper.GetString("docker-compose-file")

	cmd := exec.Command("sh", script, projectDir, dockerComposeFile, branch[1:], service)
	runCommand(c, cmd)
}

func logHandler(c *gin.Context) {
	service := c.Param("service")

	// Just to check if the service is valid
	if dir := getDirectory(service); dir == "" {
		c.String(http.StatusBadRequest, "the service does not exist")
		return
	}

	dockerComposeFile := viper.GetString("docker-compose-file")

	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "logs", service)
	runCommand(c, cmd)
}

func runCommand(c *gin.Context, cmd *exec.Cmd) {
	log.Println(cmd.String())
	outputData, err := cmd.CombinedOutput()

	output := string(outputData)

	if err != nil {
		msg := fmt.Sprintf("error: %s; output: %s", err.Error(), output)
		log.Println(msg)
		c.String(http.StatusBadRequest, "error!\n----detail----\n%s", output)
		return
	}

	c.String(http.StatusOK, "success\n----output----\n%s", output)
}

func getDirectory(serviceName string) string {
	projects := viper.GetStringMapString("projects")
	dir, ok := projects[serviceName]
	if !ok {
		return ""
	}
	return dir
}
