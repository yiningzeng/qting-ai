#!/bin/bash
versionId=0
buildTime=$(date "+%Y-%m-%d %H:%M:%S")
packageName="qting-ai"
buildPlugin=0
print_help() {
cat <<EOF
usage:  build.sh  -v 编译的版本号
                  -p 是否编译插件,默认不编译 如果输入大于0表示编译
                  -o 编译的程序名
                  -h 帮助
EOF
exit 1
}

while getopts "v:o:p:h" opt; do
  case $opt in
    v)
      versionId=$OPTARG
      ;;
    o)
      packageName=$OPTARG
      ;;
    p)
      buildPlugin=$OPTARG
      ;;
    h)
      print_help
      ;;
    \?)
      print_help
      ;;
  esac
done
rm qting-ai
if [ $buildPlugin -gt 0 ];then
  rm plugins/*.yn
  go build -buildmode=plugin -gcflags="all=-N -l" -o plugins/QTing-fast-multilabel.yn plugins/QTing-fast-multilabel/main.go plugins/QTing-fast-multilabel/func.go
  go build -buildmode=plugin -gcflags="all=-N -l" -o plugins/QTing-tiny-3l-single.yn plugins/QTing-tiny-3l-single/main.go plugins/QTing-tiny-3l-single/func.go
  go build -buildmode=plugin -gcflags="all=-N -l" -o plugins/QTing-tiny-3l-multilabel.yn plugins/QTing-tiny-3l-multilabel/main.go plugins/QTing-tiny-3l-multilabel/func.go
fi

go build -gcflags="all=-N -l" -ldflags  "-X \"qting-ai/version.BuildDate=$buildTime\" -X \"qting-ai/version.ID=$versionId\"" -o "$packageName"