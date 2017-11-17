package main

import (
	"fmt"
	"text/template"
	"io/ioutil"
	"bytes"
	"sync"
	"io"
	"log"
)

var cfgPath=`D:\go\a.txt`

type Backend struct{
	Name string `json:"Name"`
	Ip []string
	ExPort int
}
type Upstream struct{
	Backends []Backend
}

func (a *Upstream) Set(name string,ip []string,exPort int) {
	a.Backends=append(a.Backends,Backend{name,ip,exPort})
}
func UpdateNginx(upstream Upstream) error{
	tmpl,_:=template.ParseFiles("a.tmpl")
	bp:=sync.Pool{
	New: func() interface{} {
		b := bytes.NewBuffer(make([]byte, 65535))
		b.Reset()
		return b
	},
	}
	tmplBuf:=bp.Get().(*bytes.Buffer)
	defer bp.Put(tmplBuf)
	tmpl.Execute(tmplBuf,upstream)
	oputBuf:=bp.Get().(*bytes.Buffer)
	for{
		str,err:=tmplBuf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("The file end is touched.")
				break
			} else {
				return err
			}
		}
		if 0 == len(str) || str == "\r\n" {
			continue
		}
		fmt.Println(str)
		oputBuf.WriteString(str)
	}
	src, _ := ioutil.ReadFile(cfgPath)
	if !bytes.Equal(src,oputBuf.Bytes()) {
		log.Println("SSS")


	//	tmpfile, err := ioutil.TempFile("", "new-nginx-cfg")
	//	defer tmpfile.Close()
	//	//diffOutput, _ := exec.Command("diff", "-u", cfgPath, tmpfile.Name()).CombinedOutput()
	//	//glog.Infof("%v\n", string(diffOutput))
	//	err= ioutil.WriteFile(cfgPath, oputBuf.Bytes(), 0644)
	//	if err != nil {
	//		return err
	//	}
	//	o, err := exec.Command("nginx", "-s", "reload", "-c", cfgPath).CombinedOutput()
	//	if err != nil {
	//		return fmt.Errorf("%v\n%v", err, string(o))
	//	}
	//	os.Remove(tmpfile.Name())
	//	return err
	} else{
		log.Println("ggg")
	}

	return nil
}

func main(){
	fmt.Println("aa")
	name:="zx.com"
	ip:=[]string{"192.168.1.2:8090","192.168.1.2:8091","192.168.1.2:8092"}
	exPort:=81
	a:=Upstream{}
	a.Set(name,ip,exPort)
	fmt.Println(a.Backends)
	UpdateNginx(a)

}