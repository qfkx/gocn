package util

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/valyala/fasthttp"
)

// DrawExcel 绘制模板
// 下载模板的时候，增加了随机数，确保拿到最新的模板
func DrawExcel(jsonUrl, xlsxUrl string, dataInfo map[string]gjson.Result, dataDetail []gjson.Result) (string, error) {
	// 下载模板
	var filename = JoinX("", GetGUID().Hex()+".xlsx")
	Download(xlsxUrl+"?"+GetGUID().Hex(), filename, nil)
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return "", err
	}
	status, tplJson, err := fasthttp.Get(nil, jsonUrl+"?"+GetGUID().Hex())
	if err != nil {
		return "", err
	}
	if status != fasthttp.StatusOK {
		return "", err
	}

	// 模板
	tpl := gjson.Get(string(tplJson), "@this")
	tplStart, _ := strconv.Atoi(tpl.Get("start").String())
	tplCount, _ := strconv.Atoi(tpl.Get("count").String())
	tplInfo := tpl.Get("info").Map()
	tplDetail := tpl.Get("detail").Map()

	// 填写表头
	newValue := ""
	for k := range dataInfo {
		if _, ok := tplInfo[k]; ok == false { //模板不需要
			continue
		}
		v := dataInfo[k]
		pos := tplInfo[k].String()
		cell := f.GetCellValue("Sheet1", pos)
		t1 := v.Type.String()
		if t1 == "Number" {
			newValue = v.String()
		}
		if t1 == "String" {
			newValue = v.String()
		}
		newValue = strings.Replace(cell, "{"+k+"}", newValue, 1)
		f.SetCellValue("Sheet1", pos, newValue)
	}
	totals := len(dataDetail)
	// 删除多余的行
	if tplCount > totals {
		for i := totals + 1; i <= tplCount; i++ {
			f.RemoveRow("Sheet1", tplStart+i)
		}
	}
	// 补充行
	if tplCount < totals {
		for i := 0; i < totals-tplCount; i++ {
			f.DuplicateRow("Sheet1", tplStart+tplCount-1+i) //从最后一行复制行，乃不传之秘
		}
	}

	//填制表格
	id := 1
	for _, row := range dataDetail {
		r := row.Map()
		for k := range r {
			if _, ok := tplDetail[k]; ok == false { // 模板不需要这个数据
				continue
			}
			pos := tplDetail[k].String() + strconv.Itoa(tplStart+id-1)
			t1 := r[k].Type.String()
			if t1 == "Number" {
				v1 := r[k].Float()
				f.SetCellValue("Sheet1", pos, v1)
			}
			if t1 == "String" {
				f.SetCellValue("Sheet1", pos, r[k])
			}
		}
		id = id + 1
	}
	if err := f.Save(); err != nil {
		return "", err
	}
	key := "print/" + filepath.Base(filename)
	if res, err := Upload(filename, key, 7200); err != nil {
		return res, err
	}
	del := os.Remove(filename)
	if del != nil {
		return "", del
	}
	return key, nil
}

// AmountConvert 金额转为大写
func AmountConvert(pMoney float64, pRound bool) string {
	var NumberUpper = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖", "零"}
	var Unit = []string{"分", "角", "圆", "拾", "佰", "仟", "万", "拾", "佰", "仟", "亿", "拾", "佰", "仟"}
	var regex = [][]string{
		{"零拾", "零"}, {"零佰", "零"}, {"零仟", "零"}, {"零零零", "零"}, {"零零", "零"},
		{"零角零分", "整"}, {"零分", "整"}, {"零角", "零"}, {"零亿零万零元", "亿元"},
		{"亿零万零元", "亿元"}, {"零亿零万", "亿"}, {"零万零元", "万元"}, {"万零元", "万元"},
		{"零亿", "亿"}, {"零万", "万"}, {"拾零圆", "拾元"}, {"零圆", "元"}, {"零零", "零"}}
	str, DigitUpper, UnitLen, round := "", "", 0, 0
	if pMoney == 0 {
		return "零"
	}
	if pMoney < 0 {
		str = "负"
		pMoney = math.Abs(pMoney)
	}
	if pRound {
		round = 2
	} else {
		round = 1
	}
	DigitByte := []byte(strconv.FormatFloat(pMoney, 'f', round+1, 64)) //注意币种四舍五入
	UnitLen = len(DigitByte) - round

	for _, v := range DigitByte {
		if UnitLen >= 1 && v != 46 {
			s, _ := strconv.ParseInt(string(v), 10, 0)
			if s != 0 {
				DigitUpper = NumberUpper[s-1]

			} else {
				DigitUpper = "零"
			}
			str = str + DigitUpper + Unit[UnitLen-1]
			UnitLen = UnitLen - 1
		}
	}
	for i, _ := range regex {
		reg := regexp.MustCompile(regex[i][0])
		str = reg.ReplaceAllString(str, regex[i][1])
	}
	if string(str[0:3]) == "元" {
		str = string(str[3:len(str)])
	}
	if string(str[0:3]) == "零" {
		str = string(str[3:len(str)])
	}
	return str
}
