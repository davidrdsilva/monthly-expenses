package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Bill struct
type Bill struct {
	Label              string
	Price              float64
	Date               string
	PaidInCreditCard   bool
	TotalInstallments  int
	CurrentInstallment int
}

// Person struct
type Person struct {
	FullName string
	Salary   float64
}

// BillingData struct
type BillingData struct {
	Person      Person
	MonthNumber int
	MonthName   string
	Year        int
	Bills       []Bill
	Total       float64
	Paid        bool
}

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

func listExpenses(myApp fyne.App, bills []Bill) {
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

	content := container.NewVBox(tableScroll)

	expensesTableWindow.SetFixedSize(true)
	expensesTableWindow.SetContent(content)
	expensesTableWindow.Resize(fyne.NewSize(600, 400))
	expensesTableWindow.Show()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Billing Manager")

	// Create fields for entering bill information
	labelEntry := widget.NewEntry()
	labelEntry.SetPlaceHolder("Bill Label")
	priceEntry := widget.NewEntry()
	priceEntry.SetPlaceHolder("Price")
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("Date (YYYY-MM-DD)")

	// Person input fields
	personNameEntry := widget.NewEntry()
	personNameEntry.SetPlaceHolder("Person's name")

	// Add a salary field for the person
	salaryEntry := widget.NewEntry()
	salaryEntry.SetPlaceHolder("Person Salary")

	// Person data
	person := Person{}

	// Bills list
	var bills []Bill

	// Total label to display calculated amount
	totalLabel := widget.NewLabel("Total: $0.00")

	// Add Bill Button
	addBillButton := widget.NewButton("Add Bill", func() {
		// Convert price to float64
		price, err := strconv.ParseFloat(priceEntry.Text, 64)
		if err != nil {
			fmt.Println("Invalid price format")
			return
		}

		// Create a new bill
		bill := Bill{
			Label: labelEntry.Text,
			Price: price,
			Date:  dateEntry.Text,
		}
		bills = append(bills, bill)

		// Clear input fields
		labelEntry.SetText("")
		priceEntry.SetText("")
		dateEntry.SetText("")

		// Calculate total expenses and display
		total := calculateTotal(bills)
		totalLabel.SetText(fmt.Sprintf("Total: $%.2f", total))
	})

	// Save JSON Button
	saveJSONButton := widget.NewButton("Save to JSON", func() {
		// Set person data
		person.FullName = personNameEntry.Text
		salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
		person.Salary = salary

		// Create billing data
		billingData := BillingData{
			Person: person,
			Bills:  bills,
			Total:  calculateTotal(bills),
		}
		saveToJSON(billingData)
	})

	// Save CSV Button
	saveCSVButton := widget.NewButton("Save to CSV", func() {
		// Set person salary
		salary, _ := strconv.ParseFloat(salaryEntry.Text, 64)
		person.Salary = salary

		// Create billing data
		billingData := BillingData{
			Person: person,
			Bills:  bills,
			Total:  calculateTotal(bills),
		}
		saveToCSV(billingData)
	})

	// Test it
	listExpensesButton := widget.NewButton("Show Expenses", func() {
		listExpenses(myApp, bills)
	})

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Enter Bill Information:"),
		labelEntry, priceEntry, dateEntry,
		addBillButton,
		widget.NewLabel("Enter Person Information:"),
		personNameEntry,
		salaryEntry,
		totalLabel,
		listExpensesButton,
		saveJSONButton,
		saveCSVButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}
