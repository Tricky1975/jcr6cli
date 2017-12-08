/*
	JCR6CLI
	Coverter Windows Only Code
	
	
	
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

/*
 * This is NOT tested yet!
 *
 */


import (
	"path"
	"os"
	//"qff"
	)

type tpack struct{
	packexecutable string
	unpackexecutable string
	packcommand string
	unpackcommand string
}

var winpacks = map[string]tpack{}
var me = path.Base(os.Args[0])

const suffix=".exe"

func initpack(){
	winpacks["jcr"] = { packexecutable:"jcr6_add.exe",  unpackexecutable:"jcr6_extract.exe", packcommand:"jcr6_add %s",       unpackcommand:"jcr6_extract %s"}
	winpacks["7z" ] = { packexecutable:"7z.exe",        unpackexecutable:"7z.exe",           packcommand:"7z a %s *",         unpackcommand:"7z e %s *"}
	winpacks["zip"] = { packexecutable:"zip.exe",       unpackexecutable:"unzip.exe",        packcommand:"zip -9 -r %s *",    unpackcommand:"unzip %s" }
	winpacks["tar"] = { packexecutable:"7z.exe",        unpackexecutable:"jcr6_extract.exe", packcommand:"7z a -ttar %s *",   unpackcommand:"jcr6_extract %s"}
	winpacks["arj"] = { packexecutable:"arj32.exe",     unpackexecutable:"7z.exe",           packcommand:"arj32 a -r %s",     unpackcommand:"7z e %s *"}
	
	// lha
	winpacks["lha"] = { packexecutable:"lha.exe", unpackexecutable:"lha.exe",      packcommand:"lha a %s *",      unpackcommand:"lha x %s *"}
	winpacks["lzh"] = winpacks["lha"]
	
mkl.Version("JCR6 CLI (GO) - packtools_windows.go","17.12.08")
mkl.Lic    ("JCR6 CLI (GO) - packtools_windows.go","GNU General Public License 3")
}

func d(file string) bool{
	return qstr.Left(file,1)="/" || qstr.Mid(file,2)==":"
}


func checkpack(act,packer string){
	var ok bool
	var wp tpack
	var a:="pack"
	var want string
	if act=="u"{
		a:="unpack"
	}
	if wp,ok=winpacks[packer];!ok{
		fmt.Println("ERROR!\nI don't have the required datato "+a+" files of the "+packer+".\nTry one of these:")
		for n,_:=range(winpacks){
			fmt.Println("= "+n)
		}
		os.Exit(1)
	}
	if act=="u"{
		want = winpacks[packer].unpackexecutable
	} else {
		want = winpacks[packer].packexecutable
	}
	if _,e := exec.LookPath(me+"/"+want) {
		fmt.Println("ERROR!\nIn order to "+a+" anyfile of the "+packer+"type I need the program "+want+" to be present in the same folder as where the jcr6 tools are installed which is currently not the case")
		os.Exit(2)
	}
}

func pack(packer,tofile string){
}

