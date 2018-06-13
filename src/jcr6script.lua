--[[
	JCR6 CLI
	Additional Script
	
	
	
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
Version: 18.04.13
]]


----------------------------------------------------------------------
--[[


	This script must come along with the script utility of 
	JCR6 (the Go version that is). It must be present in the same
	directory as the jcr6_script binary and have the same name 
	(suffixed with .lua of course).
	
	This file contains some important definitions the 
	script utility needs. It's best NOT to modify it in any way
	unless you really know what you are doing.
	
	One important note, although this script is licensed under the
	terms of the GPL, I hereby make one exception clause for the
	scripts you load with the JCR6 scripter to which this script
	will automatically be linked. Your own scripts may have any
	license you deem right. The GPL for this script counts for its
	being a part of the JCR6 cli utilities, and when being 
	distributed apart from that utility package.


]]
----------------------------------------------------------------------



function Add(src,tgt,data)
	JCR_ResetData()
	for k,v in pairs(data or {}) do
		JCR_SetData(k,v)
	end
	JCR_Add(src,tgt)
end

function AddImport(dependency,sig)
	value = dependency
	if sig then value = value .. ";" .. sig end
	Output("IMPORT:"..value)
end

function AddRequire(dependency,sig)
	value = dependency
	if sig then value = value .. ";" .. sig end
	Output("REQUIRE:"..value)
end

function AddComment(name,cmt)
	Output(name..":"..cmt)
end

function JCR6MergeSkipPrefix(prefix)
	Output("JCRSKIPPREFIX:"..prefix)
end

function Sig(sig)
	Output("SIGNATURE:"..sig)
end

--[[
	This table is created to replace the original BlitzMax library.
	It's not needed in Go, but I set it up to ensure backward 
	compatibility
]]
JLS = {
	SetJCR6OutputFile = SetJCR6OutputFile,
	Output = JCR_Output,
	FType = filetype
}
Output = JCR_Output

function fileexists(file)
	return JLS.FType(file)~=0
end 

function isfile(file)
	return JLS.FType(file)==1
end

function isdir(file)
	return JLS.FType(file)==2
end

function GetDir(path)
	--print("test")
	local ggetdir,e = load(JCR_GetDir(path))
	if not ggetdir then print("ERROR: ") print(e) end
	--print(type(ggetdir))
	return ggetdir()
end getdir = GetDir




-- mkl 
mkl={}
function mkl.version(a,b)
	JCRMKL("VER",a,b)
end

function mkl.lic(a,b)
	JCRMKL("LIC",a,b)
end

mkl.version("JCR6 CLI (GO) - jcr6script.lua","18.04.13")
mkl.lic    ("JCR6 CLI (GO) - jcr6script.lua","GNU General Public License 3")



-- String features --
upper = string.upper
lower = string.lower
chr = string.char
printf = string.format
replace = string.gsub
rep = string.rep
substr = string.sub


function cprintf(a,b)
print(printf(a,b))
end

function len(a)
local k,v
local ret=0
if not a then return 0 end
if type(a)=="table" then
  --for k,v in ipairs(a) do
  --    ret = ret + 1
  --    end
  return #a
  end
return string.len(a.."") -- the .."" is to make sure this is string formatted! ;)  
end

function left(s,l)
return substr(s,1,l)
end

function right(s,l)
local ln = l or 1
local st = s or "nostring"
-- return substr(st,string.len(st)-ln,string.len(st))
return substr(st,-ln,-1)
end 

function mid(s,o,l)
local ln=l or 1
local of=o or 1
local st=s or ""
return substr(st,of,(of+ln)-1)
end 


function trim(s)
  return (s:gsub("^%s*(.-)%s*$", "%1"))
end
-- from PiL2 20.4

function findstuff(haystack,needle) -- BLD: Returns the position on which a substring (needle) is found inside a string or (array)table (haystrack). If nothing if found it will return nil.<p>Needle must be a string if haystack is a string, if haystack is a table, needle can be any type.
local ret = nil
local i
for i=1,len(haystack) do
    if type(haystack)=='table'  and needle==haystack[i] then ret = ret or i end
    if type(haystack)=='string' and needle==mid(haystack,i,len(needle)) then ret = ret or i end
    -- rint("finding needle: "..needle) if ret then print("found at: "..ret) end print("= Checking: "..i.. " >> "..mid(haystack,i,len(needle)))
    end
return ret    
end


function safestring(s)
local allowed = "qwertyuiopasdfghjklzxcvbnmmQWERTYUIOPASDFGHJKLZXCVBNM 12345678890-_=+!@#$%^&*()[]{};:|,.<>/?"
--local i
local safe = true
local alt = ""
assert ( type(s)=='string' , "safestring expects a string not a "..type(s) )
for i=1,#s do
    safe = safe and (findstuff(allowed,mid(s,i,1))~=nil)
    alt = alt .."\\"..string.byte(mid(s,i,1),1)
    end
-- print("DEBUG: Testing string"); if safe then print("The string "..s.." was safe") else print("The string "..s.." was not safe and was reformed to: "..alt) end    
local ret = { [true] = s, [false]=alt }
-- print("returning "..ret[safe])
return ret[safe]     
end 



-- Serializing
__serialize_work = __serialize_work or {
                ["nil"]        = function(vvalue) return "nil" end,
                ["number"]     = function(vvalue) return vvalue end,
                ["function"]   = function(vvalue) SysError("Cannot serialize functions") return "ERROR" end,
                ["string"]     = function(vvalue) return "\""..safestring(vvalue).."\"" end,
                ["boolean"]    = function(vvalue) return ({[true]="true",[false]="false"})[vvalue] end,
                ["table"]      = function(vvalue)
                                 local titype
                                 local tindex = {
                                                   ["number"]     = function(v) return v end,
                                                   ["boolean"]    = function(v) return ({[true]="true",[false]="false"})[v] end,
                                                   ["string"]     = function(v) return "\""..safestring(v).."\"" end
                                 }
                                 local wrongindex = function() SysError("Type "..titype.." can not be used as a table index in serializing") return "ERROR" end
                                 local ret = "{"
                                 local k,v
                                 local result
                                 local notfirst
                                 for k,v in pairs(vvalue) do
                                     if notfirst then ret = ret .. ",\n" else notfirst=true ret = ret .."\n" end
                                     titype = type(k)
                                     result = (tindex[titype] or wrongindex)(k)
                                     -- print(titype.."/"..k)
                                     ret = ret .. TRUE_SERIALIZE("["..result.."]",v,(tabs or 0)+1,true)                                      
                                     end
                                 --if notfirst then    
                                 --  ret = ret .."\n"    
                                 --  for i=1,tabs or 0 do ret = ret .."     " end   
                                 --  for i=1,len(vname.." = ") do ret = ret .. " " end
                                 --  end 
                                 ret = ret .. "}"  
                                 return ret  
                                 end 
                                   
             } 
                
function TRUE_SERIALIZE(vname,vvalue,tabs,noenter)
local ret = ""
--local work = __serialize_work   
--[[
if type(vvalue=='string') then                   
	Console.Write('Serialize("'..vname..'","'..vvalue..'",'..(tabs or 0)..')')
elseif type(vvalue=='number' then
   	Console.Write('Serialize("'..vname..'",'..vvalue..','..(tabs or 0)..')')
else
	Console.Write('Serialize("'..vname..'",<'..type(vvalue)..'>,'..(tabs or 0)..')')
endif	
]]	
local letsgo = __serialize_work[type(vvalue)] or function(vvalue) error("Unknown type. Cannot serialize","Variable,"..vname..";Type Value,"..type(vvalue)) end
local i
for i=1,tabs or 0 do ret = ret .."       " end
ret = ret .. vname .." = "..letsgo(vvalue) 
if not noenter then ret = ret .."\n" end
return ret
end


function serialize(vname,variableitself)
local ret = ""
local v = variableitself or _G[vname]
assert(type(vname)=='string',"First variable must be the name to return in the serializer string")
ret = TRUE_SERIALIZE(vname,v)
--JBCSYSTEM.Returner(ret)
return ret
end

