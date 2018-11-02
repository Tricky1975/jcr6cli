// License Information:
// 	JCR6CLI
// 	Add files to a JCR6 resource
// 	
// 	
// 	
// 	(c) Jeroen P. Broks, 2017, 2018, All rights reserved
// 	
// 		This program is free software: you can redistribute it and/or modify
// 		it under the terms of the GNU General Public License as published by
// 		the Free Software Foundation, either version 3 of the License, or
// 		(at your option) any later version.
// 		
// 		This program is distributed in the hope that it will be useful,
// 		but WITHOUT ANY WARRANTY; without even the implied warranty of
// 		MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 		GNU General Public License for more details.
// 		You should have received a copy of the GNU General Public License
// 		along with this program.  If not, see <http://www.gnu.org/licenses/>.
// 		
// 	Exceptions to the standard GNU license are available with Jeroen's written permission given prior 
// 	to the project the exceptions are needed for.
// Version: 18.11.02
// End License Information
/*
	JCR6CLI
	Add files to a JCR6 resource
	
	
	
	(c) Jeroen P. Broks, 2017, 2018, All rights reserved
	
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
Version: 18.06.13
*/
package main

import(
	// Go
	"os"
	"fmt"
	"log"
	"flag"
	"path"
	"math"
	"strings"
	"strconv"

	// Internal
	"jcr6cli/src/imps/ver"
_	"jcr6cli/src/imps/drv"

	// Tricky units
	"trickyunits/mkl"
	"trickyunits/qff"
	"trickyunits/qstr"
	"trickyunits/tree"
	"trickyunits/ansistring"
	"trickyunits/jcr6/jcr6main"
)


const debuglist = false


const Red = ansistring.A_Red
const Yellow = ansistring.A_Yellow
const Cyan = ansistring.A_Cyan
const Magenta = ansistring.A_Magenta
const Green = ansistring.A_Green
const Blue = ansistring.A_Blue

const Blink  = ansistring.A_Blink
const Bright = ansistring.A_Bright

var cyes = ansistring.SCol("YES",Green,Bright)
var cno  = ansistring.SCol("NO!",Red,0)

var counterrors = 0
var countwarnings = 0

var goed = 0
var goeddep = 0

func ERR(msg string, fatal ...bool) {
	f:=false
	if len(fatal)>0 { f=fatal[0] } // Optional parameters are rather clumsy in Go
	log.Print(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol(msg,Yellow,0))
	counterrors++
	if f { log.Fatal(ansistring.SCol("This error is fatal!",Magenta,0)) }
}

func WARNING(msg string) {
	log.Print(ansistring.SCol("ERROR! ",Red,Blink)+ansistring.SCol(msg,Yellow,0))
	countwarnings++
}

func plural(number int,nounsingle,nounplural string) string{
	if number==1 { return "1 "+nounsingle } else { return fmt.Sprintf("%d %s",number,nounplural) }
}

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6add.go","18.11.02")
mkl.Lic    ("JCR6 CLI (GO) - jcr6add.go","GNU General Public License 3")
}

type tAddMe struct {
	files []map[string]string
	alias map[string]string
	comments map[string]string
	config map[string]string
	cf int
	ca string
	skipprefix []string // Only taken notice of when merging JCR files
	destroyoriginal bool
	fatstorage string
	inputdir string
	aliasfrom string
	impkind map[string]string
	impsig  map[string]string
}

func (a *tAddMe) fa (k,v string) {
	if a.cf<0 { ERR("HUH?",true) } // Protection. When this error goes of we have a very serious bug, as this should be impossible!
	a.files[a.cf][k] = v
}

func (a *tAddMe) NewFileRec(file string){
	a.files = append(a.files,map[string]string{})
	a.cf = len(a.files)-1
	a.fa("!__FROMFILE",file)
}

func CreateAddMe(inputdir string) tAddMe {
	r:=tAddMe{}
	r.files = []map[string]string{}
	r.alias = map[string]string{}
	r.comments = map[string]string{}
	r.config = map[string]string{}
	r.skipprefix = []string{}
	r.impkind=map[string]string{}
	r.impsig=map[string]string{}
	r.cf=-1
	r.ca=""
	r.destroyoriginal=false
	r.inputdir = inputdir
	return r
}

func ParseJIF(JIF,inputdir string) tAddMe {
	// Let's parse a "JCR Instruction File" "JIF"
	fmt.Println(ansistring.SCol("Parsing instructions",Cyan,0))
	r:=CreateAddMe(inputdir)
	if !qff.Exists(JIF) { ERR("ParseJIF(\""+JIF+"\"): File not found",true) }
	lines:=qff.GetLines(JIF)
	for tln,tlt:=range(lines){
		ln:=tln+1
		lt:=qstr.MyTrim(tlt)
		if lt!="" && qstr.Left(lt,1)!="#"{   // Ignore empty lines and comments
			p:=strings.Index(lt,":")
			if p<1 {
				ERR(fmt.Sprintf("JIF: %s %d: Illegal line",JIF,ln))
			} else {
				cmd := qstr.MyTrim(strings.ToUpper(lt[:p]))
				prm := qstr.MyTrim(lt[p+1:])
				switch cmd {
					case "DESTROYORIGINAL","DESTROYORIGINALJCR":
						r.destroyoriginal=strings.ToUpper(qstr.Left(prm,1))=="Y"
					case "CONFIG":
						pcfg:=strings.Index(prm,",")
						if pcfg<1 {
							ERR(fmt.Sprintf("JIF: %s %d: Illegal config definition",JIF,ln))
						} else {
							r.config[prm[:pcfg]] = prm[pcfg+1:]
						}
					case "FATSTORAGE":
						r.fatstorage=prm
					case "INPUTDIR":
						r.inputdir=prm
					case "COMMENT":
						pcmt:=strings.Index(prm,",")
						if pcmt<1 {
							ERR(fmt.Sprintf("JIF: %s %d: Illegal comment definition",JIF,ln))
						} else {
							r.comments[prm[:pcmt]] = strings.Replace(prm[pcmt+1:],"\\n","\n",-10)
						}
					case "ALIAS","ALIASFROM":
						r.aliasfrom = prm
					case "AS":
						if r.aliasfrom=="" {
							ERR(fmt.Sprintf("JIF: %s %d: AS command can only be used after an ALIAS request has been done!",JIF,ln))
						} else {
							r.alias[prm]=r.aliasfrom
						}
					case "FILE":
						sfile:=strings.Replace(prm,"\\","/",-1)
						if qstr.Mid(sfile,2,1)!=":" && qstr.Left(sfile,1)!="/" { sfile = inputdir+"/"+prm }
						r.NewFileRec(sfile)
					case "TARGET","AUTHOR","NOTES","STORAGE","MERGE","MERGESTRIPEXT":
						if r.cf<0 {
							ERR(fmt.Sprintf("JIF: %s %d: %s command can only be used after a FILE request has been done!",JIF,ln,cmd))
						} else {
							r.fa("!__"+cmd,prm)
						}
					case "JCRSKIPPREFIX":
						r.skipprefix = append(r.skipprefix,prm)
					case "CD":
						os.Chdir(prm)
						WARNING("Although CD has been executed, the current setup of this tool might void its effect, and the CD command has therefore been deprecated")
					case "IMPORT","REQUIRE":
						pim:=strings.Index(prm,";")
						file:=""
						sig:=""
						if pim<1 {
							sig=""
							file=prm
						} else {
							sig=prm[:pim+1]
							file=prm[:pim]
						}
						r.impkind[file]=cmd
						r.impsig [file]=sig
					default:
						ERR(fmt.Sprintf("JIF: %s %d: I don't understand command %s",JIF,ln,cmd))
				}
			}
		}
	}
	return r
}

func ParseList(listfile,idir,cm,fc string,doj bool,author,notes string) tAddMe {
	fmt.Println(ansistring.SCol("Reading file list",Cyan,0))
	r:=CreateAddMe(idir)
	r.destroyoriginal=doj
	r.fatstorage=fc
	if !qff.Exists(listfile) { ERR("List file "+listfile+" has not been found",true) }
	for tln,tline:=range qff.GetLines(listfile) {
		line:=qstr.MyTrim(tline)
		ln:=tln+1
		if line!=""{
			sfile:=strings.Replace(line,"\\","/",-100)
			tfile:=strings.Replace(line,"\\","/",-100)
			if qstr.Mid(tfile,2,1)==":" {
				tfile=qstr.Right(tfile,len(tfile)-1)
				WARNING(fmt.Sprintf("%s %d: Windows drive letter detected. This could lead to odd results!",listfile,ln))
			}
			if qstr.Left(tfile,1)=="/" {
				tfile=qstr.Right(tfile,len(tfile)-1)
				WARNING(fmt.Sprintf("%s %d: Full path name detected. This could lead to odd results!",listfile,ln))
			}
			if qstr.Mid(sfile,2,1)!=":" && qstr.Left(sfile,1)=="/" { sfile = idir + "/" + sfile }
			if !qff.Exists(sfile) {
				ERR(fmt.Sprintf("%s %d: The file \"%s\" has not been found!",path.Base(listfile),ln,sfile))
			} else {
				r.NewFileRec(sfile)
				r.fa("!__TARGET",tfile)
				r.fa("!__STORAGE",cm)
				r.fa("!__AUTHOR",author)
				r.fa("!__NOTES",notes)
			}
		}
	}
	return r
}

func TreeDir(idir,cm,fc,pre,suf string, doj bool,author,notes string) tAddMe {
	fmt.Println(ansistring.SCol("Analysing dir: "+idir,Cyan,0))
	r:=CreateAddMe(idir)
	r.destroyoriginal=doj
	r.fatstorage=fc
	t:=tree.GetTree(idir,false)
	for _,fil:=range(t){
		ok:=true
		if pre!="" { ok = ok && qstr.Left (fil,len(pre))==pre }
		if suf!="" { ok = ok && qstr.Right(fil,len(suf))==suf }
		sfile:=idir+"/"+fil
		tfile:=fil
		if ok {
			r.NewFileRec(sfile)
			r.fa("!__TARGET",tfile)
			r.fa("!__STORAGE",cm)
			r.fa("!__AUTHOR",author)
			r.fa("!__NOTES",notes)
		}
	}
	return r
}


func main(){
	ver.CHVER()
	ansiyes:="yes"
	if !ansistring.ANSI_Use {ansiyes="no"}
	flag.CommandLine.SetOutput(os.Stdout)
	fl_idir:=flag.String("i" ,qff.PWD(),"Directory to add files from (this works RECURSIVELY -- always)")
	fl_pref:=flag.String("pr",  ""   ,"Only add files prefixed with the given value.")
	fl_suff:=flag.String("sf",  ""   ,"Only add files suffixed with the given value.")
	fl_list:=flag.String("ls",  ""   ,"Use a text file as a list of all files.\n\tUsing this flag will ignore -pr and -sf")
	fl_jif :=flag.String("jif", ""   ,"Use an instruction file to build data from.\n\tUsing this flag will ignore -pr, -sf, -cm, -fc and -ls")
	fl_encm:=flag.String("cm" ,"lzma","Use one of the supported compression methods for packing the entries.\n\tYou can also use the method 'Store' for no compression, and 'Brute' to make JCR6 try all of them to see which is best.")
	fl_ftcm:=flag.String("fc" ,"lzma", "Use one of the supported compression methods for packing the file table.\n\tYou can also use the method 'Store' for no compression, and 'Brute' to make JCR6 try all of them to see which is best.")
	fl_ansi:=flag.String("ansi",ansiyes,"Allow using ansi in output.")
	fl_nmrg:=flag.Bool ("nm",false,"When set the addition routine will not merge files recognized by the JCR6 detector as directories")
	fl_kill:=flag.Bool ("doj",false,"When set the system will always create a new JCR6 file and destroy the old one if it exists")
	fl_author:=flag.String("author","","Define the \"Author\" field in all files added (please note, this flag will be ignored when you use jif files)")
	fl_notes:=flag.String("notes","","Define the \"Notes\" field in all files added (please note, this flag will be ignored when you use jif files)")
	//fl_yes :=flag.Bool ("y",false,"When set all existing files will be overwritten")
	flag.Parse()
	//autoyes=*fl_yes
	nonflags:=flag.Args()
	if len(nonflags)<1 {
		fmt.Print(ansistring.SCol("Usage: ",Red,0),ansistring.SCol("jcr6 extract ",Yellow,0),ansistring.SCol("[ flags ] ",Magenta,ansistring.A_Dark),ansistring.SCol("<JCR6 file> ",Cyan,0),"\n\n")
		flag.PrintDefaults()
		fmt.Print("")
		fmt.Println(ansistring.SCol("\nSupported compression methods:",Yellow,0))
		for k,_:=range jcr6main.JCR6StorageDrivers{
			if k==strings.ToLower(k){
				fmt.Print(ansistring.SCol("= ",Red,Bright),ansistring.SCol(k,Yellow,0),"\n")
			}
		}
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
	updatefile:=wjcr
	// Show parsing results
	update:=(!*fl_kill) && (qff.Exists(wjcr))
	if update {
		ERR("Updating existing resources is not yet supported! Please check later versions for that. If you want to replace this resource, please add the -doj flag.",true)
		fmt.Print(ansistring.SCol("Updating JCR:       ",Yellow,0))
		idx:=0
		for qff.Exists(wjcr+".update."+fmt.Sprint(idx)+".tmp") { idx++ }
		updatefile=wjcr+".update."+fmt.Sprint(idx)+".tmp"
	} else {
		fmt.Print(ansistring.SCol("Creating JCR:       ",Yellow,0))
	}
	fmt.Println(ansistring.SCol(wjcr,Cyan,0))
	fmt.Print(ansistring.SCol("Input directory:    ",Yellow,0))
	fmt.Println(ansistring.SCol(*fl_idir,Cyan,0))
	if *fl_pref!="" && *fl_list=="" && *fl_jif=="" {
		fmt.Print(ansistring.SCol("Prefix:             ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_pref,Cyan,0))
	}
	if *fl_suff!="" && *fl_list=="" && *fl_jif=="" {
		fmt.Print(ansistring.SCol("Suffix:             ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_suff,Cyan,0))
	}
	if *fl_jif!="" {
		fmt.Print(ansistring.SCol("Instruction file:   ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_jif,Cyan,0))

	} else if *fl_list!="" {
		fmt.Print(ansistring.SCol("List:               ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_list,Cyan,0))
	}
	if *fl_jif=="" {
		fmt.Print(ansistring.SCol("Entry compression:  ",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_encm,Cyan,0))
		fmt.Print(ansistring.SCol("Filetab compression:",Yellow,0))
		fmt.Println(ansistring.SCol(*fl_ftcm,Cyan,0))
		if *fl_author!=""{
			fmt.Print(ansistring.SCol("Author:             ",Yellow,0))
			fmt.Println(ansistring.SCol(*fl_author,Cyan,0))
		}
		if *fl_notes!=""{
			fmt.Print(ansistring.SCol("Notes:              ",Yellow,0))
			if len(*fl_notes)>20 {
			  r:=strings.Replace(*fl_notes,"\n","\\n",-1)
				r=r[:20]
				fmt.Println(ansistring.SCol(r,Cyan,0))
			} else {
			  fmt.Println(ansistring.SCol(strings.Replace(*fl_notes,"\n","\\n",-1),Cyan,0))
			}
		}
		fmt.Print(ansistring.SCol("Merging             ",Yellow,0))
		if *fl_nmrg { fmt.Println(cno) } else {fmt.Println(cyes)}
	}
	fmt.Print("\n\n")

	// Let's get a clear view on which files to add and which not!
	list:=tAddMe{}
	if *fl_jif!="" { // JIF instruction file, this overrides ALL other settings, although it will use the directory set with -i as base dir unless a full pathname is present in the instruction file
		list = ParseJIF(*fl_jif,*fl_idir)
	} else if *fl_list!="" {
		list = ParseList(*fl_list,*fl_idir,*fl_encm,*fl_ftcm,*fl_kill,*fl_author,*fl_notes)
	} else {
		list = TreeDir(*fl_idir,*fl_encm,*fl_ftcm,*fl_pref,*fl_suff,*fl_kill,*fl_author,*fl_notes)
	}
	if debuglist {fmt.Println(list)} // debug line to keep Go from WHINING during development!

	if len(list.files)==0 && len(list.comments)==0 && len(list.impkind)==0 { ERR("There seems to be nothing I can do!") }

	// Create file
	jc:=jcr6main.JCR_Create(updatefile,list.fatstorage)

	// Configuration
	for k,v:=range list.config {
		tk:=k[1:]
		switch qstr.Left(k,1){
			case "!": // Ignore.... Nothing to do here, "!" is reserved for system declarations only
			case "$": jc.ConfigString(tk,v)
			case "&": jc.ConfigBool(tk,strings.ToUpper(v)=="TRUE")
			case "%":
				i,err:=strconv.ParseInt(v,32,32)
				if err!=nil { ERR("Config nummeric error while converting config var \""+k+"\": "+err.Error()); i=0}
				jc.ConfigInt(tk,int32(i))
			default:
				ERR("Unknown type for config variable "+k)
		}
	}

	// Update
	// This comes later!

	// Add comments
	fmt.Println(ansistring.SCol("Adding comments",Cyan,0))
	for k,v :=range list.comments { jc.AddComment(k,v); fmt.Println(ansistring.SCol("= ",Red,Bright)+ansistring.SCol(k,Yellow,0)) }

	// Add dependencies
	fmt.Println(ansistring.SCol("Adding dependency requests",Cyan,0))
	for f,knd :=range list.impkind {
		sig:=list.impsig[f]
		fmt.Println(ansistring.SCol("= ",Red,Bright)+ansistring.SCol(knd,Blue,Bright)+" "+ansistring.SCol(f,Yellow,0)+" "+ansistring.SCol(sig,Magenta,0))
		switch knd {
			case "IMPORT":  jc.AddImport(f,sig); goeddep++
			case "REQUIRE": jc.AddRequire(f,sig); goeddep++
			default:        ERR("Unknown dependency type: "+knd)
		}
	}

	// Add entries
	fmt.Println(ansistring.SCol("Adding raw files",Cyan,0))
	for _,entry := range list.files {
		//originalfile,entryname,algorithm,author,notes
		if _,ok:=entry["!__AUTHOR"];!ok{ entry["!__AUTHOR"]=""}
		if _,ok:=entry["!__NOTES"] ;!ok{ entry["!__NOTES"] =""}
		if _,ok:=entry["!__FROMFILE"];!ok{ ERR(fmt.Sprintf("INTERNAL ERROR! Please report! No source file received!\n%s",entry),true) }
		entry["!__FROMFILE"]=strings.Replace(entry["!__FROMFILE"],"//","/",-100)
		fmt.Print(ansistring.SCol("= ",Red,Bright)+ansistring.SCol(entry["!__FROMFILE"],Yellow,0)+ansistring.SCol(" ... ",Cyan,0))
		if !qff.Exists(entry["!__FROMFILE"]){
			fmt.Println(ansistring.SCol("FAILED",Red,0))
			ERR("File not found!")
		} else if ISJCR:=jcr6main.Recognize(entry["!__FROMFILE"]);ISJCR!="NONE" {
			fmt.Println(ansistring.SCol("Merging: "+ISJCR,Magenta,0))
			jdm:=jcr6main.Dir(entry["!__FROMFILE"])
			if jcr6main.JCR6Error!="" {
				ERR(jcr6main.JCR6Error)
			} else {
				for k,v :=range jdm.Comments { jc.AddComment(k,v); fmt.Println(ansistring.SCol("  = Comment: ",Cyan,0)+ansistring.SCol(k,Yellow,0)) }
				for _,e :=range jdm.Entries  {
					fmt.Print(ansistring.SCol("  = ",Red,Bright)+ansistring.SCol(e.Entry,Yellow,0)+ansistring.SCol(" ... ",Cyan,0))
					tm:=int64(0)
					if tmc,ok:=e.Dataint["__TIMESTAMP"];ok{ tm=int64(tmc) }
					tgt:=entry["!__TARGET"]+"/"+e.Entry
					//fmt.Print(tgt+" ") // debug
					csize,storage:=jc.AddData(jcr6main.JCR_B(jdm,e.Entry),tgt,entry["!__STORAGE"],e.UnixPerm,tm,e.Author,e.Notes)
					rsize:=e.Size
					if jcr6main.JCR6Error!="" {
						ERR(jcr6main.JCR6Error)
					} else if storage=="Store" {
						fmt.Println(ansistring.SCol("stored",ansistring.A_White,0))
						goed++
					} else {
						deel    := float64(csize)
						geheel  := float64(rsize)
						procent := (deel/geheel)*100
						ratio   := int( math.Floor(procent + .5) )
						pteken:="%"
						fmt.Println(ansistring.SCol(fmt.Sprintf("%s: Reduced to %d%s",storage,ratio,pteken),Green,Bright))
						goed++
					}
				}
			}
		} else {
			rsize,csize,storage:=jc.AddFile(entry["!__FROMFILE"],entry["!__TARGET"],entry["!__STORAGE"],entry["!__AUTHOR"],entry["!__NOTES"])
			if jcr6main.JCR6Error!="" {
				ERR(jcr6main.JCR6Error)
			} else if storage=="Store" {
				fmt.Println(ansistring.SCol("stored",ansistring.A_White,0))
				goed++
			} else {
				deel    := float64(csize)
				geheel  := float64(rsize)
				procent := (deel/geheel)*100
				ratio   := int( math.Floor(procent + .5) )
				pteken:="%"
				fmt.Println(ansistring.SCol(fmt.Sprintf("%s: Reduced to %d%s",storage,ratio,pteken),Green,Bright))
				goed++
			}
		}
	}

	// Alias requests
	fmt.Println(ansistring.SCol("Processing aliases",Cyan,0))
	for alias,from:=range list.alias{
		fmt.Print(ansistring.SCol(from,Red,0))
		fmt.Print(ansistring.SCol(" => ",Yellow,0))
		fmt.Print(ansistring.SCol(alias,Green,0))
		fmt.Print(ansistring.SCol(" ... ",Cyan,0))
		centryname:=strings.ToUpper(from)
		if _,ok:=jc.Entries[strings.ToUpper(centryname)] ; ok {
			jc.AliasFile(from,alias)
			e:=jcr6main.JCR6Error
			if e=="" {
				fmt.Println(ansistring.SCol("Done",7,Bright))
			} else {
				fmt.Println(ansistring.SCol("Failed",Red,0))
				ERR(e)
			}
		} else {
			fmt.Println(ansistring.SCol("Failed -- No original",Red,0))
			counterrors++
		}
	}

	// Closure
	fmt.Println(ansistring.SCol("Finalizing JCR6 file",Cyan,0))
	jc.Close()

	// Replace files and clean temp files in case of update
	// This comes later
	if update { fmt.Print("I'll work with "+updatefile+" later") } // not used error prevention!


	// Last
	fmt.Println("\n")
	if counterrors>0 { fmt.Println(ansistring.SCol("\t"+plural(counterrors,"error","errors"),Red,0)+" "+ansistring.SCol("occurred during the entire process!",Yellow,0)) }
	if countwarnings>0 { fmt.Println(ansistring.SCol("\t"+plural(countwarnings,"warning","warnings"),Red,0)+" "+ansistring.SCol("occurred during the entire process!",Yellow,0)) }
	if len(list.comments)>0 { fmt.Println(ansistring.SCol("\t"+plural(len(list.comments),"comment","comments"),Green,Bright)+" "+ansistring.SCol("were processed!",Yellow,0)) }
	if goeddep>0 { fmt.Println(ansistring.SCol("\t"+plural(goeddep,"dependency request","dependency requests"),Green,Bright)+" "+ansistring.SCol("were succesfully processed!",Yellow,0)) }
	if goed>0 { fmt.Println(ansistring.SCol("\t"+plural(goed,"entry","entries"),Green,Bright)+" "+ansistring.SCol("were succesfully processed!",Yellow,0)) }
	fmt.Println("\n")
}
