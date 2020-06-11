package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type fileinfo struct {
	fullpath string
	name     string
	isdir    bool
}

var listfileinfo = make([]fileinfo, 0, 10000)

func movefile(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return err
	}
	var mfileinfo fileinfo
	mfileinfo.fullpath = path
	mfileinfo.name = info.Name()
	mfileinfo.isdir = info.IsDir()
	listfileinfo = append(listfileinfo, mfileinfo)
	return nil
}

var h = flag.Bool("h", false, "this help")
var path = flag.String("p", "./", "change path")
var old = flag.String("o", "AAA", "oldname")
var new = flag.String("n", "BBB", "newname")

func main() {
	flag.Parse()
	if *h {
		flag.Usage()
	}

	f, _ := os.OpenFile("replacekey.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	mylog := log.New(f, "", log.Ldate|log.Ltime)
	filepath.Walk(*path, movefile)
	for _, value := range listfileinfo {
		if strings.Contains(value.name, *old) {
			tmpfullpath := value.fullpath
			// 因为前面路径已经被修改改直接修改全部即可
			tmpfullpath = strings.ReplaceAll(tmpfullpath, *old, *new)

			ncount := strings.Count(value.fullpath, *old)
			value.fullpath = strings.Replace(value.fullpath, *old, *new, ncount-1)

			err := os.Rename(value.fullpath, tmpfullpath)
			fmt.Println(err, value.fullpath, "===>", tmpfullpath)
			mylog.Println(err, value.fullpath, "===>", tmpfullpath)
		}
	}
	fmt.Println("repalce " + *path + " " + *old + "===>" + *new + " " + " [OK]!")
	mylog.Println("repalce " + *path + " " + *old + "===>" + *new + " " + " [OK]!\n")
}
