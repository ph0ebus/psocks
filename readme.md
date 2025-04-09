# psocks - 轻量级 SOCKS5 代理服务器

![GitHub](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/github/go-mod/go-version/ph0ebus/psocks)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows%20%7C%20macOS-brightgreen.svg)

**psocks** 是一个简单易用的 SOCKS5 代理服务器实现，专为个人使用而设计，注重轻量化和易用性，避免复杂的配置。

## 功能特性

- 🚀 纯 SOCKS5 协议实现
- ⚡ 轻量高效，资源占用低
- 🖥️ 跨平台支持 (Linux/Windows/macOS)
- 📦 简洁的代码结构，易于理解和修改

## 当前状态

⚠ **注意**: 该项目仍处于开发阶段，尚未经过大规模测试，可能存在未预料的 bug。暂不支持 IPv6。

## 快速开始

### 安装

```bash
git clone https://github.com/ph0ebus/psocks.git
cd psocks
go build main.go -o psocks # 编译
./main # 运行
```
### 使用

```
-host string
    监听地址 (default "0.0.0.0")
-port int
    监听端口 (default 1080)
-socks string
    socks5代理地址, 默认为空
```

简单开启一个 SOCKS5 正向代理服务器，监听 1080 端口

```bash
./psocks -host 0.0.0.0 -port 1080
```

sokcs5 代理链

```bash
# machine 1
./psocks -host 0.0.0.0 -port 1080 -socks 192.168.0.10:1081
```

```bash
# machine 2
./psocks -host 0.0.0.0 -port 1081
```
## 许可证
本项目基于 MIT 许可证发布，详见 [LICENSE](LICENSE) 文件。