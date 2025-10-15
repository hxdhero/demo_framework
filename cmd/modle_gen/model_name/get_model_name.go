package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ModuleInfo 存储模块信息
type ModuleInfo struct {
	ModulePath   string
	DirPath      string
	Replacements map[string]string
}

// ModulePathFinder 模块路径查找器
type ModulePathFinder struct {
	moduleRoot string
	moduleInfo *ModuleInfo
}

// FindGoModPath 从当前目录向上查找 go.mod 文件
func FindGoModPath(startDir string) (string, error) {
	currentDir := startDir
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir { // 已经到达根目录
			return "", fmt.Errorf("go.mod not found")
		}
		currentDir = parentDir
	}
}

// ParseGoMod 解析 go.mod 文件
func ParseGoMod(goModPath string) (*ModuleInfo, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	moduleInfo := &ModuleInfo{
		Replacements: make(map[string]string),
	}

	var inReplaceSection bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "module ") {
			moduleInfo.ModulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		} else if strings.HasPrefix(line, "replace ") {
			inReplaceSection = true
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				old := strings.TrimSpace(parts[1])
				newVal := strings.TrimSpace(parts[3])
				moduleInfo.Replacements[old] = newVal
			}
		} else if inReplaceSection && strings.HasPrefix(line, ")") {
			inReplaceSection = false
		} else if inReplaceSection && line != "" && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				old := strings.TrimSpace(parts[0])
				newVal := strings.TrimSpace(parts[2])
				moduleInfo.Replacements[old] = newVal
			}
		}
	}

	return moduleInfo, scanner.Err()
}

// NewModulePathFinder 创建新的模块路径查找器
func NewModulePathFinder() (*ModulePathFinder, error) {
	goModPath, err := FindGoModPath(".")
	if err != nil {
		return nil, err
	}

	moduleInfo, err := ParseGoMod(goModPath)
	if err != nil {
		return nil, err
	}

	return &ModulePathFinder{
		moduleRoot: filepath.Dir(goModPath),
		moduleInfo: moduleInfo,
	}, nil
}

// FindStructsByPath 根据模块路径查找结构体
func (m *ModulePathFinder) FindStructsByPath(relativePath string) ([]string, error) {
	fullPath := filepath.Join(m.moduleRoot, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", fullPath)
	}
	return GetStructNamesFromDir(fullPath)
}

// GetStructNamesFromDir 从指定目录读取所有.go文件中的结构体名称
func GetStructNamesFromDir(dirPath string) ([]string, error) {
	var structNames []string

	// 遍历目录下的所有.go文件
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理.go文件，排除_test.go文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			names, err := getStructNamesFromFile(path)
			if err != nil {
				return err
			}
			structNames = append(structNames, names...)
		}

		return nil
	})

	return structNames, err
}

// getStructNamesFromFile 从单个.go文件中提取结构体名称
func getStructNamesFromFile(filePath string) ([]string, error) {
	fset := token.NewFileSet()

	// 解析.go文件
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var structNames []string

	// 遍历AST，查找结构体声明
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			// 检查是否为结构体类型
			if _, ok := x.Type.(*ast.StructType); ok {
				if x.Name != nil {
					structNames = append(structNames, x.Name.Name)
				}
			}
		}
		return true
	})

	return structNames, nil
}

// GetModulePath 获取模块完整路径
func (m *ModulePathFinder) GetModulePath() string {
	return m.moduleInfo.ModulePath
}

// GetModuleRoot 获取模块根目录
func (m *ModulePathFinder) GetModuleRoot() string {
	return m.moduleRoot
}

// 打印出用于gorm自动迁移的模型名称
func main() {
	finder, err := NewModulePathFinder()
	if err != nil {
		fmt.Printf("Error creating module finder: %v\n", err)
		return
	}

	fmt.Printf("Module: %s\n", finder.GetModulePath())
	fmt.Printf("Module root: %s\n", finder.GetModuleRoot())

	// 查找不同路径下的结构体
	paths := []string{"internal/model/gen"}

	for _, path := range paths {
		structNames, err := finder.FindStructsByPath(path)
		if err != nil {
			fmt.Printf("Error finding structs in %s: %v\n", path, err)
			continue
		}

		if len(structNames) > 0 {
			fmt.Printf("\nStructs in %s:\n", path)
			for _, name := range structNames {
				fmt.Printf("gen.%s{},\n", name)
			}
		}
	}
}
