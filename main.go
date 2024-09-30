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

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Enter Bill Information:"),
		labelEntry, priceEntry, dateEntry,
		addBillButton,
		widget.NewLabel("Enter Person Information:"),
		personNameEntry,
		salaryEntry,
		totalLabel,
		saveJSONButton,
		saveCSVButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}
