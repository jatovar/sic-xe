package util

import (
	"fmt"
	"strconv"
	"strings"
)

//Assembler Estructura del ensamblador de la SIC estandar
type Assembler struct {
	table      map[string]int64 //Guarda los codigos de operacion
	code       []int64          //Guarda el codigo objeto en decimal
	hexconst   map[int][]int64  //Guarda el arreglo de bytes de una constante hexadecimal
	charconst  map[int][]int64  //Guarda el arreglo de bytes de una constante cadena
	Hexcode    []string         //Guarda el codigo transformado a hexadecimal en cadenas
	Addrs      []int64          //Va guardando el contador de programa en la direccion que va
	Progname   string           //Nombre del programa
	StartAddrs int64            //Direccion de inicio del programa
	EndAddrs   int64            //Direccion final del programa
	registers  map[string]int64 //Tabla de registros
	mReg       []string
}

//Codes Variable donde se guarda la tabla codigos de operacion
var Codes Assembler

//GetAssembler estructura que regresa la direccion de la tabla de codigos de operacion
func GetAssembler() *Assembler {
	return &Codes
}

//Allocate inicializa todas las tablas
func (p *Assembler) Allocate() {
	p.code = nil
	p.Hexcode = nil
	p.Addrs = nil
	p.EndAddrs = 0
	p.mReg = nil
	p.table = make(map[string]int64)
	p.hexconst = make(map[int][]int64)
	p.charconst = make(map[int][]int64)
	p.registers = make(map[string]int64)

}

//GetCodeEntry nos regresa un codigo de operacion segun en nemonico
func (p *Assembler) GetCodeEntry(name string) int64 {
	return p.table[name]

}

//InsertCodes inserta los codigos de operacion con su respectivo nemonico
func (p *Assembler) InsertCodes() {
	p.table["ADD"] = 0x18
	p.table["ADDF"] = 0x58
	p.table["ADDR"] = 0x90
	p.table["AND"] = 0x40
	p.table["CLEAR"] = 0xb4
	p.table["COMP"] = 0x28
	p.table["COMPF"] = 0x88
	p.table["COMPR"] = 0xa0
	p.table["DIV"] = 0x24
	p.table["DIVF"] = 0x64
	p.table["DIVR"] = 0x9c
	p.table["FIX"] = 0xc4
	p.table["FLOAT"] = 0xc0
	p.table["HIO"] = 0xf4
	p.table["J"] = 0x3c
	p.table["JEQ"] = 0x30
	p.table["JGT"] = 0x34
	p.table["JLT"] = 0x38
	p.table["JSUB"] = 0x48
	p.table["LDA"] = 0x0
	p.table["LDB"] = 0x68
	p.table["LDCH"] = 0x50
	p.table["LDF"] = 0x70
	p.table["LDL"] = 0x8
	p.table["LDS"] = 0x6c
	p.table["LDT"] = 0x74
	p.table["LDX"] = 0x4
	p.table["LPS"] = 0xd0
	p.table["MUL"] = 0x20
	p.table["MULF"] = 0x60
	p.table["MULR"] = 0x98
	p.table["NORM"] = 0xc8
	p.table["OR"] = 0x44
	p.table["RD"] = 0xd8
	p.table["RMO"] = 0xac
	p.table["RSUB"] = 0x4c
	p.table["SHIFTL"] = 0xa4
	p.table["SHIFTR"] = 0xa8
	p.table["SIO"] = 0xf0
	p.table["SSK"] = 0xec
	p.table["STA"] = 0xc
	p.table["STB"] = 0x78
	p.table["STCH"] = 0x54
	p.table["STL"] = 0x14
	p.table["STS"] = 0x7c
	p.table["STSW"] = 0xe8
	p.table["STT"] = 0x84
	p.table["STX"] = 0x10
	p.table["SUB"] = 0x1c
	p.table["SUBF"] = 0x5c
	p.table["SUBR"] = 0x94
	p.table["SVC"] = 0xb0
	p.table["TD"] = 0xe0
	p.table["TIO"] = 0xf8
	p.table["TIX"] = 0x2c
	p.table["TIXR"] = 0xb8
	p.table["WD"] = 0xdc
}

//CreateOpCodeTable crea la tabla de codigos de operacion
func (p *Assembler) CreateOpCodeTable() {
	p.Allocate()
	p.InsertCodes()
	p.InsertRegTable()
	fmt.Println("¡¡Tabla creada!!")
}

//InsertRegTable crea la tabla de registros
func (p *Assembler) InsertRegTable() {
	p.registers["A"] = 0
	p.registers["X"] = 1
	p.registers["L"] = 2
	p.registers["B"] = 3
	p.registers["S"] = 4
	p.registers["T"] = 5
	p.registers["F"] = 6
	p.registers["CP"] = 8
	p.registers["SW"] = 9
}

//Assemble ensambla el programa
func (p *Assembler) Assemble(line int, opcode int64, x int64, Addrs int64) {
	asminstr := Addrs | (x << 15) | (opcode << 16)
	res := p.code
	res = append(res, asminstr)
	p.code = res

	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

}

//AssembleByteHex Ensambla una instruccion byte con una constante hex
func (p *Assembler) AssembleByteHex(line int, bytes []byte) string {

	res1 := p.hexconst[len(p.hexconst)]
	for i := 0; i < len(bytes); i++ {
		s := string(bytes[i])
		val, _ := strconv.ParseInt(s, 16, 0)
		res1 = append(res1, val)

	}
	p.hexconst[len(p.hexconst)] = res1
	res := p.code
	res = append(res, -1)
	p.code = res
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

	var s string
	for i := 0; i < len(res1); i++ {
		s += strconv.FormatInt(res1[i], 16)
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return s
}

//AssembleByteASCII Ensambla una instruccion byte con una cadena
func (p *Assembler) AssembleByteASCII(line int, bytes []byte) string {

	res1 := p.charconst[len(p.charconst)]
	for i := 0; i < len(bytes); i++ {
		res1 = append(res1, int64(bytes[i]))
	}
	p.charconst[len(p.charconst)] = res1
	res := p.code
	res = append(res, -2)
	p.code = res
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

	var s string
	for i := 0; i < len(res1); i++ {
		s += strconv.FormatInt(res1[i], 16)
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return s
}

//AssembleWord ensambla una palabra
func (p *Assembler) AssembleWord(line int, word int64) string {
	res := p.code
	res = append(res, word)
	p.code = res
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

	s := strconv.FormatInt(word, 16)
	fmtStr := p.formaTo3byte(s)
	return fmtStr
}

//SetRegFlag Pone una bandera para indicar que se redujo un resw o resb
func (p *Assembler) SetRegFlag(line int) {
	res := p.code
	res = append(res, -3)
	p.code = res
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])
}

//SetCutRegFlag Pone una bandera para indicar que se redujo un resw o resb
func (p *Assembler) SetCutRegFlag(line int) {
	p.Hexcode = append(p.Hexcode, "flagCut")
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])
}

//SetNoCodeFlag Pone una bandera para indicar que se redujo un resw o resb
func (p *Assembler) SetNoCodeFlag(line int) {
	p.Hexcode = append(p.Hexcode, "flagNocode")
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])
}

//PrintCode Imprime el codigo objeto
func (p *Assembler) PrintCode() {
	//	res := p.code
	//fmt.Println(res)

	for key, value := range p.charconst {

		fmt.Println("Etiqueta CHAR :", key, "Valor:", value)
	}
	for key, value := range p.hexconst {

		fmt.Println("Etiqueta HEX:", key, "Valor:", value)
	}

}

func (p *Assembler) formaTo3byte(s string) string {
	for len(s) < 6 {
		s = "0" + s
	}
	return s
}

func (p *Assembler) formaTo1byte(s string) string {
	for len(s) < 2 {
		s = "0" + s
	}
	return s
}

func (p *Assembler) formaTo2byte(s string) string {
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}
func (p *Assembler) formaTo4byte(s string) string {
	for len(s) < 8 {
		s = "0" + s
	}
	return s
}

//Endinstr guarda la instruccion final
func (p *Assembler) Endinstr(symbol string) {
	p.EndAddrs = GetTabSim().GetSymbolAddrs(symbol)
}

//FormatCodeToHEX cambia el codigo objeto en decimal a hexadecimal representado en cadenas
func (p *Assembler) FormatCodeToHEX() {
	var hcont int
	var ccont int
	for i := 0; i < len(p.code); i++ {
		if p.code[i] == -1 {
			str := getBYTES(hcont, p.hexconst)
			p.Hexcode = append(p.Hexcode, str)
			hcont++
		} else if p.code[i] == -2 {
			str := getBYTES(ccont, p.charconst)
			p.Hexcode = append(p.Hexcode, str)
			ccont++
		} else if p.code[i] == -3 {
			p.Hexcode = append(p.Hexcode, "------")
		} else {
			s := strconv.FormatInt(p.code[i], 16)
			fmtStr := p.formaTo3byte(s)
			p.Hexcode = append(p.Hexcode, fmtStr)
		}
	}
	for i := 0; i < len(p.Hexcode); i++ {
		//fmt.Println(p.Hexcode[i])
	}
}

func getBYTES(index int, mapa map[int][]int64) string {
	array := mapa[index]
	var s string
	for i := 0; i < len(array); i++ {
		s += strconv.FormatInt(array[i], 16)
	}
	if len(s)%2 != 0 {
		s = "0" + s
	}
	return s
}

//ObjCode genera el codigo objeto
func (p *Assembler) ObjCode() []string {
	var registers []string
	registers = append(registers, p.hRegister())
	registers = append(registers, p.tRegisters()...)
	registers = append(registers, p.eRegister())
	//fmt.Println(registers)
	return registers
}

//ObjCodeXE genera el codigo objeto
func (p *Assembler) ObjCodeXE() []string {
	var registers []string
	registers = append(registers, p.hRegister())
	registers = append(registers, p.tRegistersXE()...)
	registers = append(registers, p.mRegistersXE()...)
	registers = append(registers, p.eRegister())
	//fmt.Println(registers)
	return registers
}

func (p *Assembler) hRegister() string {

	for len(p.Progname) < 6 {
		p.Progname = p.Progname + " "
	}
	return "H" + strings.ToUpper(p.Progname[0:6]+p.formaTo3byte(strconv.FormatInt(p.StartAddrs, 16))+p.formaTo3byte(GetTabSim().GetProgSizeHex()))
}
func (p *Assembler) eRegister() string {

	if p.EndAddrs != 0 {
		return "\nE" + strings.ToUpper(p.formaTo3byte(strconv.FormatInt(p.EndAddrs, 16))) + "\n"
	}
	return "\nE" + strings.ToUpper(p.formaTo3byte(strconv.FormatInt(p.StartAddrs, 16))) + "\n"
}

func (p *Assembler) tRegisters() []string {

	tReg := ""
	var tRegArray []string
	var dirstr string
	var tams []string
	//fmt.Println("TAM Addrs ", len(p.Addrs))
	//fmt.Println("TAM Hexcode", len(p.Hexcode))

	for i := 0; i < len(p.Hexcode); i++ {
		if p.Hexcode[i] == "------" || i == 0 || len(tReg)+len(p.Hexcode[i-1]) >= 71 {

			tRegArray = append(tRegArray, tReg)
			val := strconv.FormatInt(int64((len(tReg)-10)/2), 16)
			if len(val)%2 != 0 {
				val = "0" + val
			}
			tams = append(tams, val)
			tReg = ""
			if i == 0 {
				dirstr = p.formaTo3byte(strconv.FormatInt(p.StartAddrs, 16))
			} else {
				dirstr = p.formaTo3byte(strconv.FormatInt(p.Addrs[i-1], 16))
				if p.Hexcode[i] == "------" {
					dirstr = p.formaTo3byte(strconv.FormatInt(p.Addrs[i], 16))
				}

			}

			tReg += "\nT" + dirstr + "XX"

		}
		if p.Hexcode[i] != "------" {
			tReg += p.Hexcode[i]

		}
	}

	tRegArray = append(tRegArray, tReg)
	val := strconv.FormatInt(int64((len(tReg)-10)/2), 16)
	if len(val)%2 != 0 {
		val = "0" + val
	}
	tams = append(tams, val)
	for i := 0; i < len(tams); i++ {
		tRegArray[i] = strings.Replace(tRegArray[i], "XX", tams[i], 1)
		tRegArray[i] = strings.ToUpper(tRegArray[i])
	}

	//	fmt.Println(tRegArray)
	//fmt.Println(tams)
	return tRegArray
}

func (p *Assembler) mRegistersXE() []string {

	for i := 0; i < len(p.mReg); i++ {
		p.mReg[i] += p.Progname[0:6]
	}
	/*	for i := 0; i < len(p.Hexcode); i++ {
			if len(p.Hexcode[i]) == 8 && i != 0 {

				dirstr = strings.ToUpper(p.formaTo3byte(strconv.FormatInt(p.Addrs[i-1]+1, 16)))
				mRegArray = append(mRegArray, "\nM"+dirstr+"05"+"+"+p.Progname[0:6])
			}

		}
	*/
	/*for i, j := 0, len(p.mReg)-1; i < j; i, j = i+1, j-1 {
		p.mReg[i], p.mReg[j] = p.mReg[j], p.mReg[i]
	}*/
	return p.mReg
}
func (p *Assembler) tRegistersXE() []string {
	tReg := ""
	var tRegArray []string
	var dirstr string
	var tams []string
	//fmt.Println("TAM Addrs ", len(p.Addrs))
	//fmt.Println("TAM Hexcode", len(p.Hexcode))

	for i := 0; i < len(p.Hexcode); i++ {
		if p.Hexcode[i] == "flagCut" || i == 0 || len(tReg)+len(p.Hexcode[i-1]) >= 71 {

			tRegArray = append(tRegArray, tReg)
			val := strconv.FormatInt(int64((len(tReg)-10)/2), 16)
			if len(val)%2 != 0 {
				val = "0" + val
			}
			tams = append(tams, val)
			tReg = ""
			if i == 0 {
				dirstr = p.formaTo3byte(strconv.FormatInt(p.StartAddrs, 16))
				if i-1 < 0 && p.Hexcode[i] == "flagCut" {
					dirstr = p.formaTo3byte(strconv.FormatInt(p.Addrs[i], 16))
				}

			} else {
				dirstr = p.formaTo3byte(strconv.FormatInt(p.Addrs[i-1], 16))

				if p.Hexcode[i] == "flagCut" {

					dirstr = p.formaTo3byte(strconv.FormatInt(p.Addrs[i], 16))
				}

			}
			fmt.Println(dirstr, p.Hexcode[i], i)
			tReg += "\nT" + dirstr + "XX"

		}
		if p.Hexcode[i] != "flagCut" && p.Hexcode[i] != "flagNocode" {
			tReg += p.Hexcode[i]

		}
	}

	tRegArray = append(tRegArray, tReg)
	val := strconv.FormatInt(int64((len(tReg)-10)/2), 16)
	if len(val)%2 != 0 {
		val = "0" + val
	}
	tams = append(tams, val)
	for i := 0; i < len(tams); i++ {
		tRegArray[i] = strings.Replace(tRegArray[i], "XX", tams[i], 1)
		tRegArray[i] = strings.ToUpper(tRegArray[i])
	}

	//	fmt.Println(tRegArray)
	//fmt.Println(tams)
	return tRegArray
}

///SIC XE EXTRAS

//Assemblef1 Ensambla el formato 1
func (p *Assembler) Assemblef1(line int, opcode int64) {
	asmInst := opcode
	//p.code = append(p.code, asmInst)
	asmInstStr := strconv.FormatInt(asmInst, 16)
	asmInstStr = strings.ToUpper(asmInstStr)
	asmInstStr = p.formaTo1byte(asmInstStr)
	p.Hexcode = append(p.Hexcode, asmInstStr)
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])
}

//Assemblef2 Ensambla el formato 2
func (p *Assembler) Assemblef2(line int, opcode int64, reg1 int64, reg2 int64) {
	asmInst := reg2 | (reg1 << 4) | (opcode << 8)
	asmInstStr := strconv.FormatInt(asmInst, 16)
	asmInstStr = strings.ToUpper(asmInstStr)
	asmInstStr = p.formaTo2byte(asmInstStr)
	p.Hexcode = append(p.Hexcode, asmInstStr)
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])
}

//Assemblef3 Ensambla el formato 3
func (p *Assembler) Assemblef3(line int, opcode int64, ni int64, x int64, b int64, pp int64, e int64, desp int64) {

	opcode = opcode | ni
	asmInst := desp | (e << 12) | (pp << 13) | (b << 14) | (x << 15) | (opcode << 16)
	asmInstStr := strconv.FormatInt(asmInst, 16)
	asmInstStr = strings.ToUpper(asmInstStr)
	asmInstStr = p.formaTo3byte(asmInstStr)
	p.Hexcode = append(p.Hexcode, asmInstStr)
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

}

//Assemblef4 Ensambla el formato 4
func (p *Assembler) Assemblef4(line int, opcode int64, ni int64, x int64, b int64, pp int64, e int64, desp int64) {

	opcode = opcode | ni
	asmInst := desp | (e << 20) | (pp << 21) | (b << 22) | (x << 23) | (opcode << 24)
	asmInstStr := strconv.FormatInt(asmInst, 16)
	asmInstStr = strings.ToUpper(asmInstStr)
	asmInstStr = p.formaTo4byte(asmInstStr)
	p.Hexcode = append(p.Hexcode, asmInstStr)
	p.Addrs = append(p.Addrs, GetTabSim().Progpc[line])

	var dirstr string
	found := false

	for key, value := range GetTabSim().table {

		num, _ := strconv.ParseInt(asmInstStr, 16, 0)
		num = num & 0x000FFFFF
		if value == num {
			found = true

			if found {
				fmt.Println(key, value)

				if len(asmInstStr) == 8 {

					dirstr = strings.ToUpper(p.formaTo3byte(strconv.FormatInt(GetTabSim().Progpc[line-1]+1, 16)))
					p.mReg = append(p.mReg, "\nM"+dirstr+"05"+"+")
				}
				found = false
			}
		}

	}

}

//GetRegEntry Encuentra el registro correspondiente
func (p *Assembler) GetRegEntry(name string) int64 {
	return p.registers[name]
}

//IsValidRegister Encuentra si un registro es valido o no
func (p *Assembler) IsValidRegister(reg int64) bool {
	if reg == 0 || reg == 1 || reg == 2 || reg == 3 || reg == 4 || reg == 5 || reg == 6 || reg == 8 || reg == 9 {
		return true
	}
	return false
}

//AssembleXE ensambla la sic extendida
func (p *Assembler) AssembleXE(line int, nemonic string, indexed int64, symbol string, bRegister int64, ni int64) {
	ta := GetTabSim().GetSymbolAddrs(symbol)
	desp := ta - GetTabSim().Progpc[line]
	opcode := p.GetCodeEntry(nemonic)
	if desp <= 2047 && desp >= -2048 {
		if desp < 0 { //C'2
			desp = 4096 + desp
		}
		p.Assemblef3(line, opcode, ni, indexed, 0, 1, 0, desp)
	} else {
		desp := ta - bRegister
		if desp >= 0 && desp <= 4095 {
			p.Assemblef3(line, opcode, ni, indexed, 1, 0, 0, desp)
		} else { //error
			p.Assemblef3(line, opcode, ni, indexed, 1, 1, 0, 4095)
		}
	}
}
