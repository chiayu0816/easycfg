# EasyCfg

EasyCfg 是一個 Go 工具，用於簡化系統配置管理。它可以自動將 YAML 配置文件轉換為 Go 結構體，並使用 Viper 讀取和監控配置變更。

## 功能特點

- 自動將 YAML 配置文件轉換為 Go 結構體
- 生成相應的 Go 文件
- 使用 Viper 讀取 YAML 配置
- 支持配置熱重載
- 支持監控配置文件變更

## 安裝

```bash
go get github.com/chiayu0816/easycfg
```

## 使用方法

### 生成配置結構體

```bash
# 基本用法
go run cmd/easycfg/main.go -yaml path/to/config.yml

# 指定輸出目錄
go run cmd/easycfg/main.go -yaml path/to/config.yml -output myconfig

# 指定包名
go run cmd/easycfg/main.go -yaml path/to/config.yml -package myconfig

# 監控配置文件變更
go run cmd/easycfg/main.go -yaml path/to/config.yml -watch
```

### 在程序中使用生成的配置

```go
package main

import (
    "fmt"
    "log"

    "github.com/chiayu0816/easycfg"
)

func main() {
    // 創建配置結構體實例
    cfg := &MyConfig{}

    // 加載配置
    if err := easycfg.LoadConfig("config.yml", cfg); err != nil {
        log.Fatalf("加載配置失敗: %v", err)
    }

    // 使用配置
    fmt.Printf("配置值: %s\n", cfg.SomeField)

    // 監控配置變更
    easycfg.WatchConfig("config.yml", cfg, func() {
        fmt.Println("配置已更新")
    })
}
```

## 示例

查看 `examples/complete` 目錄中的完整示例。

運行示例：

```bash
# 運行基本示例
make run-example

# 運行完整示例
make run-complete-example
```

## 測試

EasyCfg 包含單元測試和基準測試，以確保代碼正確性和性能。

運行測試：

```bash
# 運行所有測試
make test

# 運行詳細測試（帶覆蓋率）
make test-verbose

# 運行基準測試
make benchmark
```

## 依賴

- [github.com/go-yaml/yaml](https://github.com/go-yaml/yaml)
- [github.com/spf13/viper](https://github.com/spf13/viper)

## 許可證

MIT 