package main

import (
	"encoding/json" // Import the json package
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gota/gota/dataframe" // Import dataframe
	"github.com/go-gota/gota/series"    // Import series (dataframe dependency)
)

// OHLCVData struct represents a single OHLCV data point (now primarily for internal data point generation, not as the final storage structure)
type OHLCVData struct {
	Period string // Date, format "YYYY-MM-DD"
	Open   int    // Opening price
	High   int    // Highest price
	Low    int    // Lowest price
	Last   int    // Closing price
	Volume int    // Volume
}

// generateDynamicDataFrame dynamically generates OHLCV dataset, directly returning dataframe.DataFrame
// startYear: Starting year
// startMonth: Starting month (1-12)
// numMonths: Number of months to generate
// initialOpen: Initial opening price, as a reference for the first data point's opening price
func generateDynamicDataFrame(startYear, startMonth, numMonths, initialOpen int) dataframe.DataFrame {
	// Temporary slices to store data for each column
	periods := make([]string, 0, numMonths)
	opens := make([]int, 0, numMonths)
	highs := make([]int, 0, numMonths)
	lows := make([]int, 0, numMonths)
	lasts := make([]int, 0, numMonths)
	volumes := make([]int, 0, numMonths)

	currentOpen := initialOpen // Used to track the closing price of the previous period, as a reference for the next period's opening price

	// Seed the random number generator to ensure different data each run
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < numMonths; i++ {
		// Calculate current date
		date := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC).AddDate(0, i, 0)
		period := date.Format("2006-01-02") // Format as "YYYY-MM-DD"

		// Simulate price fluctuations
		open := currentOpen + rand.Intn(21) - 10
		if open <= 0 {
			open = 1
		}

		high := open + rand.Intn(51) + 5
		if high < open {
			high = open + rand.Intn(10) + 1
		}

		low := open - rand.Intn(51) - 5
		if low < 1 {
			low = 1
		}
		if low > open {
			low = open - rand.Intn(10) - 1
			if low < 1 {
				low = 1
			}
		}

		last := low + rand.Intn(high-low+1)
		if last <= 0 {
			last = 1
		}

		if high < open {
			high = open + rand.Intn(10) + 1
		}
		if high < last {
			high = last + rand.Intn(10) + 1
		}
		if low > open {
			low = open - rand.Intn(10) - 1
			if low < 1 {
				low = 1
			}
		}
		if low > last {
			low = last - rand.Intn(10) - 1
			if low < 1 {
				low = 1
			}
		}

		volume := 500 + rand.Intn(5000)

		// Add generated data to respective column slices
		periods = append(periods, period)
		opens = append(opens, open)
		highs = append(highs, high)
		lows = append(lows, low)
		lasts = append(lasts, last)
		volumes = append(volumes, volume)

		currentOpen = last // Use the current period's closing price as the reference for the next period's opening price, forming a chained fluctuation
	}

	// Create dataframe.DataFrame using the collected column data
	df := dataframe.New(
		series.New(periods, series.String, "Period"),
		series.New(opens, series.Int, "Open"),
		series.New(highs, series.Int, "High"),
		series.New(lows, series.Int, "Low"),
		series.New(lasts, series.Int, "Last"),
		series.New(volumes, series.Int, "Volume"),
	)

	return df
}

// dfToJSONRecords converts a dataframe.DataFrame to a slice of JSON records
// Each record is a JSON object corresponding to a row in the DataFrame
func dfToJSONRecords(df dataframe.DataFrame) ([]byte, error) {
	records := make([]map[string]interface{}, 0, df.Nrow()) // Pre-allocate capacity

	colNames := df.Names() // Get all column names

	for i := 0; i < df.Nrow(); i++ { // Iterate over each row
		row := make(map[string]interface{}) // Create a map to store data for the current row
		for _, colName := range colNames {  // Iterate over each column in the current row
			s := df.Col(colName)        // Get the Series for the current column
			val, err := s.Elem(i).Val() // Get the value for the current row in this column
			if err != nil {
				// Handle error, e.g., if value is NaN or cannot be converted to a specific type
				row[colName] = nil // Or skip, or set to a default value
			} else {
				row[colName] = val
			}
		}
		records = append(records, row) // Add the current row (map) to the records slice
	}

	// Marshal the records slice to a JSON string, using Indent for pretty printing
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonData, nil
}

func main_df() {
	// Define starting date and initial opening price
	var currentYear int = 2024
	var currentMonth int = 1
	var initialOpen int = 100

	// Generate and print data every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting timed dataset generation...")

	for range ticker.C {
		// Directly generate dataframe.DataFrame
		df := generateDynamicDataFrame(currentYear, currentMonth, 15, initialOpen)

		fmt.Println("\n--- Original Dataset (printed directly from DataFrame) ---")
		fmt.Println(df) // Print the entire DataFrame
		fmt.Println("--- Original Dataset Generation Complete ---")

		// --- Operations using gota/dataframe ---
		fmt.Println("\n--- Operations using github.com/go-gota/gota/dataframe ---")

		// 1. Select columns (similar to df[['col1', 'col2']])
		fmt.Println("\nSelecting 'Period' and 'Last' columns:")
		selectedCols := df.Select([]string{"Period", "Last"})
		fmt.Println(selectedCols)

		// 2. Filter rows (similar to df[df['Volume'] > 3000])
		fmt.Println("\nFiltering rows where 'Volume' > 3000:")
		filteredDF := df.Filter(
			dataframe.F{
				Colname:    "Volume",
				Comparator: series.Greater,
				Comparando: 3000,
			},
		)
		fmt.Println(filteredDF)

		// 3. Sort (similar to df.sort_values('Volume', ascending=False))
		fmt.Println("\nSorting by 'Volume' in descending order:")
		sortedDF := df.Arrange(
			dataframe.RevSort("Volume"), // RevSort means descending order
		)
		fmt.Println(sortedDF)

		fmt.Println("\nSorting by 'Period' in ascending order:")
		sortedDFByPeriod := df.Arrange(
			dataframe.Sort("Period"), // Sort means ascending order
		)
		fmt.Println(sortedDFByPeriod)

		// 4. Get data for a specific row/column (similar to df.loc[index, 'col'])
		// Get data for a specific date by filtering (similar to df[df['Period'] == '2024-03-01'])
		targetPeriod := "2024-03-01"
		specificDayDF := df.Filter(
			dataframe.F{
				Colname:    "Period",
				Comparator: series.Eq,
				Comparando: targetPeriod,
			},
		)
		fmt.Printf("\nGetting data for date %s (using DataFrame filter):\n", targetPeriod)
		fmt.Println(specificDayDF)
		if specificDayDF.Nrow() > 0 {
			// If found, extract specific values
			// For example, get 'Last' value
			lastVal, _ := specificDayDF.Col("Last").Elem(0).Int() // Elem(0) gets the first element, Int() gets its integer value
			fmt.Printf("Closing price for %s is: %d\n", targetPeriod, lastVal)
		} else {
			fmt.Printf("No data found for date %s\n", targetPeriod)
		}

		// Example: Search for a non-existent date
		targetPeriodNotFound := "2025-01-01" // Assume we want to find data for this date
		specificDayDFNotFound := df.Filter(
			dataframe.F{
				Colname:    "Period",
				Comparator: series.Eq,
				Comparando: targetPeriodNotFound,
			},
		)
		fmt.Printf("\nGetting data for date %s (using DataFrame filter):\n", targetPeriodNotFound)
		fmt.Println(specificDayDFNotFound)
		if specificDayDFNotFound.Nrow() > 0 {
			fmt.Printf("Found data: Period: %s, Open: %d, High: %d, Low: %d, Last: %d, Volume: %d\n",
				specificDayDFNotFound.Col("Period").Elem(0).String(),
				specificDayDFNotFound.Col("Open").Elem(0).Int(),
				specificDayDFNotFound.Col("High").Elem(0).Int(),
				specificDayDFNotFound.Col("Low").Elem(0).Int(),
				specificDayDFNotFound.Col("Last").Elem(0).Int(),
				specificDayDFNotFound.Col("Volume").Elem(0).Int(),
			)
		} else {
			fmt.Printf("No data found for date %s.\n", targetPeriodNotFound)
		}

		fmt.Println("--- Gota DataFrame Operations Complete ---")

		// --- Convert DataFrame to JSON records ---
		fmt.Println("\n--- Converting DataFrame to JSON Records ---")
		jsonData, err := dfToJSONRecords(df)
		if err != nil {
			fmt.Printf("Failed to convert to JSON: %v\n", err)
		} else {
			fmt.Println(string(jsonData))
		}
		fmt.Println("--- JSON Conversion Complete ---")

		// --- Concatenating and Deduplicating DataFrames Example ---
		fmt.Println("\n--- Concatenating and Deduplicating DataFrames Example ---")

		// Generate df1 (e.g., 10 months starting from Jan 2024)
		df1 := generateDynamicDataFrame(2024, 1, 10, 100)
		fmt.Println("\ndf1 (2024-01-01 to 2024-10-01):")
		fmt.Println(df1)

		// Generate df2 (e.g., 10 months starting from Aug 2024, with overlap with df1)
		// To ensure overlap, we set df2's initialOpen to the last 'Last' value of df1
		df1LastVal, err := df1.Col("Last").Elem(df1.Nrow() - 1).Int()
		if err != nil {
			df1LastVal = 100 // fallback
		}
		df2 := generateDynamicDataFrame(2024, 8, 10, df1LastVal)
		fmt.Println("\ndf2 (2024-08-01 to 2025-05-01):")
		fmt.Println(df2)

		// 1. Concatenate df1 and df2
		// dataframe.Concat vertically concatenates two or more DataFrames
		combinedDF := dataframe.Concat([]dataframe.DataFrame{df1, df2})
		fmt.Println("\nConcatenated DataFrame (may contain duplicate data):")
		fmt.Println(combinedDF)

		// 2. Deduplicate, keeping the last occurrence (i.e., from df2 if overlapped)
		// DropDuplicates method removes duplicate rows based on specified columns (here, "Period")
		// By setting the 'keep' argument to series.Last, we ensure that if there are duplicates,
		// the last occurrence (which would be from df2 in the concatenated order) is retained.
		deduplicatedDF := combinedDF.DropDuplicates("Period", series.Last)
		fmt.Println("\nDeduplicated DataFrame (deduplicated by 'Period' column, keeping last occurrence):")
		fmt.Println(deduplicatedDF)

		// Verify row counts after deduplication
		fmt.Printf("\ndf1 rows: %d, df2 rows: %d, Combined rows: %d, Deduplicated rows: %d\n",
			df1.Nrow(), df2.Nrow(), combinedDF.Nrow(), deduplicatedDF.Nrow())

		fmt.Println("--- Concatenating and Deduplicating DataFrames Example Complete ---")

		// Update the starting date and initial opening price for the next generation cycle
		// Need to get the last closing price from the DataFrame
		if df.Nrow() > 0 {
			// Get the last element of the 'Last' column
			lastSeries := df.Col("Last")
			lastVal, err := lastSeries.Elem(lastSeries.Len() - 1).Int()
			if err == nil {
				initialOpen = lastVal
			} else {
				fmt.Printf("Error getting last value: %v\n", err)
				// If fetching fails, set a default value or handle the error in another way
				initialOpen = 100 // fallback
			}
		}

		currentMonth += 15
		if currentMonth > 12 {
			currentYear += (currentMonth - 1) / 12
			currentMonth = (currentMonth-1)%12 + 1
		}
	}
}
