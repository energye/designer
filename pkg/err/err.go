package err

// 简单的错误管理

// 返回状态类型
type ResultStatus int32

const (
	RsSuccess       ResultStatus = iota // 成功 通用
	RsFail                              // 失败 通用
	RsIgnoreProp                        // 忽略的属性
	RsNotValid                          // 对象无效
	RsDuplicateName                     // 组件名重复
)

// 检测 err
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
