package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", nil)
	})

	r.POST("/run", func(c *gin.Context) {
		// 1. Get input data
		source := c.PostForm("source")

		// 2. Write data to file
		file, err := ioutil.TempFile(".", "run-*.go")
		defer os.Remove(file.Name())
		defer file.Close()

		if err != nil {
			c.HTML(500, "error.tmpl", gin.H{
				"Error": err,
			})
			return
		}

		fmt.Fprintf(file, source)
		if err = file.Sync(); err != nil {
			c.HTML(500, "error.tmpl", gin.H{
				"Error": err,
			})
			return
		}

		// 3. Try to run file
		cmd := exec.Command("go", "run", file.Name())
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			c.HTML(400, "index.tmpl", gin.H{
				"Source":         source,
				"Error":          err,
				"ComplilerError": &stderr,
			})
			return
		}

		// 4. Return result
		c.HTML(200, "index.tmpl", gin.H{
			"Source": source,
			"Result": &stdout,
		})
	})

	r.Run()
}
