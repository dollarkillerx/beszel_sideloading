import redis
import json
import time

# 连接 Redis（默认 127.0.0.1:6379）
r = redis.Redis(host='100.120.27.93', port=6779, db=0, password='Kw7dQwFnLGAe2')

# 构造你要插入的对象
data = {
    "name": "移动联通深港IEPL11-X-02",
    "id": 383,
    "type": "trojan",
    "online": 120,
    "last_update": int(time.time())  # 当前 Unix 时间戳
}

# 将对象编码成 JSON 字符串
json_str = json.dumps(data, ensure_ascii=True)

# 插入 Redis，key 可以自定义
key = "v2board_database_AGENT_移动联通深港IEPL11-X-02"
r.set(key, json_str)

print(f"成功写入 Redis：{key} -> {json_str}")
