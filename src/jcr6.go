/*
	JCR6 CLI
	Main program
	
	
	
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
Version: 17.12.16
*/
package main

import(
	// Go
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"path"
	"runtime"
	//"os/exec"
	//"strings"
	
	// Internal
	"jcr6cli/src/imps/ver"
	
	// Tricky
	"trickyunits/mkl"
	"trickyunits/qstr"
	"trickyunits/ansistring"
	"trickyunits/shell"
)

const Yellow =ansistring.A_Yellow
const Cyan   =ansistring.A_Cyan
const Magenta=ansistring.A_Magenta
const Red    =ansistring.A_Red

const Bright=ansistring.A_Bright

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6.go","17.12.16")
mkl.Lic    ("JCR6 CLI (GO) - jcr6.go","GNU General Public License 3")
}


func main()																						{
	ver.CHVER()
	me,_:=os.Executable()
	dir:=path.Dir(me)
	fmt.Println(ansistring.SCol("JCR6 - Command Line Tools!",Yellow,Bright))
	fmt.Println(ansistring.SCol("Coded by: Jeroen P. Broks",Cyan,0))
	fmt.Println(ansistring.SCol("(c) Jeroen P. Broks 2016-2017",Magenta,0))
	fmt.Println()
	if len(os.Args)>1 {
		a:=[]string{}
		for i:=2;i<len(os.Args);i++ { a=append(a,os.Args[i]) }
		cmd:=shell.ArrayCommand(dir+"/jcr6_"+os.Args[1],a)
		cmd.Stdout = os.Stdout
		cmd.Run()
		os.Exit(0)
	}
	fmt.Println(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 ",Yellow,0),ansistring.SCol("<command> ",Magenta,ansistring.A_Dark),ansistring.SCol("<command parameters> ",Cyan,0),"\n\n")
	fmt.Println()
	fmt.Println(ansistring.SCol("Giving up a command without any parameters will show you the help page of that specific command",Magenta,0))
	fmt.Println()
	fmt.Println(ansistring.SCol("Available commands",Cyan,0))	
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, ifile := range files {
		//fmt.Println(file.Name()) // debug
		file:=ifile.Name()
		//fmt.Println("Looing at: "+file) // debug
		npext:=""
		if runtime.GOOS=="windows" { npext="exe" }
		if qstr.Left(file,5)=="jcr6_" && qstr.Right(file,4)!=".lua" && path.Ext(file)==npext {
			fmt.Print(ansistring.SCol("= ",Red,Bright))
			fmt.Println(ansistring.SCol(qstr.Right(file,len(file)-5),Yellow,0))
		}
	}
	fmt.Println()
	fmt.Println("The JCR6 modules have been licenced under the terms of the MPL 2.0")
	fmt.Println("These command line utilities have been licenced under the termso of the GNU GPL 3")
	fmt.Println()
	fmt.Println("The extentable nature of JCR6 makes it possible to add more commands.\nThe copyrights will always belong its respective authors, and they do have the right to chose their own licenses as far as the libraries they used allow those")
}
