package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	luajson "layeh.com/gopher-json"
	"reflect"
)
var path = "TraitBuddy.lua"
var luaFunc = "throughtTable.lua"
var excelFileName  = "sample.xlsx"
var styleNameArray = "styleName.json"
var newExcelName = "motif_list.xlsx"

var styles []string

func main() {	
	var charID = askCharId()
	var dataArray = getMotifTable(charID)
	jsonFile, _ := ioutil.ReadFile(styleNameArray)

	_ = json.Unmarshal(jsonFile, &styles)

	processToExcel(dataArray)
	
}

func processToExcel(motifs []interface{}) {
	if motifs != nil {
		excelFile, _ := xlsx.OpenFile(excelFileName)
		
		for i, value := range motifs {
			generateNewRow(i, excelFile)
			
			// Style name
			var row *xlsx.Row = excelFile.Sheets[0].Rows[i+1]
			row.Cells[0].Value = styles[i]

			// If the motif has only a book
			if reflect.TypeOf(value) == reflect.TypeOf(true) {
				writeWholeLine(i, value, excelFile)
			} else {
				valueAsArray, _ := value.([]interface{})
				writeSpecific(i, valueAsArray, excelFile)
			}

			// Delete last row of the excel file
			if i == len(motifs)-1 {
				_ = excelFile.Sheets[0].RemoveRowAtIndex(i+2)
			}
			// Save the file each iteration
			excelFile.Save(newExcelName)
		}
	}
}

func generateNewRow(index int, excelFile *xlsx.File) {
	var sheet = excelFile.Sheets[0]
	var originalRow = sheet.Rows[index+1]
	var originalCell = originalRow.Cells[0]
	var row *xlsx.Row
	row, err := sheet.AddRowAtIndex(index + 1)

	if err != nil {
		panic(err)
	}

	var cell = row.AddCell()
	cell.Value = originalCell.Value
	cell.SetStyle(originalCell.GetStyle())
}
func writeSpecific(styleID int, values []interface{}, excelFile *xlsx.File) {

	var sheet = excelFile.Sheets[0]
	
	//16 - 17 is the index of the stylish column
	 var green = sheet.Rows[0].Cells[16]
	var red = sheet.Rows[0].Cells[17]
		// 14 is the number of different style available
		for _ , value := range values {
			var newCell = sheet.Rows[styleID+1].AddCell()
			if value == true {
				newCell.Value = ""
				newCell.SetStyle(green.GetStyle())
			} else {
				newCell.Value = ""
				newCell.SetStyle(red.GetStyle())
			}
		}
}
func writeWholeLine(styleID int, value interface{}, excelFile *xlsx.File) {

	var sheet = excelFile.Sheets[0]
	
	//16 - 17 is the index of the stylish column
	var green = sheet.Rows[0].Cells[16]
	var red = sheet.Rows[0].Cells[17]
		// 14 is the number of different style available
		for i := 1; i<15; i++ {
			var newCell = sheet.Rows[styleID+1].AddCell()
			if value == true {
				newCell.Value = ""
				newCell.SetStyle(green.GetStyle())
			} else {
				newCell.Value = ""
				newCell.SetStyle(red.GetStyle())
			}
		}
}


func askCharId() string {
	fmt.Print("Enter the character name you want to export motifs from: ")
    var charID string
	_, _ = fmt.Scanln(&charID)
	
	return charID
}
func getMotifTable(charID string) []interface{} {
	var motifArray []interface{}

	L := lua.NewState()
	luajson.Preload(L)
	defer L.Close()

	// Get traitbuddy savedvariable
	if err := L.DoFile(path); err != nil {
		panic(err)
	}
	// get the pre-done function
	if err := L.DoFile(luaFunc); err != nil {
		panic(err)
	}
	
	// Setting variable inside LUA VM
	L.SetGlobal("table", L.GetGlobal("TraitBuddySettings").(*lua.LTable))
	L.SetGlobal("charID", lua.LString(charID))
	L.SetGlobal("motifTable", L.NewTable())
	
	// Execute motif search function
	if err := L.DoString(`findCharacterTable(charID)`); err != nil {
		panic(err)
	}

	var luaMotifTable, _ = luajson.Encode(L.GetGlobal("motifTable").(*lua.LTable))

	if err := json.Unmarshal(luaMotifTable, &motifArray); err != nil {
        panic(err)
	}
	
    return motifArray
}