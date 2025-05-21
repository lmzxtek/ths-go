package main

import (
	"fmt"
	"strings"

	"github.com/go-gota/gota/dataframe"
)

// TestReadCSV is a test function to read CSV file and print the dataframe
func TestReadCSV() {
	fmt.Println("\n >>> Start read dataframe... ")

	csvStr := `
Country,Date,Age,Amount,Id
"United States",2012-02-01,50,112.1,01234
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,17,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United Kingdom",2012-02-01,NA,18.2,12345
"United States",2012-02-01,32,321.31,54320
"United States",2012-02-01,32,321.31,54320
Spain,2012-02-01,66,555.42,00241
`
	df := dataframe.ReadCSV(strings.NewReader(csvStr))
	// fmt.Println(df.Col("Country"))
	fmt.Println(df.Types())
	fmt.Println(df.Describe())
	fmt.Println(df)
}

// TestReadJSON is a test function to read JSON file and print the dataframe
func TestReadJSON() {
	jsonStr := `[{"COL.2":1,"COL.3":3},{"COL.1":5,"COL.2":2,"COL.3":2},{"COL.1":6,"COL.2":3,"COL.3":1}]`
	df := dataframe.ReadJSON(strings.NewReader(jsonStr))
	fmt.Println(df)

}

// TestloadRecords is a test function to load records and print the dataframe
func TestloadRecords() {
	rec := [][]string{
		{"A", "B", "C", "D"},
		{"a", "4", "5.1", "true"},
		{"k", "5", "7.0", "true"},
		{"k", "4", "6.0", "true"},
		{"a", "2", "7.1", "false"},
	}
	fmt.Println(rec)

	df := dataframe.LoadRecords(rec)
	fmt.Println(df)
}

// TestloadStructs is a test function to load structs and print the dataframe
func TestloadStructs() {
	type User struct {
		Name     string
		Age      int
		Accuracy float64
	}
	users := []User{
		{"Aram", 17, 0.2},
		{"Juan", 18, 0.8},
		{"Ana", 22, 0.5},
	}
	df := dataframe.LoadStructs(users)
	fmt.Println(df)
}
