/*
	JCR6 CLI
	Type
	
	
	
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
Version: 17.12.08
*/
package main



import (
    "trickyunits/mkl"
    "trickyunits/ansistring"
    "trickyunits/jcr6/jcr6main"
	"jcr6cli/src/imps/ver"
_	"jcr6cli/src/imps/drv"
	"fmt"
	"os"
)

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6type.go","17.12.08")
mkl.Lic    ("JCR6 CLI (GO) - jcr6type.go","GNU General Public License 3")
}


func main(){
	ver.CHVER()
	//fmt.Printf("Not fully functional yet, but that will soon come, I guess ;)\n\n")
	if len(os.Args)<3{
		fmt.Print(ansistring.SCol("Usage:",ansistring.A_Cyan,0),ansistring.SCol("jcr6 type ",ansistring.A_Yellow,0),ansistring.SCol("<JCR6 Resource File> <entry>",ansistring.A_Magenta,0),"\n")
		os.Exit(0)
	}
	jd:=jcr6main.Dir(os.Args[1])
	if jcr6main.JCR6Error!=""{
		fmt.Println(ansistring.SCol("!! ERROR !!",ansistring.A_Red,ansistring.A_Blink))
		fmt.Println(ansistring.SCol(jcr6main.JCR6Error,ansistring.A_Yellow,0))
		os.Exit(120)
	}
	s:=jcr6main.JCR_String(jd,os.Args[2])
	if jcr6main.JCR6Error!=""{
		fmt.Println(ansistring.SCol("!! ERROR !!",ansistring.A_Red,ansistring.A_Blink))
		fmt.Println(ansistring.SCol(jcr6main.JCR6Error,ansistring.A_Yellow,0))
		os.Exit(120)
	}
	fmt.Println(s)
}
