#!/bin/sh
# Adds the git-hook described below. Appends to the hook file
# if it already exists or creates the file if it does not.
# Note: CWD must be inside target repository


if curl http://nexus.hyperchain.cn/repository/blocface/blocface/githooks/githook.tar.gz -o githook.tar.gz ; then
    echo "download githook file  "
    else
      echo "download githook file failed"
      exit 1
fi

if tar -zxf githook.tar.gz ; then
    echo "extract githook file"
    else
      echo "extract githook file failed"
      exit 1
fi

SAMPLE_DIR=$(tar -tzf githook.tar.gz | head -1 | cut -f1 -d"/")
rm -rf githook.tar.gz
HOOK_DIR=$(git rev-parse --show-toplevel)/.git/hooks

for template_file in "$SAMPLE_DIR"/*; do

  HOOK_FILE="$HOOK_DIR"/${template_file##*/}
  echo "update ${HOOK_FILE}"
  # Create script file if doesn't exist
  if [ ! -e "$HOOK_FILE" ] ; then
          echo "#!/bin/sh" >> "$HOOK_FILE"
          chmod 700 "$HOOK_FILE"
  fi
  sed '/#START BLOCFACE/,/#END BLOCFACE/d' "$HOOK_FILE" > "$HOOK_DIR"/temp.txt
  cat "$HOOK_DIR"/temp.txt > "$HOOK_FILE"
  cat "$template_file" >> "$HOOK_FILE"
  rm "$HOOK_DIR"/temp.txt
done
rm -rf "$SAMPLE_DIR"