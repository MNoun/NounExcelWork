package main

import (
	"fmt"
	"log"
	"strings"
)
import "github.com/xuri/excelize/v2"

func main() {

	//get slices of excel data
	county_rows, state_rows, popChange_rows := loadExcelData()

	//get rid of empty strings in slices
	county_slice := sanitizeData(county_rows)
	state_slice := sanitizeData(state_rows)
	popChange_slice := sanitizeData(popChange_rows)

	//make a slice of index where county = 0
	index_list := makeIndexList(county_slice)
	index_slice := sanitizeIndex(index_list)

	//create slices of test data
	state_testSlice := make([]string, 50)
	popChange_testSclice := make([]string, 50)

	//populate slices
	for _, index := range index_slice {
		state_testSlice = append(state_testSlice, state_slice[index])
		popChange_testSclice = append(popChange_testSclice, popChange_slice[index])
	}

	//sanitize test data for display
	state_displaySlice := sanitizeData(state_testSlice)
	popChange_displaySlice := sanitizeData(popChange_testSclice)

	print("length = ", len(state_displaySlice))
	print("\n")
	print("length = ", len(popChange_displaySlice))

}

func loadExcelData() ([]string, []string, []string) {
	excelFile, err := excelize.OpenFile("countyPopChange2020-2021.xlsx")
	if err != nil {
		log.Fatalln(err)
	}
	all_rows, err := excelFile.GetRows("co-est2021-alldata") //returns all rows of excel sheet
	if err != nil {
		log.Fatalln(err)
	}

	//creating slices for needed rows
	county_rows := make([]string, 3196)
	state_rows := make([]string, 3196)
	popChange_rows := make([]string, 3196)

	for _, row := range all_rows {
		temp_string := fmt.Sprintln(row[4])
		temp_slice := strings.Split(temp_string, "\n")
		for _, s := range temp_slice {
			county_rows = append(county_rows, s) //returns slice of the county column
		}
	}
	for _, row := range all_rows {
		temp_string := fmt.Sprintln(row[5])
		temp_slice := strings.Split(temp_string, "\n")
		for _, s := range temp_slice {
			state_rows = append(state_rows, s) //returns slice of the state column
		}
	}
	for _, row := range all_rows {
		temp_string := fmt.Sprintln(row[11])
		temp_slice := strings.Split(temp_string, "\n")
		for _, s := range temp_slice {
			popChange_rows = append(popChange_rows, s) //returns slice of the popchange2021 column
		}
	}

	return county_rows, state_rows, popChange_rows
}

func makeIndexList(county_slice []string) []int {
	index_list := make([]int, 50)
	for index, row := range county_slice {
		if row == "0" {
			index_list = append(index_list, index)
		}
	}
	return index_list
}

func sanitizeData(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func sanitizeIndex(s []int) []int {
	var r []int
	for _, x := range s {
		if x != 0 {
			r = append(r, x)
		}
	}
	return r
}

func makeUIWindow() {

}
