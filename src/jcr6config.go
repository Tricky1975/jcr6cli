// License Information:
// JCR6 CLI
// Config manager
// 
// 
// 
// (c) Jeroen P. Broks, 
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// 
// Please note that some references to data like pictures or audio, do not automatically
// fall under this licenses. Mostly this is noted in the respective files.
// 
// Version: 19.03.20
// End License Information


package main



import (
        "os"
        "fmt"
		"jcr6cli/src/imps/ver"
        "trickyunits/ansistring"
        "trickyunits/mkl"
)



const Yellow =ansistring.A_Yellow
const Cyan   =ansistring.A_Cyan
const Magenta=ansistring.A_Magenta
const Red    =ansistring.A_Red




func init(){
mkl.Version("JCR6 CLI (GO) - jcr6config.go","19.03.20")
mkl.Lic    ("JCR6 CLI (GO) - jcr6config.go","GNU General Public License 3")
}



func main() {
	if len(os.Args)<3 {
		ver.CHVER()
		fmt.Println(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 config ",Yellow,0),ansistring.SCol("<variable> ",Magenta,ansistring.A_Dark),ansistring.SCol("<value>",Cyan,0),"\n\n")
		os.Exit(0)
	}
	ver.Config.Set(os.Args[1],os.Args[2])
}

