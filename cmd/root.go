package cmd

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/spf13/cobra"
)

var httpClient *http.Client
var RootCmd = &cobra.Command{
	Use:          "smts",
	Short:        "Sign Me This Shit",
	Long:         "SMTS is a tool to generate and sign attendance sheets for FIP 3A students.",
	SilenceUsage: true,
}

func init() {
	jar, _ := cookiejar.New(nil)
	httpClient = &http.Client{Jar: jar}
}
