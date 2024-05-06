package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

var (
	tempFolder     = filepath.Join(os.TempDir(), "tyson_tap_testdata")
	testFolder     = "../testdata"
	dummyTestFiles []string
)

func init() {
	dummyTestFiles = createDummyFiles()
}

func TestMakeJoinerIndex(t *testing.T) {
	validateMakeJoinerIndex(tempFolder, "created", true, dummyTestFiles, t)
	validateMakeJoinerIndex(tempFolder, "created", false, dummyTestFiles, t)
	sort.Strings(dummyTestFiles)
	validateMakeJoinerIndex(tempFolder, "lastmod", true, dummyTestFiles, t)
	validateMakeJoinerIndex(tempFolder, "lastmod", false, dummyTestFiles, t)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "size", true, loadSortJSONValidtor("size.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "size", false, loadSortJSONValidtor("size.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "path", true, loadSortJSONValidtor("path.json"), t,
	)
	validateMakeJoinerIndex(
		nf("dump/yaml"), "path", false, loadSortJSONValidtor("path.json"), t,
	)
}

func validateMakeJoinerIndex(folder, sortBy string, ascending bool, exp []string, t *testing.T) {
	idx := makeJoinerIndex(newTestParams(folder, sortBy, ascending))
	if !ascending {
		exp = reverseArr(exp)
	}
	if !orderOK(idx, exp) {
		order := "asc"
		if ascending == false {
			order = "desc"
		}
		t.Errorf(
			"sort failed: %s %s, exp: %v, got: %v",
			sortBy, order, fmt.Sprintf("%v", exp), shortprintJI(idx),
		)
	}
}

func orderOK(idx tJoinerIndex, exp []string) bool {
	for i, el := range idx {
		if el.Path != exp[i] {
			return false
		}
	}
	return true
}

func reverseArr(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func newTestParams(folder, sortBy string, ascending bool) (p tIDXParams) {
	p.Endpoint.Folder = folder
	p.Endpoint.ReturnValues.Content = true
	p.Endpoint.ReturnValues.Created = true
	p.Endpoint.ReturnValues.LastMod = true
	p.Endpoint.ReturnValues.Metadata = true
	p.Endpoint.ReturnValues.Size = true
	p.Threads = 8
	p.Ascending = ascending
	p.SortBy = sortBy
	return
}

func loadSortJSONValidtor(s string) (arr []string) {
	file := filepath.Join(testFolder, "validate_sort", s)
	by, _, err := readFile(file)
	if err == nil {
		err := json.Unmarshal(by, &arr)
		if err == nil {
			return
		}
	}
	return
}

func nf(s string) string {
	return filepath.Join(testFolder, s)
}

func shortprintJI(ji tJoinerIndex) (s string) {
	for _, el := range ji {
		s += fmt.Sprintf("%s ", el.Path)
	}
	return
}

func createDummyFiles() (arr []string) {
	os.MkdirAll(tempFolder, os.ModePerm)
	for i := 1; i <= 3; i++ {
		name := filepath.Join(tempFolder, fmt.Sprintf("%03d", i)+".tmp")
		f, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		arr = append(arr, filepath.Base(name))
		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
	return
}
