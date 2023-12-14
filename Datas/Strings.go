package Datas

func StringToUint32(str string) (int32Slice []int32) {
	// 将字符串转换为rune类型
	runes := []rune(str)

	// 创建一个int32类型的切片
	//var int32Slice []int32

	// 将rune类型的值转换为int32，并存储到切片中
	for _, r := range runes {
		int32Slice = append(int32Slice, int32(r))
	}
	return
}

func StringToUint16(str string) (int32Slice []uint16) {
	// 将字符串转换为rune类型
	runes := []rune(str)

	// 创建一个int32类型的切片
	//var int32Slice []int32

	// 将rune类型的值转换为int32，并存储到切片中
	for _, r := range runes {
		int32Slice = append(int32Slice, uint16(r))
	}
	return
}
