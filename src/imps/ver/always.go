// License Information:
// JCR6 cli central unit
// Central stuff
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


package ver

import(
	"os"
	"fmt"
	"strings"
	"trickyunits/qstr"
	"trickyunits/qff"
	"trickyunits/mkl"
	"trickyunits/gini"
	"trickyunits/dirry"
	"trickyunits/ansistring"
)


type tconfig struct{
	config gini.TGINI	
}

func (s tconfig) Set(key,value string) {
	s.config.D(key,value)	
	s.config.SaveSource(dirry.Dirry("$AppSupport$/jcr6cli_config.gini"))
}

func (s tconfig) Get(key string) string {
	return s.config.C(key)
}

func (s tconfig) Yes(key string) bool {
	v:=strings.ToUpper(s.config.C(key))
	return v=="TRUE" || v=="JA" || v=="YES" || v=="T" || v=="Y" || v=="J"
}





var Config tconfig


func WANTVER() bool {
	if len(os.Args)>1{
		return os.Args[1]==jvstring || os.Args[1]==fvstring
	}
	return false
}

func CHVER(){
	if len(os.Args)>1{
		if os.Args[1]==jvstring || os.Args[1]==fvstring{
			fmt.Println(qstr.StripAll(os.Args[0])+"\t"+"v"+mkl.Newest()+" go")
			if os.Args[1]==fvstring{
				fmt.Print(mkl.ListAll())
			}
			os.Exit(0)
		}
	}
}

func init(){
	if qff.IsFile(dirry.Dirry("$AppSupport$/jcr6cli_config.gini")) {   
		Config.config = gini.ReadFromFile(dirry.Dirry("$AppSupport$/jcr6cli_config.gini"))
	}
	if Config.Get("ANSI")!="" {ansistring.ANSI_Use=Config.Yes("ANSI")}
}
