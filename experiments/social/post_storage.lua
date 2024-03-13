local socket = require("socket")
local JSON = require("JSON")
local UUID = require("uuid")
local utility = require('utility')
time = socket.gettime() * 1000
UUID.randomseed(time)
math.randomseed(time)
math.random();
math.random();
math.random()


local post_ids = {}
local number_of_post_ids = 0

function read_post_id_file(path)
    local file = io.open(path, "r");
    for line in file:lines() do
       table.insert (post_ids, line);
    end
    number_of_post_ids = table.getn(post_ids)
end

-- TODO: Read from somewhere
read_post_id_file("post_ids.txt")

local function uuid()
    return UUID():gsub('-', '')
end

request = function()
    
    -- TODO: We want some actual distribution here!
    local post_id_index = math.random(1,number_of_post_ids)
    local post_id = post_ids[post_id_index]

    local method = "POST"
    local headers = {}
    local param = {
        post_id = post_id
    }
    local body = JSON:encode(param)
    headers["Content-Type"] = "application/json"

    return wrk.format(method, path, headers, body)
end