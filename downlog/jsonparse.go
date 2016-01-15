package downlog

import (
    "io/ioutil"
    "os"
    json "encoding/json"
    )

type JsonStruct struct{
}    
 
func NewJsonStruct() *JsonStruct{
    return &JsonStruct{}
}

func (self *JsonStruct) Load(filename string,v interface{}){
    fi,err := os.Open(filename)   
    if err != nil{panic(err)}   
    defer fi.Close()   
    data,err := ioutil.ReadAll(fi) 
    if err != nil{panic(err)}   
     
    datajson:=[]byte(data)
    err = json.Unmarshal(datajson,v)
    if err!=nil{return}  
}    

