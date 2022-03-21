/**
 * @Title
 * @Description
 * @Author hyman
 * @Date 2022-03-20
 **/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func publish(t *TargetFile, to string) error {
	err := Copy(t, to)
	if err != nil {
		return err
	}
	for _, c := range t.Children {
		err = publish(c, filepath.Join(to, t.Name))
		if err != nil {
			return err
		}
	}
	return err
}

func Copy(t *TargetFile, to string) error {
	dst := filepath.Join(to, t.Name)
	if len(dst) < 10 { // 防御，
		return fmt.Errorf("path %s too short", dst)
	}
	_, err := os.Lstat(dst)
	if !os.IsNotExist(err) {
		return fmt.Errorf("path %s exist", dst)
	}
	if t.IsDir() { // 文件夹
		fmt.Println("make dir: " + dst)
		return os.MkdirAll(dst, t.Info().Mode()) // 创建文件夹
	}
	fmt.Println(fmt.Sprintf("copy file: from %s to %s", t.Path(), dst))
	// 文件的话需要先创建文件夹，在创建文件
	err = os.MkdirAll(filepath.Dir(dst), t.Info().Mode()) // 创建文件夹
	if err != nil {
		return err
	}
	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	srcFile, err := os.Open(t.Path())
	if err != nil {
		return err
	}
	defer srcFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}
