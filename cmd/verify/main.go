package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Paul-weqe/orac-data/models"

	"github.com/Paul-weqe/orac-verify/internal/encoder"
	"github.com/Paul-weqe/orac-verify/internal/solver"
)

func main() {
	var (
		dataPath = flag.String("data", "", "path to normalized interface data JSON")
		mtuMin   = flag.Int("mtu-min", 0, "minimum MTU every interface must have")
	)
	flag.Parse()

	if *dataPath == "" {
		log.Fatal("--data is required")
	}
	if *mtuMin <= 0 {
		log.Fatal("--mtu-min must be positive")
	}

	data, err := loadInterfaceData(*dataPath)
	if err != nil {
		log.Fatalf("load: %v", err)
	}

	constraint := encoder.MTUConstraint{Min: *mtuMin}
	rule, err := encoder.EncodeMTU(data, constraint)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}

	result, err := solver.Check(rule)
	if err != nil {
		log.Fatalf("solve: %v", err)
	}

	if result.Satisfiable {
		fmt.Printf("VIOLATION: at least one interface has MTU < %d\n", *mtuMin)
		for name, value := range result.Model {
			fmt.Printf("  %s = %s\n", name, value)
		}
		os.Exit(1)
	}

	fmt.Printf("OK: all interfaces have MTU >= %d\n", *mtuMin)
}

func loadInterfaceData(path string) (*models.InterfaceData, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data models.InterfaceData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &data, nil
}
