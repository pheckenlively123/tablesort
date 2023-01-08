package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
)

// Make a test roster for use in developing and testing tablesort.
// (Yes, I know this uses math/rand.  We are not doing crypto work here,
// so this use should be fine.)

func main() {

	opts := getOpt()

	ng := newNameGerator(opts.minNameLen, opts.maxNameLen)

	// Start by generating a list of family names.
	families := make([]string, 0)
	for i := 0; i < int(opts.families); i++ {
		nextFam := ng.makeName()
		families = append(families, nextFam)
	}

	// Start assembling the roster
	roster := make([][]string, 0)
	header := []string{"LastName", "FirstName", "Grade"}
	roster = append(roster, header)

	for grade := 0; grade <= int(opts.maxGrades); grade++ {

		classSize := randomNumberRange(opts.minClassSize, opts.maxClassSize)
		for classmember := 0; classmember <= int(classSize); classmember++ {
			lastName := pickLastName(families)
			firstName := ng.makeName()
			roster = append(roster, []string{
				lastName,
				firstName,
				fmt.Sprintf("%d", grade),
			})
		}
	}

	csvFile, err := os.OpenFile(opts.outPutCsv, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		// No point in continuing.
		panic("error opening output csv: " + err.Error())
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	err = csvWriter.WriteAll(roster)
	if err != nil {
		panic("error writing to csv writer: " + err.Error())
	}
}

func pickLastName(families []string) string {
	index := rand.Int63n(int64(len(families)))
	return families[index]
}

func randomNumberRange(min uint64, max uint64) uint64 {
	numLen := max - min
	diffSize := rand.Int63n(int64(numLen + 1))
	return min + uint64(diffSize)
}

type nameGenerator struct {
	alphabet []byte
	minLen   uint64
	maxLen   uint64
}

func newNameGerator(minLen uint64, maxLen uint64) *nameGenerator {

	rv := &nameGenerator{
		alphabet: make([]byte, 0),
		minLen:   minLen,
		maxLen:   maxLen,
	}

	for i := 65; i < 91; i++ {
		rv.alphabet = append(rv.alphabet, byte(i))
	}

	for i := 97; i < 123; i++ {
		rv.alphabet = append(rv.alphabet, byte(i))
	}

	return rv
}

func (n *nameGenerator) makeName() string {

	nameLen := randomNumberRange(n.minLen, n.maxLen)

	rv := make([]byte, 0)
	for i := 0; i < int(nameLen); i++ {
		charPos := rand.Int63n(int64(len(n.alphabet)))
		rv = append(rv, n.alphabet[charPos])
	}

	return string(rv)
}

type ops struct {
	families     uint64
	maxGrades    uint64
	minClassSize uint64
	maxClassSize uint64
	minNameLen   uint64
	maxNameLen   uint64
	outPutCsv    string
}

func getOpt() ops {

	families := flag.Uint64("families", 20, "Number of families for test file.")
	maxGrades := flag.Uint64("maxgrades", 8, "Number of the oldest grade (K == 0).")
	minClassSize := flag.Uint64("minclasssize", 20, "Smallest class size.")
	maxClassSize := flag.Uint64("maxclasssize", 25, "Largest class size.")
	minNameLen := flag.Uint64("minnamelen", 4, "Min name length.")
	maxNameLen := flag.Uint64("maxnamelen", 15, "Max name length.")
	outPutCsv := flag.String("outputcsv", "testroster.csv", "Name of the test roster csv.")

	flag.Parse()

	return ops{
		families:     *families,
		maxGrades:    *maxGrades,
		minClassSize: *minClassSize,
		maxClassSize: *maxClassSize,
		minNameLen:   *minNameLen,
		maxNameLen:   *maxNameLen,
		outPutCsv:    *outPutCsv,
	}
}
