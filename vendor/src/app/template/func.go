package template

import (
	"fmt"
	tmplhtml "html/template"
	"regexp"
	"strconv"
	"strings"
)

type ByteSize float64

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

func initFuncMap(m map[string]string) tmplhtml.FuncMap {
	return map[string]interface{}{
		"existMapKey": existMapKey,
		"toUpper":     strings.ToUpper,
		"toLower":     strings.ToLower,
		"title":       strings.Title,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"strFormatMeasureUnit": strFormatMeasureUnit,
		"strFormatDate":        strFormatDate,
		"instanceMapping":      instanceMapping(m),
		"measurePrecision":     measurePrecision,
	}
}

func (b ByteSize) String() string {
	switch {
	case b >= EB:
		return fmt.Sprintf("%.1fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.1fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.1fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.1fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.1fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.1fKB", b/KB)
	}
	return fmt.Sprintf("%fB", b)
}

func instanceMapping(m map[string]string) func(string) string {
	return (func(instance string) string {

		str := strings.Split(instance, ":")
		if len(str) != 2 {
			return fmt.Sprintf("[ %s ]", instance)
		}
		ip := str[0]
		if v, ok := m[ip]; ok {
			return fmt.Sprintf("[%s | %s ]", ip, v)
		}
		return fmt.Sprintf("[ %s ]", ip)
	})
}

func existMapKey(dict map[string]interface{}, key_search string) bool {
	if _, ok := dict[key_search]; ok {
		return true
	}
	return false
}

func strFormatMeasureUnit(MeasureUnit string, value string) string {
	return value
}

func measurePrecision(value string) string {
	successFloat64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value
	}
	return fmt.Sprintf("%.2f", successFloat64)
}

/*
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
log.Info(Round(123.555555, .5, 2))
*/

func strFormatDate(toformat string) string {

	return toformat
}
