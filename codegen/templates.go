package codegen

import (
	"strings"
	"text/template"
)

// 自定义模板函数
var templateFuncs = template.FuncMap{
	"ToCamelCase": toCamelCase,
}

// 将字符串转换为驼峰命名法（首字母大写）
func toCamelCase(s string) string {
	if s == "" {
		return ""
	}

	// 分割字符串
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || r == '/'
	})

	// 首字母大写
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}

	return strings.Join(parts, "")
}

var gameTemplate = template.Must(template.New("game").Funcs(templateFuncs).Parse(`
package internal

import (
	"context"
	"slp/rpc/server/internal/room_game/state/internal/{{.HandlerPackage}}"
)

// {{.GameStructName}} {{.GameName}}游戏结构体
type {{.GameStructName}} struct {
	BaseGameStateMachine
}

// 业务参数
type {{.GameStructName}}Param struct {
}

// GetGameKey 实现StateMachine接口
func (g *{{.GameStructName}}) GetGameKey() string {
	return "{{.GameKey}}"
}

{{- if .HasBefore}}
// BeforeTransition 全局前置处理函数
func (g *{{.GameStructName}}) Before(ctx context.Context, from, event string, gameId int64, val ...interface{}) error {
	return {{$.HandlerPackage}}.Before(ctx, from, event, gameId, val...)
}
{{end}}

{{- if .HasAfter}}
// AfterTransition 全局后置处理函数
func (g *{{.GameStructName}}) After(ctx context.Context, from, to, event string, gameId int64, err error, val ...interface{}) error {
	return {{$.HandlerPackage}}.After(ctx, from, to, event, gameId, val...)
}
{{end}}

func (g *{{.GameStructName}}) Transitions() map[string][]Transition {
	data := map[string][]Transition{
		{{range $state, $transitions := .State}}
		"{{$state}}": {
			{{range $i, $transition := $transitions}}
			{
				Event:   "{{$transition.Event}}",
				To:      "{{$transition.To}}",
				Handler: {{$.HandlerPackage}}.{{$.GameStructName}}{{$state | ToCamelCase}}{{$transition.Event | ToCamelCase}}Handler,
			},
			{{end}}
		},
		{{end}}
	}
	return data
}
`))

var handlerTemplate = template.Must(template.New("handler").Funcs(templateFuncs).Parse(`
package {{.HandlerPackage}}

import (
	"context"
)

func {{.HandlerName}}(ctx context.Context, gameKey string, gameId int64, val ...interface{}) error {
	// TODO: 实现{{.State}}状态下的{{.Event}}事件处理逻辑
	// 可以通过val获取事件相关参数
	return nil
}
`))

var handlerBeforeTemplate = template.Must(template.New("handler").Funcs(templateFuncs).Parse(`
package {{.HandlerPackage}}

import (
	"context"
)

func Before(ctx context.Context, gameKey string, gameId int64, val ...interface{}) error {
	// TODO: 实现{{.State}}状态下的{{.Event}}事件处理逻辑
	// 可以通过val获取事件相关参数
	return nil
}
`))

var handlerAfterTemplate = template.Must(template.New("handler").Funcs(templateFuncs).Parse(`
package {{.HandlerPackage}}

import (
	"context"
)

func After(ctx context.Context, gameKey string, gameId int64, val ...interface{}) error {
	// TODO: 实现{{.State}}状态下的{{.Event}}事件处理逻辑
	// 可以通过val获取事件相关参数
	return nil
}
`))
