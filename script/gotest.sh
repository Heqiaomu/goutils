#!/usr/bin/env bash
# cannot set -e to show curl error
### set -e

# WORKPATH=`pwd`
WORKPATH=""

if [ "$CI_PROJECT_DIR" = "" ]; then
current_path=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
WORKPATH="${current_path}/.."
else
WORKPATH=$CI_PROJECT_DIR # k8s ci
fi

#### ===============================================================

# {WORKPATH}/output/cover/
if [ -d "${WORKPATH}/output/cover/"  ]; then
    rm -rf "${WORKPATH}/output/cover"
fi
# mkdir
mkdir -p "${WORKPATH}/output/cover/"

# TODO include k8s and move *.pb.go & *.pb.gw.go to /proto
PKG_LIST=$(go list "${WORKPATH}/..." | grep -v /third | grep -v /proto)

echo "test coverage counting..."
st=0     # 返回值
fail=0   # 失败数
pass=0   # 成功数
for package in ${PKG_LIST}; do
	echo "cover ==> $package"

if ! go test -v -gcflags=-l -covermode=count -coverprofile "${WORKPATH}/output/cover/${package##*/}.cov" "$package" ;
then
  fail=$((fail+1)) && echo -e "\033[31mFAIL	${package#*./}\033[0m" && continue;
fi
    
pass=$((pass+1))
echo -e "\033[32mPASS	$package\033[0m"

done
echo "cover..."
tail -q -n +2 "${WORKPATH}"/output/cover/*.cov >> "${WORKPATH}"/output/cover/coverage.out 2> /dev/null

sys="$(uname)"
if [ "${sys}" == "Darwin" ]; then
## Mac OS sed
## i\之后有空格
echo "MacOSX..."
sed -i '' "1 i\ 
mode: count
" "${WORKPATH}"/output/cover/coverage.out
else
echo "Linux..."
## Linux sed
sed -i '1imode: count' "${WORKPATH}"/output/cover/coverage.out
fi

go tool cover -func="${WORKPATH}"/output/cover/coverage.out

echo "generate report..."
go tool cover -html="${WORKPATH}"/output/cover/coverage.out -o "${WORKPATH}"/output/cover/coverage.html

if [ -f "${WORKPATH}"/output/cover/coverage.html ]; then
    echo report  ==> "${WORKPATH}"/output/cover/coverage.html
    ## clean *.cov
    rm -rf "${WORKPATH}"/output/cover/*.cov
else
    echo "report generate fail"
fi

[ $fail -ne 0 ] && st=1 
echo ================================== 
[ $fail -ne 0 ] && echo -e "\033[31m$fail packages FAIL\033[0m"
[ $pass -ne 0 ] && echo -e "\033[32m$pass packages PASS\033[0m"
echo ================================== 
exit $st
