package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"
)

const sampleSize int = 500
const dataPath string = "data/2020-21-student-directory.csv"
const outputPath string = "data/student-sample.csv"

func loadData(name string) ([][]string, error) {
	csvfile, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(csvfile)

	r.Read()
	data, err := r.ReadAll()
	return data, err
}

func main() {
	data, err := loadData(dataPath)
	if err != nil {
		log.Fatalln("Failed to load data: ", err)
	}

	rand.Seed(time.Now().UnixNano())
	used := map[int]bool{}

	outputfile, err := os.Create(outputPath)
	defer outputfile.Close()
	defer outputfile.Sync()
	if err != nil {
		log.Fatalln("failed to create output file: ", err)
	}
	w := csv.NewWriter(outputfile)
	defer w.Flush()
	for len(used) < sampleSize {
		n := rand.Intn(len(data))
		if _, ok := used[n]; !ok {
			w.Write(data[n])
			used[n] = true
		}
	}
}
