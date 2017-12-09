/*
	JCRCLI 6 
	Alt Command for Converter
	
	
	
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
Version: 17.12.09
*/
package main

import "os/exec"
//import "trickyunits/qstr"
//import "path/filepath"
import "os"
import "trickyunits/ansistring"
import "fmt"
import "trickyunits/mkl"


/* This is the original Command function from the original Go libraries
   Let's alter it below
func Command(name string, arg ...string) *Cmd {
  	cmd := &Cmd{
  		Path: name,
  		Args: append([]string{name}, arg...),
  	}
  	if filepath.Base(name) == name {
  		if lp, err := LookPath(name); err != nil {
  			cmd.lookPathErr = err
  		} else {
  			cmd.Path = lp
  		}
  	}
  	return cmd
  }
   
*/


func AltCommand(shit []string) *exec.Cmd {
	cmd:= &exec.Cmd{
		Path: shit[0],
		Args: shit,
	}
	/*
	if qstr.Left(cmd.Path,5)=="jcr6_" {
		cmd.Path = filepath.Dir(os.Args[0])+"/"+cmd.Path
	} else {
	*/ 
		if lp,err:=exec.LookPath(cmd.Path); err!= nil {
			fmt.Println(ansistring.SCol("ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
			os.Exit(50)
		} else {
			cmd.Path = lp
		}
	//}
	return cmd
}

func DoIt(shit []string) {
	cmd:=AltCommand(shit)
	o,err := cmd.Output()
	outputstring:=fmt.Sprintf("%s",o)
	fmt.Println(outputstring)
	if err!=nil{
		fmt.Println(ansistring.SCol("EXECUTION ERROR!",ansistring.A_Red,ansistring.A_Blink)+"\n"+ansistring.SCol(err.Error(),ansistring.A_Yellow,0))
		os.Exit(51)
	}
}

func AC_Vers(){
mkl.Version("JCR6 CLI (GO) - alt_command.go","17.12.09")
mkl.Lic    ("JCR6 CLI (GO) - alt_command.go","GNU General Public License 3")
}
