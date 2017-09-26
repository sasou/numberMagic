// numbeMagic project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const phoneExp = `(010[2-9]\d{6,7})|(02[0-9][2-9]\d{6,7})|(0[3-9]\d{9,10})`
const mobileExp = `(1[3-8]\d{9})`
const phoneAndMobileExp = mobileExp + `|` + phoneExp

func ReadAllStr(fnames []string, dataStr *string) {
	for _, f := range fnames {
		f, e := os.Open(f)
		if e != nil {
			fmt.Println(e, f)
		}
		if fdata, err := ioutil.ReadAll(f); err != nil {
			fmt.Println(e, f.Name())
		} else {
			*dataStr = *dataStr + string(fdata)
		}
	}
	return
}

func FindMobiles(source *string) []string {
	r, e := regexp.Compile(mobileExp)
	if e != nil {
		fmt.Println(e, r)
		return nil
	}
	return r.FindAllString(*source, -1)
}

func checkCreateFile(fileName string) {
	_, e := os.Stat(fileName)
	if e != nil {
		ioutil.WriteFile(fileName, []byte{}, 0644)
	}
}

func SaveSliceToFile(Slice *[]string, fileName string) {
	if len(*Slice) == 0 {
		return
	}
	checkCreateFile(fileName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	f.WriteString(strings.Join(*Slice, "\n"))
	defer f.Close()
}

func SaveStrToFile(strData *string, fileName string) {
	if len(*strData) == 0 {
		return
	}
	checkCreateFile(fileName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	f.WriteString(*strData)
	defer f.Close()
}

func SaveDistinctStrToFile(strData *string, fileName string) {
	if len(*strData) == 0 {
		return
	}
	checkCreateFile(fileName)
	strMap := make(map[string]bool)
	for _, str := range strings.Split(*strData, "\n") {
		strMap[str] = true
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for str, _ := range strMap {
		f.WriteString(str + "\n")
	}

}

func splitToFileByMaxLines(strData *string, MaxLines int, name string) {
	lines := strings.Split(*strData, "\n")
	for i := 0; i < len(lines)/MaxLines+1; i++ {
		fileName := name + fmt.Sprintf("%06d", i+1) + ".txt"
		//os.Remove(fileName)
		checkCreateFile(fileName)
		upper := i*MaxLines + MaxLines
		if upper > len(lines) {
			upper = len(lines)
		}
		slice := lines[i*MaxLines : upper]
		SaveSliceToFile(&slice, fileName)
	}
}

func splitToFileByMobileLocation(strData *string, name string) {
	mobiles := FindMobiles(strData)

	if len(mobiles) == 0 {
		return
	}
	segs := Segs{}
	segs.Init()
	segs.ImportSegs()
	locNums := make(map[string][]string)

	for _, mobile := range mobiles {
		location := segs.segsmap[mobile[0:7]].province + segs.segsmap[mobile[0:7]].city
		if location == "" {
			location = "未知"
		}
		locNums[location] = append(locNums[location], mobile)
	}

	for l, m := range locNums {
		SaveSliceToFile(&m, l+fmt.Sprintf("(%d个)", len(m))+".txt")
	}
}

func cardToISP(card string) (ISP string) {
	if strings.Contains(card, "移动") {
		return "移动"
	}
	if strings.Contains(card, "联通") {
		return "联通"
	}
	if strings.Contains(card, "电信") {
		return "电信"
	}

	return "其它"

}

func splitToFileByISP(strData *string, name string) {
	mobiles := FindMobiles(strData)

	if len(mobiles) == 0 {
		return
	}
	segs := Segs{}
	segs.Init()
	segs.ImportSegs()
	ispNums := make(map[string][]string)

	for _, mobile := range mobiles {
		isp := cardToISP(segs.segsmap[mobile[0:7]].card)
		ispNums[isp] = append(ispNums[isp], mobile)
	}

	for p, m := range ispNums {
		SaveSliceToFile(&m, p+fmt.Sprintf("(%d个)", len(m))+".txt")
	}
}

func ShowChoiceMessage() {
	fmt.Println("这是您选择的文件:\n")
	for i, f := range os.Args[1:] {
		fmt.Println(i+1, ":", f)
	}
	fmt.Println("\n请选择您要做的操作:\n")
	fmt.Println("1: 提取手机号")
	fmt.Println("2: 合并文件")
	fmt.Println("3: 分割文件")
	fmt.Println("4: 提取手机号并区分归属地")
	fmt.Println("5: 提取手机号并区分运营商")
	fmt.Println("6: 号码去重")
	fmt.Printf("\n我需要: ")
}

func GetFileName() string {
	fileName := ""
	for i := 1; i < len(os.Args); i++ {
		if fileName == "" {
			fileName = strings.TrimSuffix(filepath.Base(os.Args[i]), ".txt")
		} else {
			fileName = fileName + "+" + strings.TrimSuffix(filepath.Base(os.Args[i]), ".txt")
		}

	}
	return fileName
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(dir)
	os.Chdir(dir)
	fileName := GetFileName()
	dataStr := ""
	if len(os.Args) == 1 {
		fmt.Println("请将要处理的文件选中,拖到本程序的图标上.")
		fmt.Scanln()
		os.Exit(1)
	}
	ShowChoiceMessage()
	c := ""
	fmt.Scanln(&c)

	switch c {
	case "1":
		ReadAllStr(os.Args[1:], &dataStr)
		mobiles := FindMobiles(&dataStr)
		cnt := len(mobiles)
		SaveSliceToFile(&mobiles, fileName+"(手机号"+strconv.Itoa(cnt)+"个).txt")
	case "2":
		ReadAllStr(os.Args[1:], &dataStr)
		SaveStrToFile(&dataStr, fileName+"(合并).txt")
	case "3":
		ReadAllStr(os.Args[1:], &dataStr)
		maxLines := 0
		fmt.Printf("\n请输入单文件最大行数: ")
		_, e := fmt.Scanf("%d", &maxLines)
		if e != nil || maxLines < 1 {
			fmt.Println(e)
			os.Exit(2)
		}
		splitToFileByMaxLines(&dataStr, maxLines, fileName+"_子文件")
	case "4":
		ReadAllStr(os.Args[1:], &dataStr)
		splitToFileByMobileLocation(&dataStr, "")
	case "5":
		ReadAllStr(os.Args[1:], &dataStr)
		splitToFileByISP(&dataStr, fileName+"_")
	case "6":
		ReadAllStr(os.Args[1:], &dataStr)
		SaveDistinctStrToFile(&dataStr, fileName+"(去重).txt")
	}

	fmt.Println("处理完成! 请按任意键退出.")
	fmt.Scanln()
	fmt.Scanln()
}
