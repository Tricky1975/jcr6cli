/*
	JCR6 CLI
	Converter Mac Only code
	
	
	
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

/* This source only compiles on mac
 * On mac the converter is able to install homebrew for you
 * in order to install missing dependencies
 * It will however always ask permission for this!!!
 * 
 * Installing homebrew will happen through the installation
 * scripts as provided by the homebrew crew, so any update will
 * happen eventually.
 * 
 * As I have most packages installed from the start the homebrew
 * support is hard for me to check, so I must assume it works.
 */ 

import (
	//"path"
	"os"
	"trickyunits/ansistring"
	"trickyunits/mkl"
	"trickyunits/qstr"
	"strings"
	"os/exec"
	"fmt"
	"trickyunits/shell"
	//"qff"
	)

type tpack struct{
	packexecutable string
	unpackexecutable string
	packcommand string
	unpackcommand string
}

var macpacks = map[string]*tpack{}
//var me = path.Base(os.Args[0])
var brew = map[string] string{}

const suffix=""

func initpack(){
	macpacks["jcr"] = &tpack{ packexecutable:"jcr6_add",  unpackexecutable:"jcr6_extract", packcommand:"jcr6_add %s",       unpackcommand:"jcr6_extract %s -y"}
	macpacks["7z" ] = &tpack{ packexecutable:"7z",        unpackexecutable:"7z",           packcommand:"7z a %s *",         unpackcommand:"7z x %s *"}
	macpacks["zip"] = &tpack{ packexecutable:"zip",       unpackexecutable:"unzip",        packcommand:"zip -9 -r %s *",    unpackcommand:"unzip %s" }
	macpacks["tar"] = &tpack{ packexecutable:"7z",        unpackexecutable:"jcr6_extract", packcommand:"7z a -ttar %s *",   unpackcommand:"jcr6_extract %s -y"}
	macpacks["arj"] = &tpack{ packexecutable:"",          unpackexecutable:"7z",           packcommand:"",                  unpackcommand:"7z x %s *"}
	macpacks["rar"] = &tpack{ packexecutable:"rar",       unpackexecutable:"rar",          packcommand:"rar a -r -m5 %s *", unpackcommand:"rar x %s *"}
	
	// lha
	macpacks["lha"] = &tpack{ packexecutable:"lha", unpackexecutable:"lha",      packcommand:"lha a %s *",      unpackcommand:"lha x %s *"}
	macpacks["lzh"] = macpacks["lha"]
	
	// brew packages
	brew["7z"] = "p7zip"
	brew["lha"] = "lha"
	
mkl.Version("JCR6 CLI (GO) - packtools_darwin.go","17.12.09")
mkl.Lic    ("JCR6 CLI (GO) - packtools_darwin.go","GNU General Public License 3")
}

func BrewIt(barrel string){
	if _,e:=exec.LookPath("brew") ; e!= nil {
		fmt.Println(ansistring.SCol("\nYIKES",Red,0))
		fmt.Println(ansistring.SCol("You don't have homebrew. No worries, I can install it for you!",Yellow,0))
		fmt.Print(ansistring.SCol("Do you want me to install homebrew for you",Yellow,0))
		if strings.ToUpper(qstr.RawInput(ansistring.SCol(" ? ",Yellow,ansistring.A_Blink)+ansistring.SCol("(Y/N) ",Yellow,0)))!="Y" {
			fmt.Println(ansistring.SCol("Then unfortunately the road ends here.... :-/",Magenta,0))
			os.Exit(61)
		}
		ibrew:=[]string{"bash","-e","/usr/bin/ruby -e \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)\""}
		DoIt(ibrew)
	}
	ibarrel:=[]string{"brew","install",barrel}
	DoIt(ibarrel)
}


func checkpack(act,packer string){
	var ok bool
	var wp *tpack
	a:="pack"
	var want string
	if act=="u"{
		a="unpack"
	}
	if wp,ok=macpacks[packer];!ok{
		fmt.Println("ERROR!\nI don't have the required datato "+a+" files of the "+packer+".\nTry one of these:")
		for n,_:=range(macpacks){
			fmt.Println("= "+n)
		}
		os.Exit(1)
	}
	if act=="u"{
		want = wp.unpackexecutable
	} else {
		want =wp.packexecutable
	}
	if _,e := exec.LookPath(want) ; e!=nil {
		/*
		fmt.Println("ERROR!\nIn order to "+a+" anyfile of the "+packer+"type I need the program "+want+" to be present in the same folder as where the jcr6 tools are installed which is currently not the case")
		os.Exit(2)
		*/
		fmt.Println(ansistring.SCol("In order to "+a+" anyfile of the "+packer+" type I need the program "+want+" to be present which is currently not the case",Yellow,0))
		if barrel,ok:=brew[want];ok{
			fmt.Println(ansistring.SCol("But never fear. I can use Homebrew in order to install this program for you.",Yellow,0))
			fmt.Print(ansistring.SCol("Do you want me to ? (Y/N)",Yellow,0))
			if strings.ToUpper(qstr.RawInput(" "))=="Y" {
				BrewIt(barrel)
			} else {
				fmt.Println(ansistring.SCol("Then unfortunately the road ends here.... :-/",Magenta,0))
				os.Exit(60)
			}
		}
	}
}

func pack(packer,tofile string){
	var eline string
	if want,ok:=macpacks[packer];!ok{
		fmt.Println("PACK: FATAL INTERNAL ERROR!")
		os.Exit(255)
	} else {
		eline=want.packcommand
	}
	/*
	cutup:=strings.Split(eline," ")
	for i:=0;i<len(cutup);i++{
		if cutup[i]=="%s" {
			cutup[i]=tofile
		}
	}
	DoIt(cutup)
	*/
	shell.Shell(fmt.Sprintf(eline,"\""+tofile+"\""))
}

func isd(file string) bool{
	return qstr.Left(file,1)=="/"
}

func unpack(packer,fromfile string){
	var eline string
	if want,ok:=macpacks[packer];!ok{
		fmt.Println("UNPACK: FATAL INTERNAL ERROR!")
		os.Exit(255)
	} else {
		eline=want.unpackcommand
	}
	/*
	cutup:=strings.Split(eline," ")
	for i:=0;i<len(cutup);i++{
		if cutup[i]=="%s" {
			cutup[i]=fromfile
		}
	}
	DoIt(cutup)
	*/
	shell.Shell(fmt.Sprintf(eline,"\""+fromfile+"\"")) 
}

