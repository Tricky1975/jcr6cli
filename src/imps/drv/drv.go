// License Information:
// 	JCR6 CLI
// 	Driver list
// 	
// 	
// 	
// 	(c) Jeroen P. Broks, 2017, 2019
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
// 	Please note that some references to data like pictures or audio, do not automatically
// 	fall under this licenses. Mostly this is noted in the respective files.
// 		
// Version: 19.02.28
// End License Information
/*
	JCR6 CLI
	Driver list
	
	
	
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
Version: 17.12.13
*/
package drv
import(
	// compression
	_ "trickyunits/jcr6/jcr6zlib"
	_ "trickyunits/jcr6/jcr6lzma"
	_ "trickyunits/jcr6/jcr6lzw"
	_ "trickyunits/jcr6/jcr6flate"
	_ "trickyunits/jcr6/jcr6jxsrcca"
//	_ "trickyunits/jcr6/jcr6xz"

	// file formats
	_ "trickyunits/jcr6/jcr6realdir"
	_ "trickyunits/jcr6/jcr6wad"
)	
