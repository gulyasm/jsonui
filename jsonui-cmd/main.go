package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rmasci/jsonui"
	"github.com/spf13/pflag"
)

func main() {
	var inFile string
	var jsonByte []byte
	var err error
	pflag.StringVarP(&inFile, "file", "f", "", "JSON File to parse")
	pflag.Parse()
	if inFile != "" {
		jsonByte, err = ioutil.ReadFile(inFile)
		if err != nil {
			fmt.Println("Could not read file", err)
			os.Exit(1)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			jsonByte, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
	if len(jsonByte) <= 1 {
		fmt.Println("No Data")
		os.Exit(1)
	}
	jp := jsonui.Interactive(jsonByte)
	fmt.Println("JSON Path:", jp)
}
