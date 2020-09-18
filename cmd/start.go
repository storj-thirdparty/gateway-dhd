package cmd

import (
	"context"
	"fmt"
	"io"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/:bucketName/*objectPath", func(c *gin.Context) {

		bucketName := c.Param("bucketName")
		objectPathWithSlash := c.Param("objectPath")
		objectPath := objectPathWithSlash[1:]

		reqToken := c.GetHeader("Authorization")
		
		index := strings.Index(reqToken, "Storj ")
		size := len("Storj ")

		if index == -1 {
			//Entered Authorization Header is in improper format
			return
		}
		
		serializedKey := reqToken[index+size:]

		access, err := uplink.ParseAccess(serializedKey)
		if err != nil {
			//Error in getting access handle from Storj uplink
			return
		}
		//Access handle obtained from Storj uplink

		project, err := uplink.OpenProject(context.Background(), access)
		if err != nil {
			//Error in getting project handle from Storj uplink
			return
		}
		//Project handle obtained from Storj uplink

		download, err := project.DownloadObject(context.Background(), bucketName, objectPath, nil)
		if err != nil {
			//Error in download object
			return
		}

		//Downloading the object
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%", objectPath))
		c.Writer.Header().Add("Content-Type", c.GetHeader("Content-Type"))
		_, err = io.Copy(c.Writer, download)
		if err != nil {
			//Some error occurred while downloading the object
			return
		}
		//Object downloaded successfully

		err = download.Close()
		if err != nil {
			//Error in closing the download of the object
			return
		}
		//Object download completed & closed successfully

	})

	// Listen and serve
	router.Run(":" + port)

}
