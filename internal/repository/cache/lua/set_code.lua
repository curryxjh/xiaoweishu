-- 发送到的key, 也就是 phone_code:biz:phone
-- example: code:login:159xxxxxxxx
local key = KEYS[1]
-- 验证次数, 一个验证码最多重复三次, 超过三次就不再验证, 记录还可以验证几次
-- example: code:login:159xxxxxxxx:cnt
local countKey = key..":cnt"
-- 验证码
-- example: 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在, 但是没有过期时间
   return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", countKey, 3)
    redis.call("expire", countKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end

