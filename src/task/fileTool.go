package task

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// 当前目录 ./data
var RootDir string

func SetPath(path string) {
	length := len(path)
	if length > 0 {
		// tmp
		if strings.HasSuffix(path, string(os.PathSeparator)) {
			RootDir = string([]rune(path)[0 : length-1])
		}
		fmt.Println(RootDir)
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	RootDir = dir + string(os.PathSeparator) + "data"
	fi, err := os.Stat(RootDir)
	if err != nil {
		os.Mkdir(RootDir, os.ModePerm)
	}
	fmt.Println(fi)
}

// 去掉最后的路径分隔符如：/
func GetPath(dir string) string {
	var root string
	if os.IsPathSeparator('\\') {
		i := strings.LastIndex(dir, "\\")
		root = string([]rune(dir)[0:i])
	} else {
		i := strings.LastIndex(dir, "/")
		root = string([]rune(dir)[0:i])
	}
	return root
}

func CreateFile(name string) (*os.File, error) {
	filename := RootDir + string(os.PathSeparator) + name

	if checkFileIsExist(name) {
		return os.Open(filename)
	} else {
		return os.Create(filename)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 生产文件名
func filename(id int) string {
	name := "ZX" + fmt.Sprint(id) + "X_"
	// 时间戳
	s := Timestamp(time.Now())
	//	fmt.Println(s)
	strings.TrimSpace(s)
	ss := strings.Split(s, " ")
	name += ss[0] + "-" + ss[1]

	return name
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}

		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}

		return nil
	})

	return files, err
}

var FileCount uint32 = 0

// 异步进行 将 .csb.bak 文件改为 .bak 文件
func WritedAndRename() {
	fmt.Println("rename file task open")
	atomic.SwapUint32(&FileCount, 0)

	dir, _ := ioutil.ReadDir(RootDir)
	for _, fi := range dir {
		name := RootDir + string(os.PathSeparator) + fi.Name()
		reName := strings.TrimSuffix(name, ".bak")
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(name, ".csv.bak") { //匹配文件
			os.Rename(name, reName)
		}
	}
}

func MiTask(filer *time.Ticker, datas []string, i int) {
	s := rand.New(rand.NewSource(time.Now().UnixNano())).Uint32()
	sleep := time.Duration(s >> 7) //12ms
	time.Sleep(sleep)
	var id = "ZX" + fmt.Sprint(i) + "X,"
	//	fmt.Println("开始每分钟写文件任务" + id)
	go func() {
		dater := time.NewTicker(time.Second * 1)
		task(dater, datas)
		for _ = range filer.C {
			WorkTime()
			name := filename(i)
			name = strings.Replace(name, ":", "_", -1)
			f, err := CreateFile(name + ".csv.bak")
			if err != nil {
				fmt.Println(err)
			}
			w := bufio.NewWriter(f)
			for _, s := range datas {
				if len(s) != 0 {
					w.WriteString(id + s + "\n")
				}
			}
			w.Flush()
			//			defer f.Close()
			f.Close()

			atomic.AddUint32(&FileCount, 1)
			//linux
			os.Rename(RootDir+string(os.PathSeparator)+name+".csv.bak", RootDir+string(os.PathSeparator)+name+".csv")
			// windows 可以这么干
			//			if FileCount > 5000 {
			//				go WritedAndRename()
			//			}
		}
	}()
}
