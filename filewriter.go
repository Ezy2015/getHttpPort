package main

import (
	"fmt"
	"os"
)

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)         // 查找文件末尾的偏移量
		_, err = f.WriteAt([]byte(content), n) // 从末尾的偏移量开始写入内容
	}
	defer f.Close()
	return err
}
