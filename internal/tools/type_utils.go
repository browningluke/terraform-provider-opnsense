package tools

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

// Ints

func Int64ToString(i int64) string {
	return fmt.Sprintf("%d", i)
}

func StringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return i
	}
	return -1
}

func StringToInt64Null(s string) types.Int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return types.Int64Value(i)
	}
	return types.Int64Null()
}

func Int64ToStringNegative(i int64) string {
	s := fmt.Sprintf("%d", i)
	if i == -1 {
		s = ""
	}
	return s
}

// Bools

func BoolToString(b bool) string {
	if b {
		return "1"
	} else {
		return "0"
	}
}

func StringToBool(s string) bool {
	return s == "1"
}

// Strings

func StringOrNull(s string) types.String {
	if s != "" {
		return types.StringValue(s)
	} else {
		return types.StringNull()
	}
}
