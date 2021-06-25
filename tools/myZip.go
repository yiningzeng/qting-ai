package tools

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func ZipV1(dst, src string, taskId string, label string, width string, height string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			logrus.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(p string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}
		if fi.IsDir() || !strings.Contains(fi.Name(), taskId) {
			return nil
		}
		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}
		// 替换文件信息中的文件名
		fh.Name = filepath.Base(p)
		if path.Ext(fh.Name) == ".weights" {
			fh.Name	= label + ".weights"
		} else if path.Ext(fh.Name) == ".names" {
			fh.Name	= "labels.names"
		} else if path.Ext(fh.Name) == ".suggest" {
			fh.Name	= "suggest_score.suggest"
		} else if path.Ext(fh.Name) == ".modelInfo" {
			fh.Name	= "model_info.txt"
		} else if path.Ext(fh.Name) == ".cfg" {
			fh.Name	= label + ".cfg"
			if bytes, err := afero.ReadFile(afero.NewOsFs(), p); err == nil {
				newWidth := fmt.Sprintf("width=%s #", width)
				newHeight := fmt.Sprintf("height=%s #", height)
				res := strings.Replace(string(bytes), "width=", newWidth, 1)
				res = strings.Replace(res, "width =", newWidth, 1)

				res = strings.Replace(res, "height=", newHeight, 1)
				res = strings.Replace(res, "height =", newHeight, 1)
				if err = afero.WriteFile(afero.NewOsFs(), path.Join("/tmp/", fh.Name), []byte(res), 0755); err == nil {
					p = path.Join("/tmp/", fh.Name)
				} else {
					logrus.Error(fmt.Sprintf("%+v", errors.Wrap(err, "zip")))
				}
			}
		}
		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}
		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}
		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}


		// 打开要压缩的文件
		fr, err := os.Open(p)
		defer fr.Close()
		if err != nil {
			return err
		}
		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return err
		}
		// 输出压缩的内容
		logrus.Info(fmt.Sprintf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", p, n))
		return nil
	})
}