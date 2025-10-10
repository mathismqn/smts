package cmd

import (
	"bufio"
	"fmt"
	"os"
	"smts/internal/sso"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const service = "imt-pass"

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure your credentials for PASS",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Username : ")
		username, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		username = strings.TrimSpace(username)

		fmt.Print("Password : ")
		bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}
		password := strings.TrimSpace(string(bytePassword))

		session := sso.NewSession()
		if err := session.Login(username, password); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		if err := keyring.Set(service, "username", username); err != nil {
			return err
		}
		if err := keyring.Set(service, "password", password); err != nil {
			return err
		}

		fmt.Println("Credentials saved successfully")

		return nil
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)
}
