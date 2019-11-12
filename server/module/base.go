package module

type Module interface {
	// 创建模块对象
	Create()

	// 初始化模块
	Init()

	// 关闭模块
	Close()

	// 调用模块
	Call(interface{})

	// 信号处理
	Signal(signal uint32)
}

type Modules struct {
	// 当前模块数
	Count uint32 `json:"count"`

	// module映射
	ModuleMap map[string]*Module `json:"module_map"`
}

func (this *Modules) AddModule() {}
