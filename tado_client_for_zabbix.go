package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SebastiaanKlippert/go-tado"
	"os"
	"strconv"
)

type ZoneInfo struct {
	Id   int    `json:"{#TADO_ZONE_ID}"`
	Name string `json:"{#TADO_ZONE_NAME}"`
}

type ZoneState struct {
	Temp         float64 `json:"temp"`
	Humidity     float64 `json:"humidity"`
	HeatingPower float64 `json:"heating_power_percentage"`
}

func main() {
	tadoUser := flag.String("u", "dummy@example.com", "Tado username")
	tadoPasswd := flag.String("p", "this_is_dummy_password", "Tado password")

	flag.Parse()

	command := flag.Args()[0]

	client := tado.NewClient(*tadoUser, *tadoPasswd)
	me, err := client.GetMe()
	if err != nil {
		fmt.Printf("Get Me Error: %s\n", err.Error())
		os.Exit(3)
	}
	homeId := me.Homes[0].ID

	switch command {
	case "outside_temp":
		getWeather(client, homeId)
	case "list_zones":
		listZones(client, homeId)
	case "zone_status":
		zoneId, err := strconv.Atoi(flag.Args()[1])
		if err != nil {
			fmt.Printf("Parse zone id Error: %s\n", err.Error())
			os.Exit(1)
		}

		getZoneState(client, homeId, zoneId)
	default:
		fmt.Println("command not supported.")
		os.Exit(2)
	}
}

func getWeather(client *tado.Client, homeId int) {
	weather, err := client.GetWeather(&tado.GetWeatherInput{HomeID: homeId})
	if err != nil {
		fmt.Printf("Get weather Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Print(weather.OutsideTemperature.Celsius)
}

func listZones(client *tado.Client, homeId int) {
	zones, err := client.GetZones(&tado.GetZonesInput{HomeID: homeId})
	if err != nil {
		fmt.Printf("Get zones Error: %s\n", err.Error())
		os.Exit(1)
	}

	var zoneList []ZoneInfo
	for _, zone := range zones {
		zoneList = append(zoneList, ZoneInfo{Id: zone.ID, Name: zone.Name})
	}
	outputJson, err := json.Marshal(zoneList)
	if err != nil {
		fmt.Printf("Marshal zone list json error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s", outputJson)
}

func getZoneState(client *tado.Client, homeId int, zoneId int) {
	zoneDetailState, err := client.GetZoneState(&tado.GetZoneStateInput{HomeID: homeId, ZoneID: zoneId})
	if err != nil {
		fmt.Printf("Get zone state error: %s\n", err.Error())
		os.Exit(1)
	}

	zoneState := ZoneState{
		Temp:         zoneDetailState.SensorDataPoints.InsideTemperature.Celsius,
		Humidity:     zoneDetailState.SensorDataPoints.Humidity.Percentage,
		HeatingPower: zoneDetailState.ActivityDataPoints.HeatingPower.Percentage,
	}

	outputJson, err := json.Marshal(zoneState)
	if err != nil {
		fmt.Printf("Marshal zone state json error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s", outputJson)
}
