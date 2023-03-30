# jupyter官方镜像提供，为jupyter所需最小docker
FROM jupyter/base-notebook:latest

USER root

# 安装时不interative
ARG DEBIAN_FRONTEND=noninteractive

# 安装系统基础命令和依赖
RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends figlet curl git git-lfs vim curl wget \
    && apt-get install ffmpeg libsm6 libxext6  -y \
    && rm -rf /var/lib/apt/lists/*

# 修改原有用户
# /home/<原有用户>用户文件夹没删除，需再分析
RUN usermod -l mindspore -d /home/mindspore -m ${NB_USER} \
    && groupmod -n mindspore users \
    && sudo rm -rf /home/${NB_USER}

# 更改原有镜像env
ENV NB_USER=mindspore 
ENV CHOWN_HOME=yes
ENV DOCKER_STACKS_JUPYTER_CMD=lab
ENV HOME=/home/mindspore
ENV PWD=/home/mindspore

# 切换用户
USER $NB_UID

# 安装依赖
COPY --chown=mindspore:mindspore  requirements.txt  ./
RUN pip install --upgrade -i https://pypi.mirrors.ustc.edu.cn/simple pip
RUN pip install --default-timeout=1200 --user -i https://pypi.mirrors.ustc.edu.cn/simple -r requirements.txt
RUN rm requirements.txt

# 原有镜像ls为ls -a --color，这里去掉-a
RUN echo "figlet MindSpore" >> ~/.bashrc \
    && echo "alias ls='ls --color'" >> ~/.bashrc \
    && echo "alias ll='ls -l'" >> ~/.bashrc

# 更改hostname，没生效，需更改/etc/hostname
ENV HOSTNAME=xihe

RUN jupyter lab --generate-config
RUN echo 'c.ServerApp.token="f9a3bd4e9f2c3be01cd629154cfb224c2703181e050254b5"' >> ~/.jupyter/jupyter_lab_config.py 

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
# http://127.0.0.1:8080/lab?token=f9a3bd4e9f2c3be01cd629154cfb224c2703181e050254b5
CMD jupyter lab  --notebook-dir=/home/mindspore/work --ip='*' --port=8080 --no-browser
