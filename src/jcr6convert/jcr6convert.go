/*
	JCR6 CLI
	Converter main source
	
	
	
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


// Only "ver" needed for versioning.
// Since the conversion tool does neither pack nor extract itself 
// the drivers will not be needed here.
import(
	"jcr6cli/src/imps/ver"
	"flag"
	"fmt"
	"trickyunits/ansistring"
	"path"
	"strings"
	"trickyunits/qstr"
	"trickyunits/qff"
	"os"
)

const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green


// Yes or no!
func yn(b bool) string {
	if b {
		return ansistring.SCol("YES",Green,0)
	} else {
		return ansistring.SCol("NO",Red,0)
	}
}
func init() {
	AC_Vers()
	initpack()
}

func main(){
	// Version show if asked
	ver.CHVER()
	// Flags
	inputfile := flag.String("i","","File to be converted (required)")
	outputfile := flag.String("o","","Output file. When not set it will be generated out of the -i tag. Don't use extensions. They will be added automatically!")
	convertto := flag.String("t","jcr","Convert to format. Multiple targets possible, separate with semicolons and put the targets between quotes!")
	dontverify := flag.Bool("dv",false,"When set verification will be skipped")
	flag.Parse()
	verify:=!*dontverify
	// 	fmt.Printf("i = %s\no = %s\nt = %s\ndv= %s\n\n",*inputfile,*outputfile,*convertto,*dontverify) // debug

	// No input no action
	if *inputfile == "" {
		fmt.Print(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 convert ",Yellow,0),ansistring.SCol("-i <input> ",Cyan,0),ansistring.SCol("[ -o output -t targetformats ]",Magenta,ansistring.A_Dark),"\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	
	// Parsing flags into the way to go
	iname:=strings.Replace(*inputfile,"\\","/",-10)
	if !isd(iname) {
		iname = qff.PWD()+"/"+iname
		iname=strings.Replace(iname,"//","/",-10)
	}
	oname:=strings.Replace(*outputfile,"\\","/",-10)
	if oname=="" {
		oname=qstr.StripExt(iname)
	}
	if !isd(oname) {
		oname = qff.PWD()+"/"+oname
		oname=strings.Replace(oname,"//","/",-10)
	}
		
	targets := strings.Split(strings.ToLower(*convertto),";")
	source  := strings.ToLower(path.Ext(*inputfile))
	source   = qstr.Right(source,len(source)-1)
	
	// Parse results
	fmt.Print(ansistring.SCol("Input file:      ",Yellow,0))
	fmt.Println(ansistring.SCol(iname,Cyan,0))
	fmt.Print(ansistring.SCol("Input type:      ",Yellow,0))
	fmt.Println(ansistring.SCol(source,Cyan,0))
	if len(targets)==1 {
		fmt.Print(ansistring.SCol("Output file:     ",Yellow,0))
		fmt.Println(ansistring.SCol(oname+"."+targets[0],Cyan,0))
	} else {
		for i,f:=range targets {
			if i==1 {
				fmt.Print(ansistring.SCol("Output files:    ",Yellow,0))
			} else {
				fmt.Print(                "                 ")
			}
			fmt.Println(ansistring.SCol(oname+"."+f,Cyan,0))
		}
	}
	fmt.Print(ansistring.SCol("Output type:     ",Yellow,0))
	for _,tn:=range targets {
		fmt.Print(ansistring.SCol(tn+"  ",Cyan,0))
	}
	fmt.Print("\n")
	fmt.Print(ansistring.SCol("Verification:    ",Yellow,0))
	fmt.Println(yn(verify))
	
	// All required dependencies accounted for?
	checkpack("u",source)
	for _,tf := range targets{
		checkpack("p",tf)
	}
}
