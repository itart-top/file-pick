/**
 * @Title
 * @Description
 * @Author hyman
 * @Date 2022-03-19
 **/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RunCmd() *cobra.Command {
	var tags []string
	var env string
	runCmd := &cobra.Command{
		Use:   "publish",
		Short: "将文件发布到指定文件夹",
		Long:  "根据规则将最终文件发布到指定的文件夹",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			publishRun(args, env, tags)
		},
	}
	runCmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "文件标识，从左到右优先级递增")
	runCmd.Flags().StringVarP(&env, "env", "e", "", "部署环境")
	err := runCmd.MarkFlagRequired("env")
	if err != nil {
		panic(err)
	}
	return runCmd
}

func publishRun(args []string, env string, tags Tags) {
	info, err := os.Stat(filepath.Join(args[0], env))
	if err != nil {
		panic(err)
	}
	root := &TargetFile{
		Name: info.Name(),
		Default: &File{
			Name: info.Name(),
			Dir:  args[0],
			Info: info,
		},
	}
	root.Children = GenTarget(root, tags)
	for _, sub := range root.Children {
		if err = publish(sub, args[1]); err != nil {
			fmt.Println("publish error:", err)
			panic(err)
		}
	}
}

func GenTarget(t *TargetFile, tags []string) TargetFiles {
	pp, err := ioutil.ReadDir(t.Path())
	if err != nil {
		panic(err)
	}
	pick := NewFilePick()
	for _, p := range pp {
		pick.Add(t.Path(), p)
	}
	targets := pick.Pick(tags)
	for _, t := range targets {
		if t.IsDir() {
			t.Children = GenTarget(t, tags)
		}
	}
	return targets
}

func init() {
	rootCmd.AddCommand(RunCmd())
}
