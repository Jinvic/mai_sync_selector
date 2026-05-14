# mai_sync_selector

舞萌 DX（MaiMai DX）歌曲筛选工具。

## 功能

- 从 [diving-fish.com](https://www.diving-fish.com) API 同步舞萌DX歌曲数据
- 支持 1P/2P 双玩家筛选，可按版本、曲包、定数、等级筛选歌曲
- 定时自动同步（每周一凌晨3点）
- 响应式 Web 界面，支持分页、列显示控制

## 技术栈

- **后端**: Go + Gin
- **数据库**: SQLite + GORM
- **前端**: 原生 HTML/CSS/JS

## 快速启动

```bash
go build -o mai_sync_selector .
./mai_sync_selector
```

访问 http://localhost:8080

## Docker 部署

```bash
docker-compose up -d
```

访问 http://localhost:8080

数据存储在 `./data` 目录。

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_PATH` | `sqlite.db` | 数据库路径 |
| `SERVER_PORT` | `8080` | 服务端口 |

## API

| 接口 | 方法 | 说明 |
|------|------|------|
| `/sync-data` | GET | 手动同步歌曲数据 |
| `/select-song` | POST | 筛选歌曲 |
| `/from` | GET | 获取版本列表 |
| `/genre` | GET | 获取曲包列表 |
| `/level` | GET | 获取等级列表 |

### 筛选请求示例

```json
{
  "filter1": {
    "from_list": [],
    "genre_list": ["niconico & VOCALOID"],
    "min_ds": 1.0,
    "max_ds": 15.9,
    "min_level": "11+",
    "max_level": "12+"
  },
  "filter2": {
    "from_list": ["maimai でらっくす PRiSM"],
    "genre_list": [],
    "min_ds": 12.9,
    "max_ds": 13.9,
    "min_level": "",
    "max_level": ""
  },
  "page": 1,
  "page_size": 20
}
```

## 数据来源

歌曲数据来自 [diving-fish/maimaidxprober](https://github.com/diving-fish/maimaidxprober) API。