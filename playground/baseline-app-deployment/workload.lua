local socket = require("socket")
local JSON = require("JSON")
local UUID = require("uuid")
time = socket.gettime() * 1000
UUID.randomseed(time)
math.randomseed(time)
math.random();
math.random();
math.random()

local HOTEL_NUM = 100
local USER_NUM = 100

local charset = { 'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', 'a', 's',
                  'd', 'f', 'g', 'h', 'j', 'k', 'l', 'z', 'x', 'c', 'v', 'b', 'n', 'm', 'Q',
                  'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P', 'A', 'S', 'D', 'F', 'G', 'H',
                  'J', 'K', 'L', 'Z', 'X', 'C', 'V', 'B', 'N', 'M', '1', '2', '3', '4', '5',
                  '6', '7', '8', '9', '0' }

function string.random(length)
    if length > 0 then
        return string.random(length - 1) .. charset[math.random(1, #charset)]
    else
        return ""
    end
end

local function uuid()
    return UUID():gsub('-', '')
end

request = function()
    local userid = tostring(math.random(0, USER_NUM - 1))
    local username = "username_" .. tostring(user_id)
    local hotelid = tostring(math.random(0, HOTEL_NUM - 1))

    local method = "POST"
    local headers = {}
    local param = {
        user_name = username,
        user_id = userid,
        hotel_id = hotelid
    }
    local body = JSON:encode(param)
    headers["Content-Type"] = "application/json"
    -- headers["Host"] = "caller1.default.example.com"

    return wrk.format(method, path, headers, body)
end

-- {"user_name":"User Name","user_id":"2", "hotel_id":"10"}
-- ./wrk2/wrk -t1 -c1 -d20 -R1 --latency http://localhost:5002/book -s workload.lua 
