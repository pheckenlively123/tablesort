package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {

	opts := getOpts()

	rosterData, err := os.ReadFile(opts.inputCsv)
	if err != nil {
		// No point in continuing.
		panic("error reading input roster: " + err.Error())
	}
	rosterReader := strings.NewReader(string(rosterData))
	rosterCsv := csv.NewReader(rosterReader)
	roster, err := rosterCsv.ReadAll()
	if err != nil {
		// No point in continuing
		panic("error parsing roster csv input: " + err.Error())
	}

	shuffleList := make([]student, 0)
	for i, rosterRow := range roster {

		if i == 0 {
			// Ignore the header row.
			continue
		}

		if len(rosterRow) != 3 {
			// We got a mangled input csv.  Punch out.
			panic("unexpected number of columns in csv row: " + strings.Join(rosterRow, ","))
		}
		grade, err := strconv.ParseUint(rosterRow[2], 10, 64)
		if err != nil {
			panic("error converting grade to uint64: " + rosterRow[2])
		}

		shuffleList = append(shuffleList, student{
			lastName:  rosterRow[0],
			firstName: rosterRow[1],
			grade:     grade,
		})
	}

	// Randomize the roster
	rand.Shuffle(len(shuffleList), func(i, j int) {
		shuffleList[i], shuffleList[j] = shuffleList[j], shuffleList[i]
	})

	// Assign the students to tables evenly.
	tableNum := 1
	for i := 0; i < len(shuffleList); i++ {
		shuffleList[i].table = uint64(tableNum)
		tableNum++
		if tableNum > int(opts.maxTables) {
			tableNum = 1
		}
	}

	// Now that we have assigned tables, sort the students by table,
	// so the output is eaiser to use.
	sort.Slice(shuffleList, func(i, j int) bool {
		return shuffleList[i].table < shuffleList[j].table
	})

	// Write out the results, starting with the header
	tableSort := make([][]string, 0)
	tableSort = append(tableSort, []string{
		"LastName",
		"FirstName",
		"Grade",
		"Table",
	})

	for _, st := range shuffleList {
		row := []string{
			st.lastName,
			st.firstName,
			fmt.Sprintf("%d", st.grade),
			fmt.Sprintf("%d", st.table),
		}
		tableSort = append(tableSort, row)
	}

	csvFile, err := os.OpenFile(opts.outputCsv, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		// We were so close...  :-(
		panic("error opening output file: " + err.Error())
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	err = csvWriter.WriteAll(tableSort)
	if err != nil {
		// Even closer... :-(
		panic("error writing to csv writer: " + err.Error())
	}
}

type student struct {
	firstName string
	lastName  string
	grade     uint64
	table     uint64
}

type opts struct {
	maxTables uint64
	inputCsv  string
	outputCsv string
}

func getOpts() opts {

	maxTables := flag.Uint64("maxtables", 8, "Number of lunch tables")
	inputCsv := flag.String("roster", "testroster.csv", "Input roster in CSV format.")
	outputCsv := flag.String("tablelist", "tablelist.csv", "Table list in CSV format.")

	flag.Parse()

	return opts{
		maxTables: *maxTables,
		inputCsv:  *inputCsv,
		outputCsv: *outputCsv,
	}
}
