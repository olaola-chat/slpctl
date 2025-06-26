package codegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 游戏配置结构体
type GameConfig struct {
	State    map[string][]StateTransition `json:"state"`
	GameKey  string                       `json:"game_key"`
	GameName string                       `json:"game_name"`
}

// 状态转换结构
type StateTransition struct {
	Event string `json:"Event"`
	To    string `json:"To"`
}

type GameGenerator struct {
	config       GameConfig
	outputDir    string
	handlerDir   string
	gameFilePath string
}

func NewGameGenerator(configPath, outputDir string) (*GameGenerator, error) {
	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置
	var config GameConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 准备目录和文件路径
	handlerDir := filepath.Join(outputDir, "state", "internal", strings.ToLower(config.GameKey)+"_handler")
	gameFilePath := filepath.Join(outputDir, "state", "internal", strings.ToLower(config.GameKey)+"_game.go")

	return &GameGenerator{
		config:       config,
		outputDir:    outputDir,
		handlerDir:   handlerDir,
		gameFilePath: gameFilePath,
	}, nil
}

func (g *GameGenerator) Generate() error {
	if err := g.createDirectories(); err != nil {
		return err
	}

	// 生成游戏文件
	if err := g.generateGameFile(); err != nil {
		return err
	}

	// 生成处理函数文件
	if err := g.generateHandlerFiles(); err != nil {
		return err
	}

	return nil
}

func (g *GameGenerator) createDirectories() error {
	if err := os.MkdirAll(g.handlerDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	return nil
}

func (g *GameGenerator) generateGameFile() error {
	// 检查文件是否存在
	if fileExists(g.gameFilePath) {
		fmt.Printf("文件已存在，跳过生成: %s\n", g.gameFilePath)
		return nil
	}

	// 准备模板数据
	gameStructName := toCamelCase(g.config.GameKey) + "Game"
	handlerPackage := strings.ToLower(g.config.GameKey) + "_handler"

	// 查找初始状态（默认使用配置中的第一个状态）
	var initialState string
	for state := range g.config.State {
		initialState = state
		break
	}

	tmplData := struct {
		GameKey        string
		GameName       string
		GameStructName string
		InitialState   string
		State          map[string][]StateTransition
		HandlerPackage string
	}{
		GameKey:        g.config.GameKey,
		GameName:       g.config.GameName,
		GameStructName: gameStructName,
		InitialState:   initialState,
		State:          g.config.State,
		HandlerPackage: handlerPackage,
	}

	var buf bytes.Buffer
	if err := gameTemplate.Execute(&buf, tmplData); err != nil {
		return fmt.Errorf("执行游戏文件模板失败: %v", err)
	}

	if err := ioutil.WriteFile(g.gameFilePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入游戏文件失败: %v", err)
	}

	fmt.Printf("生成游戏文件: %s\n", g.gameFilePath)
	return nil
}

func (g *GameGenerator) generateHandlerFiles() error {
	gameStructName := toCamelCase(g.config.GameKey) + "Game"
	handlerPackage := strings.ToLower(g.config.GameKey) + "_handler"

	for state, transitions := range g.config.State {
		for _, transition := range transitions {
			handlerName := generateHandlerName(gameStructName, state, transition.Event)
			handlerPath := filepath.Join(g.handlerDir, fmt.Sprintf("%s.go", strings.ToLower(strings.ReplaceAll(handlerName, "Handler", ""))))

			if err := g.generateHandlerFile(handlerPath, handlerName, state, transition.Event, handlerPackage); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *GameGenerator) generateHandlerFile(filePath, handlerName, state, event, handlerPackage string) error {
	// 检查文件是否存在
	if fileExists(filePath) {
		fmt.Printf("文件已存在，跳过生成: %s\n", filePath)
		return nil
	}

	tmplData := struct {
		HandlerPackage string
		HandlerName    string
		State          string
		Event          string
	}{
		HandlerPackage: handlerPackage,
		HandlerName:    handlerName,
		State:          toCamelCase(state),
		Event:          toCamelCase(event),
	}

	var buf bytes.Buffer
	if err := handlerTemplate.Execute(&buf, tmplData); err != nil {
		return fmt.Errorf("执行处理函数模板失败: %v", err)
	}

	if err := ioutil.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入处理函数文件失败: %v", err)
	}

	fmt.Printf("生成处理函数文件: %s\n", filePath)
	return nil
}

// 生成处理函数名称
func generateHandlerName(gameStructName, state, event string) string {
	stateTitle := toCamelCase(state)
	eventTitle := toCamelCase(event)

	return fmt.Sprintf("%s%s%sHandler", gameStructName, stateTitle, eventTitle)
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
