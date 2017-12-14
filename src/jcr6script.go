/*
	
	
	
	
	
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
Version: 17.12.14
*/
package main


import (
	// Go libs
	"os"
	"log"
	"fmt"
	"strings"
	"runtime"
	"os/exec"
	
//	"io"
//	"bufio"
	
	// 3rd party libs
	"github.com/shopify/go-lua"
	
	// Tricky's Units
	"trickyunits/mkl"
	"trickyunits/qff"
	"trickyunits/qstr"
	"trickyunits/tree"
	"trickyunits/ansistring"
	
	// Internal libs
	"jcr6cli/src/imps/ver"
	)
	
const AddDebug = false
	
const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green
const Blue = ansistring.A_Blue

const Blink  = ansistring.A_Blink
const Bright = ansistring.A_Bright

// Errors
func ERR(msg string, fatal ...bool) {
	f:=false
	if len(fatal)>0 { f=fatal[0] } // Optional parameters are rather clumsy in Go
	log.Print(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol(msg,Yellow,0))
	if f { log.Fatal(ansistring.SCol("This error is fatal!",Magenta,0)) }
}


// LUA API
var API_Data map[string]string
var JIF_File string
var sl *lua.State
var jif = "# GENERATED WITH JCR6_SCRIPT\n\n" // Here the contents saved to the jif file will be stored!
var ffields = []string{"Notes","Author","Storage"}
var jcrfile = ""


func AddChat(msg string) { if AddDebug { fmt.Println(ansistring.SCol("Add: ",Cyan,0)+" "+ansistring.SCol(msg,Yellow,0)) }}


func TrueAdd(src,tgt string){
	if qff.IsDir(src) {
		AddChat(src+" is a directory, so I'm gonna add all its contents as "+tgt)
		trl:=tree.GetTree(src,false)
		for _,f:=range trl {
			TrueAdd(src+"/"+f,tgt+"/"+f)
		}
	} else if qff.Exists(src) {
		rtgt:=strings.Replace(tgt,"\\","/",-1)
		if rtgt=="" {rtgt=strings.Replace(src,"\\","/",-1) }
		if qstr.Mid(rtgt,2,1)==":" { rtgt = qstr.Right(rtgt,len(rtgt)-2) }
		if qstr.Mid(rtgt,1,1)=="/" { rtgt = qstr.Right(rtgt,len(rtgt)-1) }
		if qstr.Right(rtgt,1)=="/" { rtgt += qstr.StripDir(src) }
		AddChat("Adding file: "+src+" as "+rtgt)
		jif+="FILE:"+src+"\nTARGET:"+rtgt+"\n"
		for _,k:=range(ffields) {
			if v,ok:=API_Data[k];ok{
				jif+=k+":"+v+"\n"
			}
		}
	} else {
		ERR("Raw file "+src+" does not exist")
	}
}
func API_Add(l *lua.State) int{
	src:=lua.CheckString(l,1)
	tgt:=""
	if !l.IsNil(2) { tgt=lua.CheckString(l,2) }
	TrueAdd(src,tgt)
	return 0
}

func API_ResetData(l *lua.State) int{
	API_Data = map[string]string{}
	return 0
}

func API_SetData(l *lua.State) int {
	k:=lua.CheckString(l,1)
	v:=lua.CheckString(l,2)
	API_Data[k]=v
	return 0
}

func API_Platform(l *lua.State) int {
	l.PushString(runtime.GOOS)
	return 1
}

func PANIEK(l *lua.State) int {
	ERR("Error in Lua Script!")
	return 0
}

func API_Use(l *lua.State) int {
	u:=lua.CheckString(l,1)
	DoLua(u)
	return 0
}

func API_SetJCR6OutputFile(l *lua.State) int {
	jcrfile = lua.CheckString(l,1)
	return 0
}

func API_JIF(l *lua.State) int{
	JIF_File = lua.CheckString(l,1)
	return 0
}

var oldalias = ""
func API_Alias(l *lua.State) int {
	src:=lua.CheckString(l,1)
	as :=lua.CheckString(l,2)
	if src!=oldalias {
		jif+="ALIAS:"+src+"\n"
	}
	oldalias=src
	jif+="AS:"+as+"\n"
	return 0
}

func API_Output(l *lua.State) int {
	jif+=lua.CheckString(l,1)+"\n"
	return 0
}

func API_MKL(l * lua.State) int {
	what:=lua.CheckString(l,1)
	file:=lua.CheckString(l,2)
	data:=lua.CheckString(l,3)
	switch what{
		case "LIC": mkl.Lic(file,data)
		case "VER": mkl.Version(file,data)
	}
	return 0
}

func SysError(l *lua.State) int {
	ERR(lua.CheckString(l,1),!l.IsNoneOrNil(2))
	return 0
}

var used map[string]bool = map[string]bool{}
// Lua processing

func DoLua(file string) {
	if d,ok:=used[file];ok{
		if d { return }
	}
	used[file]=true
	fmt.Print  (ansistring.SCol("Compiling: ",Yellow,0))
	fmt.Println(ansistring.SCol(file,Cyan,0))
	lua.LoadFile(sl,file,"")
	sl.Call(0,0)
}

func DoLuaSilent(file string) {
	if d,ok:=used[file];ok{
		if d { return }
	}
	used[file]=true
	//fmt.Print  (ansistring.SCol("Compiling: ",Yellow,0))
	//fmt.Println(ansistring.SCol(file,Cyan,0))
	lua.LoadFile(sl,file,"")
	sl.Call(0,0)
}

func DoLuaString(scriptstring,desc string) {
	fmt.Print  (ansistring.SCol("Compiling: ",Yellow,0))
	fmt.Println(ansistring.SCol(desc,Magenta,0))
	lua.LoadString(sl,scriptstring)
	sl.Call(0,0)
}

// Saved variables parsing
func savedvars() (string,string,string){
	sfile:=qstr.StripExt(os.Args[1])+".lsv"
	cfile:=qstr.StripExt(os.Args[1])+".savedvars.lua"
	if !qff.Exists(sfile) { return "","","" }
	script:="function veryverysecret__savedvariables__yeahyeah()\nret = ''\n"
	lines:=qff.GetLines(sfile)
	for i,l := range(lines) {
		ln:=qstr.MyTrim(l)
		if qstr.Left(strings.ToUpper(ln),5)=="SAVE " {
			sv:=strings.Split(ln," ")
			if len(sv)!=2 {
				ERR(fmt.Sprintf("Invalid lsv instruction in line #%d",i+1))
			} else {
				script+="\tret = ret .. serialize('"+sv[1]+"',"+sv[1]+")..'\\n'\n"
			}
		}
	}
	script+="\n\nJCR_ResetData() JCR_SetData('serializer',ret)\n"
	script+="\treturn ret\nend\n"
	return script,sfile,cfile
}


// Main
func init(){
	sl = lua.NewState()
	lua.BaseOpen(sl)
	lua.OpenLibraries(sl)
	lua.SetFunctions(sl, []lua.RegistryFunction{ {"JCR_Add",API_Add},
		                                         {"JCR_ResetData",API_ResetData},
		                                         {"JCR_SetData",API_SetData},
		                                         {"Platform",API_Platform},       // JCR_ is a prefix only used for commands the user scripts should not DIRECTLY call (unless they don't give a damn about backward compatibility that is). Platform can just be called directly, although I cannot guarantee the results might be the same in all underlying programming languages used to code the APIs.
		                                         {"Use",API_Use},
		                                         {"Alias",API_Alias},
		                                         {"SetJCR6OutputFile",API_SetJCR6OutputFile},
		                                         {"JIF",API_JIF},
		                                         {"Output",API_Output},
		                                         {"JCRMKL",API_MKL},
		                                       },0 )
	lua.AtPanic(sl,PANIEK)
mkl.Version("JCR6 CLI (GO) - jcr6script.go","17.12.14")
mkl.Lic    ("JCR6 CLI (GO) - jcr6script.go","GNU General Public License 3")
}

func main(){
	fexe,_:=os.Executable() //(string, error)
	flua:=qstr.StripExt(fexe)+".lua"
	if !qff.Exists(flua) { ERR(flua+" not found! It must be in the same directory as where the "+qstr.StripAll(fexe)+" binary is located. Please note that JCR6 should NOT be called through symlinks!",true) }
	if ver.WANTVER(){
		DoLuaSilent(flua)
	}
	ver.CHVER()
	if len(os.Args)<2 {
		fmt.Print(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 script ",Yellow,0),ansistring.SCol("<Lua Script> ",Cyan,0),ansistring.SCol("[ script-args ] ",Magenta,ansistring.A_Dark),"\n\n")
		fmt.Println("If you have configurations variables inside your script you can create the scriptstring with the same name as your Lua Script,\nbut with .lua replaced ith .lsv and put the line \"SAVE myconfigvar\" in it. You can put in as many as you want!")
		fmt.Println("Based on what you put in the script the scripter will create a JCR6 Instruction File (jif)\nand launch the 'Add' features to run it and create the actual JCR6 scriptstring")
		os.Exit(0)
	}
	slua:=os.Args[1]
	jcrfile  = qstr.StripExt(slua)+".jcr"
	JIF_File = qstr.StripExt(slua)+".jif"
	alua := "Args = {}\n"
	for i:=1;i<len(os.Args);i++ { alua += "Args[#Args+1]=\"" + strings.Replace(os.Args[i],"\"","\\\"",-10) + "\"\n" }
	svlua,svluaf,svluafsv := savedvars()
	DoLuaString(alua,"Arguments definitions")
	if svlua!="" { DoLuaString(svlua, svluaf) }
	if svluafsv!="" && qff.Exists(svluafsv) { DoLua(svluafsv) }
	DoLua(flua)
	DoLua(slua)
	//fmt.Println("Result:\n"+jif) // debug
	fmt.Println(ansistring.SCol("  Writing: ",Yellow,0)+ansistring.SCol(JIF_File,Cyan,0))
	bt,err:=os.Create(JIF_File)
	if err!=nil {
		ERR(err.Error(),true)
	}
	qff.RawWriteString(bt,jif)
	bt.Close()
	if svluafsv!="" {
		fmt.Println(ansistring.SCol("  Writing: ",Yellow,0)+ansistring.SCol(svluafsv,Cyan,0))
		sl.Global("veryverysecret__savedvariables__yeahyeah")
		sl.Call(0,0)
		if serialized,ok:=API_Data["serializer"];ok{
			bt,err=os.Create(svluafsv)
			if err!=nil {
				ERR(err.Error(),true)
			} else {
				qff.RawWriteString(bt,serialized)
			}
			bt.Close()			
		} else {
			ERR("No serializer data received, so nothing's saved!")
		}
	}
	fmt.Printf("\n\n")
	// Now let's make the add program do the rest
	addutil := strings.Replace(fexe,"jcr6_script","jcr6_add",-1)
	c:=exec.Command(addutil,"-doj","-jif",JIF_File,jcrfile)
	c.Stdout = os.Stdout
	c.Run()
}
/*
	stdout, _ := c.StdoutPipe()

	
    c.Start()
    //go print(stdout)
	oneByte := make([]byte, 100)
	//num := 1
	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		fmt.Println(string(line))
		//num = num + 1
		//if num > 3 {
		//	os.Exit(0)
		//}
	}
    c.Wait()

}

// to print the processed information when stdout gets a new line
func print(stdout io.ReadCloser) {
     r := bufio.NewReader(stdout)
     //line, _, err := r.ReadLine()
     //fmt.Println("line: %s err %s", line, err)
     l,_,e:=r.ReadLine()
     fmt.Printf("%s\n",l)
     if e!=nil {
		 ERR(e.Error())
	 }
}
*/

