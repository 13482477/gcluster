package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"reflect"
	"time"
)

type (
	BaseExcel struct {
		File  *xlsx.File
		Sheet *xlsx.Sheet
	}
	ExportExcel struct {
		MyExcel  *BaseExcel
		downPath string
		savePath string
		fileName string
		Header   []string
	}
)

func NewExportExcel() *ExportExcel {
	baseExcel := new(BaseExcel)
	baseExcel.File = xlsx.NewFile()
	baseExcel.Sheet, _ = baseExcel.File.AddSheet("Sheet1")
	exportExcel := ExportExcel{
		MyExcel: baseExcel,
	}
	return &exportExcel
}
func (ex *ExportExcel) SetOptions(options map[string]string) {
	if savePath, ok := options["savePath"]; ok {
		ex.savePath = savePath
	}
	if downPath, ok := options["downPath"]; ok {
		ex.downPath = downPath
	}
}

func (ex *ExportExcel) SetFileName(name string, addTimeSuffix bool) {
	ex.fileName = name
	if addTimeSuffix {
		date := time.Now().Format("2006-01-02")
		ex.fileName += "_" + date + ".xlsx"
	} else {
		ex.fileName += ".xlsx"
	}
}

func (ex *ExportExcel) SetHeader(header []string) {
	ex.Header = header
}

func (ex *ExportExcel) Export(data interface{}) (string, error) {
	sheet := ex.MyExcel.Sheet
	logrus.Debug(data)
	f := reflect.ValueOf(data)
	if ex.Header == nil {
		return "", fmt.Errorf("导出EXCEL时Herder信息错误")
	}
	y := reflect.TypeOf(data)
	z := reflect.TypeOf(data).Elem()
	if y.Kind() != reflect.Slice || z.Kind() != reflect.Ptr {
		return "", fmt.Errorf("导出EXCEL时Data格式错误")
	}
	var indexs []string
	for i := 0; i <= f.Len(); i++ {
		row := sheet.AddRow()
		if i == 0 {
			for _, title := range ex.Header {
				cell := row.AddCell()
				cell.SetValue(title)
				indexs = append(indexs, title)
			}
			logrus.Debug(indexs)
		} else {
			object := f.Index(i - 1).Interface()
			t := reflect.ValueOf(object).Elem()
			x := reflect.TypeOf(object).Elem()
			for j := 0; j < t.NumField(); j++ {
				for _, index := range indexs {
					cname := x.Field(j).Tag.Get("cname")
					if index == cname {
						cell := row.AddCell()
						f := t.Field(j)
						switch f.Kind() {
						case reflect.String:
							cell.Value = f.String()
						case reflect.Int32:
							cell.SetInt64(int64(f.Int()))
						case reflect.Float32:
							cell.SetFloat(float64(f.Float()))
						default:
							cell.SetValue(f.InterfaceData())
						}
					}
				}
			}
		}
	}
	err := ex.MyExcel.File.Save(ex.savePath + ex.fileName)
	if err != nil {
		return "", err
	} else {
		return ex.downPath + ex.fileName, nil
	}
}
