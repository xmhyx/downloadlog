package downlog

import (
    "fmt"
    "io/ioutil"
    "os"
    "time"
    "strings"
    )

 
const logconfigfile string ="log_config.txt"
const taskconfigfile string ="task_config.txt"

func Maketask()(workpool int,downrec string,allfile map[string]map[string]interface{}){

    allfiles:=make(map[string]map[string]interface{})
    jsonParse:=NewJsonStruct()
    
    var logfromconfig map[string]map[string]interface{}  
    jsonParse.Load(logconfigfile,&logfromconfig)
    
    var taskconfig map[string]interface{}
    jsonParse.Load(taskconfigfile,&taskconfig)
    
    dayoffset:=taskconfig["dayoffset"].(float64)
    workpool=int(taskconfig["concurrentpool"].(float64))
    
    the_day:=time.Now().AddDate(0,0,int(dayoffset))
    arr:=strings.Split(the_day.Format("2006-01-02"),"-")
    y,m,d:=arr[0],arr[1],arr[2]
    // y,m,d="2015","09","08"
    
    downrec = "downloaded_"+y+m+d+".txt"
	downing:= "~"+downrec
    if _, err := os.Stat(downing); err == nil || os.IsExist(err) {
       fmt.Println("download is runing ,record file:",downrec) 
       return 0,"",nil
    }
    var dr *os.File
    
    _, err := os.Stat(downrec)
    if err == nil ||os.IsExist(err){
      dr,_=os.Open(downrec)
    }else {
      dr,_= os.Create(downrec)
    }
    filedata , err :=ioutil.ReadAll(dr) 
    if err!=nil{
        fmt.Println(err.Error())   
    }
    dr.Close()
    downfiles:=strings.Split(string(filedata),"\n")

    for k,v :=range logfromconfig {
        if strings.HasPrefix(k,"log_" )&&v["isenabled"].(bool)==true {
            path_local:=v["path_local"].(string)
            logname:=v["logname"].(string)

            logname=strings.Replace(logname,"[Y]",y,-1)
            logname=strings.Replace(logname,"[m]",m,-1)
            logname=strings.Replace(logname,"[d]",d,-1)
            
            var s,tmp string
            
            if  strings.Contains(logname,"[H]"){   
               s=""
               for i:=0;i<24 ; i++ {
                 tmp=strings.Replace(logname,"[H]",fmt.Sprintf("%02d",i),-1)
                 s=k+","+path_local+","+tmp
                 allfiles[s]=v
               }
             } else {
               s=k+","+path_local+","+logname
               allfiles[s]=v
           }
       }  
     }
     for _,d:= range downfiles {
        d=strings.Replace(d, "\r", "", -1) 
        delete(allfiles,d) 
     }  
     return workpool,downrec,allfiles
}
