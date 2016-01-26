package downlog

import (
    "fmt"
    "os"
    "io"
    "archive/tar"
    "compress/gzip"
    "path"
)

func untargzfile(srcfile,dstpath string,) error {
    f,err:=os.Open(srcfile)
    if err!=nil {
       return err
    }
    defer f.Close()
    gz,err:=gzip.NewReader(f)
    if err!=nil {
       return err
    }
    defer gz.Close()
    tf:=tar.NewReader(gz)
    for {
       head,err:=tf.Next()
       if err!=nil {
         return err
       } 
       if err == io.EOF {
            // End of tar archive
            break
        }
       fmt.Println("untargzing file..."+head.Name)
       if head.Typeflag!=tar.TypeDir{
          df,err:=os.Create(dstpath+"/"+path.Base(head.Name))
          fmt.Println("create file..."+head.Name+" base file:"+path.Base(head.Name))
          if err!=nil {
             fmt.Println("error :",err)
             return err
          } 
          _,err=io.Copy(df,tf) 
          if err!=nil {
             return err
          } 
       }
    }
    return nil
}

func ungzfile(srcfile,dstpath string) error{
    f,err:=os.Open(srcfile)
    if err!=nil {
       return err
    }
    defer f.Close()
    
    gz,err:=gzip.NewReader(f)
    if err!=nil {
       return err
    }
    defer gz.Close()
	//dstfilename=gz.Name
	dstfilename:=srcfile[:len(srcfile)-3]
    fmt.Println("ungzing file..."+dstfilename)
    newfilename := dstpath+"/"+dstfilename
    
    df, err := os.Create(newfilename)
    if err!=nil {
       return err
    }
    defer df.Close()
    
    if _, err = io.Copy(df, gz); err != nil {
       fmt.Println(err)
       return err
    }
    return nil
}
