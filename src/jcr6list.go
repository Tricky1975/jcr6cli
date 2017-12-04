/*
	JCR6 CLI
	List out JCR6 file
	
	
	
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
Version: 17.12.04
*/
package main



import(
	"jcr6cli/src/imps/ver"
_	"jcr6cli/src/imps/drv"
	"fmt"
	"os"
	"sort"
	"math"
	"trickyunits/mkl"
	"trickyunits/ansistring"
	"trickyunits/qstr"
	"trickyunits/jcr6/jcr6main"
)

const A_Magenta=ansistring.A_Magenta
const A_Cyan   =ansistring.A_Cyan
const A_Yellow =ansistring.A_Yellow

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6list.go","17.12.04")
mkl.Lic    ("JCR6 CLI (GO) - jcr6list.go","GNU General Public License 3")
}

func main(){
	ver.CHVER()
	//fmt.Printf("Not fully functional yet, but that will soon come, I guess ;)\n\n")
	if len(os.Args)<2{
		fmt.Print(ansistring.SCol("Usage:",ansistring.A_Cyan,0),ansistring.SCol("jcr6 list ",ansistring.A_Yellow,0),ansistring.SCol("<JCR6 Resource File>",ansistring.A_Magenta,0),"\n")
		os.Exit(0)
	}
	jcrfilename:=os.Args[1]
	fmt.Printf("%s %s\n\n",ansistring.SCol("Reading:",A_Cyan,0),ansistring.SCol(jcrfilename,A_Magenta,0))
	jcr:=jcr6main.Dir(jcrfilename)
	if jcr6main.JCR6Error!="" {
		fmt.Printf("%s\n%s\n",ansistring.SCol("ERROR",ansistring.A_Red,ansistring.A_Blink),ansistring.SCol(jcr6main.JCR6Error,ansistring.A_Yellow,0))
		os.Exit(20)
	}
	// Comments
	for k,v := range jcr.Comments{
		fmt.Println(ansistring.SCol(k,A_Cyan,0))
		for i:=0;i<len(k);i++{
			fmt.Print(ansistring.SCol("=",A_Yellow,0))
		}
		fmt.Print("\n")
		fmt.Println(v)
		fmt.Print("\n")
	}
	// Main files analysis
	maincodes := make(map[string] string)
	maintypes := make(map[string] string)
	maincount := make(map[string] int)
	jent:=jcr6main.EntryList(jcr)
	for i:=0;i<len(jent);i++{
		mf:=jcr6main.Entry(jcr,jent[i]).Mainfile
		c:=qstr.Left(qstr.StripAll(mf)+"________",8)
		maincodes [mf] = c
		maintypes [mf] = jcr6main.Recognize(mf)
		if _,ok:=maincount[mf];!ok{
			maincount[mf]=0
		}
		maincount[mf]++
	}
	mainorder := make([]string,len(maincodes))
	mci := 0
	for k,_ := range maincodes{
		mainorder[mci]=k
		mci++
	}
	sort.Strings(mainorder)
	fmt.Println(ansistring.SCol("MainCode  Type                  Entries  Main File",A_Cyan,0))
	fmt.Println(ansistring.SCol("========  ====================  =======  =========",A_Yellow,0))
	for i:=0;i<len(mainorder);i++{
		fn:=mainorder[i]
		fmt.Print  (ansistring.SCol(maincodes[fn]+"  ",1,0))
		fmt.Print  (ansistring.SCol(qstr.Left(jcr6main.JCR6Drivers[maintypes[fn]].Drvname+"                           ",20)+"  ",2,0))
		fmt.Print  (ansistring.SCol(qstr.Right(fmt.Sprintf("       %d",maincount[fn]),7),3,0)+"  ")
		fmt.Println(ansistring.SCol(fn,7,0))
	}
	fmt.Print(ansistring.SCol("\tThis resource has ",A_Yellow,0),ansistring.SCol(fmt.Sprintf("%d",len(mainorder)),A_Cyan,0))
	if len(mainorder)==1{
		fmt.Println(ansistring.SCol(" main file\n",A_Yellow,0))
	} else {
		fmt.Println(ansistring.SCol(" main files\n",A_Yellow,0))
	}
	// Entry list
	fmt.Println(ansistring.SCol(" Real Size Comp. Size Ratio   Offset MainCode    Storage Entry",A_Cyan,0))
	fmt.Println(ansistring.SCol("========== ========== ===== ======== ======== ========== =====",A_Yellow,0))
	for i:=0;i<len(jent);i++{
		e:=jcr6main.Entry(jcr,jent[i])
		fmt.Print(ansistring.SCol(qstr.Right(fmt.Sprintf("          %d",e.Size),10),1,0)," ")
		fmt.Print(ansistring.SCol(qstr.Right(fmt.Sprintf("          %d",e.Compressedsize),10),2,0)," ")
		ratio:=0
		if e.Size>0 {
			deel    := float64(e.Compressedsize)
			geheel  := float64(e.Size)
			procent := (deel/geheel)*100
			ratio   = int( math.Floor(procent + .5) )
		}
		fmt.Print(ansistring.SCol(qstr.Right(fmt.Sprintf("    %d%s",ratio,"%"),5),3,0)," ")
		fmt.Print(ansistring.SCol(qstr.Right(fmt.Sprintf("        %X",e.Offset),8),4,0)," ")
		fmt.Print(ansistring.SCol(maincodes[e.Mainfile],5,0)," ")
		fmt.Print(ansistring.SCol(qstr.Right("               "+e.Storage,10),6,0)," ")
		fmt.Println(ansistring.SCol(e.Entry,7,0))
	}
	fmt.Print(ansistring.SCol("\tThis resource has ",A_Yellow,0),ansistring.SCol(fmt.Sprintf("%d",len(jent)),A_Cyan,0))
	if len(jent)==1{
		fmt.Println(ansistring.SCol(" entry\n",A_Yellow,0))
	} else {
		fmt.Println(ansistring.SCol(" entries\n",A_Yellow,0))
	}

}
