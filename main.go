package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Transaction struct {
	TxDate   Date
	Amount   Money
	Summary  string
	Category string
}

var autoCategories map[string]string

// Parse a QIF file and report key details
func main() {
	fromDateS := flag.String("df", "01/01/1980", "Ignore transactions before this date")
	toDateS := flag.String("dt", "31/12/2199", "Ignore transactions after this date")
	writeData := flag.Bool("wd", false, "Write transactions as pipe-separated values")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: moneyview {args} filename.")
		fmt.Println("Catgory mappings will be loaded from automapping.txt if it exists.")
		fmt.Println("Arguments...")
		flag.PrintDefaults()
		os.Exit(-1)
	}

	dataFile := flag.Args()[0]

	log.Println("MoneyView")
	log.Println("Processing data file " + dataFile)

	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatalf("Unable to open data file %s", dataFile)
	}
	defer file.Close()

	txs := []Transaction{}
	tx := Transaction{}

	autoCategories = LoadMappings()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch line[0] {
		case 'D':
			tx.TxDate = DateFrom(line[1:])
		case 'P':
			tx.Summary = line[1:]
			tx.Category = categorise(tx.Summary)
		case 'T':
			tx.Amount = MoneyFrom(line[1:])
		case '^':
			{
				txs = append(txs, tx)
				tx = Transaction{}
			}
		}
	}

	start := DateFrom(*fromDateS)
	end := DateFrom(*toDateS)
	log.Printf("Processing transactions from %v to %v, inclusive", start, end)

	categories := make(MoneyMap)

	total := Money{}
	txCount := 0
	for _, tx := range txs {
		if tx.TxDate.InRange(start, end) {
			txCount++
			total.Add(tx.Amount)

			if tx.Category == "" {
				if tx.Amount.IsPositive() {
					categories.Add("in", tx.Amount)
				} else {
					categories.Add("out", tx.Amount)
				}
			} else {
				categories.Add(tx.Category, tx.Amount)
			}

			if *writeData {
				fmt.Printf("%v | %v | %s | %s\n", tx.TxDate, tx.Amount, tx.Category, tx.Summary)
			}
		}
	}

	log.Printf("Processed %d out of %d transactions in file", txCount, len(txs))

	for key, value := range categories {
		log.Printf("  %s = %s", key, value)
	}
}

func categorise(text string) string {
	for marker, category := range autoCategories {
		ucSummary := strings.ToUpper(text)
		if strings.Contains(ucSummary, marker) {
			return category
		}
	}
	return ""
}

func LoadMappings() map[string]string {
	autocat := make(map[string]string)

	file, err := os.Open("./automapping.txt")
	if os.IsNotExist(err) {
		log.Println("No automapping.txt file found")
		return autocat
	} else if err != nil {
		log.Fatalf("Unable to open automapping file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, ":")
		autocat[strings.ToUpper(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
	}

	log.Printf("Loaded %d automapping rules from automapping.txt\n", len(autocat))

	return autocat
}
