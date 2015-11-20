package gui

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"parser"
	"path/filepath"
	"util"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//Gui es la Estructura para crear la instancia del paquete GUI
type Gui struct {
}

//New es una funcion para crear una referencia a gui
func New() *Gui {
	return &Gui{}
}

//InitializeGUI es la funcion que inicializa y ejecuta todo el entorno gráfico de la aplicacion
func (p *Gui) InitializeGUI(errorstr string, file *os.File) {
	gtk.Init(nil)

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.Maximize()
	window.SetTitle("Go-SIC ASM/SIM!")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	//--------------------------------------------------------
	//	GTK statusbar
	//--------------------------------------------------------
	statusbar := gtk.NewStatusbar()
	context_id := statusbar.GetContextId("hola")
	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	//--------------------------------------------------------
	// GtkMenuBar
	//--------------------------------------------------------
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)

	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("Código Fuente/Archivo intermedio/TABSIM")
	framebox1 := gtk.NewVBox(false, 1)
	frame1.Add(framebox1)
	frame1.SetSizeRequest(300, 300)

	frame2 := gtk.NewFrame("Debug/Código Objeto")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)

	vpaned.Pack1(frame1, false, false)
	vpaned.Pack2(frame2, false, false)

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	ventanasPrincipales := gtk.NewHBox(false, 1)
	//----------------------------------------------

	label := gtk.NewLabel("Ensamblador SIC SIC/XE")
	label.ModifyFontEasy("DejaVu Serif 15")
	framebox1.PackStart(label, false, true, 0)
	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin := gtk.NewScrolledWindow(nil, nil)

	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()

	textview.ModifyFontEasy("Sans 10")
	var start, end gtk.TextIter
	buffer := textview.GetBuffer()
	buffer.GetStartIter(&start)
	swin.Add(textview)
	ventanasPrincipales.Add(swin)
	//framebox1.Add(swin)
	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin4 := gtk.NewScrolledWindow(nil, nil)

	swin4.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin4.SetShadowType(gtk.SHADOW_IN)
	textview4 := gtk.NewTextView()

	textview4.ModifyFontEasy("Sans 10")
	var start4, end4 gtk.TextIter
	buffer4 := textview4.GetBuffer()
	buffer4.GetStartIter(&start4)
	swin4.Add(textview4)
	ventanasPrincipales.Add(swin4)
	//framebox1.Add(swin)
	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin5 := gtk.NewScrolledWindow(nil, nil)

	swin5.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin5.SetShadowType(gtk.SHADOW_IN)
	textview5 := gtk.NewTextView()

	textview5.ModifyFontEasy("Sans 10")
	var start5, end5 gtk.TextIter
	buffer5 := textview5.GetBuffer()
	buffer5.GetStartIter(&start5)
	swin5.Add(textview5)
	ventanasPrincipales.Add(swin5)
	framebox1.PackStart(ventanasPrincipales, true, true, 1)
	//framebox1.Add(swin)
	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	ventanas := gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin2 := gtk.NewScrolledWindow(nil, nil)

	swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin2.SetShadowType(gtk.SHADOW_IN)
	textview2 := gtk.NewTextView()
	textview2.SetEditable(false)
	var start2, end2 gtk.TextIter
	buffer2 := textview2.GetBuffer()
	buffer2.GetStartIter(&start2)
	swin2.Add(textview2)
	//framebox2.Add(swin2)
	ventanas.Add(swin2)
	//--------------------------------------------------------
	// GtkTextView
	//--------------------------------------------------------
	swin3 := gtk.NewScrolledWindow(nil, nil)

	swin3.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin3.SetShadowType(gtk.SHADOW_IN)
	textview3 := gtk.NewTextView()
	textview3.SetEditable(false)
	var start3, end3 gtk.TextIter
	buffer3 := textview3.GetBuffer()
	buffer3.GetStartIter(&start3)
	swin3.Add(textview3)
	//framebox2.Add(swin2)
	ventanas.Add(swin3)
	framebox2.PackStart(ventanas, true, true, 1)

	//--------------------------------------------------------
	// GtkEntry
	//--------------------------------------------------------
	entry := gtk.NewEntry()
	entry.SetText("Para comenzar, favor de escoger un archivo del directorio....")
	//entry.SetSensitive(false)
	entry.SetEditable(false)
	framebox2.Add(entry)

	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons := gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkButton //**OPEN FILE****//
	//--------------------------------------------------------
	button := gtk.NewButtonWithLabel("Elegir archivo...")
	var filename string
	var isXE bool
	button.Clicked(func() {
		fmt.Println("button clicked:", button.GetLabel())
		messagedialog := gtk.NewMessageDialog(
			button.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK,
			"Escoja un archivo .*s o *.x de su directorio")
		messagedialog.Response(func() {
			fmt.Println("Dialog OK!")

			//--------------------------------------------------------
			// GtkFileChooserDialog
			//--------------------------------------------------------
			filechooserdialog := gtk.NewFileChooserDialog(
				"Choose File...",
				button.GetTopLevelAsWindow(),
				gtk.FILE_CHOOSER_ACTION_OPEN,
				gtk.STOCK_OK,
				gtk.RESPONSE_ACCEPT)
			filter := gtk.NewFileFilter()
			filter.AddPattern("*.s")
			filter.AddPattern("*.x")
			filechooserdialog.AddFilter(filter)
			filechooserdialog.Response(func() {
				/*aqui va el dialogo */
				filename = filechooserdialog.GetFilename()
				var extension = filepath.Ext(filename)
				if extension == ".x" {
					isXE = true
				} else {
					isXE = false
				}
				fmt.Println(filechooserdialog.GetFilename())
				statusbar.Push(context_id, filename)
				dat, err := ioutil.ReadFile(filechooserdialog.GetFilename())
				if err == nil {
					buffer.GetStartIter(&start)
					buffer.GetEndIter(&end)
					buffer.Delete(&start, &end)

					buffer.Insert(&start, string(dat))
					entry.SetText("Ahora haz click en el boton '¡Analizar!'")
					//fmt.Print(string(dat))
				}
				filechooserdialog.Destroy()
			})
			filechooserdialog.Run()
			messagedialog.Destroy()
		})
		messagedialog.Run()
	})
	buttons.Add(button)
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	button2 := gtk.NewButtonWithLabel("Ensamblar!")
	button2.Clicked(func() {
		fmt.Println("button clicked:", button2.GetLabel())
		messagedialog2 := gtk.NewMessageDialog(
			button.GetTopLevelAsWindow(),
			gtk.DIALOG_MODAL,
			gtk.MESSAGE_INFO,
			gtk.BUTTONS_OK,
			"Analisis lexico y sintactico presiona OK para continuar...")
		messagedialog2.Response(func() {
			fmt.Println(filename)
			statusbar.Push(context_id, filename)
			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("Error al cargar el archivo: %s", err.Error())
			} else {
				//*********************************************************
				util.GetTabSim().Allocate()
				if isXE == false {
					fmt.Println("IS SIC")
					errorstr = parser.New().Parse(file, true, false)
					file, err = os.Open(filename)

					util.GetAssembler().CreateOpCodeTable()
					errorstr = parser.New().Parse(file, false, false)
					file.Close()
					//util.GetAssembler().PrintCode()
					util.GetAssembler().FormatCodeToHEX()
					buffer2.GetStartIter(&start2)
					buffer2.GetEndIter(&end2)
					buffer2.Delete(&start2, &end2)
					if len(errorstr) > 0 {
						buffer2.Insert(&start2, string(errorstr+util.GetTabSim().Errores))

						strArray := strings.Split(filename, ".s")
						f, err := os.OpenFile(strArray[0]+".ts", os.O_WRONLY|os.O_CREATE, 0600)
						if err != nil {
							panic(err)
						}
						defer f.Close()
						if _, err = f.WriteString(errorstr); err != nil {
							panic(err)
						}
					} else {
						buffer2.Insert(&start2, "Correcto!"+util.GetTabSim().Errores)
					}
					///*******************************************************
					///*******************************************************
					tabsimlines := util.GetTabSim().ReturnTable()
					buffer5.GetStartIter(&start5)
					buffer5.GetEndIter(&end5)
					buffer5.Delete(&start5, &end5)
					if len(tabsimlines) > 0 {
						buffer5.Insert(&start5, string(tabsimlines))
					}
					///*******************************************************
					///*******************************************************
					util.GetTabSim().DecToHex()
					objectCode := util.GetAssembler().ObjCode()
					var obj string
					for i := 0; i < len(objectCode); i++ {

						if len(objectCode[i]) > 10 {

							obj += objectCode[i]
						}
						if i == len(objectCode)-1 {
							obj += objectCode[i]
						}
					}
					buffer3.GetStartIter(&start3)
					buffer3.GetEndIter(&end3)
					buffer3.Delete(&start3, &end3)
					if len(obj) > 0 {
						buffer3.Insert(&start3, string(obj))
					}
					f2 := newFile(strings.Split(filename, ".s")[0], ".os")
					if _, err = f2.WriteString(obj); err != nil {
						panic(err)
					}
					///*******************************************************
					///*******************************************************
					var todostr string
					lines, err := readLines(filename)
					if err == nil {
						addrsstr := util.GetTabSim().GetProgPcStr()
						//fmt.Println(addrsstr)
						for i := 0; i < len(lines); i++ {
							if addrsstr[i+1] == "" {
								addrsstr[i+1] = addrsstr[i]
								todostr += addrsstr[i] + "\t" + lines[i] + "\n"
							} else {
								todostr += addrsstr[i] + "\t" + lines[i] + "\t" + util.GetTabSim().Progobjstr[i] + "\n"
							}
						}
						//util.GetTabSim().PrintProgpc()
						//fmt.Print(todostr)
						addrsstr = nil
					}
					buffer4.GetStartIter(&start4)
					buffer4.GetEndIter(&end4)
					buffer4.Delete(&start4, &end4)
					if len(obj) > 0 {
						buffer4.Insert(&start4, string(todostr))
					}
				} else {
					fmt.Println("IS SIC XE")
					errorstr = parser.New().Parse(file, true, true)
					file, err = os.Open(filename)

					util.GetAssembler().CreateOpCodeTable()
					parser.New().Parse(file, false, true)
					file.Close()
					//util.GetAssembler().PrintCode()
					//util.GetAssembler().FormatCodeToHEX()
					buffer2.GetStartIter(&start2)
					buffer2.GetEndIter(&end2)
					buffer2.Delete(&start2, &end2)
					if len(errorstr) > 0 {
						buffer2.Insert(&start2, string(errorstr+util.GetTabSim().Errores))

						strArray := strings.Split(filename, ".x")
						f, err := os.OpenFile(strArray[0]+".tx", os.O_WRONLY|os.O_CREATE, 0600)
						if err != nil {
							panic(err)
						}
						defer f.Close()
						if _, err = f.WriteString(errorstr); err != nil {
							panic(err)
						}
					} else {
						buffer2.Insert(&start2, "Correcto!"+util.GetTabSim().Errores)
					}
					///*******************************************************
					///*******************************************************
					tabsimlines := util.GetTabSim().ReturnTable()
					buffer5.GetStartIter(&start5)
					buffer5.GetEndIter(&end5)
					buffer5.Delete(&start5, &end5)
					if len(tabsimlines) > 0 {
						buffer5.Insert(&start5, string(tabsimlines))
					}
					///*******************************************************
					///*******************************************************

					objectCode := util.GetAssembler().ObjCodeXE()
					util.GetTabSim().PcToHex()
					var obj string
					for i := 0; i < len(objectCode); i++ {

						if len(objectCode[i]) > 10 {

							obj += objectCode[i]
						}
						if i == len(objectCode)-1 {
							obj += objectCode[i]
						}
					}
					buffer3.GetStartIter(&start3)
					buffer3.GetEndIter(&end3)
					buffer3.Delete(&start3, &end3)
					if len(obj) > 0 {
						buffer3.Insert(&start3, string(obj))
					}
					f2 := newFile(strings.Split(filename, ".x")[0], ".ox")
					if _, err = f2.WriteString(obj); err != nil {
						panic(err)
					}
					///*******************************************************
					///*******************************************************
					var todostr string
					lines, err := readLines(filename)
					if err == nil {
						addrsstr := util.GetTabSim().GetProgPcStr()
						//fmt.Println(addrsstr)
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
					buffer4.GetStartIter(&start4)
					buffer4.GetEndIter(&end4)
					buffer4.Delete(&start4, &end4)
					if len(obj) > 0 {
						buffer4.Insert(&start4, string(todostr))
					}
				}
			}
			entry.SetText("En la ventana debug se encuentran los errores lexicos y sintacticos del programa")
			fmt.Println("Dialog OK!")

			messagedialog2.Destroy()
		})
		messagedialog2.Run()
	})
	buttons.Add(button2)
	framebox2.PackStart(buttons, false, false, 0)

	/////////////////////////////////////////////

	button3 := gtk.NewButtonWithLabel("Cargador...")
	button3.Clicked(func() {
		strArray2 := strings.Split(filename, ".s")
		f2, err := os.OpenFile(strArray2[0]+".o", os.O_WRONLY|os.O_CREATE, 0600)
		var obj string
		objectCode := util.GetAssembler().ObjCode()
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
		//nombre := strings.Split(strArray2[0], "work")

		cmd := exec.Command("script.sh", strArray2[0]+".o")

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		//err = cmd.Wait()
		log.Printf("Command finished with error: %v", err)
	})

	buttons.Add(button3)

	////////////////////////////

	buttons = gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkMenuItem ///***********SAVE FILE*****************///
	//--------------------------------------------------------
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	var menuitem2 *gtk.MenuItem
	menuitem2 = gtk.NewMenuItemWithMnemonic("G_uardar")
	menuitem2.Connect("activate", func() {
		strArray := strings.Split(filename, ".s")
		if len(strArray) <= 1 {
			strArray = strings.Split(filename, ".x")
		}
		if len(strArray) > 1 && len(filename) > 0 {
			statusbar.Push(context_id, filename)
			var s string
			buffer.GetStartIter(&start)
			buffer.GetEndIter(&end)
			s = buffer.GetText(&start, &end, true)
			fmt.Println(filename)
			err := ioutil.WriteFile(filename, []byte(s), 0644)
			if err != nil {
				panic(err)
			}
		} else {

			filechooserdialog := gtk.NewFileChooserDialog(
				"Choose File...",
				button.GetTopLevelAsWindow(),
				gtk.FILE_CHOOSER_ACTION_SAVE,
				gtk.STOCK_OK,
				gtk.RESPONSE_ACCEPT)
			filter := gtk.NewFileFilter()
			///***ALLOWS SIC AND SIC-XE EXTENSION***///
			filter.AddPattern("*.s")
			filter.AddPattern("*.x")
			filechooserdialog.AddFilter(filter)
			filechooserdialog.Response(func() {
				/*aqui va el dialogo */
				statusbar.Push(context_id, filename)
				filename = filechooserdialog.GetFilename()
				if len(filename) > 0 {
					fmt.Println(filechooserdialog.GetFilename())
					var s string
					buffer.GetStartIter(&start)
					buffer.GetEndIter(&end)
					s = buffer.GetText(&start, &end, true)

					err := ioutil.WriteFile(filename, []byte(s), 0644)
					if err != nil {
						panic(err)
					}
					entry.SetText("Haz click en el boton analizar")
				}
				filechooserdialog.Destroy()
			})
			filechooserdialog.Run()

		}
	})
	submenu.Append(menuitem2)

	var menuitem *gtk.MenuItem
	menuitem = gtk.NewMenuItemWithMnemonic("S_alir")
	menuitem.Connect("activate", func() {
		gtk.MainQuit()
	})
	submenu.Append(menuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_View")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	checkmenuitem := gtk.NewCheckMenuItemWithMnemonic("_Disable")
	checkmenuitem.Connect("activate", func() {
		textview.SetSensitive(!checkmenuitem.GetActive())
		textview2.SetSensitive(!checkmenuitem.GetActive())
	})
	submenu.Append(checkmenuitem)

	cascademenu = gtk.NewMenuItemWithMnemonic("_Help")
	menubar.Append(cascademenu)
	submenu = gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem = gtk.NewMenuItemWithMnemonic("_About")
	menuitem.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("Go-SIC sim!")
		dialog.SetProgramName("Go-SIC sim")
		dialog.SetLicense("The library is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License). The idea is that if you can use Go (and Gtk) in a project, you should also be able to use go-gtk.")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})
	submenu.Append(menuitem)

	//--------------------------------------------------------
	// GtkStatusbar
	//--------------------------------------------------------

	statusbar.Push(context_id, "No hay archivo seleccionado")

	framebox2.PackStart(statusbar, false, false, 0)

	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(600, 600)
	window.ShowAll()
	gtk.Main()
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

func newFile(name string, extension string) *os.File {
	f, err := os.OpenFile(name+extension, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	return f
}
