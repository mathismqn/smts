package cmd

import (
	"errors"
	"fmt"
	"os"
	"smts/internal/auth"
	"smts/internal/pdf"
	"time"

	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Generate and sign the attendance sheet for the current week",
	RunE: func(cmd *cobra.Command, args []string) error {
		signature, _ := cmd.Flags().GetString("signature")
		if _, err := os.Stat(signature); errors.Is(err, os.ErrNotExist) {
			return err
		}

		username, err := keyring.Get(service, "username")
		if err != nil {
			return fmt.Errorf("failed to get username: %w", err)
		}
		password, err := keyring.Get(service, "password")
		if err != nil {
			return fmt.Errorf("failed to get password: %w", err)
		}

		session := auth.NewSession()
		if err := session.Login(username, password); err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}

		cookies, reqURL, err := session.GetAgendaSession()
		if err != nil {
			return fmt.Errorf("failed to get agenda: %w", err)
		}

		_, week := time.Now().ISOWeek()
		outputPath := fmt.Sprintf("%s %s – FIPA%d%s – S%d.pdf", session.User.LastName, session.User.FirstName, session.User.Year, session.User.Campus[0:1], week)
		myPDF := pdf.New(outputPath)
		if err := myPDF.Generate(cookies, reqURL); err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}

		if err := myPDF.AddWatermark("Certifie sur l’honneur avoir été présent(e) sur les créneaux indiqués dans le planning.", 80, 80); err != nil {
			return fmt.Errorf("failed to sign PDF: %w", err)
		}
		if err := myPDF.AddWatermark(fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName), 100, 60); err != nil {
			return fmt.Errorf("failed to sign PDF: %w", err)
		}
		if err := myPDF.AddSignature(signature); err != nil {
			return fmt.Errorf("failed to sign PDF: %w", err)
		}

		fmt.Printf("Signed attendance sheet generated: %s\n", outputPath)

		return nil
	},
}

func init() {
	signCmd.Flags().StringP("signature", "s", "signature.png", "Path to the signature image")
	RootCmd.AddCommand(signCmd)
}
