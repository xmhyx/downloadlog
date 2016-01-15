package main
import (
       "./downlog"
       )

func main() {  
    pool,downrec,task:=downlog.Maketask()
    downlog.Runtask(pool,downrec,task) 
}
