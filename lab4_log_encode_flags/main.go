/*
CMPT 315 (Winter 2018)
Program: Lab #04: log, flag, encoding/xml, and encoding/json
Author: Nicholas M. Boers

This file implements the user interface.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	//    "encoding/json"
	//    "encoding/xml"
	"io/ioutil"

	"github.com/kr/pretty"
)

var minOnly bool
var avgOnly bool
var quiet bool
var returnType string
var logger *log.Logger

func init() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: %s [-a|-m] <address>
       %s <address> [new value]

Options:
`, path.Base(os.Args[0]), path.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.BoolVar(&minOnly, "m", false, "show the minimum value")
	flag.BoolVar(&avgOnly, "a", false, "show only average value")
	flag.BoolVar(&quiet, "q", false, "hide log messages")
	flag.StringVar(&returnType, "f", "json", "choose return type of api")

	flag.Parse()

	if quiet {
		log.New(ioutil.Discard, "Status log: ", log.LstdFlags)
	} else {
		log.New(os.Stdout, "Status log: ", log.LstdFlags)
	}
}

func main() {
	db, err := OpenDatabase()
	if err != nil {
		log.Fatalf("OpenDatabase: %v", err)
	}
	defer db.Close()

	otherArgs := flag.Args()
	if len(otherArgs) == 0 ||
		len(otherArgs) == 1 && avgOnly && minOnly ||
		len(otherArgs) == 2 && (avgOnly || minOnly) ||
		len(otherArgs) > 2 {
		flag.Usage()
		os.Exit(1)
	}
	/*
	   encoder := json.NewEncoder(os.Stdout)
	   if returnType == "xml" {
	       encoder = xml.NewEncoder(os.Stdout)
	   }
	*/
	// add wildcard characters to each side of the pattern
	otherArgs[0] = "%" + otherArgs[0] + "%"

	if minOnly {
		val, err := db.GetMinimum(otherArgs[0])
		if err != nil {
			log.Fatalf("GetMinimum: %v", err)
		}

		fmt.Println("Minimum:", val)
	} else if avgOnly {
		avg, err := db.GetAverage(otherArgs[0])
		if err != nil {
			log.Fatalf("GetAverage: %v", err)
		}

		fmt.Println("Average:", avg)
	} else if len(otherArgs) == 1 {
		assessments, err := db.GetAssessments(otherArgs[0])
		if err != nil {
			log.Fatalf("GetAssessments: %v", err)
		}

		fmt.Println("Assessments:")
		pretty.Println(assessments)
	} else {
		newValue, err := strconv.Atoi(otherArgs[1])
		if err != nil {
			log.Fatalf("Atoi (%v): %v", otherArgs[1], err)
		}

		count, err := db.UpdateAssessments(otherArgs[0], newValue)
		if err != nil {
			log.Fatalf("UpdateAssessments: %v", err)
		}

		fmt.Println("Updated records:", count)
	}

	//    encoder.Encode(assessments)//write to writer when this is called
}
