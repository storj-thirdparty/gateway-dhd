package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"storj.io/uplink"
)

// port represents the port no. of REST server to be setup
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the REST server",
	Long:  `start the REST server`,
	Run:   startServer,
}

func init() {

	// Setup the start command with its flags.
	rootCmd.AddCommand(startCmd)
	var defaultPortNo string
	startCmd.Flags().StringVarP(&defaultPortNo, "port", "p", "8080", "Port number of the REST server")

}

func startServer(cmd *cobra.Command, args []string) {

	port, _ := cmd.Flags().GetString("port")

	// Set gin to Release Mode before initializing the gin router
	gin.SetMode(gin.ReleaseMode)

	// Initialize the gin router
	router := gin.Default()

	// GET request
	router.GET("/:bucketName/*objectPath", func(c *gin.Context) {

		// Create a custom header
		c.Header("x-storj-metadata", "")

		// Get bucket name and object path from GET request
		bucketName := c.Param("bucketName")
		objectPathWithSlash := c.Param("objectPath")
		objectPath := objectPathWithSlash[1:]

		// Get serialized access key from the Authorization Header entered
		reqToken := c.GetHeader("Authorization")
		if reqToken == "" {
			c.Status(http.StatusBadRequest)
			log.Print("Entered Authorization Header is in improper format")
			return
		}

		index := strings.Index(reqToken, "Storj ")
		size := len("Storj ")

		if index == -1 {
			c.Status(http.StatusBadRequest)
			log.Print("Entered Authorization Header is in improper format")
			return
		}

		serializedKey := reqToken[index+size:]

		access, err := uplink.ParseAccess(serializedKey)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			log.Print(err)
			return
		}
		// Access handle obtained from Storj uplink

		project, err := uplink.OpenProject(context.Background(), access)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Print(err)
			return
		}
		// Project handle obtained from Storj uplink

		obj, err := project.StatObject(context.Background(), bucketName, objectPath)
		if err != nil {
			c.Status(http.StatusNotFound)
			log.Print(err)
			return
		}
		// Stat information of the object obtained from Storj uplink

		metadataObject, err := json.Marshal(obj.System)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Print(err)
			return
		}
		c.Writer.Header().Add("x-storj-metadata", string(metadataObject))

		// Downloading the object
		download, err := project.DownloadObject(context.Background(), bucketName, objectPath, nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Print(err)
			return
		}

		_, err = io.Copy(c.Writer, download)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Print(err)
			return
		}

		err = download.Close()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			log.Print(err)
			return
		}
		
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%", objectPath))
		c.Writer.Header().Add("Content-Type", c.GetHeader("Content-Type"))
		// Object download completed & closed successfully
	})

	// Listen and serve
	router.Run(":" + port)

}