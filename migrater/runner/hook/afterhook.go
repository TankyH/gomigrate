package hook

/*
	用于展示After Hook 检查或者执行的时候记录错误的方式
*/
type HookErr struct {
	ModuleName string
	FolderName string
	Code       int
}

//用于检查AfterHook的运行是否有报错，如果有，则返回错误码
func (e HookErr) HasError() bool {
	if e.Code == 0 {
		return false
	}
	return true
}
