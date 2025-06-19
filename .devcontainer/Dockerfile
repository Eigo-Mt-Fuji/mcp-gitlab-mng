# .devcontainer/Dockerfile

# ベースイメージとしてGoの公式イメージを使用
FROM golang:1.23

# 必要なツールのインストール
RUN apt-get update && apt-get install -y \
    git \
    curl \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y apt-transport-https ca-certificates gnupg2 lsb-release

RUN curl -sS https://download.docker.com/linux/debian/gpg | gpg --dearmor > /usr/share/keyrings/docker-ce.gpg && \
    chmod a+r /usr/share/keyrings/docker-ce.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-ce.gpg] https://download.docker.com/linux/debian $(lsb_release -sc) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Docker CLIのインストール
RUN apt-get update && \
    apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

RUN apt-get install -y jq sl

# ワークスペースディレクトリの設定
WORKDIR /workspace
