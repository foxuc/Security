package main
import(
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	inix "github.com/go-ini/ini"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"path"
	"log"
)
var p = fmt.Println

/*请按任意键继续*/
func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}

func main() {
	PthSep := string(os.PathSeparator)
	NowTime :=time.Now().Format("2006_01_02_15_04_05")
	p(NowTime)
	p("Ver: 2017.11.3.01")
	p("NowPath:",getCurrentDirectory())
	iniPath := filepath.Join( getCurrentDirectory(),"Bat_EncryptPNG.ini")
	OUT_ciphercode_dir :=filepath.Join( getCurrentDirectory(),"OUT_ciphercode")
	if IsDirFileExist(OUT_ciphercode_dir) {
		os.RemoveAll(OUT_ciphercode_dir)
	}
	os.MkdirAll(OUT_ciphercode_dir,0666)

	if IsDirFileExist(iniPath) !=true {
		p(iniPath,"找不到配置文件，请检查!")
		pause()
		return
	}
	cfg,_ := inix.LoadSources(inix.LoadOptions{IgnoreInlineComment: true}, iniPath)
	timeLayout := "2006-01-02 15:04:05"
	isPAUSE := cfg.Section("BAT").Key("isPAUSE").String()
	var boolisPAUSE bool
	if isPAUSE == "true" {
		boolisPAUSE = true
	}

	startTime :=cfg.Section("BAT").Key("startTime").String()
	endTime :=cfg.Section("BAT").Key("endTime").String()
	inPath :=cfg.Section("BAT").Key("inPath").String()
	outPath :=cfg.Section("BAT").Key("outPath").String()
	EncryptPNGextName :=cfg.Section("BAT").Key("EncryptPNGextName").String()
	loc1, _ := time.LoadLocation("Local")
	startTime1, _ := time.ParseInLocation(timeLayout, startTime,loc1)
	n1startTime :=startTime1.Unix()

	loc2, _ := time.LoadLocation("Local")
	endTime1, _ := time.ParseInLocation(timeLayout, endTime,loc2)
	n1endTime :=endTime1.Unix()

	p("#-----------print ",iniPath,"---------s--#")
	p("isPAUSE --> ",isPAUSE)
	p("EncryptPNGextName --> ",EncryptPNGextName)
	p("startTime --> ",startTime,"--> ",n1startTime)
	p("endTime   --> ",endTime,"--> ",n1endTime)
	p("inPath --> ",inPath)
	p("outPath --> ",outPath)
	p("#-----------print ",iniPath,"---------e--#")

	if boolisPAUSE {
		pause()
	}

	saveBatCmdFile := filepath.Join( getCurrentDirectory(),EncryptPNGextName+".log.bat")
	if IsDirFileExist(saveBatCmdFile) {
		p(saveBatCmdFile,"配置文件存在，自动删除!")
		os.Remove(saveBatCmdFile)
	}

	k := time.Now()
	n1 :=k.Unix()
	n1Format :=k.Format("2006-01-02 15:04:05")
	p(n1,"当前时间-->",n1Format)

	//一天之前
	d, _ := time.ParseDuration("-24h")
	dx1 :=k.Add(d)
	d1 :=dx1.Unix()
	d1Format :=dx1.Format("2006-01-02 15:04:05")
	p(d1,"一天之前-->",d1Format)

	p("#--------inPath---------------------#")

	outfiles, _ := WalkDir(inPath, ".png")
	 for _, f1 := range outfiles {

	 	 f1timeUnix :=GetFileModTime(f1)
		 f1timeUnixString:=GetFileModTimestring(f1)
		 if int64(f1timeUnix) >= int64(n1startTime) && int64(f1timeUnix) <= int64(n1endTime) {
			 p(f1timeUnixString,"-->",f1timeUnix)
			 cmd1 := EncryptPNGextName+" "+ f1
			 len1 := len(f1)
			 LENinPath :=len(inPath)
			 strtemp1 := f1[LENinPath+1:len1]
			 cmd :=cmd1+" "+outPath+PthSep+strtemp1+"\r\n"
			 p(cmd)
			 appendToFile(saveBatCmdFile,cmd)
			 }


	}

	mvCMD1 := "rem ::xcopy /y /e /q "+ OUT_ciphercode_dir +" "+ inPath+"\r\n"
	appendToFile(saveBatCmdFile,mvCMD1)

	mvCMD := "xcopy /y /e /q "+ OUT_ciphercode_dir +" "+ inPath+"\r\n"
	appendToFile(saveBatCmdFile,mvCMD)
   p("#--------inPath2---------------------#")
	if boolisPAUSE {
		pause()
	}
}//--main End----

//获取文件修改时间 返回unix时间戳
func GetFileCreatTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}


//获取文件修改时间 返回unix时间戳
func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	//return fi.ModTime().Unix()
	return fi.ModTime().Unix()
}

//获取文件修改时间 返回string时间戳
func GetFileModTimestring (path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().String()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().String()
	}

	return fi.ModTime().String()
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

func appendToFile(fileName string, content string) error {

	f, err := os.OpenFile(fileName,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
	}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}


//获取指定目录下".png"的,所有文件和目录
func ListDir(dirPth string ) (files []string,files1 []string, err error) {
	//fmt.Println(dirPth)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil,nil, err
	}
	PthSep := string(os.PathSeparator)

	for _, fi := range dir {

		if fi.IsDir() { // 忽略目录
			files1 = append(files1, dirPth+PthSep+fi.Name())
			ListDir(dirPth + PthSep + fi.Name())
		}else {

			ext1 := strings.ToLower(dirPth + PthSep + fi.Name())
			fisize := fi.Size()

			if fisize > 0 && path.Ext(ext1) == ".png"  {
			files = append(files, dirPth+PthSep+fi.Name())
		}
		}
	}
	return files,files1, nil
}


func pause() {
	fmt.Println("请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:

			break Loop
		}
	}
}

//判断文件或文件夹是否存在
func IsDirFileExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

/*获取程序运行路径*/
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		p(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
