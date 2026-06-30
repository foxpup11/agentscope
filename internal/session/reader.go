package session

// Reader 解析不同 Agent 的会话日志
type Reader interface {
	// Read 从给定路径读取会话
	Read(path string) (*Session, error)
	// DetectFormat 检测日志格式是否匹配
	DetectFormat(path string) bool
}
