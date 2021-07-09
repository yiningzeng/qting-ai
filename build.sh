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
rm qting-ai
rm plugins/*.yn
go build -buildmode=plugin -gcflags="all=-N -l" -o plugins/QTing-tiny-3l-single.yn plugins/QTing-tiny-3l-single/main.go plugins/QTing-tiny-3l-single/func.go
go build -buildmode=plugin -gcflags="all=-N -l" -o plugins/QTing-tiny-3l-multilabel.yn plugins/QTing-tiny-3l-multilabel/main.go plugins/QTing-tiny-3l-multilabel/func.go
go build -gcflags="all=-N -l" -ldflags  "-X \"qting-ai/version.BuildDate=$buildTime\" -X \"qting-ai/version.ID=$versionId\"" -o "$packageName"