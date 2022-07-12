package cmd

import (
	"bytes"
	"csv-splitter/splitter"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "csv-splitter",
	Short: "Splits a CSV file into smaller files",
	Long:  "Splits a CSV file into small files for each of a unique column value.",
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceCsv, err := os.Open(csvFile)
		if err != nil {
			return err
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(sourceCsv)

		// output filename
		sourceExtension := filepath.Ext(csvFile)
		sourceBasename := filepath.Base(csvFile)
		outputBasename := strings.TrimSuffix(sourceBasename, sourceExtension) + "_{suffix}" + sourceExtension
		outputFilenameTemplate := outputDirectory + string(os.PathSeparator) + outputBasename

		// column separator
		separator := '0'
		if len(csvSeparator) > 0 {
			separator = bytes.Runes([]byte(csvSeparator))[0]
		}

		// splitter columns
		columns := strings.Split(csvColumns, string(separator))

		if err := splitter.CsvSplit(sourceCsv, separator, columns, csvPreserveColumn, outputFilenameTemplate); err != nil {
			return err
		}

		return nil
	},
}

var csvFile string
var csvSeparator string
var csvColumns string
var csvPreserveColumn bool
var outputDirectory string

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&csvFile,
		"file",
		"f",
		"",
		"Source CSV/TSV file.",
	)
	_ = rootCmd.MarkPersistentFlagRequired("file")
	_ = rootCmd.MarkPersistentFlagFilename("file", "csv", "tsv")

	defaultOutputDirectory, _ := os.Getwd()

	rootCmd.PersistentFlags().StringVarP(
		&outputDirectory,
		"output-dir",
		"d",
		defaultOutputDirectory,
		"Output directory to store output files.",
	)
	_ = rootCmd.MarkPersistentFlagDirname("output-dir")

	rootCmd.PersistentFlags().StringVarP(
		&csvColumns,
		"columns",
		"c",
		"",
		"Columns which values will be used to split file. Multiple columns can be specified using same --separator value.",
	)
	_ = rootCmd.MarkPersistentFlagRequired("column")

	rootCmd.PersistentFlags().StringVarP(
		&csvSeparator,
		"separator",
		"s",
		",",
		"Source CSV/TSV column separator.",
	)
	rootCmd.PersistentFlags().BoolVar(
		&csvPreserveColumn,
		"preserve-columns",
		false,
		"Weather the split columns should be written to the output files.",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
