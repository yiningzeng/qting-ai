#!/bin/bash
versionId=0
buildTime=$(date "+%Y-%m-%d %H:%M:%S")
packageName="qting-ai"

print_help() {
cat <<EOF
usage:  build.sh  -v 编译的版本号
                  -o 编译的程序名
                  -h 帮助
EOF
exit 1
}

while getopts "v:o:h" opt; do
  case $opt in
    v)
      versionId=$OPTARG
      ;;
    o)
      packageName=$OPTARG
      ;;
    h)
      print_help
      ;;
    \?)
      print_help
      ;;
  esac
done
go build -buildmode=plugin -o plugins/QTing-tiny-3l-single.yn plugins/QTing-tiny-3l-single/main.go
go build -buildmode=plugin -o plugins/QTing-tiny-3l-multilabel.yn plugins/QTing-tiny-3l-multilabel/main.go
go build -ldflags  "-X \"qting-ai/version.BuildDate=$buildTime\" -X \"qting-ai/version.ID=$versionId\"" -o "$packageName"