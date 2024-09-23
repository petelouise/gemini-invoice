package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
)

var currencySymbols = map[string]string{
	"USD": "$",
	// Add other currency symbols as needed
}

//go:embed "fonts/Inter.ttf"
var interFont []byte

//go:embed "fonts/Inter-Bold.ttf"
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
	myApp.Settings().SetTheme(NewPinkTheme())
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

	var outputDir string
	outputDirButton := widget.NewButton("Select Output Directory", nil)

	title := canvas.NewText("Invoice Generator", color.NRGBA{R: 219, G: 112, B: 147, A: 255})
	title.TextSize = 24
	title.Alignment = fyne.TextAlignCenter

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idEntry},
			{Text: "To", Widget: toEntry},
			{Text: "Item Name", Widget: itemNameEntry},
			{Text: "Item Price", Widget: itemPriceEntry},
			{Text: "Output Directory", Widget: outputDirButton},
		},
	}

	generateButton := widget.NewButton("Generate Invoice", nil)
	generateButton.Importance = widget.HighImportance

	outputDirButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error selecting directory:", err)
				return
			}
			if uri == nil {
				return
			}
			outputDir = uri.Path()
			outputDirButton.SetText(outputDir)
		}, myWindow)
	}

	generateButton.OnTapped = func() {
		if outputDir == "" {
			dialog.ShowInformation("Error", "Please select an output directory", myWindow)
			return
		}

		inv.Id = idEntry.Text
		inv.To = toEntry.Text
		inv.Items = []string{itemNameEntry.Text}
		price, err := strconv.ParseFloat(itemPriceEntry.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error parsing price: %v", err), myWindow)
			return
		}
		inv.Rates = []float64{price}
		inv.Quantities = []int{1}
		
		output := filepath.Join(outputDir, "invoice.pdf")
		err = GenerateInvoice(inv, output)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error generating invoice: %v", err), myWindow)
		} else {
			dialog.ShowInformation("Success", "Invoice generated successfully!", myWindow)
		}
	}

	content := container.NewVBox(
		title,
		layout.NewSpacer(),
		form,
		layout.NewSpacer(),
		generateButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.ShowAndRun()
}
type PinkTheme struct{}

var _ fyne.Theme = (*PinkTheme)(nil)

func (m PinkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 255, G: 240, B: 245, A: 255} // Light pink background
	case theme.ColorNameButton:
		return color.NRGBA{R: 219, G: 112, B: 147, A: 255} // Medium pink for buttons
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 139, G: 0, B: 139, A: 255} // Dark pink for text
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 182, B: 193, A: 255} // Light pink for hover
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 199, G: 21, B: 133, A: 255} // Medium violet red for placeholders
	case theme.ColorNamePressed:
		return color.NRGBA{R: 199, G: 21, B: 133, A: 255} // Medium violet red for pressed state
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 219, G: 112, B: 147, A: 255} // Medium pink as primary color
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 255, G: 182, B: 193, A: 255} // Light pink for scrollbar
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m PinkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m PinkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m PinkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func NewPinkTheme() fyne.Theme {
	return &PinkTheme{}
}
