clear
cat txt/preinst.txt
read -r -p "Do you wish to continue? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY]) 
        mkdir -p bin
        echo Compiling
        echo = Main; go build -o bin/jcr6 src/jcr6.go
        echo = Script; go build -o bin/jcr6_script src/jcr6script.go;cp src/jcr6script.lua bin/jcr6_script.lua
        echo = Add;  go build -o bin/jcr6_add src/jcr6add.go
#       echo = Delete;  go build -o bin/jcr6_delete src/jcr6delete.go
        echo = List; go build -o bin/jcr6_list src/jcr6list.go
        echo = Type; go build -o bin/jcr6_type src/jcr6type.go
        echo = Extract; go build -o bin/jcr6_extract src/jcr6extract.go
        echo = Convert; go build -o bin/jcr6_convert jcr6cli/src/jcr6convert.go
        echo = Config;  go build -o bin/jcr6_config jcr6cli/src/jcr6config.go
        echo = Version; go build -o bin/jcr6_version src/jcr6version.go
        echo "Ready!"

        ;;
    *)
        echo Then I guess I\'ll see you back when you DO want to compile this.
        ;;
esac

echo "Ok"
