package invoice

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/signintech/gopdf"
)

//go:embed "Inter/Inter Variable/Inter.ttf"
var interFont []byte

//go:embed "Inter/Inter Hinted for Windows/Desktop/Inter-Bold.ttf"
var interBoldFont []byte

type Invoice struct {
	Id         string    `json:"id" yaml:"id"`
	Title      string    `json:"title" yaml:"title"`
	Logo       string    `json:"logo" yaml:"logo"`
	From       string    `json:"from" yaml:"from"`
	To         string    `json:"to" yaml:"to"`
	Date       string    `json:"date" yaml:"date"`
	Due        string    `json:"due" yaml:"due"`
	Items      []string  `json:"items" yaml:"items"`
	Quantities []int     `json:"quantities" yaml:"quantities"`
	Rates      []float64 `json:"rates" yaml:"rates"`
	Tax        float64   `json:"tax" yaml:"tax"`
	Discount   float64   `json:"discount" yaml:"discount"`
	Currency   string    `json:"currency" yaml:"currency"`
	Note       string    `json:"note" yaml:"note"`
}

func DefaultInvoice() Invoice {
	return Invoice{
		Id:         time.Now().Format("20060102"),
		Title:      "INVOICE",
		Rates:      []float64{25},
		Quantities: []int{2},
		Items:      []string{"Paper Cranes"},
		From:       "Project Folded, Inc.",
		To:         "Untitled Corporation, Inc.",
		Date:       time.Now().Format("Jan 02, 2006"),
		Due:        time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"),
		Tax:        0,
		Discount:   0,
		Currency:   "USD",
	}
}

func GenerateInvoice(inv Invoice, output string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})
	pdf.SetMargins(40, 40, 40, 40)
	pdf.AddPage()
	err := pdf.AddTTFFontData("Inter", interFont)
	if err != nil {
		return err
	}

	err = pdf.AddTTFFontData("Inter-Bold", interBoldFont)
	if err != nil {
		return err
	}

	writeLogo(&pdf, inv.Logo, inv.From)
	writeTitle(&pdf, inv.Title, inv.Id, inv.Date)
	writeBillTo(&pdf, inv.To)
	writeHeaderRow(&pdf)
	subtotal := 0.0
	for i := range inv.Items {
		q := 1
		if len(inv.Quantities) > i {
			q = inv.Quantities[i]
		}

		r := 0.0
		if len(inv.Rates) > i {
			r = inv.Rates[i]
		}

		writeRow(&pdf, inv.Items[i], q, r)
		subtotal += float64(q) * r
	}
	if inv.Note != "" {
		writeNotes(&pdf, inv.Note)
	}
	writeTotals(&pdf, subtotal, subtotal*inv.Tax, subtotal*inv.Discount)
	if inv.Due != "" {
		writeDueDate(&pdf, inv.Due)
	}
	writeFooter(&pdf, inv.Id)

	err = pdf.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}

// Add the following functions from cli_example.go:
// writeLogo, writeTitle, writeBillTo, writeHeaderRow, writeRow, writeNotes, writeTotals, writeDueDate, writeFooter

// You'll need to copy these functions from cli_example.go and adjust them to work within this package.
