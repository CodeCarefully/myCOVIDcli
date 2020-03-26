package main

import (
	"bufio"
	"coronaVirus/renderfloat"
	"encoding/csv"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	currentTime := time.Now() //get current time/date
	currentTime = currentTime.AddDate(0, 0, -1) //go to yesterday, this source updates only daily
	strcurrentdate := currentTime.Format("01-02-2006") //reformat for URL format
	COVIDurl := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"+strcurrentdate+".csv"

	//fmt.Println(COVIDurl)

	filepath := "covid.csv" // path to save the CSV to.
	if err := DownloadFile(filepath, COVIDurl); err != nil {
		panic(err)
	}

	csvFile, _ := os.Open(filepath)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	//zero out all variables
	usconfirmed := 0
	usdeaths := 0
	region := ""
	state := ""
	confirmed := ""
	deaths:=""

	data := [][]string{}

	for {
		line, error := reader.Read() //read in a line
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		region = line[3]
		state = line[1]
		confirmed = line[7]
		deaths = line[8]

		if region == "US" {

			confirmed, err := strconv.Atoi(line[7]) //convert confirmed to int
			if err != nil {
				// handle error
				fmt.Println(err)
				os.Exit(2)
			}
			deaths, err := strconv.Atoi(line[8]) //convert deaths to int
			if err != nil {
				// handle error
				fmt.Println(err)
				os.Exit(2)
			}
			usconfirmed = usconfirmed + confirmed
			usdeaths = usdeaths + deaths
		}

		//create array for table output
		if state == "New York City" {
			f, _ := strconv.ParseFloat(confirmed, 8)
			data = append(data, []string{state, renderfloat.RenderFloat("#,###.", f), deaths})
			//fmt.Println(state + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}

		if region == "Israel" {
			f, _ := strconv.ParseFloat(confirmed, 8)
			data = append(data, []string{region, renderfloat.RenderFloat("#,###.", f), deaths})
			//fmt.Println(region + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}

		if region == "Italy" {
			f, _ := strconv.ParseFloat(confirmed, 8)
			d, _ := strconv.ParseFloat(deaths, 8)
			data = append(data, []string{region, renderfloat.RenderFloat("#,###.", f), renderfloat.RenderFloat("#,###.", d)})
			//fmt.Println(region + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}

		if region == "Estonia" {
			f, _ := strconv.ParseFloat(confirmed, 8)
			data = append(data, []string{region, renderfloat.RenderFloat("#,###.", f), deaths})
			//fmt.Println(region + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}
	}

		usconfirmed1, _ := strconv.ParseFloat(strconv.Itoa(usconfirmed), 8)
		usdeaths1, _ := strconv.ParseFloat(strconv.Itoa(usdeaths), 8)
		data = append(data, []string{"USA", renderfloat.RenderFloat("#,###.", usconfirmed1), renderfloat.RenderFloat("#,###.", usdeaths1) })

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Area", "Confirmed", "Deaths"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output

	}

