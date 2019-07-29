// trace_schema_validator validates individual lines of an input file against a
// provided JSON-Schema for git trace2 event output.
//
// Note that this expects each object to validate to be on its own line in the
// input file (AKA JSON-Lines format). This is what Git natively writes with
// GIT_TRACE2_EVENT enabled.
//
// Traces can be collected by setting the GIT_TRACE2_EVENT environment variable
// to an absolute path and running any Git command; traces will be appended to
// the file.
//
// Traces can then be verified like so:
//   trace_schema_validator \
//     --trace2-event-file /path/to/trace/output \
//     --schema-file /path/to/schema
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// Required flags
var schemaFile = flag.String("schema-file", "", "JSON-Schema filename")
var trace2EventFile = flag.String("trace2-event-file", "", "trace2 event filename")
var progress = flag.Int("progress", 0, "Print progress message each time we have validated this many lines. --progress=0 means no messages are printed")

func main() {
	flag.Parse()
	if *schemaFile == "" || *trace2EventFile == "" {
		log.Fatal("Both --schema-file and --trace2-event-file are required.")
	}
	schemaURI, err := filepath.Abs(*schemaFile)
	if err != nil {
		log.Fatal("Can't get absolute path for schema file: ", err)
	}
	schemaURI = "file://" + schemaURI

	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		log.Fatal("Problem loading schema: ", err)
	}

	tracesFile, err := os.Open(*trace2EventFile)
	if err != nil {
		log.Fatal("Problem opening trace file: ", err)
	}
	defer tracesFile.Close()

	scanner := bufio.NewScanner(tracesFile)

	count := 0
	for ; scanner.Scan(); count++ {
		if *progress != 0 && count%*progress == 0 {
			log.Print("Validated items: ", count)
		}
		event := gojsonschema.NewStringLoader(scanner.Text())
		result, err := schema.Validate(event)
		if err != nil {
			log.Fatal(err)
		}
		if !result.Valid() {
			log.Printf("Trace event line %d is invalid: %s", count+1, scanner.Text())
			for _, desc := range result.Errors() {
				log.Print("- ", desc)
			}
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Scanning error: ", err)
	}

	log.Print("Validated events: ", count)
}
