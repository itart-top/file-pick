/**
 * @Title
 * @Description
 * @Author hyman
 * @Date 2022-03-19
 **/
package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

const RWFlag = ".rw.for."

type File struct {
	Name string      // 文件名
	Dir  string      // 目录名
	Info fs.FileInfo // 文件的信息
}

func (f File) Path() string {
	return filepath.Join(f.Dir, f.Name)
}

// 目标文件
type TargetFile struct {
	Name     string        // 文件名
	Default  *File         // 缺省文件
	RW       *TagFile      // 重写文件
	Children []*TargetFile // 子文件
}

func (r *TargetFile) File() *File {
	if r.RW != nil {
		return &r.RW.File
	}
	return r.Default
}

// 文件真实路径
func (r *TargetFile) Path() string {
	return r.File().Path()
}

func (r *TargetFile) Info() fs.FileInfo {
	return r.File().Info
}
func (r *TargetFile) IsDir() bool {
	return r.Info().IsDir()
}

func (r *TargetFile) Dir() string {
	return r.File().Dir
}

func (r *TargetFile) Label() string {
	var buf bytes.Buffer
	buf.WriteString(r.Name)
	if r.RW != nil { // 如果重写了
		buf.WriteString(" -> ")
		buf.WriteString(r.RW.Info.Name())
	}
	format := "%s"
	if r.Default == nil { // 不存在默认文件，指定重写
		format = "\u001b[33;1m%s\u001b[0m"
	} else if r.RW != nil { // 文件重写
		format = "\u001b[32m%s\u001b[0m"
	}
	return fmt.Sprintf(format, buf.String())
}

type TargetFiles []*TargetFile

type Tags []string

/**
标识文件
*/
type TagFile struct {
	File
	tag string // 标识
	src string // 对应源文件名
}

func TageFileParse(dir string, info fs.FileInfo) (*TagFile, bool) {
	name := info.Name()
	i := strings.Index(name, RWFlag)               // 是否包含重写标识
	if i <= 0 || strings.HasSuffix(name, RWFlag) { // 没有tag标识：1. 标识是第一个，2. 标识不存， 3. 标识最后一个
		return nil, false
	}
	return &TagFile{
		File: File{
			Name: info.Name(),
			Dir:  dir,
			Info: info,
		},
		tag: name[i+len(RWFlag):],
		src: name[:i],
	}, true
}

type TagFiles []*TagFile

type FilePick struct {
	targetFiles map[string]*TargetFile
	tagFiles    map[string]TagFiles
}

func NewFilePick() *FilePick {
	return &FilePick{
		targetFiles: make(map[string]*TargetFile),
		tagFiles:    make(map[string]TagFiles),
	}
}
func (f *FilePick) Add(dir string, info fs.FileInfo) {
	if tagPath, ok := TageFileParse(dir, info); ok {
		f.tagFiles[tagPath.tag] = append(f.tagFiles[tagPath.tag], tagPath)
		return
	}
	f.targetFiles[info.Name()] = &TargetFile{
		Name: info.Name(),
		Default: &File{
			Name: info.Name(),
			Dir:  dir,
			Info: info,
		},
	}
}

func (f *FilePick) Pick(tags Tags) TargetFiles {
	for _, t := range tags {
		for _, tagFile := range f.tagFiles[t] {
			if target, ok := f.targetFiles[tagFile.src]; ok {
				target.RW = tagFile
				continue
			}
			// 只有tag的文件也要加入TargetFile
			f.targetFiles[tagFile.src] = &TargetFile{
				Name: filepath.Base(tagFile.src),
				RW:   tagFile,
			}
		}
	}
	var ff TargetFiles
	for _, f := range f.targetFiles {
		ff = append(ff, f)
	}
	return ff
}
