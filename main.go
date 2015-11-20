package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"parser"
	"util"
	"gui"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("entra")
		gui.New().InitializeGUI("", nil)
	}
	for i := 1; i < len(os.Args); i++ {
		//var errorstr string

		fmt.Println("----------------------------------------")
		fmt.Println("Ruta del programa: " + os.Args[i])

		file, err := os.Open(os.Args[i])

		if err != nil {
			log.Fatalf("Error al cargar el archivo: %s", err.Error())
		}

		util.GetTabSim().Allocate()
		finfo, err := os.Stat(os.Args[i])
		name := finfo.Name()

		var extension = filepath.Ext(name)
		name = name[0 : len(name)-len(extension)]

		switch extension {
		case ".s":
			errorstr := parser.New().Parse(file, true, false)
			fmt.Println(errorstr)
			fmt.Println("El archivo fuente \"" + name + extension + "\" es de la arquitectura SIC Estandar")
			sic(os.Args[i], name)
			break
		case ".x":
			errorstr := parser.New().Parse(file, true, true)
			fmt.Println(errorstr)
			fmt.Println("El archivo fuente \"" + name + "\" es de la arquitectura SIC-XE (SIC Extendida)")
			sicXE(os.Args[i], name, errorstr)
			break
		}

	}

}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func sicXE(Arg string, name string, errorstr string) {
	f := newFile(name, ".tx")
	f2 := newFile(name, ".ox")
	var todostr string
	file, err := os.Open(Arg)
	util.GetAssembler().CreateOpCodeTable()
	//util.GetAssembler().Hexcode = append(util.GetAssembler().Hexcode, "flagNocode")
	parser.New().Parse(file, false, true)
	file.Close()

	objectCode := util.GetAssembler().ObjCodeXE()

	util.GetTabSim().PcToHex()
	//fmt.Println(util.GetTabSim().GetProgPcStr())
	var obj string
	for i := 0; i < len(objectCode); i++ {

		if len(objectCode[i]) > 10 {

			obj += objectCode[i]
		}
		if i == len(objectCode)-1 {
			obj += objectCode[i]
		}
	}
	if _, err = f2.WriteString(obj); err != nil {
		panic(err)
	}
	//fmt.Println(obj)
	tabsimlines := util.GetTabSim().ReturnTable()
	fmt.Println("----------------------------------------")
	fmt.Println("Tabla de simbolos: ")
	fmt.Println("----------------------------------------")
	fmt.Print(tabsimlines)
	progSize := util.GetTabSim().GetProgSizeHex()
	fmt.Println("----------------------------------------")
	fmt.Println("Tamaño del programa = " + progSize + " bytes")
	fmt.Println("----------------------------------------")

	lines, err := readLines(Arg)
	//fmt.Println(len(util.GetAssembler().Hexcode))
	if err == nil {
		addrsstr := util.GetTabSim().GetProgPcStr()
		for i := 0; i < len(lines); i++ {
			if addrsstr[i+1] == "" {
				addrsstr[i+1] = addrsstr[i]
				todostr += addrsstr[i] + "\t" + lines[i] + "\tError" + "\n"
			} else {
				todostr += addrsstr[i] + "\t" + lines[i] + "\t \t"
				if i != 0 && i < len(lines)-1 {
					todostr += util.GetAssembler().Hexcode[i-1] + "\n"
				} else {
					todostr += "\n"
				}
			}
		}
		//util.GetTabSim().PrintProgpc()
		//fmt.Print(todostr)
		addrsstr = nil
	}
	defer f.Close()

	errorstr += todostr
	errorstr += tabsimlines

	if _, err = f.WriteString(errorstr); err != nil {
		panic(err)
	}
}

func sic(Arg string, name string) {
	fmt.Println("----------------------------------------")
	util.GetAssembler().CreateOpCodeTable()
	file, err := os.Open(Arg)
	errorstr := parser.New().Parse(file, false, false)
	file.Close()
	//util.GetAssembler().PrintCode()
	util.GetAssembler().FormatCodeToHEX()
	objectCode := util.GetAssembler().ObjCode()
	//fmt.Println(errorstr)

	f := newFile(name, ".ts")
	f2 := newFile(name, ".os")

	var obj string
	for i := 0; i < len(objectCode); i++ {

		if len(objectCode[i]) > 10 {

			obj += objectCode[i]
		}
		if i == len(objectCode)-1 {
			obj += objectCode[i]
		}
	}
	if _, err = f2.WriteString(obj); err != nil {
		panic(err)
	}
	//	fmt.Println("¡Análisis léxico/sintactico terminado!")
	//util.GetTabSim().PrintTable()
	//util.GetTabSim().SortTable()
	util.GetTabSim().DecToHex()
	tabsimlines := util.GetTabSim().ReturnTable()
	fmt.Println("----------------------------------------")
	fmt.Println("Tabla de simbolos: ")
	fmt.Println("----------------------------------------")
	fmt.Print(tabsimlines)
	progSize := util.GetTabSim().GetProgSizeHex()
	fmt.Println("----------------------------------------")
	fmt.Println("Tamaño del programa = " + progSize + " bytes")
	fmt.Println("----------------------------------------")

	var todostr string
	lines, err := readLines(Arg)
	if err == nil {
		addrsstr := util.GetTabSim().GetProgPcStr()
		for i := 0; i < len(lines); i++ {
			if addrsstr[i+1] == "" {
				addrsstr[i+1] = addrsstr[i]
				todostr += addrsstr[i] + "\t" + lines[i] + "\tError" + "\n"
			} else {

				todostr += addrsstr[i] + "\t" + lines[i] + "\t" + util.GetTabSim().Progobjstr[i] + "\n"
			}
		}
		//util.GetTabSim().PrintProgpc()
		//fmt.Print(todostr)
		addrsstr = nil
	}
	defer f.Close()
	defer f2.Close()
	errorstr += todostr
	errorstr += tabsimlines

	if _, err = f.WriteString(errorstr); err != nil {
		panic(err)
	}

}

func newFile(name string, extension string) *os.File {
	f, err := os.OpenFile(name+extension, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	return f
}
