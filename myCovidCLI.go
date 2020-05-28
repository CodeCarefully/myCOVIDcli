package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"myCOVIDcli/renderfloat"
)

func main() {
	maxloops := 3
	currentTime := time.Now() //get current time/date
	bodydata := ""
	for i := 0; i < maxloops; i++ {

		strcurrentdate := currentTime.Format("01-02-2006") //reformat for URL format
		COVIDurl := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/" + strcurrentdate + ".csv"
		//fmt.Println(COVIDurl)
		resp, err := http.Get(COVIDurl)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			bodydata = string(body)
			fmt.Println("Last Updated: " + strcurrentdate)
			break
		}

		currentTime = currentTime.AddDate(0, 0, -1) //go to yesterday, this source updates only daily
	}

	//fmt.Println(bodydata)
	reader := csv.NewReader(strings.NewReader(bodydata))

	//zero out all variables
	usconfirmed := 0
	usdeaths := 0
	usrecovered := 0
	region := ""
	state := ""
	confirmed := ""
	deaths := ""
	recovered := ""
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
		recovered = line[9]

		if region == "US" {

			confirmed, err := strconv.Atoi(line[7]) //convert confirmed to int
			if err != nil {
				// handle error
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
			deaths, err := strconv.Atoi(line[8]) //convert deaths to int
			if err != nil {
				// handle error
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			recovered, err := strconv.Atoi(line[9]) //convert recovered to int
			if err != nil {
				// handle error
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}

			usconfirmed = usconfirmed + confirmed
			usdeaths = usdeaths + deaths
			usrecovered = usrecovered + recovered
		}

		//create array for table output
		if state == "New York City" || region == "Israel" || region == "Estonia" {
			f_confirmed, _ := strconv.ParseFloat(confirmed, 8)
			f_deaths, _ := strconv.ParseFloat(deaths, 8)
			f_recovered, _ := strconv.ParseFloat(recovered, 8)

			data = append(data, []string{region, state,
				renderfloat.RenderFloat("#,###.", f_confirmed),
				renderfloat.RenderFloat("#,###.", f_deaths),
				renderfloat.RenderFloat("#,###.", f_recovered),
				renderfloat.RenderFloat("#,###.", f_confirmed-f_recovered-f_deaths)})
			//fmt.Println(state + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}

	}

	f_confirmed, _ := strconv.ParseFloat(strconv.Itoa(usconfirmed), 8)
	f_deaths, _ := strconv.ParseFloat(strconv.Itoa(usdeaths), 8)
	f_recovered, _ := strconv.ParseFloat(strconv.Itoa(usdeaths), 8)
	data = append(data, []string{"US", "Total",
		renderfloat.RenderFloat("#,###.", f_confirmed),
		renderfloat.RenderFloat("#,###.", f_deaths),
		renderfloat.RenderFloat("#,###.", f_recovered),
		renderfloat.RenderFloat("#,###.", f_confirmed-f_recovered-f_deaths)})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region", "Area", "Confirmed", "Deaths", "Recovered", "Still sick"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
