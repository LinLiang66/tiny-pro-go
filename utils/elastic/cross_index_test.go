package elastic

import (
	"strings"
	"testing"
	"time"
)

// 测试模型
type TestUser struct {
	BaseModel
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// 带自定义索引名的测试模型
type CustomIndexModel struct {
	BaseModel
	Name  string `json:"name"`
	Index string `es:"custom_index"`
}

// 按月分索引的测试模型
type MonthlyPartitionModel struct {
	BaseModel
	Name  string `json:"name"`
	Index string `es:"monthly_model,partition=monthly"`
}

// 按年分索引的测试模型
type YearlyPartitionModel struct {
	BaseModel
	Name  string `json:"name"`
	Index string `es:"yearly_model,partition=yearly"`
}

// 不分索引的测试模型
type NoPartitionModel struct {
	BaseModel
	Name  string `json:"name"`
	Index string `es:"no_partition_model,partition=none"`
}

// 使用默认索引名的按月分索引测试模型
type DefaultPartitionModel struct {
	BaseModel
	Name string `json:"name"`
}

// TestGetIndexPattern 测试获取索引模式
func TestGetIndexPattern(t *testing.T) {
	// 测试默认索引名
	user := TestUser{}
	pattern := GetIndexPattern(user)
	expectedPattern := "testuser_*"
	if pattern != expectedPattern {
		t.Errorf("GetIndexPattern() = %v, expected %v", pattern, expectedPattern)
	}

	// 测试自定义索引名
	custom := CustomIndexModel{}
	pattern = GetIndexPattern(custom)
	expectedPattern = "custom_index_*"
	if pattern != expectedPattern {
		t.Errorf("GetIndexPattern() = %v, expected %v", pattern, expectedPattern)
	}
}

// TestGetIndexNamesByDateRange 测试根据日期范围获取索引名列表
func TestGetIndexNamesByDateRange(t *testing.T) {
	user := TestUser{}

	// 测试同一年同一月
	startTime := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	indexNames := GetIndexNamesByDateRange(user, startTime, endTime, "200601")
	expectedNames := []string{"testuser_202601"}
	if len(indexNames) != len(expectedNames) {
		t.Errorf("GetIndexNamesByDateRange() length = %d, expected %d", len(indexNames), len(expectedNames))
	}
	for i, name := range indexNames {
		if name != expectedNames[i] {
			t.Errorf("GetIndexNamesByDateRange()[%d] = %v, expected %v", i, name, expectedNames[i])
		}
	}

	// 测试跨月
	startTime = time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	endTime = time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC)
	indexNames = GetIndexNamesByDateRange(user, startTime, endTime, "200601")
	expectedNames = []string{"testuser_202601", "testuser_202602"}
	if len(indexNames) != len(expectedNames) {
		t.Errorf("GetIndexNamesByDateRange() length = %d, expected %d", len(indexNames), len(expectedNames))
	}
	for i, name := range indexNames {
		if name != expectedNames[i] {
			t.Errorf("GetIndexNamesByDateRange()[%d] = %v, expected %v", i, name, expectedNames[i])
		}
	}

	// 测试跨年
	startTime = time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
	endTime = time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	indexNames = GetIndexNamesByDateRange(user, startTime, endTime, "200601")
	expectedNames = []string{"testuser_202512", "testuser_202601"}
	if len(indexNames) != len(expectedNames) {
		t.Errorf("GetIndexNamesByDateRange() length = %d, expected %d", len(indexNames), len(expectedNames))
	}
	for i, name := range indexNames {
		if name != expectedNames[i] {
			t.Errorf("GetIndexNamesByDateRange()[%d] = %v, expected %v", i, name, expectedNames[i])
		}
	}

	// 测试自定义时间格式（按年）
	startTime = time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	endTime = time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	indexNames = GetIndexNamesByDateRange(user, startTime, endTime, "2006")
	expectedNames = []string{"testuser_2024", "testuser_2025", "testuser_2026"}
	if len(indexNames) != len(expectedNames) {
		t.Errorf("GetIndexNamesByDateRange() length = %d, expected %d", len(indexNames), len(expectedNames))
	}
	for i, name := range indexNames {
		if name != expectedNames[i] {
			t.Errorf("GetIndexNamesByDateRange()[%d] = %v, expected %v", i, name, expectedNames[i])
		}
	}
}

// TestGetIndexNamesByDateRangeReverse 测试日期范围反转的情况
func TestGetIndexNamesByDateRangeReverse(t *testing.T) {
	user := TestUser{}

	// 测试startTime > endTime的情况
	startTime := time.Date(2026, 2, 15, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	indexNames := GetIndexNamesByDateRange(user, startTime, endTime, "200601")
	expectedNames := []string{"testuser_202601", "testuser_202602"}
	if len(indexNames) != len(expectedNames) {
		t.Errorf("GetIndexNamesByDateRange() length = %d, expected %d", len(indexNames), len(expectedNames))
	}
	for i, name := range indexNames {
		if name != expectedNames[i] {
			t.Errorf("GetIndexNamesByDateRange()[%d] = %v, expected %v", i, name, expectedNames[i])
		}
	}
}

// TestIndexConfig 测试实体注解配置的分索引策略
func TestIndexConfig(t *testing.T) {
	// 测试默认配置（应该按月分索引）
	defaultModel := DefaultPartitionModel{}
	config := GetIndexConfig(defaultModel)
	if config.BaseIndexName != "defaultpartitionmodel" {
		t.Errorf("Default partition model base index name = %v, expected %v", config.BaseIndexName, "defaultpartitionmodel")
	}
	if config.PartitionType != "monthly" {
		t.Errorf("Default partition model type = %v, expected %v", config.PartitionType, "monthly")
	}
	if config.TimeSuffixFormat != "200601" {
		t.Errorf("Default partition model time suffix format = %v, expected %v", config.TimeSuffixFormat, "200601")
	}

	// 测试按月分索引的配置
	monthlyModel := MonthlyPartitionModel{}
	config = GetIndexConfig(monthlyModel)
	if config.BaseIndexName != "monthly_model" {
		t.Errorf("Monthly partition model base index name = %v, expected %v", config.BaseIndexName, "monthly_model")
	}
	if config.PartitionType != "monthly" {
		t.Errorf("Monthly partition model type = %v, expected %v", config.PartitionType, "monthly")
	}
	if config.TimeSuffixFormat != "200601" {
		t.Errorf("Monthly partition model time suffix format = %v, expected %v", config.TimeSuffixFormat, "200601")
	}

	// 测试按年分索引的配置
	yearlyModel := YearlyPartitionModel{}
	config = GetIndexConfig(yearlyModel)
	if config.BaseIndexName != "yearly_model" {
		t.Errorf("Yearly partition model base index name = %v, expected %v", config.BaseIndexName, "yearly_model")
	}
	if config.PartitionType != "yearly" {
		t.Errorf("Yearly partition model type = %v, expected %v", config.PartitionType, "yearly")
	}
	if config.TimeSuffixFormat != "2006" {
		t.Errorf("Yearly partition model time suffix format = %v, expected %v", config.TimeSuffixFormat, "2006")
	}

	// 测试不分索引的配置
	noPartitionModel := NoPartitionModel{}
	config = GetIndexConfig(noPartitionModel)
	if config.BaseIndexName != "no_partition_model" {
		t.Errorf("No partition model base index name = %v, expected %v", config.BaseIndexName, "no_partition_model")
	}
	if config.PartitionType != "none" {
		t.Errorf("No partition model type = %v, expected %v", config.PartitionType, "none")
	}
	if config.TimeSuffixFormat != "" {
		t.Errorf("No partition model time suffix format = %v, expected %v", config.TimeSuffixFormat, "")
	}
}

// TestGetIndexNameWithPartition 测试根据配置获取索引名
func TestGetIndexNameWithPartition(t *testing.T) {
	// 测试按月分索引的模型
	monthlyModel := MonthlyPartitionModel{}
	monthlyIndex := GetIndexName(monthlyModel)
	if !strings.Contains(monthlyIndex, "monthly_model_") {
		t.Errorf("Monthly partition model index name = %v, expected to contain %v", monthlyIndex, "monthly_model_")
	}
	if len(monthlyIndex) != len("monthly_model_")+6 { // "200601"格式有6位
		t.Errorf("Monthly partition model index name length = %d, expected %d", len(monthlyIndex), len("monthly_model_")+6)
	}

	// 测试按年分索引的模型
	yearlyModel := YearlyPartitionModel{}
	yearlyIndex := GetIndexName(yearlyModel)
	if !strings.Contains(yearlyIndex, "yearly_model_") {
		t.Errorf("Yearly partition model index name = %v, expected to contain %v", yearlyIndex, "yearly_model_")
	}
	if len(yearlyIndex) != len("yearly_model_")+4 { // "2006"格式有4位
		t.Errorf("Yearly partition model index name length = %d, expected %d", len(yearlyIndex), len("yearly_model_")+4)
	}

	// 测试不分索引的模型
	noPartitionModel := NoPartitionModel{}
	noPartitionIndex := GetIndexName(noPartitionModel)
	if noPartitionIndex != "no_partition_model" {
		t.Errorf("No partition model index name = %v, expected %v", noPartitionIndex, "no_partition_model")
	}
}

// TestGetIndexPatternWithPartition 测试根据配置获取索引模式
func TestGetIndexPatternWithPartition(t *testing.T) {
	// 测试按月分索引的模型
	monthlyModel := MonthlyPartitionModel{}
	monthlyPattern := GetIndexPattern(monthlyModel)
	if monthlyPattern != "monthly_model_*" {
		t.Errorf("Monthly partition model index pattern = %v, expected %v", monthlyPattern, "monthly_model_*")
	}

	// 测试按年分索引的模型
	yearlyModel := YearlyPartitionModel{}
	yearlyPattern := GetIndexPattern(yearlyModel)
	if yearlyPattern != "yearly_model_*" {
		t.Errorf("Yearly partition model index pattern = %v, expected %v", yearlyPattern, "yearly_model_*")
	}

	// 测试不分索引的模型
	noPartitionModel := NoPartitionModel{}
	noPartitionPattern := GetIndexPattern(noPartitionModel)
	if noPartitionPattern != "no_partition_model" {
		t.Errorf("No partition model index pattern = %v, expected %v", noPartitionPattern, "no_partition_model")
	}
}
