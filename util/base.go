package util

import (
	"bufio"
	"fmt"
	"os"
	"errors"
)

func WriteFile(path string, name string, content string) (int, error) {
    if len(path) == 0 || len(name) == 0 {
	return 0, errors.New("path or name is null.")
    }
	
    f, err := os.OpenFile(path + name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
    if err != nil {
        fmt.Println("Open file error")
	return 0, err
    }
    defer f.Close()
    w := bufio.NewWriter(f)
    len, err := w.WriteString(content)
    //fmt.Println("size : ", len)
    w.Flush()
    return len, err
}

