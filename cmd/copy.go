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
	dstStat, err := os.Lstat(dst)
	if !os.IsNotExist(err) { // 文件存在
		if dstStat.IsDir() != t.IsDir() { // 如果类型不匹配
			return fmt.Errorf("目标类型不匹配：%s(dir=%t) != %s(dir=%t)", dst, dstStat.IsDir(), t.Path(), t.IsDir())
		}
		if dstStat.IsDir() && dstStat.IsDir() {
			fmt.Println("dir exist: " + dst)
			return nil
		}
		// 删除文件
		err = os.Remove(dst)
		if err != nil {
			return err
		}
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
