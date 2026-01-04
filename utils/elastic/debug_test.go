package elastic

import (
	"fmt"
	"testing"
)

// TestDebugIndexConfig 调试索引配置解析
func TestDebugIndexConfig(t *testing.T) {
	// 调试CustomIndexModel的索引配置
	customModel := CustomIndexModel{}
	config := GetIndexConfig(customModel)
	fmt.Printf("CustomIndexModel config: %+v\n", config)
	fmt.Printf("CustomIndexModel index pattern: %s\n", GetIndexPattern(customModel))

	// 调试NoPartitionModel的索引配置
	noPartitionModel := NoPartitionModel{}
	config = GetIndexConfig(noPartitionModel)
	fmt.Printf("NoPartitionModel config: %+v\n", config)
	fmt.Printf("NoPartitionModel index pattern: %s\n", GetIndexPattern(noPartitionModel))
}
