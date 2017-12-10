echo off
cls
type txt\preinst.txt
echo "If you are not ready to go press ctrl-C and say "yes" if Windows asks to terminate the batch job."
pause
echo Do not mind the error if the the directory already exists
md bin
        echo Compiling
rem        go build -o bin/jcr6.exe src/jcr6.go
rem        go build -o bin/jcr6_add.exe src/jcr6add.go
rem        go build -o bin/jcr6_delete.exe src/jcr6delete.go
        echo = List
        go build -o bin/jcr6_list.exe src/jcr6list.go
        echo = Type
        go build -o bin/jcr6_type.exe src/jcr6type.go
        go build -o bin/jcr6_extract.exe src/jcr6extract.go
        echo = Convert
        go build -o bin/jcr6_convert.exe jcr6cli/src/jcr6convert
       echo "Ready!"
