package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Remove ANSI escape codes from strings
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var lastTime time.Time = time.Now().Add(time.Hour * -1)
var re = regexp.MustCompile(ansi)

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
		sendResponse(c, http.StatusBadRequest, "access too frequent", "")
		return
	}

	service := c.Param("service")
	branch := c.Param("branch")

	projectDir := getDirectory(service)

	if projectDir == "" {
		sendResponse(c, http.StatusBadRequest, "the service does not exist", "")
		return
	}

	if branch == "/" || branch == "" {
		sendResponse(c, http.StatusBadRequest, "invalid branch name", "")
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
		sendResponse(c, http.StatusBadRequest, "the service does not exist", "")
		return
	}

	dockerComposeFile := viper.GetString("docker-compose-file")

	cmd := exec.Command("docker-compose", "-f", dockerComposeFile, "logs", service)
	runCommand(c, cmd)
}

func runCommand(c *gin.Context, cmd *exec.Cmd) {
	log.Println(cmd.String())

	var stderr, stdout bytes.Buffer

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()

	if err != nil {
		escapedErr := re.ReplaceAllString(stderr.String(), "")
		msg := fmt.Sprintf("error: %s; output: %s", err.Error(), escapedErr)
		log.Println(msg)
		sendResponse(c, http.StatusBadRequest, err.Error(), escapedErr)
		return
	}

	escapedOut := re.ReplaceAllString(stdout.String(), "")

	sendResponse(c, http.StatusOK, "success", escapedOut)
}

func getDirectory(serviceName string) string {
	projects := viper.GetStringMapString("projects")
	dir, ok := projects[serviceName]
	if !ok {
		return ""
	}
	return dir
}

func sendResponse(c *gin.Context, code int, msg string, output string) {
	data := map[string]interface{}{
		"message": msg,
	}
	if output != "" {
		data["output"] = output
	}
	c.JSON(code, data)
}
