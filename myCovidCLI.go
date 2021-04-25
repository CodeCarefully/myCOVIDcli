package main

import (
	"encoding/csv"
	"fmt"
	"git.tilde.institute/kneezle/mycovidcli/renderfloat"
	"github.com/olekukonko/tablewriter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var bodydata string //where we'll store the raw data when we import it

	maxloops := 3 //how far "back" to try to get a valid file
	currentTime := time.Now() //get current time/date

	for i := 0; i < maxloops; i++ {

		strcurrentdate := currentTime.Format("01-02-2006") //reformat for URL format
		COVIDurl := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/" + strcurrentdate + ".csv"
		fmt.Println(COVIDurl)
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

	//at this point we have the CSV data from github. we now need to parse it into something recognizable.


	//fmt.Println(bodydata)

	type searchPlaces struct {
		placename string
		confirmed int
		deaths int
		recovered int
	}

	//at some point, I'll read these in from the CLI or a config
	var placesIwantToSee [6]searchPlaces
	placesIwantToSee[0].placename = "Israel"
	placesIwantToSee[1].placename = "Estonia"
	placesIwantToSee[2].placename = "US"
	placesIwantToSee[3].placename = "Italy"
	placesIwantToSee[4].placename = "Spain"
	placesIwantToSee[5].placename = "Netherlands"

	reader := csv.NewReader(strings.NewReader(bodydata))


	// Initialize variables.
	var (
		region,

		confirmed,
		deaths,
		recovered string
	)
	data := [][]string{}

	for {
		line, error := reader.Read() //read in a line
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		//this is running on a line-by-line basis now.
		//we need to determine if we're in "Province_State" mode, or "Country_Region" mode
		//israel is populated on 3, but not before.
		//others are populated on 2, then grouped on 3

		//lets get the raw data

		region = line[3] //like US, Israel, etc...
		//incident_rate = line[12] //if theres multiple in a region
		confirmed = line[7]
		deaths = line[8]
		recovered = line[9]

		//now lets get funky.
		//essentially we want to just group on region

		for i, s := range placesIwantToSee {
			if region == s.placename {
				iconfirmed, err := strconv.Atoi(confirmed) //convert confirmed to int
				if err != nil {
					// handle error
					fmt.Fprintln(os.Stderr, err)
					os.Exit(2)
				}
				ideaths, err := strconv.Atoi(deaths) //convert deaths to int
				if err != nil {
					// handle error
					fmt.Fprintln(os.Stderr, err)
					ideaths = 0
				}

				irecovered, err := strconv.Atoi(recovered) //convert recovered to int
				if err != nil {
					// handle error
					fmt.Fprintln(os.Stderr, err)
					irecovered = 0
				}

				placesIwantToSee[i].confirmed = iconfirmed + s.confirmed
				placesIwantToSee[i].deaths = ideaths + s.deaths
				placesIwantToSee[i].recovered = irecovered + s.recovered

			}
		}

	}

	for _, s := range placesIwantToSee {

			//f_confirmed, _ := strconv.ParseFloat(confirmed, 8)
			//f_deaths, _ := strconv.ParseFloat(deaths, 8)
			//f_recovered, _ := strconv.ParseFloat(recovered, 8)

			f_confirmed := float64(s.confirmed)
			f_deaths := float64(s.deaths)
			f_recovered := float64(s.recovered)


			data = append(data, []string{s.placename,
				renderfloat.RenderFloat("#,###.", f_confirmed),
				renderfloat.RenderFloat("#,###.", f_deaths),
				renderfloat.RenderFloat("#,###.", f_recovered),
				renderfloat.RenderFloat("#,###.", f_confirmed-f_recovered-f_deaths)})
			//fmt.Println(state + "   Deaths: " + deaths + " Confirmed: " + confirmed )
		}


	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Region",  "Confirmed", "Deaths", "Recovered", "Still sick"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
	}

