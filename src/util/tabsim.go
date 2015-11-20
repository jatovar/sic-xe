package util

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//Tabsim definicion del tipo
type Tabsim struct {
	table       map[string]int64 /*Tabla de simbolos alojada en  memoria*/
	Progpc      []int64          /*Arreglo para guardar el arreglo del contador del programa*/
	progpcstr   []string         /*Arreglo para pasar a hexadecimal el arreglo de contador de programa*/
	keys        []string         /*Arreglo para guardar las llaves de la tabla en orden*/
	Progobjstr  []string         /*Programa objeto en hexadecimal*/
	Errores     string           /*Cadena para guardar los errores que se vayan generando*/
	CurrentLine int64            /*Linea actual del contador de programa*/
}

//SymbT Variable donde se guarda la tabla de simbolos
var SymbT Tabsim
var lastPC int64

//GetTabSim estructura que regresa la direccion de tabsim
func GetTabSim() *Tabsim {
	return &SymbT
}

//Allocate hace una tabla de forma dinamica
func (p *Tabsim) Allocate() {
	//if p.table == nil {
	p.Progpc = make([]int64, 2048)
	for i := 0; i < len(p.Progpc); i++ {
		p.Progpc[i] = -1
	}
	p.table = make(map[string]int64)
	p.progpcstr = make([]string, 2048)
	p.Errores = ""
	p.keys = make([]string, 2048)
	p.Progobjstr = make([]string, 2048)
	//}

}

//AddSymbol Agrega un nuevo simbolo a la tabla
func (p *Tabsim) AddSymbol(s string, i int64) {
	if val, ok := p.table[s]; ok {
		str := "Simbolo duplicado en la direccion: "
		str += strconv.FormatInt(val, 16)
		fmt.Print(str)
		fmt.Print(" Nombre: ")
		fmt.Print(s + "\n")
		p.Errores = p.Errores + str + " Nombre: " + s + "\n"

	} else {
		p.table[s] = i
	}
}

//GetSymbolAddrs Obtiene la direccion del simbolo
func (p *Tabsim) GetSymbolAddrs(s string) int64 {
	if val, ok := p.table[s]; ok {
		return val
	}
	fmt.Println("Simbolo: " + s + " no encontrado")
	p.Errores = p.Errores + "Simbolo: " + s + " no encontrado" + "\n"
	return 32767 //Regresa -1
}

//AddProgpc Agrega una direccion al programa
func (p *Tabsim) AddProgpc(i int64, pc int64) {
	p.CurrentLine = i
	p.Progpc[i] = pc
	lastPC = pc
}

//PrintTable imprime la tabla
func (p *Tabsim) PrintTable() {
	for key, value := range p.table {
		s := strconv.FormatInt(value, 16)
		//s = strings.ToUpper(s)

		fmt.Println("Etiqueta:", key, "Valor: 0x"+s)
	}
	fmt.Print("Map entries: ")
	fmt.Println(len(p.table))
}

//ReturnTable regresa la tabla en formato cdena para el archivo
func (p *Tabsim) ReturnTable() string {
	var cadenas string
	for key, value := range p.table {
		s := strconv.FormatInt(value, 16)
		//s = strings.ToUpper(s)
		cadenas += "Etiqueta:  "
		cadenas += string(key)
		cadenas += "  Valor  0x" + s + "\n"
	}
	//	sort.Sort(data sort.Interface)
	//fmt.Println(cadenas)
	return cadenas
}

//PrintProgpc imprime los contadores
func (p *Tabsim) PrintProgpc() {
	for i := 0; i < len(p.Progpc); i++ {
		if p.Progpc[i] != 0 {
			fmt.Println(p.Progpc[i])
		}
	}
	fmt.Print("PC entries: ")
	fmt.Println(len(p.table))
}

//DecToHex Funcion que convierte de Decimal a hexadecimal
func (p *Tabsim) DecToHex() {
	for i := 0; i < len(p.Progpc); i++ {
		if p.Progpc[i] != -1 {
			p.progpcstr[i] = strconv.FormatInt(p.Progpc[i], 16)
			p.progpcstr[i] = strings.ToUpper(p.progpcstr[i])
			for j := 0; j < len(GetAssembler().Addrs); j++ {
				if p.progpcstr[i] == strings.ToUpper(strconv.FormatInt(GetAssembler().Addrs[j], 16)) && p.Progpc[i] != -1 {
					p.Progobjstr[j+1] = strings.ToUpper(GetAssembler().Hexcode[j])
				}
			}
			p.progpcstr[i] = GetAssembler().formaTo3byte(p.progpcstr[i])
		}
	}
}

//PcToHex Funcion que convierte de Decimal a hexadecimal
func (p *Tabsim) PcToHex() {
	for i := 0; i < len(p.Progpc); i++ {
		if p.Progpc[i] != -1 {
			p.progpcstr[i] = strconv.FormatInt(p.Progpc[i], 16)
			p.progpcstr[i] = strings.ToUpper(p.progpcstr[i])
			p.progpcstr[i] = GetAssembler().formaTo3byte(p.progpcstr[i])
		}
	}
}

//GetProgPcStr Funcion que convierte de Decimal a hexadecimal
func (p *Tabsim) GetProgPcStr() []string {
	return p.progpcstr
}

//GetProgSize nos regresa el tamaño del programa
func (p *Tabsim) GetProgSize() int64 {
	//fmt.Println(lastPC - 3)
	//fmt.Println(p.Progpc[0])

	return (lastPC) - p.Progpc[0]
}

//GetProgSizeHex nos regresa el tamaño del programa en hexadecimal
func (p *Tabsim) GetProgSizeHex() string {
	var progsize string
	progsize = strconv.FormatInt(p.GetProgSize(), 16)
	//progsize = strings.ToUpper(progsize)
	return progsize
}

//SortTable nos ordena la ttabla
func (p *Tabsim) SortTable() {

	for k := range p.table {
		p.keys = append(p.keys, k)
	}
	fmt.Print("Simbolos Ordenados: ")
	sort.Strings(p.keys)
	fmt.Println(p.keys)

}
