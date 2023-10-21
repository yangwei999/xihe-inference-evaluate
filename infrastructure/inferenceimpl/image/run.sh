#!/bin/bash
set +o history
set -euo pipefail

project_name=${PROJECT_NAME}
app=app.py
work_dir=/workspace
inference_dir=$work_dir/$project_name/inference

# run user define app
cd $inference_dir

f=requirements.txt
if [ -e "$f" -a -s "$f" ]; then
    pip install --upgrade -i https://pypi.tuna.tsinghua.edu.cn/simple pip
    pip install -r ./$f -i https://pypi.tuna.tsinghua.edu.cn/simple
fi

python3 ./$app
