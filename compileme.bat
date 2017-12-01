cls
type txt/preinst.txt
echo "If you are not ready to go press ctrl-C and say "yes" if Windows asks to terminate the batch job."
pause
echo Do not mind the error if the the directory already exists
md bin
        echo Compiling
rem        go build -o bin/jcr6 src/jcr6.go
rem        go build -o bin/jcr6_add src/jcr6add.go
rem        go build -o bin/jcr6_delete src/jcr6delete.go
rem        go build -o bin/jcr6_list src/jcr6list.go
rem        go build -o bin/jcr6_type src/jcr6type.go
rem        go build -o bin/jcr6_extract src/jcr6extract.go
rem        go build -o bin/jcr6_convert src/jcr6convert.go
       echo "Ready!"
