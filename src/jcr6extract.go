/*
	JCR6CLI
	Extractor
	
	
	
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
Version: 17.12.10
*/
package main

import (
	// internal
	"jcr6cli/src/imps/ver"
_	"jcr6cli/src/imps/drv"

	// Go Libraries
	"flag"
	"log"
	"fmt"
	"os"
	"path"
	"strings"

	// My own libraries
	"trickyunits/jcr6/jcr6main"
	"trickyunits/qff"
	"trickyunits/ansistring"
	"trickyunits/mkl"
	"trickyunits/qstr"
	)

const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green
const Blue = ansistring.A_Blue

const Blink  = ansistring.A_Blink
const Bright = ansistring.A_Bright

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6extract.go","17.12.10")
mkl.Lic    ("JCR6 CLI (GO) - jcr6extract.go","GNU General Public License 3")
}

var autoyes bool

func YES(question string) bool{
	if autoyes { return true }
	fmt.Print(question) // debug to prevent non-compile due to unused variable.
	return false
}

func ERR(msg string, fatal ...bool) {
	f:=false
	if len(fatal)>0 { f=fatal[0] } // Optional parameters are rather clumsy in Go
	log.Print(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol(msg,Yellow,0))
	if f { log.Fatal(ansistring.SCol("This error is fatal!",Magenta,0)) }
}

func FromList(d jcr6main.TJCR6Dir,lst string) []string {
	ret:=[]string{}
	if !qff.Exists(lst) { ERR("List file '"+lst+"' not found!",true) }
	tl:=qff.GetLines(lst)
	for ln,lfi:=range tl{
		lf:=qstr.MyTrim(lfi)
		if lf!="" {
			if jcr6main.HasEntry(d,lf) {
				ret=append(ret,lf)
			} else {
				ERR(fmt.Sprintf("Requested file \"%s\" in line %d has NOT been found in this resource/archive",lf,ln+1))
			}
		}
	}
	return ret
}

func FromPrefixAndSuffix(d jcr6main.TJCR6Dir, apre, asuf string) []string {
	el:=jcr6main.EntryList(d)
	rt:=[]string{}
	pre:=strings.ToUpper(apre)
	suf:=strings.ToUpper(asuf)
	for _,entry:=range el {
		ok:=true
		centry:=strings.ToUpper(entry)
		ok = ok && (pre=="" || qstr.Left (centry,len(pre))==pre)
		ok = ok && (suf=="" || qstr.Right(centry,len(suf))==suf)
		if ok {rt=append(rt,entry) }		
	}
	return rt
}

func failed(fmsg string) {
	fmt.Println(ansistring.SCol("failed!",Red,0))
	ERR(fmsg)
}

func main() {
	ver.CHVER()
	// Init flags
	ansiyes:="yes"
	if !ansistring.ANSI_Use {ansiyes="no"}
	fl_odir:=flag.String("o" ,qff.PWD(),"Directory to extract the content of the JCR6 file to")
	fl_pref:=flag.String("pr","","Only extract files prefixed with the given value.")
	fl_suff:=flag.String("sf","","Only extract files suffixed with the given value.")
	fl_list:=flag.String("ls","","Use a text file as a list of all files.\n\tUsing this flag will ignore -pr and -sf")
	fl_ansi:=flag.String("ansi",ansiyes,"Allow using ansi in output.")
	fl_yes :=flag.Bool ("y",false,"When set all existing files will be overwritten")
	flag.Parse()
	autoyes=*fl_yes
	nonflags:=flag.Args()
	if len(nonflags)<1 {
		fmt.Print(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 extract ",Yellow,0),ansistring.SCol("[ flags ] ",Magenta,ansistring.A_Dark),ansistring.SCol("<JCR6 file> ",Cyan,0),"\n\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *fl_ansi=="yes" {
		ansistring.ANSI_Use=true
	} else if *fl_ansi=="no" {
		ansistring.ANSI_Use=false
	} else if *fl_ansi!="" {
		log.Print(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol("Invalid value for ansi. Only 'yes' or 'no' would do!",Yellow,Bright))
	}	
	wjcr:=nonflags[0]
	
	// Show parsing results
	fmt.Print(ansistring.SCol("Analysing JCR:      ",Yellow,0))
	fmt.Println(ansistring.SCol(wjcr,Cyan,0))
	wtype:=jcr6main.Recognize(wjcr)
	if wtype=="NONE" {
		log.Fatal(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol("Requested JCR6 file unrecognized",Yellow,Bright))
	}
	fmt.Print(ansistring.SCol("Resource Type:      ",Yellow,0))
	fmt.Println(ansistring.SCol(wtype,Cyan,0))
	if *fl_pref!="" && *fl_list=="" {
		fmt.Print(ansistring.SCol("Prefix:             ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_pref,Cyan,0))
	}
	if *fl_suff!="" && *fl_list=="" {
		fmt.Print(ansistring.SCol("Suffix:             ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_suff,Cyan,0))
	}
	if *fl_list!="" {
		fmt.Print(ansistring.SCol("List:               ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_list,Cyan,0))
	}
	fmt.Print(ansistring.SCol("Output directory:   ",Yellow,0))
	fmt.Println(ansistring.SCol(*fl_odir,Cyan,0))

	if (!qff.Exists(wjcr)) { ERR("JCR file not found!",true) } // this one is basically checked above, but this is an extra security thing to make sure NOTHING can go wrong on this department!
	jd:=jcr6main.Dir(wjcr)
	if jcr6main.JCR6Error!="" { ERR(jcr6main.JCR6Error,true) }
	
	var list[] string // Which files do we want to unpack
	if *fl_list!="" { 
		list = FromList(jd,*fl_list) 
	} else if *fl_pref!="" || *fl_suff!="" {
		list = FromPrefixAndSuffix(jd,*fl_pref,*fl_suff)
	} else {
		list = jcr6main.EntryList(jd)
	}
	if len(list)==0 { ERR("Nothing to do! No files seem to match any request!",true) }
	fmt.Print(ansistring.SCol("Files to extract    ",Yellow,0))
	fmt.Println(ansistring.SCol(fmt.Sprintf("%d",len(list)),Cyan,0))

	// fmt.Println(list) // debug line

	// Comments
	for k,v := range jd.Comments{
		fmt.Println(ansistring.SCol(k,Cyan,0))
		for i:=0;i<len(k);i++{
			fmt.Print(ansistring.SCol("=",Yellow,0))
		}
		fmt.Print("\n")
		fmt.Println(v)
		fmt.Print("\n")
	}


	// Let's extract all this crap
	fmt.Println(ansistring.SCol("\n\nExtracting files:",Cyan,0))
	goed:=0
	fout:=0
	for _,efile:=range list {
		fmt.Print(ansistring.SCol("= ",Red,Bright))
		fmt.Print(ansistring.SCol(efile,Yellow,0))
		fmt.Print(ansistring.SCol(" ... ",Cyan,0))
		needdir:=path.Dir((*fl_odir)+"/"+efile)
		err:=os.MkdirAll(needdir,0777)
		if err!=nil{
			failed(err.Error())
			fout++
		} else if qff.Exists((*fl_odir)+"/"+efile) && !*fl_yes {
				failed("File already exists. Set -y if you want to overwrite it!")
				fout++
		} else {
			jcr6main.JCR_Extract(jd,efile,(*fl_odir)+"/"+efile)
			if jcr6main.JCR6Error=="" {
				fmt.Println(ansistring.SCol("extracted",Green,Bright))
				goed++
			} else {
				failed(jcr6main.JCR6Error)
				fout++
			}
		}
	}
	if goed>0{
		fmt.Print(ansistring.SCol("Files extracted:    ",Yellow,0))
		fmt.Println(ansistring.SCol(fmt.Sprint(goed),Green,Bright))
	} 
	if fout>0{
		fmt.Print(ansistring.SCol("Files failed:       ",Yellow,0))
		fmt.Println(ansistring.SCol(fmt.Sprint(fout),Red,0))
	}

}

