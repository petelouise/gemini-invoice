package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
)

var currencySymbols = map[string]string{
	"USD": "$",
	// Add other currency symbols as needed
}

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

func GenerateInvoice(invoice Invoice, output string) error {
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

	WriteLogo(&pdf, invoice.Logo, invoice.From)
	WriteTitle(&pdf, invoice.Title, invoice.Id, invoice.Date)
	WriteBillTo(&pdf, invoice.To)
	WriteHeaderRow(&pdf)
	subtotal := 0.0

	for i := range invoice.Items {
		q := 1
		if len(invoice.Quantities) > i {
			q = invoice.Quantities[i]
		}

		r := 0.0
		if len(invoice.Rates) > i {
			r = invoice.Rates[i]
		}

		WriteRow(&pdf, invoice.Items[i], q, r, invoice.Currency)
		subtotal += float64(q) * r
	}

	if invoice.Note != "" {
		WriteNotes(&pdf, invoice.Note)
	}
	WriteTotals(&pdf, subtotal, subtotal*invoice.Tax, subtotal*invoice.Discount, invoice.Currency)
	if invoice.Due != "" {
		WriteDueDate(&pdf, invoice.Due)
	}
	WriteFooter(&pdf, invoice.Id)
	output = strings.TrimSuffix(output, ".pdf") + ".pdf"
	err = pdf.WritePdf(output)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s\n", output)

	return nil
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Invoice Generator")

	inv := DefaultInvoice()

	idEntry := widget.NewEntry()
	idEntry.SetText(inv.Id)
	titleEntry := widget.NewEntry()
	titleEntry.SetText(inv.Title)
	fromEntry := widget.NewEntry()
	fromEntry.SetText(inv.From)
	toEntry := widget.NewEntry()
	toEntry.SetText(inv.To)
	itemEntry := widget.NewEntry()
	itemEntry.SetText(strings.Join(inv.Items, ", "))
	outputEntry := widget.NewEntry()
	outputEntry.SetText("invoice.pdf")

	content := container.NewVBox(
		widget.NewLabel("Invoice Generator"),
		widget.NewLabel("ID:"),
		idEntry,
		widget.NewLabel("Title:"),
		titleEntry,
		widget.NewLabel("From:"),
		fromEntry,
		widget.NewLabel("To:"),
		toEntry,
		widget.NewLabel("Items (comma-separated):"),
		itemEntry,
		widget.NewLabel("Output filename:"),
		outputEntry,
		widget.NewButton("Generate Invoice", func() {
			inv.Id = idEntry.Text
			inv.Title = titleEntry.Text
			inv.From = fromEntry.Text
			inv.To = toEntry.Text
			inv.Items = strings.Split(itemEntry.Text, ",")
			for i := range inv.Items {
				inv.Items[i] = strings.TrimSpace(inv.Items[i])
			}
			output := outputEntry.Text
			if !strings.HasSuffix(output, ".pdf") {
				output += ".pdf"
			}
			err := GenerateInvoice(inv, output)
			if err != nil {
				fmt.Println("Error generating invoice:", err)
			} else {
				fmt.Println("Invoice generated successfully!")
			}
		}),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.ShowAndRun()
}
