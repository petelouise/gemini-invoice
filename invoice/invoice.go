package invoice

import (
	"fmt"
	"time"

	"gemini-invoice/examples/pdf"
)

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
	p := pdf.New()

	pdf.WriteLogo(p, inv.Logo, inv.From)
	pdf.WriteTitle(p, inv.Title, inv.Id, inv.Date)
	pdf.WriteBillTo(p, inv.To)
	pdf.WriteHeaderRow(p)
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

		pdf.WriteRow(p, inv.Items[i], q, r)
		subtotal += float64(q) * r
	}
	if inv.Note != "" {
		pdf.WriteNotes(p, inv.Note)
	}
	pdf.WriteTotals(p, subtotal, subtotal*inv.Tax, subtotal*inv.Discount)
	if inv.Due != "" {
		pdf.WriteDueDate(p, inv.Due)
	}
	pdf.WriteFooter(p, inv.Id)

	err := p.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}
