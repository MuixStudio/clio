package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

type Column struct {
	Name    string
	GoName  string
	GoType  string
	GormTag string
}

type Table struct {
	Name    string
	GoName  string
	Columns []Column
}

var typeMapping = map[string]string{
	"bigint":    "int64",
	"int":       "int",
	"varchar":   "string",
	"text":      "string",
	"timestamp": "time.Time",
}

func GenModel(sql_file string) {
	sqlBytes, _ := os.ReadFile(sql_file)
	sql := string(sqlBytes)

	// 找表名
	tableRe := regexp.MustCompile(`(?i)CREATE TABLE (\w+)`)
	tableMatch := tableRe.FindStringSubmatch(sql)
	if len(tableMatch) < 2 {
		panic("no table found")
	}
	tableName := tableMatch[1]

	// 找字段
	colRe := regexp.MustCompile(`(?m)^\s*(\w+)\s+([a-zA-Z0-9()]+)`)
	matches := colRe.FindAllStringSubmatch(sql, -1)

	var columns []Column
	for _, m := range matches {
		colName := m[1]
		sqlType := strings.ToLower(m[2])
		goType := "string"
		for k, v := range typeMapping {
			if strings.HasPrefix(sqlType, k) {
				goType = v
				break
			}
		}
		columns = append(columns, Column{
			Name:    colName,
			GoName:  snakeToCamel(colName),
			GoType:  goType,
			GormTag: fmt.Sprintf("column:%s", colName),
		})
	}

	table := Table{
		Name:    tableName,
		GoName:  snakeToCamel(tableName),
		Columns: columns,
	}

	// 渲染模板
	_, filename, _, _ := runtime.Caller(0)
	tmplPath := filepath.Join(filepath.Dir(filename), "templates/model.tmpl")
	tmpl := template.Must(template.ParseFiles(tmplPath))
	tmpl.Execute(os.Stdout, table)
}

// 下划线转驼峰
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		parts[i] = strings.Title(p)
	}
	return strings.Join(parts, "")
}
