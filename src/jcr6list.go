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
Version: 17.12.02
*/
package main



import(
	"jcr6cli/src/imps/ver"
_	"jcr6cli/src/imps/drv"
	"fmt"
	"trickyunits/mkl"
)

func init(){
mkl.Version("JCR6 CLI (GO) - jcr6list.go","17.12.02")
mkl.Lic    ("JCR6 CLI (GO) - jcr6list.go","GNU General Public License 3")
}

func main(){
	ver.CHVER()
	fmt.Printf("Not functional yet, but that will soon come, I guess ;)\n\n")
}
