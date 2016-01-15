package downlog

import (
    "fmt"
    "time"
    "os"
    "io"
    "strings"
    "sort"
)
    
func task_dispatch(task map[string]map[string]interface{},pool chan bool,result_chan chan string) {
    var keys []string
    for k := range task {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
      print("waiting pool ...")
      <-pool
      fmt.Println("distpatch task:",k)
      go worker(k,task[k],result_chan)
    }
}

func worker(log string,logconfig map[string]interface{}, result_chan chan string){
/*    
    fmt.Printf("server: %s, user: %s, pass:%s \n",logconfig["hostname"].(string),
                                                  logconfig["username"].(string),
                                                  logconfig["password"].(string))
*/
    arr:=strings.Split(log,",")
    filename:=arr[2]
    srcfile:=logconfig["path_remote"].(string)+"/"+filename
    dstfile:=logconfig["path_local"].(string)+"/~"+filename
    fmt.Printf("from: %s, to: %s\n",srcfile,dstfile)
    time.Sleep(time.Second * 2)
    result_chan <-log  
}

func writedownloaded(filename, result string){
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
       return 
    }
    io.WriteString(f, result+"\n") 
    f.Close()
}

func Runtask(workpool int,downrec string,task map[string]map[string]interface{}) {
    if downrec=="" {
       return
    }

    var pool = make(chan bool, workpool)  
    var result_chan = make(chan string, workpool)
    err:=os.Rename(downrec,"~"+downrec )
    fmt.Println("main rename ",downrec," to ","~"+downrec,"err :",err) 
    
    for i:=0;i<workpool;i++{
        pool<-true
    }

    go task_dispatch(task,pool,result_chan)   

    for i:=0;i<len(task);i++{
        r:=<-result_chan   
        writedownloaded("~"+downrec, r)
        arr:=strings.Split(r,",")
        tmpfile:=arr[1]+"/~"+arr[2]
        dstfile:=arr[1]+"/"+arr[2]
        //err := os.Rename(tmpfile,dstfile )
        fmt.Println(i," rename ",tmpfile," to ",dstfile) 
        pool<-true
    }
    err=os.Rename("~"+downrec,downrec )
    fmt.Println("main rename ","~"+downrec," to ",downrec,"err :",err) 
}
