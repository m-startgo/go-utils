package mmath

// Sum 多数求和，传入 0 个参数时返回 Zero
func Sum(nums ...Decimal) Decimal {
	if len(nums) == 0 {
		return Zero
	}
	sum := Zero
	for _, n := range nums {
		sum = sum.Add(n)
	}
	return sum
}

// Mean 计算算术平均数，传入 0 个参数时返回 Zero
func Mean(nums ...Decimal) Decimal {
	if len(nums) == 0 {
		return Zero
	}
	return Sum(nums...).Div(NewFromInt(int64(len(nums))))
}
