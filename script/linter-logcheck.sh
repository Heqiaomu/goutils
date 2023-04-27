#!/bin/bash

logStands=(
  "^Before[[:space:]].*\.?$"
  "^After[[:space:]].*\.?$"
  "^Start[[:space:]][a-z]+ing[[:space:]].*\.?$"
  "^Finish[[:space:]][a-z]+ing[[:space:]].*\.?$"
  "^Call[[:space:]].*\.?$"
  "^Receive[[:space:]].*\.?$"
  "^Currently[[:space:]].*\.?$"
  "^[A-Z][a-zA-Z]*ing[[:space:]].*\.?$"
  "^Manage[[:space:]]err[[:space:]]when[[:space:]][a-z]+ing[[:space:]].*\.[[:space:]]Now[[:space:]]to[[:space:]].*\.[[:space:]]Err:[[:space:]]\[%v\]\.?$"
  "^Manage[[:space:]]err[[:space:]]when[[:space:]][a-z]+ing[[:space:]].*\,[[:space:]]because[[:space:]].*\.[[:space:]]Please[[:space:]].*\.?$"
  "^Set[[:space:]].*[[:space:]]to[[:space:]].*\.?$"
  "^Currently[[:space:]].*\.[[:space:]]Now enable schema:[[:space:]].*\.?$"
  "^Currently[[:space:]].*\.[[:space:]]There may be risk:[[:space:]].*\.[[:space:]]Please[[:space:]].*\.?$"
  "^Fail[[:space:]]to[[:space:]].*\,[[:space:]]err:[[:space:]]\[%v\]\.?$"
  "^Fail[[:space:]]to[[:space:]].*\,[[:space:]]because[[:space:]].*\.[[:space:]]Please[[:space:]].*\.?$"
  "^.*$"
)

errStands=(
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]because[[:space:]].*\.[[:space:]]please[[:space:]].*$"
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]err:[[:space:]]\[%v\]$"
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]because[[:space:]]\.[[:space:]]for help:[[:space:]].*$"
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]err:[[:space:]]\[%v\]$"
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]because[[:space:]].*$"
  "^fail[[:space:]]to[[:space:]].*,[[:space:]]err:[[:space:]]\[%v\]$"
  "^.*$"
)

function checkArgFormat() {
  Des=$1
  file=$2
  Line=$3

  HasErr=0
  if [[ ${Des} =~ [=:a-zA-Z*^%$]\[%[a-z]\] ]]; then
    HasErr=1
  fi

  if [[ ${Des} =~ ([^[]%[a-z]$)|([^[]%[a-z][^]]) ]]; then
    HasErr=1
  fi

  # 如果不符合任意一种日志格式，就报错！
#  if [[ $HasErr == 1 ]]; then
#    echo "This description does not meet argument standard. Please use '[]' to tag args and there is no other symbols before it. $file:$Line"
#    echo "$Des"
#    echo "Please see http://thoughts.hyperchain.cn:8099/workspaces/5ebd00c48db9ae00116a46c5/docs/60b6331cbe825b00014d9c00 for more information"
#    echo ""
#    HasErr=1
#  fi
#
#  if [[ $HasErr == 1 ]]; then
#    return 1
#  else
#    return 0
#  fi
  return 0
}

function checkLogDesFormat() {
  Des=$1
  file=$2
  Line=$3

  HasErr=0
  IsRight=1
  for element in ${logStands[*]}; do
    # 如果符合其中一种日志格式，就通过检测
    if [[ ${Des} =~ ${element} ]]; then
      IsRight=0
      break
    fi
  done

  # 如果不符合任意一种日志格式，就报错！
  if [[ $IsRight == 1 ]]; then
    echo "This log description does not meet the general scenario. $file:$Line"
    echo "$Des"
    echo "Please see http://thoughts.hyperchain.cn:8099/workspaces/5ebd00c48db9ae00116a46c5/docs/60b6331cbe825b00014d9c00 for more information"
    echo ""
    HasErr=1
  fi

  if [[ $HasErr == 1 ]]; then
    return 1
  else
    return 0
  fi
}

function checkErrDesFormat() {
  Des=$1
  file=$2
  Line=$3

  HasErr=0
  IsErrRight=1
  for element in ${errStands[*]}; do
    # 如果符合其中一种错误格式，就通过检测
    if [[ ${Des} =~ ${element} ]]; then
      IsErrRight=0
      break
    fi
  done

  # 如果不符合任意一种日志格式，就报错！
  if [[ $IsErrRight == 1 ]]; then
    echo "This err description does not meet the general scenario. $file:$Line"
    echo "$Des"
    echo "Please see http://thoughts.hyperchain.cn:8099/workspaces/5ebd00c48db9ae00116a46c5/docs/60b71222be825b0001e3dce3 for more information"
    echo ""
    HasErr=1
  fi

  if [[ $HasErr == 1 ]]; then
    return 1
  else
    return 0
  fi
}

ProjectPath=$1

hasErr=0
# 获取当前新增代码
touch newfiles.txt
touch newcodes.txt

gitDiffType=0
newfs=$(git diff --name-only | grep ".go$" || true)
if [[ $newfs == "" ]]; then
  newfs=$(git diff HEAD^ HEAD --name-only | grep ".go$" || true)
  gitDiffType=1
fi

echo "$newfs" >newfiles.txt

# 读取文件
while read -r file; do
  # 过滤pd自动生成的文件
  if [[ "$file" =~ ^.*\.pb\.go$ ]]; then
    contine
  fi
  # 过滤test文件
  if [[ "$file" =~ ^.*_test\.go$ ]]; then
      contine
  fi

  if [[ $gitDiffType == 0 ]]; then
    git diff "$ProjectPath/$file" >newcodes.txt
  fi
  if [[ $gitDiffType == 1 ]]; then
    git diff HEAD^ HEAD "$ProjectPath/$file" >newcodes.txt
  fi

  Line=0
  while read -r code; do
    if [[ "$code" =~ ^@@[^-]-[0-9]+,[0-9]+[^+]\+([0-9]+), ]]; then
      Line=${BASH_REMATCH[1]}
      continue
    fi

    # 跳过注释
    if [[ "$code" =~ ^[^/]*//.*$ ]]; then
      Line=$((Line + 1))
      continue
    fi

    # 提取日志描述信息
    if [[ "$code" =~ ^\+[^+].*log\.(Debug|Error|Info|Warn|Fatal)f?\((fmt.Sprintf\()?\"([^\"]*)\" ]]; then
      LogDes="${BASH_REMATCH[3]}"
      if ! checkLogDesFormat "$LogDes" "$file" "$Line"; then
        hasErr=1
      fi
      if ! checkArgFormat "$LogDes" "$file" "$Line"; then
        hasErr=1
      fi
    fi

    # log.Debug(err.Error())
    if [[ "$code" =~ ^\+[^+].*(log\.(Debug|Error|Info|Warn|Fatal)f?\([^\"]*\))$ ]]; then
      logStatement="${BASH_REMATCH[1]}"
      echo "This log statement does not meet the general scenario. $file:$Line"
      echo "$logStatement"
      echo "Please see http://thoughts.hyperchain.cn:8099/workspaces/5ebd00c48db9ae00116a46c5/docs/60b6331cbe825b00014d9c00 for more information"
      echo ""
      hasErr=1
    fi

    # 提取错误描述信息
    if [[ "$code" =~ ^\+[^+].*((fmt)|(errors))\.Errorf?\(\"([^\"]*)\" ]]; then
      ErrDes="${BASH_REMATCH[4]}"
      if ! checkErrDesFormat "$ErrDes" "$file" "$Line"; then
        hasErr=1
      fi
      if ! checkArgFormat "$ErrDes" "$file" "$Line"; then
        hasErr=1
      fi
    fi

    # 判断错误返回前是否wrap
    if [[ "$code" =~ ^\+[^+][[:space:]]*(return.*[[:space:]]err)$ ]]; then
      RawErrReturn="${BASH_REMATCH[1]}"
      echo "This err return does not use errors.Wrap(). $file:$Line"
      echo "$RawErrReturn"
      echo "Please see http://thoughts.hyperchain.cn:8099/workspaces/5ebd00c48db9ae00116a46c5/docs/60b71222be825b0001e3dce3 for more information"
      echo ""
      hasErr=1
    fi

    if [[ "$code" =~ (^[^-].*)|(^$) ]]; then
      #      echo "$code"
      Line=$((Line + 1))
    fi
  done \
    <newcodes.txt
done <newfiles.txt

rm -rf newfiles.txt
rm -rf newcodes.txt

if [[ $hasErr == 1 ]]; then
  exit 1
fi
