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
Version: 17.12.16
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
	"trickyunits/dirry"
	"trickyunits/tree"
	"os"
)

const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green
const Blue = ansistring.A_Blue

const Blink  = ansistring.A_Blink
const Bright = ansistring.A_Bright

var workdir=dirry.Dirry("$AppSupport$/$LinuxDot$JCR6G/Convert")
var swapdir=workdir+"/Swap"
var verdir=workdir+"/Verify"
var sessiondir,verifydir,sessionnum = NewSessionDir()
var launchdir=qff.PWD()

var sizes  = make (map[string]int )
var hashes = make (map[string]string )
var skipped= make (map[string]bool )
var verfail= make (map[string]bool )


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

func NewSessionDir() (string,string,int) {
	ret:=""
	retver:=""
	n:=0
	err:=os.MkdirAll(swapdir, 0777)
	if err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
		os.Exit(9)
	}
	for{
		n++
		ret=swapdir+fmt.Sprintf("/SESSION%d",n)
		retver=verdir+fmt.Sprintf("/SESSION%d",n)
		//fmt.Println("Sessiondir: "+ret)
		if (!qff.IsDir(ret)) && (!qff.Exists(ret)) && (!qff.IsDir(retver)) && (!qff.Exists(retver)) {
			return ret,retver,n
		}
	}
}

func MakeWorkDirs(){
	err:=os.MkdirAll(sessiondir, 0777)
	if err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
		os.Exit(8)
	}
	err=os.MkdirAll(verifydir, 0777)
	if err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
		os.Exit(8)
	}
	err=os.Chdir(sessiondir)
	if err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
		os.Exit(8)
	}
}

func fout(file,note string){
	fmt.Print  (ansistring.SCol("= ",Red,Bright))
	fmt.Print  (ansistring.SCol(file,Cyan,0)+ansistring.SCol(": ",Blue,0))
	fmt.Println(ansistring.SCol(note,Magenta,0))
}

func main(){
	// Version show if asked
	ver.CHVER()
	// Flags
	flag.CommandLine.SetOutput(os.Stdout)
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
	fmt.Print(ansistring.SCol("Session ID:      ",Yellow,0))
	fmt.Println(ansistring.SCol(fmt.Sprintf("%d",sessionnum),Cyan,0))
	fmt.Print(ansistring.SCol("Launched from:   ",Yellow,0))
	fmt.Println(ansistring.SCol(launchdir,Cyan,0))
	fmt.Print(ansistring.SCol("Verification:    ",Yellow,0))
	fmt.Println(yn(verify))
	
	// All required dependencies accounted for?
	checkpack("u",source)
	for _,tf := range targets{
		checkpack("p",tf)
	}
	
	// Let's create a temp directory then (this function will also make
	// the system change the directory to this directory
	MakeWorkDirs()
	
	// Let's unpack
	fmt.Println(ansistring.SCol("\n\nUnpacking: ",Cyan,0),ansistring.SCol(iname,Magenta,0))
	unpack(source,iname)
	
	// Let's get all file data
	ctree:=tree.GetTree(qff.PWD(),true)
	fmt.Println(ansistring.SCol("\n\nChecking unpacked data",Cyan,0))
	for i,v:=range ctree{
		fmt.Printf("  Checking %d/%d\r",i,len(ctree))
		sizes[v]=qff.FileSize(v)
		if verify { 
			hashes[v] = qff.MD5File(v) 
		}
	}
	
	// Let's pack in all formats
	for _,tar:=range targets {
		tarname := oname + "." + tar
		fmt.Println(ansistring.SCol("\n\nPacking: ",Cyan,0),ansistring.SCol(tarname,Magenta,0))
		ok:=!qff.Exists(tarname)
		if (!ok){
			fmt.Println(ansistring.SCol("WARNING!",Red,Blink))
			fmt.Println(ansistring.SCol("That file already exists.",Yellow,0))
			ok=strings.ToUpper(qstr.RawInput(ansistring.SCol("Do you want to overwrite it ",Yellow,0)+ansistring.SCol("?",Cyan,Blink)+ansistring.SCol(" (Y/N) ",Cyan,0)))=="Y"
			if ok {
				fmt.Println(ansistring.SCol("Deleting original: ",Cyan,0),ansistring.SCol(tarname,Magenta,0))
				err:=os.Remove(tarname)
				if err!=nil{
					fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
					fmt.Println(ansistring.SCol("I could not delete: "+tarname,Yellow,0))
					fmt.Println(ansistring.SCol(err.Error(),Blue,ansistring.A_Bright))
					ok=false
				}
			}
		}
		// Ok got redefined on the way, so a new if and no else statements!
		if ok {
			pack(tar,tarname)
			skipped[tar]=false
		} else {
			fmt.Println(ansistring.SCol("Request skipped!",Red,0))
			skipped[tar]=true
		}
	}
	// Cleanup
	fmt.Println(ansistring.SCol("\n\nCleaning up conversion swap",Cyan,0))
	if err:=os.RemoveAll(sessiondir);err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol("I could not delete: "+sessiondir,Yellow,0))
		fmt.Println(ansistring.SCol(err.Error(),Blue,ansistring.A_Bright))
		fmt.Println(ansistring.SCol("This directory can therefore linger on your system taking up space, so you'll have to do this manually!",Yellow,ansistring.A_Bright) )
	}
	// Verifydir
	if verify {
		for _,tar:=range targets{
			tarname := oname + "." + tar
			fmt.Println(ansistring.SCol("\n\nVerifying: ",Cyan,0),ansistring.SCol(tarname,Magenta,0))
			if !skipped[tar]{
				err:=os.Mkdir(verifydir+"/"+tar,0777)
				if err!=nil{
						fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
					fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
					os.Exit(18)
				}
				err=os.Chdir(verifydir+"/"+tar)
				if err!=nil{
					fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
					fmt.Println(ansistring.SCol(err.Error(),Yellow,0))
					os.Exit(18)
				}
				unpack(tar,tarname)
				ctree=tree.GetTree(qff.PWD(),true)
				fail:=false
				for _,fl:=range ctree{
					if _,hok:=hashes[fl];!hok{ fout(fl,"No hash found in memory! Did this file exist before?") ; fail=true }
					if _,hok:=sizes[fl];!hok { fout(fl,"No size found in memory! Did this file exist before?") ; fail=true }
				}
				for fl,h:=range hashes{
					if !qff.Exists(fl) { 
						fout(fl,"File not found!")
						fail=true
					} else if qff.MD5File(fl)!=h {
						fout(fl,"Hash mismatch    "+qff.MD5File(fl)+" != "+h)
						fail=true
					}
				}
				for fl,h:=range sizes{
					if !qff.Exists(fl) { 
						fout(fl,"File not found!")
						fail=true
					} else if qff.FileSize(fl)!=h {
						fout(fl,fmt.Sprintf("Size mismatch    %d !=%d",qff.FileSize(fl),h) )
						fail=true
					}
				}
				verfail[tar]=fail
			}
		}
	}
	fmt.Println(ansistring.SCol("\n\nCleaning up verification swap",Cyan,0))
	if err:=os.RemoveAll(verifydir);err!=nil{
		fmt.Println(ansistring.SCol("ERROR!",Red,Blink))
		fmt.Println(ansistring.SCol("I could not delete: "+sessiondir,Yellow,0))
		fmt.Println(ansistring.SCol(err.Error(),Blue,ansistring.A_Bright))
		fmt.Println(ansistring.SCol("This directory can therefore linger on your system taking up space, so you'll have to do this manually!",Yellow,ansistring.A_Bright) )
	}
	
	// The end
	osize:=qff.FileSize(iname)
	fmt.Println("\n\n\n")
	fmt.Print(ansistring.SCol("Original file:         ",Yellow,0)); fmt.Println(ansistring.SCol(iname,Cyan,0))
	fmt.Print(ansistring.SCol("Size:                  ",Yellow,0)); fmt.Println(ansistring.SCol(fmt.Sprintf("%d",osize),Cyan,0))
	for _,tar:=range targets{
		tarname := oname + "." + tar
		fmt.Print(ansistring.SCol("Created file:          ",Yellow,0)); fmt.Println(ansistring.SCol(tarname,Cyan,0))
		fmt.Print(ansistring.SCol("Status:                ",Yellow,0))
		//sok:=true
		if skipped[tar]{
			fmt.Println(ansistring.SCol("SKIPPED",Red,0))
			//sok=false
		} else if !qff.Exists(tarname) {
			fmt.Println(ansistring.SCol("CONVERSION FAILED",Red,0))
			//sok=flase
		} else if verify && verfail[tar] {
			fmt.Println(ansistring.SCol("VERIFICATION FAILED",Red,0))
		} else {
			fmt.Println(ansistring.SCol("DONE",Green,0))
			fmt.Print(ansistring.SCol("Size:                  ",Yellow,0))
			tsize:=qff.FileSize(tarname)
			fmt.Println(ansistring.SCol(fmt.Sprintf("%d",tsize),Cyan,0))
			tdiff:=osize-tsize
			if tsize<osize{
				fmt.Print(ansistring.SCol("Shrunk:                ",Yellow,0))
				fmt.Println(ansistring.SCol(fmt.Sprintf("%d",tdiff),Cyan,0))
			} else if tsize>osize{
				fmt.Print(ansistring.SCol("Grown:                 ",Yellow,0))
				fmt.Println(ansistring.SCol(fmt.Sprintf("%d",tdiff*(-1)),Cyan,0))
			}
		}
	}
}
