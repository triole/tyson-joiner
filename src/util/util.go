package util

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/triole/logseal"
	yaml "gopkg.in/yaml.v2"
)

func (util Util) AbsPath(str string) (p string, err error) {
	p, err = filepath.Abs(str)
	util.Lg.IfErrFatal("invalid file path", logseal.F{"path": str, "error": err})
	return p, err
}

func (util Util) GetPathDepth(pth string) int {
	return len(strings.Split(pth, string(filepath.Separator))) - 1
}

func (util Util) Find(basedir string, rxFilter string) []string {
	inf, err := os.Stat(basedir)
	if err != nil {
		util.Lg.IfErrFatal(
			"unable to access md folder", logseal.F{
				"path": basedir, "error": err,
			},
		)
	}
	if !inf.IsDir() {
		util.Lg.Fatal(
			"not a folder, please provide a directory to look for md files.",
			logseal.F{"path": basedir},
		)
	}

	filelist := []string{}
	rxf, _ := regexp.Compile(rxFilter)

	err = filepath.Walk(basedir, func(path string, f os.FileInfo, err error) error {
		if rxf.MatchString(path) {
			inf, err := os.Stat(path)
			if err == nil {
				if !inf.IsDir() {
					filelist = append(filelist, path)
				}
			} else {
				util.Lg.IfErrInfo("stat file failed", logseal.F{"path": path})
			}
		}
		return nil
	})
	util.Lg.IfErrFatal("find files failed", logseal.F{"path": basedir, "error": err})
	return filelist
}

func (util Util) GetFileSize(filename string) (siz uint64) {
	file, err := os.Open(filename)
	util.Lg.IfErrError(
		"can not open file to get file size",
		logseal.F{"path": filename, "error": err},
	)
	if err == nil {
		defer file.Close()
		stat, err := file.Stat()
		util.Lg.IfErrError(
			"can not stat file to get file size",
			logseal.F{"path": filename, "error": err},
		)
		if err == nil {
			siz = uint64(stat.Size())
		}
	}
	return
}

func (util Util) GetFileLastMod(filename string) (uts int64) {
	fil, err := os.Stat(filename)
	if err != nil {
		util.Lg.Error("can not stat file", logseal.F{"path": filename, "error": err})
		return
	}
	uts = fil.ModTime().Unix()
	return
}

func (ut Util) ReadFile(filename string) (by []byte, isTextfile bool, err error) {
	fn, err := ut.AbsPath(filename)
	if err == nil {
		by, err = os.ReadFile(fn)
		isTextfile = utf8.ValidString(string(by))
		ut.Lg.IfErrError(
			"can not read file", logseal.F{"path": filename, "error": err},
		)
	}
	return
}

func (ut Util) ReadYAMLFile(filepath string) (r map[string]interface{}) {
	by, _, err := ut.ReadFile(filepath)
	if err != nil {
		return
	} else {
		_ = yaml.Unmarshal(by, &r)
	}
	return
}

func (util Util) RxFind(rx string, content string) string {
	temp, _ := regexp.Compile(rx)
	return temp.FindString(content)
}

func (util Util) RxMatch(rx string, content string) bool {
	temp, _ := regexp.Compile(rx)
	return temp.MatchString(content)
}

func (util Util) StringifySliceOfInterfaces(itf []interface{}) (r []string) {
	for _, el := range itf {
		r = append(r, el.(string))
	}
	return
}

func (util Util) ToFloat(inp interface{}) (fl float64) {
	switch val := inp.(type) {
	case float32:
		fl = float64(val)
	case float64:
		fl = val
	case int:
		fl = float64(val)
	case int8:
		fl = float64(val)
	case int16:
		fl = float64(val)
	case int32:
		fl = float64(val)
	case int64:
		fl = float64(val)
	case uint:
		fl = float64(val)
	case uint8:
		fl = float64(val)
	case uint16:
		fl = float64(val)
	case uint32:
		fl = float64(val)
	case uint64:
		fl = float64(val)
	}
	return
}

func (util Util) ToString(inp interface{}) string {
	return fmt.Sprintf("%s", inp)
}

func (util Util) Trace() (r string) {
	pc, fullfile, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	file := util.RxFind("src.*", fullfile)
	r = fmt.Sprintf("%s:%d %s", file, line, fn.Name())
	return
}

func (util Util) FromTestFolder(s string) (r string) {
	t, err := filepath.Abs("../../testdata")
	if err == nil {
		r = filepath.Join(t, s)
	}
	return
}
