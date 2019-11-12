package handle

type Handle struct {
	// 消息处理句柄名,非必要
	Name string `json:"name"`

	//
	Module string `json:"module"`
}

func Init() {}
