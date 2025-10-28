package cmd

import (
	"errors"
	"fmt"
	"os"
	"smts/internal/cas"
	"smts/internal/creds"
	"smts/internal/pass"
	"smts/internal/pdf"
	"time"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Generate and sign the attendance sheet for the current week",
	RunE: func(cmd *cobra.Command, args []string) error {
		signature, _ := cmd.Flags().GetString("signature")
		if _, err := os.Stat(signature); errors.Is(err, os.ErrNotExist) {
			return err
		}

		credentials, err := creds.Load()
		if err != nil {
			return fmt.Errorf("failed to load credentials: %w", err)
		}

		casClient := cas.NewClient(httpClient)
		if err := casClient.Login(credentials.Username, credentials.Password); err != nil {
			return fmt.Errorf("failed to login to CAS: %w", err)
		}

		passClient := pass.NewClient(httpClient)
		if err := passClient.Authenticate(); err != nil {
			return fmt.Errorf("failed to authenticate to PASS: %w", err)
		}

		session, err := passClient.GetAgendaSession()
		if err != nil {
			return fmt.Errorf("failed to get agenda: %w", err)
		}

		_, week := time.Now().ISOWeek()
		outputPath := fmt.Sprintf("%s %s – FIPA3%s – S%d.pdf", session.User.LastName, session.User.FirstName, session.User.Campus[0:1], week)
		myPDF := pdf.New(outputPath)
		if err := myPDF.Generate(session.Cookies, session.URL); err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}

		if err := myPDF.AddWatermark("Certifie sur l'honneur avoir été présent(e) sur les créneaux indiqués dans le planning.", 80, 80); err != nil {
			return fmt.Errorf("failed to add watermark to PDF: %w", err)
		}
		if err := myPDF.AddWatermark(fmt.Sprintf("%s %s", session.User.FirstName, session.User.LastName), 100, 60); err != nil {
			return fmt.Errorf("failed to add watermark to PDF: %w", err)
		}
		if err := myPDF.AddSignature(signature); err != nil {
			return fmt.Errorf("failed to add signature to PDF: %w", err)
		}

		fmt.Printf("Signed attendance sheet generated: %s\n", outputPath)

		return nil
	},
}

func init() {
	signCmd.Flags().StringP("signature", "s", "signature.png", "Path to the signature image")
	RootCmd.AddCommand(signCmd)
}
