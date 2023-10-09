# jupyter官方镜像提供，为jupyter所需最小docker
FROM jupyter/base-notebook:ed2908bbb62e

USER root

# 安装时不interative
ARG DEBIAN_FRONTEND=noninteractive

# 安装系统基础命令和依赖
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends figlet git git-lfs vim wget \
    && apt-get install ffmpeg libsm6 libxext6  -y \
    && rm -rf /var/lib/apt/lists/*

# 修改原有用户
RUN usermod -l mindspore -d /home/mindspore -m ${NB_USER} \
    && groupmod -n mindspore users

WORKDIR /home/mindspore

# 更改原有镜像env
ENV NB_USER=mindspore 
ENV CHOWN_HOME=yes
ENV DOCKER_STACKS_JUPYTER_CMD=lab
ENV HOME=/home/mindspore
ENV PWD=/home/mindspore

# 切换用户
USER $NB_UID

# 安装依赖
RUN pip install -i https://pypi.mirrors.ustc.edu.cn/simple mindspore
# RUN pip install https://ms-release.obs.cn-north-4.myhuaweicloud.com/1.10.1/MindSpore/cpu/x86_64/mindspore-1.10.1-cp39-cp39-linux_x86_64.whl --trusted-host ms-release.obs.cn-north-4.myhuaweicloud.com -i https://pypi.tuna.tsinghua.edu.cn/simple
COPY --chown=mindspore:mindspore  requirements.txt  ./
RUN pip install --upgrade -i https://pypi.mirrors.ustc.edu.cn/simple pip 
RUN pip install --default-timeout=1200 -i https://pypi.mirrors.ustc.edu.cn/simple -r requirements.txt
RUN rm ./requirements.txt

# 原有镜像ls为ls -a --color，这里去掉-a
RUN echo "figlet MindSpore" >> ~/.bashrc \
    && echo "alias ls='ls --color'" >> ~/.bashrc \
    && echo "alias ll='ls -l'" >> ~/.bashrc

# 更改hostname
ENV HOSTNAME=mindspore

RUN jupyter lab --generate-config
RUN echo "c.ServerApp.token=\"$JUPYTER_TOKEN\"" >> ~/.jupyter/jupyter_lab_config.py 

RUN mkdir -p /home/mindspore/.jupyter/lab/user-settings/@jupyterlab/terminal-extension \
    && echo '{ "theme": "dark" }' > /home/mindspore/.jupyter/lab/user-settings/@jupyterlab/terminal-extension/plugin.jupyterlab-settings

# 设置工作路径
WORKDIR /home/mindspore/work

# 将demo复制到镜像中
COPY --chown=mindspore:mindspore  傻瓜式教程.ipynb  ./
COPY --chown=mindspore:mindspore  Jupyter操作文档.ipynb  ./

# 暴露端口
EXPOSE 8080 

# 真实执行命令
CMD jupyter lab  --notebook-dir=/home/mindspore/work --ip='*' --port=8080 --no-browser
