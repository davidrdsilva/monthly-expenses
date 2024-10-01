package main

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
