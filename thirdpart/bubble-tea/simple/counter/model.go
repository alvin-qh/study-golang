package counter

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// 定义模型类型
//
// 该类型实现了 `tea.Model` 接口, 包括 `Init`, `Update` 和 `View` 方法
type Model struct {
	count int // 计数值
}

// 创建模型对象
func NewModel() tea.Model {
	return &Model{
		count: 0,
	}
}

// 初始化模型
func (m *Model) Init() tea.Cmd {
	return nil
}

// 用于根据用户操作更新模型状态
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit // 按下 'q' 退出程序
		case "+", "=":
			m.count += 1
		case "-", "_":
			m.count -= 1
		}
		return m, nil
	}
	return m, nil
}

// 用于在控制台显式界面内容
func (m *Model) View() string {
	return fmt.Sprintf("Count: %d\nPress + to increment, - to decrement, q to quit.\n", m.count)
}
