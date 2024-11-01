package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Function to calculate total bills
func calculateTotal(bills []Bill) float64 {
	var total float64
	for _, bill := range bills {
		total += bill.Price
	}
	return total
}

// Function to save data to JSON
func saveToJSON(data BillingData) {
	file, err := os.Create("billing_data.json")
	if err != nil {
		log.Fatal("Cannot create JSON file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		log.Fatal("Cannot write JSON to file:", err)
	}
}

// Function to save data to CSV
func saveToCSV(data BillingData) {
	file, err := os.Create("billing_data.csv")
	if err != nil {
		log.Fatal("Cannot create CSV file:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Label", "Price", "Date", "PaidInCreditCard", "TotalInstallments", "CurrentInstallment"})
	for _, bill := range data.Bills {
		writer.Write([]string{
			bill.Label,
			strconv.FormatFloat(bill.Price, 'f', 2, 64),
			bill.Date,
			strconv.FormatBool(bill.PaidInCreditCard),
			strconv.Itoa(bill.TotalInstallments),
			strconv.Itoa(bill.CurrentInstallment),
		})
	}
}

// Prevent text from overflowing into adjacent cells
func truncateText(text string, length int) string {
	if len(text) > length {
		return text[:length] + "..." // Truncate with ellipsis
	}
	return text
}

func listExpenses(myApp fyne.App, bills []Bill, salaryEntry *widget.Entry) {
	expensesTableWindow := myApp.NewWindow("All expenses")

	// Create a table to display bills
	expensesTable := widget.NewTable(
		// Define the size of the table based on the bills slice
		func() (int, int) {
			return len(bills) + 1, 4 // +1 for header row
		},
		// Create each cell (headers in row 0)
		func() fyne.CanvasObject {
			return widget.NewLabel("Paid In Credit Card")
		},
		// Populate the table
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			if id.Row == 0 { // Header row
				switch id.Col {
				case 0:
					label.SetText("Label")
				case 1:
					label.SetText("Price")
				case 2:
					label.SetText("Date")
				case 3:
					label.SetText("Paid In Credit Card")
				}
			} else { // Data rows
				bill := bills[id.Row-1] // Adjust for header
				switch id.Col {
				case 0:
					label.SetText(truncateText(bill.Label, 15))
				case 1:
					label.SetText(fmt.Sprintf("$%.2f", bill.Price))
				case 2:
					label.SetText(bill.Date)
				case 3:
					label.SetText(strconv.FormatBool(bill.PaidInCreditCard))
				}
			}
		},
	)

	// Wrap the table in a scroll container and set a specific height
	tableScroll := container.NewScroll(expensesTable)
	tableScroll.SetMinSize(fyne.NewSize(600, 400)) // Width: 500, Height: 300

	totalExpensesLabel := canvas.NewText(fmt.Sprintf("R$ %.2f of total expenses", calculateTotal(bills)), color.RGBA{100, 0, 0, 255})
	totalExpensesLabel.TextStyle = fyne.TextStyle{Bold: true}

	salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
	totalAvailable := salary - calculateTotal(bills)

	totalAvailableLabel := canvas.NewText(fmt.Sprintf("R$ %.2f available", totalAvailable), color.RGBA{0, 100, 0, 255})
	totalAvailableLabel.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewVBox(
		tableScroll,
		totalExpensesLabel,
		totalAvailableLabel,
	)

	expensesTableWindow.SetFixedSize(true)
	expensesTableWindow.SetContent(content)
	expensesTableWindow.Resize(fyne.NewSize(600, 400))
	expensesTableWindow.Show()
}
