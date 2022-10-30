package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "This is gin server",
	Long:  `ignore`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "success",
			})
		})
		r.Run()
	},
}
