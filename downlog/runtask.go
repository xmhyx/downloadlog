package downlog

import (
    "fmt"
    "os"
 //   "io"
    "strings"
    "sort"
    "ftp"
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

    arr:=strings.Split(log,",")
    filename:=arr[2]
    srcfile:=logconfig["path_remote"].(string)+"/"+filename
    dstfile:=logconfig["path_local"].(string)+"/~"+filename
    fmt.Printf("from: %s, to: %s\n",srcfile,dstfile)
	ftp:= new(ftp.FTP)                                                         
    ftp.Debug = true
    ftp.Connect(logconfig["hostname"].(string), int(logconfig["port"].(float64)))  
    ftp.Login(logconfig["username"].(string), logconfig["password"].(string))
	ftp.Request("TYPE I") 
    err:=ftp.Retr(srcfile,dstfile) 
	if err!=nil  {
	  result_chan<-""
	}else{
      result_chan <-log  
	}
}

func writedownloaded(filename, result string){
    f, err := os.OpenFile(filename ,os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
	   fmt.Println("open downloaded file failed ",err)
	   return
    }
    _,e:=f.WriteString( result+"\n")
	if e != nil {
	   fmt.Println("write downloaded file failed ",e)
    }else{
	   fmt.Println("write downloaded ok : ",result)
	}
    f.Close()
	
}

func Runtask(workpool int,downrec string,task map[string]map[string]interface{}) {
    if downrec=="" {
       return
    }

    var pool = make(chan bool, workpool)  
    var result_chan = make(chan string, workpool)
	var downing="~"+downrec 
    err:=os.Rename(downrec,downing )
    fmt.Println("main rename ",downrec," to ",downing,"err :",err) 
    
    for i:=0;i<workpool;i++{
        pool<-true
    }

    go task_dispatch(task,pool,result_chan)   

    for i:=0;i<len(task);i++{
        r:=<-result_chan   
		if r!=""{
          writedownloaded(downing, r)
          arr:=strings.Split(r,",")
          tmpfile:=arr[1]+"/~"+arr[2]
          dstfile:=arr[1]+"/"+arr[2]
          os.Rename(tmpfile,dstfile )
		  if strings.HasSuffix( dstfile ,".tar.gz") {
            untargzfile(dstfile,arr[1])
          }else if strings.HasSuffix( dstfile ,"gz") {
            ungzfile(dstfile,arr[1])
          }
          //fmt.Println(i," rename ",tmpfile," to ",dstfile) 
	      }
        pool<-true
    }
    err=os.Rename(downing,downrec )
    fmt.Println("main rename ",downing," to ",downrec,"err :",err) 
}
