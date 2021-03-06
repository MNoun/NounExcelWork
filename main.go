/************************
Created by: Mitchell Noun
Date created: 4/18/22
Class: COMP415 Emerging Languages
Assignment: Project 4
*************************/
package main

import (
	"embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"image/png"
	"log"
	"strconv"
	"strings"
)
import "github.com/xuri/excelize/v2"

//go:embed graphics/*
var EmbeddedAssets embed.FS

var demoApp GuiApp

func main() {

	//get slices of excel data
	county_rows, state_rows, pop_rows, popChange_rows := loadExcelData()

	//get rid of empty strings in slices
	county_slice := sanitizeData(county_rows)
	state_slice := sanitizeData(state_rows)
	pop_slice := sanitizeData(pop_rows)
	popChange_slice := sanitizeData(popChange_rows)

	//make a slice of index where county = 0
	index_list := makeIndexList(county_slice)
	index_slice := sanitizeIndex(index_list)

	//create slices of test data
	state_testSlice := make([]string, 50)
	pop_testSlice := make([]string, 50)
	popChange_testSclice := make([]string, 50)

	//populate slices
	for _, index := range index_slice {
		state_testSlice = append(state_testSlice, state_slice[index])
		pop_testSlice = append(pop_testSlice, pop_slice[index])
		popChange_testSclice = append(popChange_testSclice, popChange_slice[index])
	}

	//sanitize test data for display
	state_displaySlice := sanitizeData(state_testSlice)
	pop_displaySlice := sanitizeData(pop_testSlice)
	popChange_displaySlice := sanitizeData(popChange_testSclice)

	//population percent change slice
	percent_slice := make([]float64, 50)

	for index, og := range pop_displaySlice {
		intPopChange, err := strconv.Atoi(popChange_displaySlice[index])
		if err != nil {
			log.Fatalln(err)
		}
		intOG, err := strconv.Atoi(og)
		if err != nil {
			log.Fatalln(err)
		}
		percentChange := (intPopChange / intOG) * 100
		percent_slice = append(percent_slice, float64(percentChange))
	}

	//make a slice of States
	var sliceOfStates []State

	for index, str := range state_displaySlice {
		stateNameAndPop := fmt.Sprintf("%s: %s", str, popChange_displaySlice[index])
		percent := percent_slice[index]
		s := State{
			FullName:      stateNameAndPop,
			PercentChange: percent,
		}
		sliceOfStates = append(sliceOfStates, s)
	}

	//setup GUI
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Excel List")

	demoApp = GuiApp{AppUI: MakeUIWindow(state_displaySlice, popChange_displaySlice, sliceOfStates)}

	err := ebiten.RunGame(&demoApp)
	if err != nil {
		log.Fatalln("Error running User Interface Demo", err)
	}

}

func (g GuiApp) Update() error {
	//TODO finish me
	g.AppUI.Update()
	return nil
}

func (g GuiApp) Draw(screen *ebiten.Image) {
	//TODO finish me
	g.AppUI.Draw(screen)
}

func (g GuiApp) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

type GuiApp struct {
	AppUI *ebitenui.UI
}

func MakeUIWindow(state_displaySlice, popChange_displaySlice []string, sliceOfStates []State) (GUIhandler *ebitenui.UI) {
	background := image.NewNineSliceColor(color.Gray16{})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(background))

	resources, err := newListResources()
	if err != nil {
		log.Println(err)
	}

	dataAsGeneric1 := make([]interface{}, len(state_displaySlice))
	for position, s := range state_displaySlice {
		dataAsGeneric1[position] = s
	}

	dataAsGeneric2 := make([]interface{}, len(popChange_displaySlice))
	for position, x := range popChange_displaySlice {
		dataAsGeneric2[position] = x
	}

	listWidget := widget.NewList(
		widget.ListOpts.Entries(dataAsGeneric1),
		widget.ListOpts.Entries(dataAsGeneric2),
		widget.ListOpts.EntryLabelFunc(func(e interface{}) string {
			name := ""
			for _, state := range sliceOfStates {
				name = state.FullName
			}
			return name
		}),
		widget.ListOpts.ScrollContainerOpts(widget.ScrollContainerOpts.Image(resources.image)),
		widget.ListOpts.SliderOpts(
			widget.SliderOpts.Images(resources.track, resources.handle),
			widget.SliderOpts.HandleSize(resources.handleSize),
			widget.SliderOpts.TrackPadding(resources.trackPadding)),
		widget.ListOpts.EntryColor(resources.entry),
		widget.ListOpts.EntryFontFace(resources.face),
		widget.ListOpts.EntryTextPadding(resources.entryPadding),
		widget.ListOpts.HideHorizontalSlider(),
		widget.ListOpts.EntrySelectedHandler(func(args *widget.ListEntrySelectedEventArgs) {
			var percentFloat float64
			for _, state := range sliceOfStates {
				percentFloat = state.PercentChange
			}
			text := fmt.Sprintf("%f", percentFloat)
			percentChangeText := widget.TextOptions{}.Text(text, basicfont.Face7x13, color.White)
			textWidget := widget.NewText(percentChangeText)
			rootContainer.AddChild(textWidget)
		}))
	rootContainer.AddChild(listWidget)
	GUIhandler = &ebitenui.UI{Container: rootContainer}
	return GUIhandler
}

func loadImageNineSlice(path string, centerWidth int, centerHeight int) (*image.NineSlice, error) {
	i := loadPNGImageFromEmbedded(path)

	w, h := i.Size()
	return image.NewNineSlice(i,
			[3]int{(w - centerWidth) / 2, centerWidth, w - (w-centerWidth)/2 - centerWidth},
			[3]int{(h - centerHeight) / 2, centerHeight, h - (h-centerHeight)/2 - centerHeight}),
		nil
}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func loadExcelData() ([]string, []string, []string, []string) {
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
	pop_rows := make([]string, 3196)
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
		temp_string := fmt.Sprintln(row[8])
		temp_slice := strings.Split(temp_string, "\n")
		for _, s := range temp_slice {
			pop_rows = append(pop_rows, s) //returns slice of the state column
		}
	}
	for _, row := range all_rows {
		temp_string := fmt.Sprintln(row[11])
		temp_slice := strings.Split(temp_string, "\n")
		for _, s := range temp_slice {
			popChange_rows = append(popChange_rows, s) //returns slice of the popchange2021 column
		}
	}

	return county_rows, state_rows, pop_rows, popChange_rows
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
