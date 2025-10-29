package cmd

import (
	"errors"
	"fmt"
	"os"
	"smts/internal/cas"
	"smts/internal/creds"
	"smts/internal/pass"
	"smts/internal/pdf"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Generate and sign the attendance sheet for the current week",
	RunE: func(cmd *cobra.Command, args []string) error {
		signature, _ := cmd.Flags().GetString("signature")
		if _, err := os.Stat(signature); errors.Is(err, os.ErrNotExist) {
			return err
		}

		campus, _ := cmd.Flags().GetString("campus")
		if campus != "" {
			campusLower := strings.ToLower(campus)
			if campusLower != "brest" && campusLower != "rennes" && campusLower != "nantes" {
				return fmt.Errorf("invalid campus '%s'; must be one of: Brest, Rennes, Nantes", campus)
			}
		}

		firstName, _ := cmd.Flags().GetString("firstname")
		lastName, _ := cmd.Flags().GetString("lastname")
		if (firstName != "" && lastName == "") || (firstName == "" && lastName != "") {
			return fmt.Errorf("both --firstname and --lastname must be provided together")
		}
		if firstName != "" && strings.TrimSpace(firstName) == "" {
			return fmt.Errorf("firstname cannot be empty")
		}
		if lastName != "" && strings.TrimSpace(lastName) == "" {
			return fmt.Errorf("lastname cannot be empty")
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

		if campus == "" {
			campus = session.User.Campus
		}
		if campus == "unknown" {
			return fmt.Errorf("failed to determine campus; please provide it using the --campus flag")
		}
		campus = cases.Title(language.French).String(campus)

		if firstName == "" {
			firstName = session.User.FirstName
			lastName = session.User.LastName
		}
		if firstName == "unknown" || lastName == "unknown" {
			return fmt.Errorf("failed to determine name; please provide it using --firstname and --lastname flags")
		}
		firstName = cases.Title(language.French).String(firstName)
		lastName = strings.ToUpper(lastName)

		_, week := time.Now().ISOWeek()
		outputPath := fmt.Sprintf("%s %s – FIPA3%s – S%d.pdf", lastName, firstName, cases.Title(language.French).String(campus[0:1]), week)
		myPDF := pdf.New(outputPath)
		if err := myPDF.Generate(session.Cookies, session.URL); err != nil {
			return fmt.Errorf("failed to generate PDF: %w", err)
		}

		if err := myPDF.AddWatermark("Certifie sur l'honneur avoir été présent(e) sur les créneaux indiqués dans le planning.", 80, 80); err != nil {
			return fmt.Errorf("failed to add watermark to PDF: %w", err)
		}
		if err := myPDF.AddWatermark(fmt.Sprintf("%s %s", firstName, lastName), 100, 60); err != nil {
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
	signCmd.Flags().StringP("signature", "s", "signature.png", "path to signature image file")
	signCmd.Flags().String("campus", "", "campus (Brest, Rennes, or Nantes) (auto-detected if not provided)")
	signCmd.Flags().String("firstname", "", "first name (auto-detected if not provided, requires --lastname)")
	signCmd.Flags().String("lastname", "", "last name (auto-detected if not provided, requires --firstname)")
	RootCmd.AddCommand(signCmd)
}
