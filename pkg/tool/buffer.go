package tool

import "bytes"

// 简单的buffer实现

type Buffer struct {
	bytes.Buffer
}

func (m *Buffer) WriteString(s ...string) *Buffer {
	for _, v := range s {
		m.Buffer.WriteString(v)
	}
	return m
}
