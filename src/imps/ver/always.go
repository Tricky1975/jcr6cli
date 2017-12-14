package ver

import(
	"os"
	"fmt"
	"trickyunits/qstr"
	"trickyunits/mkl"
)


func WANTVER() bool {
	if len(os.Args)>1{
		return os.Args[1]==jvstring || os.Args[1]==fvstring
	}
	return false
}

func CHVER(){
	if len(os.Args)>1{
		if os.Args[1]==jvstring || os.Args[1]==fvstring{
			fmt.Println(qstr.StripAll(os.Args[0])+"\t"+"v"+mkl.Newest()+" go")
			if os.Args[1]==fvstring{
				fmt.Print(mkl.ListAll())
			}
			os.Exit(0)
		}
	}
}

