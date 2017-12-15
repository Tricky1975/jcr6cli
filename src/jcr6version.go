/*
	JCR6CLI
	Version information
	
	
	
	(c) Jeroen P. Broks, 2017, All rights reserved
	
		This program is free software: you can redistribute it and/or modify
		it under the terms of the GNU General Public License as published by
		the Free Software Foundation, either version 3 of the License, or
		(at your option) any later version.
		
		This program is distributed in the hope that it will be useful,
		but WITHOUT ANY WARRANTY; without even the implied warranty of
		MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
		GNU General Public License for more details.
		You should have received a copy of the GNU General Public License
		along with this program.  If not, see <http://www.gnu.org/licenses/>.
		
	Exceptions to the standard GNU license are available with Jeroen's written permission given prior 
	to the project the exceptions are needed for.
Version: 17.12.15
*/
package main

import (
	// Go internal libs
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"os/exec"
	"strings"
	
	// JCRCLI's internal mods
	"jcr6cli/src/imps/ver"
	
	// Tricky's modules
	"trickyunits/mkl"
	"trickyunits/qstr"
	"trickyunits/ansistring"
)

const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta


func init(){
mkl.Version("JCR6 CLI (GO) - jcr6version.go","17.12.15")
mkl.Lic    ("JCR6 CLI (GO) - jcr6version.go","GNU General Public License 3")
}

func main() {
	ver.CHVER()
	ext:=""
	npext:=""
	teken:="*"
	full:=false
	param:="VERSION"
	if len(os.Args)>1 { if os.Args[1]=="FULL" { full=true; param="FULLVERSION" }}
	if runtime.GOOS=="windows" { ext=".exe"; npext="exe"; teken="-" }
	exe,_:=os.Executable()
	dir:=path.Dir(exe)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		//fmt.Println(file.Name()) // debug
		file:=file.Name()
		//fmt.Println("Looing at: "+file) // debug
		if qstr.Left(file,4)=="jcr6" && qstr.Right(file,4)!=".lua" && path.Ext(file)==npext {
			cmd:=exec.Command(dir+"/"+file+ext,teken+param+teken)
			o,err := cmd.Output()
			outputstring:=fmt.Sprintf("%s",o)
			if err!=nil { 
				log.Println(err.Error())
			} else {
				olines:=strings.Split(outputstring,"\n")
				mainline:=olines[0]
				pdat:=strings.Split(mainline,"\t")
				pdat=append(pdat,"error")
				fmt.Println(ansistring.SCol(qstr.Left(pdat[0]+"                             ",20),Yellow,0)+ansistring.SCol(pdat[1],Cyan,0))
				if full {
					for i:=1;i<len(olines);i++{
						fmt.Println("\t"+ansistring.SCol(olines[i],Magenta,0))
					}
				}
			}
		}
	}
}

