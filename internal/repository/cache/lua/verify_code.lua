local key = KYES[1]
-- 用户输入的 code
local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"
local cnt = tonumber(redis.call("get", cntKey))
if cnt == nil or c nt <= 0 then
    -- 说明用户一直输入错误
    -- 说明验证码已经使用过了
    return -1
elseif expectedCode == code then
    redis.call("set", cntKey, -1)
    return 0
else
    -- 说明用户输入错误
    -- 记录错误次数, 每次错误次数减一, 最多验证三次
    redis.call("decr", cntKey)
    return -2
end
