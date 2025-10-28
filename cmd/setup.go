package cmd

import (
	"bufio"
	"fmt"
	"os"
	"smts/internal/cas"
	"smts/internal/creds"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure CAS credentials to access PASS",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSpace(username)

		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}
		password := strings.TrimSpace(string(bytePassword))

		casClient := cas.NewClient(httpClient)
		if err := casClient.Login(username, password); err != nil {
			return fmt.Errorf("failed to login to CAS: %w", err)
		}

		credentials := creds.New(username, password)
		if err := credentials.Save(); err != nil {
			return fmt.Errorf("failed to save credentials: %w", err)
		}

		fmt.Println("Credentials saved successfully")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
