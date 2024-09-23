package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
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

type Config struct {
	Title string `yaml:"title"`
	From  string `yaml:"from"`
	Logo  string `yaml:"logo"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultInvoice(config *Config) Invoice {
	return Invoice{
		Id:         time.Now().Format("20060102"),
		Title:      config.Title,
		Logo:       config.Logo,
		Rates:      []float64{0},
		Quantities: []int{1},
		Items:      []string{""},
		From:       config.From,
		To:         "",
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
	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("Invoice Generator")

	inv := DefaultInvoice(config)

	idEntry := widget.NewEntry()
	idEntry.SetText(inv.Id)
	toEntry := widget.NewEntry()
	toEntry.SetPlaceHolder("Customer Name")
	itemNameEntry := widget.NewEntry()
	itemNameEntry.SetPlaceHolder("Item Name")
	itemPriceEntry := widget.NewEntry()
	itemPriceEntry.SetPlaceHolder("Item Price")
	outputEntry := widget.NewEntry()
	outputEntry.SetText("invoice.pdf")

	content := container.NewVBox(
		widget.NewLabel("Invoice Generator"),
		widget.NewLabel("ID:"),
		idEntry,
		widget.NewLabel("To:"),
		toEntry,
		widget.NewLabel("Item Name:"),
		itemNameEntry,
		widget.NewLabel("Item Price:"),
		itemPriceEntry,
		widget.NewLabel("Output filename:"),
		outputEntry,
		widget.NewButton("Generate Invoice", func() {
			inv.Id = idEntry.Text
			inv.To = toEntry.Text
			inv.Items = []string{itemNameEntry.Text}
			price, err := strconv.ParseFloat(itemPriceEntry.Text, 64)
			if err != nil {
				fmt.Println("Error parsing price:", err)
				return
			}
			inv.Rates = []float64{price}
			inv.Quantities = []int{1}
			output := outputEntry.Text
			if !strings.HasSuffix(output, ".pdf") {
				output += ".pdf"
			}
			err = GenerateInvoice(inv, output)
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
