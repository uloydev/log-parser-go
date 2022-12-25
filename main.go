package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	input := kingpin.Arg("", "input file").Required().String()

	outputType := kingpin.Flag("t", "output type [json or text] ex: -t json").Short('t').Required().String()
	output := kingpin.Flag("o", "output parsed log ex: -o /home/uloydev/parsed.json").Short('o').Default("json").String()

	kingpin.Parse()

	finput, err := os.Open(*input)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	defer finput.Close()

	foutput, err := os.Create(*output)

	if err != nil {
		panic(err)
	}
	defer foutput.Close()

	fileScanner := bufio.NewScanner(finput)

	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		logLine := fileScanner.Text()
		if err == io.EOF {
			break
		}
		// your defined log format
		logsFormat := `$ip -$_- \[$time_stamp\] \"$http_method $request_path $http_version\" $status_code $_ \"$_\" \"$user_agent\"`

		// transform all the defined variable into a regex-readable named format
		regexFormat := regexp.MustCompile(`\$([\w_]*)`).ReplaceAllString(logsFormat, `(?P<$1>.*)`)

		// compile the result
		re := regexp.MustCompile(regexFormat)

		// find all the matched data from the logsExample
		matches := re.FindStringSubmatch(logLine)

		if *outputType == "json" {
			logMap := map[string]string{}

			for index, key := range re.SubexpNames() {
				// ignore the first and the $_
				if index == 0 || key == "_" {
					continue
				}

				logMap[key] = matches[index]
			}

			jsonParsed, err := json.Marshal(logMap)

			if err != nil {
				panic(err)
			}

			if *output != "" {

				_, err = foutput.WriteString(string(jsonParsed) + "\n")
				if err != nil {
					panic(err)
				}

			} else {
				fmt.Println(string(jsonParsed))
			}

		} else if *outputType == "text" {
			for index, key := range re.SubexpNames() {
				// ignore the first and the $_
				if index == 0 || key == "_" {
					continue
				}

				// print the defined variable
				data := fmt.Sprintf("%-15s => %s\n", key, matches[index])
				if *output != "" {

					_, err = foutput.WriteString(data)
					if err != nil {
						panic(err)
					}

				} else {
					fmt.Println(data)
				}
			}
		} else {
			panic("Invalid output type flag must be json or text. ex : -t json ")
		}

	}
}
