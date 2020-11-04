package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"
)

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

func srs(source [][]string, used map[int]bool, n int) [][]string {
	sample := make([][]string, 0)
	for len(sample) < n {
		r := rand.Intn(len(source))
		if _, ok := used[r]; !ok {
			sample = append(sample, source[r])
			used[r] = true
		}
	}
	return sample
}

func main() {
	const sampleSize int = 500
	const dataPath string = "data/2020-21-student-directory.csv"
	const outputPath string = "data/treatments"
	treatments := []string{"control", "experimental"}
	treatmentSize := sampleSize / len(treatments)

	data, err := loadData(dataPath)
	if err != nil {
		log.Fatalln("Failed to load data: ", err)
	}

	rand.Seed(time.Now().UnixNano())
	used := map[int]bool{}

	for _, treatment := range treatments {
		outputfile, err := os.Create(outputPath + "/" + treatment + ".csv")
		w := csv.NewWriter(outputfile)
		if err != nil {
			log.Fatalln("failed to create output file: ", err)
		}
		sample := srs(data, used, treatmentSize)
		w.WriteAll(sample)
		w.Flush()
		outputfile.Sync()
		outputfile.Close()
	}
	// for len(used) < sampleSize {
	// 	n := rand.Intn(len(data))
	// 	if _, ok := used[n]; !ok {
	// 		w.Write(data[n])
	// 		used[n] = true
	// 	}
	// }
}
