/**
 * @Title
 * @Description
 * @Author hyman
 * @Date 2022-03-19
 **/
package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ViewCmd() *cobra.Command {
	var tags []string
	var env string
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "最终文件替换效果",
		Long:  "显示根据规则最终文件替换结果的文件树",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			viewHandle(args, env, tags)
		},
	}
	viewCmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "文件标识，从左到右优先级递增")
	viewCmd.Flags().StringVarP(&env, "env", "e", "", "部署环境")
	err := viewCmd.MarkFlagRequired("env")
	if err != nil {
		panic(err)
	}
	return viewCmd
}

func viewHandle(args []string, env string, tags Tags) {
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
	root.Children = ComputeTarget(root, tags)
	PrintFileTree(nil, 2, root)
}

func ComputeTarget(t *TargetFile, tags []string) TargetFiles {
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
	rootCmd.AddCommand(ViewCmd())
}
