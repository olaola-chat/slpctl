package helmgen

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// Service 表示一个服务的元数据
type Service struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Name     string `json:"name"`
	Path     string `json:"path"`
}

// HTML模板内容 - 已修复所有语法问题
const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>helm重启tag服务管理仪表盘</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            margin-bottom: 30px;
        }
        .controls {
            display: flex;
            gap: 15px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        .filter-group {
            flex: 1;
            min-width: 200px;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        select, input {
            width: 100%;
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
        }
        .global-controls {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
            flex-wrap: wrap;
        }
        .global-controls button {
            padding: 8px 15px;
            background-color: #4a6fa5;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            transition: background-color 0.2s;
        }
        .global-controls button:hover {
            background-color: #3a5a8c;
        }
        .category {
            margin-bottom: 20px;
            border: 1px solid #e0e0e0;
            border-radius: 5px;
            overflow: hidden;
        }
        .category-header {
            background-color: #f8f9fa;
            padding: 12px 15px;
            font-weight: bold;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
            user-select: none;
        }
        .category-header:hover {
            background-color: #e9ecef;
        }
        .category-title {
            color: #2c3e50;
        }
        .category-count {
            background-color: #6c757d;
            color: white;
            border-radius: 12px;
            padding: 2px 10px;
            font-size: 0.85em;
        }
        .services-list {
            display: none;
        }
        .services-list.active {
            display: block;
        }
        .service-item {
            display: flex;
            align-items: center;
            padding: 12px 15px;
            border-bottom: 1px solid #eee;
            transition: background 0.2s;
        }
        .service-item:hover {
            background-color: #f8f9fa;
        }
        .service-select {
            display: block;
            position: relative;
            padding-left: 28px;
            margin-right: 15px;
            cursor: pointer;
            user-select: none;
        }
        .service-select input {
            position: absolute;
            opacity: 0;
            cursor: pointer;
            height: 0;
            width: 0;
        }
        .checkmark {
            position: absolute;
            top: 0;
            left: 0;
            height: 20px;
            width: 20px;
            background-color: #eee;
            border-radius: 4px;
            transition: all 0.2s;
        }
        .service-select:hover .checkmark {
            background-color: #ddd;
        }
        .service-select input:checked ~ .checkmark {
            background-color: #4a6fa5;
        }
        .checkmark:after {
            content: "";
            position: absolute;
            display: none;
            left: 7px;
            top: 3px;
            width: 5px;
            height: 10px;
            border: solid white;
            border-width: 0 2px 2px 0;
            transform: rotate(45deg);
        }
        .service-select input:checked ~ .checkmark:after {
            display: block;
        }
        .service-info {
            flex: 1;
            min-width: 0;
        }
        .service-name {
            font-weight: bold;
            margin-bottom: 3px;
            color: #343a40;
        }
        .service-path {
            color: #6c757d;
            font-size: 0.85em;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .service-type {
            padding: 4px 10px;
            background-color: #e9ecef;
            border-radius: 12px;
            font-size: 0.8em;
            color: #495057;
            margin-left: 15px;
        }
        .no-results {
            text-align: center;
            padding: 30px;
            color: #6c757d;
            display: none;
        }
        .command-dialog {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0,0,0,0.5);
            display: flex;
            justify-content: center;
            align-items: center;
            z-index: 1000;
        }
        .dialog-content {
            background: white;
            padding: 25px;
            border-radius: 8px;
            width: 90%;
            max-width: 700px;
            box-shadow: 0 4px 20px rgba(0,0,0,0.15);
        }
        .dialog-content h3 {
            margin-top: 0;
            color: #2c3e50;
        }
        #command-output {
            width: 100%;
            height: 250px;
            margin: 15px 0;
            font-family: 'Courier New', monospace;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            resize: none;
            font-size: 14px;
            line-height: 1.5;
            white-space: pre;
            overflow-x: auto;
        }
        .dialog-buttons {
            display: flex;
            justify-content: flex-end;
            gap: 10px;
        }
        .dialog-buttons button {
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
            transition: background-color 0.2s;
        }
        #copy-commands {
            background-color: #28a745;
            color: white;
        }
        #copy-commands:hover {
            background-color: #218838;
        }
        #close-dialog {
            background-color: #6c757d;
            color: white;
        }
        #close-dialog:hover {
            background-color: #5a6268;
        }
        .toast {
            position: fixed;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            background-color: #333;
            color: white;
            padding: 12px 24px;
            border-radius: 4px;
            opacity: 0;
            transition: opacity 0.3s;
            z-index: 1001;
        }
        .toast.show {
            opacity: 1;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>服务管理仪表盘</h1>

        <div class="controls">
            <div class="filter-group">
                <label for="category-filter">服务类型</label>
                <select id="category-filter">
                    <option value="all">所有类型</option>
                    <option value="HTTP">HTTP</option>
                    <option value="WEB">WEB</option>
                    <option value="RPC">RPC</option>
                    <option value="CMD">CMD</option>
                    <option value="TASK">TASK</option>
                </select>
            </div>

            <div class="filter-group">
                <label for="search-filter">搜索服务</label>
                <input type="text" id="search-filter" placeholder="输入服务名称或路径...">
            </div>
        </div>

        <div class="global-controls">
            <button id="select-all-btn">全选</button>
            <button id="deselect-all-btn">取消全选</button>
            <button id="show-selected-btn">显示已选 (<span id="selected-count">0</span>)</button>
            <button id="execute-selected-btn">生成执行命令</button>
        </div>

        <div id="services-container"></div>
        <div id="no-results" class="no-results">没有找到匹配的服务</div>
    </div>

    <div id="toast" class="toast"></div>

    <script>
        // 服务数据
        window.serviceMetadata = {{.ServicesJSON}};

        // 显示Toast通知
        function showToast(message, duration = 2000) {
            const toast = document.getElementById('toast');
            toast.textContent = message;
            toast.classList.add('show');

            setTimeout(() => {
                toast.classList.remove('show');
            }, duration);
        }

        // 初始化应用
        document.addEventListener('DOMContentLoaded', function() {
            const state = {
                selectedServices: new Set(),
                currentServices: [...window.serviceMetadata]
            };

            // 渲染服务列表
            function renderServices(services) {
                const container = document.getElementById('services-container');
                container.innerHTML = '';

                if (services.length === 0) {
                    document.getElementById('no-results').style.display = 'block';
                    return;
                }

                document.getElementById('no-results').style.display = 'none';

                // 按类别分组
                const categories = services.reduce(function(acc, service) {
                    if (!acc[service.category]) acc[service.category] = [];
                    acc[service.category].push(service);
                    return acc;
                }, {});

                // 渲染每个类别
                Object.entries(categories).forEach(function([category, services]) {
                    const categoryElement = document.createElement('div');
                    categoryElement.className = 'category';

                    const header = document.createElement('div');
                    header.className = 'category-header';
                    header.innerHTML = 
                        '<span class="category-title">' + category + '</span>' +
                        '<span class="category-count">' + services.length + '个服务</span>';

                    const list = document.createElement('div');
                    list.className = 'services-list active';

                    services.forEach(function(service) {
                        const serviceId = service.type + '.' + service.name;
                        const isChecked = state.selectedServices.has(serviceId);

                        const serviceItem = document.createElement('div');
                        serviceItem.className = 'service-item';
                        serviceItem.innerHTML = 
                            '<label class="service-select">' +
                                '<input type="checkbox" class="service-checkbox" data-service-id="' + serviceId + '" ' + (isChecked ? 'checked' : '') + '>' +
                                '<span class="checkmark"></span>' +
                            '</label>' +
                            '<div class="service-info">' +
                                '<div class="service-name">' + service.name + '</div>' +
                                '<div class="service-path">' + service.path + '</div>' +
                            '</div>' +
                            '<div class="service-type">' + service.type + '</div>';
                        list.appendChild(serviceItem);
                    });

                    categoryElement.appendChild(header);
                    categoryElement.appendChild(list);
                    container.appendChild(categoryElement);

                    // 点击类别标题切换展开/折叠
                    header.addEventListener('click', function() {
                        list.classList.toggle('active');
                    });
                });
            }

            // 更新选中计数
            function updateSelectedCount() {
                document.getElementById('selected-count').textContent = state.selectedServices.size;
            }

            // 设置事件监听
            function setupEventListeners() {
                // 服务选择
                document.addEventListener('change', function(e) {
                    if (e.target.classList.contains('service-checkbox')) {
                        const serviceId = e.target.dataset.serviceId;
                        if (e.target.checked) {
                            state.selectedServices.add(serviceId);
                        } else {
                            state.selectedServices.delete(serviceId);
                        }
                        updateSelectedCount();
                    }
                });

                // 全选/取消全选
                document.getElementById('select-all-btn').addEventListener('click', function() {
                    document.querySelectorAll('.service-checkbox').forEach(function(checkbox) {
                        checkbox.checked = true;
                        state.selectedServices.add(checkbox.dataset.serviceId);
                    });
                    updateSelectedCount();
                    showToast('已全选所有服务');
                });

                document.getElementById('deselect-all-btn').addEventListener('click', function() {
                    document.querySelectorAll('.service-checkbox').forEach(function(checkbox) {
                        checkbox.checked = false;
                        state.selectedServices.delete(checkbox.dataset.serviceId);
                    });
                    updateSelectedCount();
                    showToast('已取消全选');
                });

                // 显示已选服务
                document.getElementById('show-selected-btn').addEventListener('click', function() {
                    if (state.selectedServices.size === 0) {
                        showToast('当前没有选中任何服务');
                        return;
                    }

                    const dialog = document.createElement('div');
                    dialog.className = 'command-dialog';

                    let message = "<h4>已选服务:</h4><ul>";
                    state.selectedServices.forEach(function(id) {
                        const parts = id.split('.', 2);
                        const type = parts[0];
                        const name = parts[1];
                        message += "<li>" + type + "/" + name + "</li>";
                    });
                    message += "</ul><p>共 " + state.selectedServices.size + " 个服务</p>";

                    dialog.innerHTML = 
                        '<div class="dialog-content">' +
                            '<h3>已选服务列表</h3>' +
                            '<div>' + message + '</div>' +
                            '<div class="dialog-buttons">' +
                                '<button id="close-selected-dialog">关闭</button>' +
                            '</div>' +
                        '</div>';
                    document.body.appendChild(dialog);

                    // 关闭对话框
                    document.getElementById('close-selected-dialog').addEventListener('click', function() {
                        document.body.removeChild(dialog);
                    });
                });

                // 生成执行命令
                document.getElementById('execute-selected-btn').addEventListener('click', function() {
                    if (state.selectedServices.size === 0) {
                        showToast('请先选择要执行的服务');
                        return;
                    }

                    let commands = "";
                    Array.from(state.selectedServices).sort().forEach(function(id) {
                        const parts = id.split('.', 2);
                        const type = parts[0];
                        const name = parts[1];
                        commands += "./yq.sh " + type + " " + name + "\n";
                    });

                    // 显示命令对话框
                    const dialog = document.createElement('div');
                    dialog.className = 'command-dialog';
                    dialog.innerHTML = 
                        '<div class="dialog-content">' +
                            '<h3>执行命令</h3>' +
                            '<textarea id="command-output" readonly>' + commands + '</textarea>' +
                            '<div class="dialog-buttons">' +
                                '<button id="copy-commands">复制命令</button>' +
                                '<button id="close-dialog">关闭</button>' +
                            '</div>' +
                        '</div>';
                    document.body.appendChild(dialog);

                    // 复制命令
                    document.getElementById('copy-commands').addEventListener('click', function() {
                        const textarea = document.getElementById('command-output');
                        textarea.select();
                        document.execCommand('copy');
                        showToast('命令已复制到剪贴板');
                    });

                    // 关闭对话框
                    document.getElementById('close-dialog').addEventListener('click', function() {
                        document.body.removeChild(dialog);
                    });
                });

                // 服务过滤
                document.getElementById('category-filter').addEventListener('change', function() {
                    filterServices();
                });

                document.getElementById('search-filter').addEventListener('input', function() {
                    filterServices();
                });

                function filterServices() {
                    const searchTerm = document.getElementById('search-filter').value.toLowerCase();
                    const category = document.getElementById('category-filter').value;

                    state.currentServices = window.serviceMetadata.filter(function(service) {
                        return (category === 'all' || service.category === category) &&
                               (searchTerm === '' ||
                                service.name.toLowerCase().includes(searchTerm) ||
                                service.path.toLowerCase().includes(searchTerm));
                    });

                    renderServices(state.currentServices);
                }
            }

            // 初始化渲染和事件
            renderServices(state.currentServices);
            setupEventListeners();
            updateSelectedCount();
        });
    </script>
</body>
</html>`

type FunctionHelp struct {
	jsonFolder string
	jsonFile   string
	outputDir  string
}

// 主函数
func GenHelmCode() {
	// 收集服务数据
	services := collectServices()

	// 生成HTML文件
	outputFile := "/tmp/services_dashboard.html"
	if err := generateHTML(services, outputFile); err != nil {
		fmt.Printf("生成HTML失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("服务仪表盘已生成到 %s\n", outputFile)

	// 尝试自动打开浏览器
	if err := openBrowser(outputFile); err != nil {
		fmt.Printf("自动打开浏览器失败，您可以手动打开文件: %v\n", err)
		fmt.Printf("请打开浏览器访问该文件，或使用命令: open %s\n", outputFile)
	}
}

// 收集所有服务信息
func collectServices() []Service {
	var services []Service

	// 添加单文件服务
	addServiceIfExists := func(path, typ, category string) {
		if _, err := os.Stat(path); err == nil {
			name := strings.TrimSuffix(filepath.Base(path), ".yaml")
			services = append(services, Service{
				Type:     typ,
				Category: category,
				Name:     name,
				Path:     path,
			})
		}
	}

	// 单文件服务
	addServiceIfExists("./deploy/helm/http/values.yaml", "http", "HTTP")
	addServiceIfExists("./deploy/helm/web/values.yaml", "web", "WEB")

	// 目录服务扫描
	scanDir := func(dir, typ, category string) {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return
		}

		files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
		if err != nil {
			fmt.Printf("扫描目录 %s 时出错: %v\n", dir, err)
			return
		}

		for _, f := range files {
			name := strings.TrimSuffix(filepath.Base(f), ".yaml")
			services = append(services, Service{
				Type:     typ,
				Category: category,
				Name:     name,
				Path:     f,
			})
		}
	}

	// 扫描目录服务
	scanDir("./deploy/helm/rpc/rpcs", "rpc", "RPC")
	scanDir("./deploy/helm/cmd/cmds", "cmd", "CMD")
	scanDir("./deploy/helm/cli/tasks", "task", "TASK")

	// 按分类和名称排序
	sort.Slice(services, func(i, j int) bool {
		if services[i].Category != services[j].Category {
			return services[i].Category < services[j].Category
		}
		return services[i].Name < services[j].Name
	})

	return services
}

// 生成HTML文件
func generateHTML(services []Service, outputFile string) error {
	// 将服务数据转换为JSON格式
	servicesJSON, err := json.Marshal(services)
	if err != nil {
		return fmt.Errorf("转换服务数据为JSON失败: %w", err)
	}

	// 创建模板数据
	data := struct {
		ServicesJSON template.JS
	}{
		ServicesJSON: template.JS(servicesJSON),
	}

	// 解析模板
	tmpl, err := template.New("dashboard").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// 创建输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer file.Close()

	// 执行模板并写入文件
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// 自动打开浏览器
func openBrowser(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// 根据操作系统选择合适的命令
	var cmd string
	var args []string

	switch os := runtime.GOOS; os {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", absPath}
	case "darwin": // macOS
		cmd = "open"
		args = []string{absPath}
	case "linux":
		cmd = "xdg-open"
		args = []string{absPath}
	default:
		return fmt.Errorf("不支持的操作系统: %s", os)
	}

	return exec.Command(cmd, args...).Start()
}
