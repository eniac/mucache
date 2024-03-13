local socket = require("socket")
local JSON = require("JSON")
local UUID = require("uuid")
local time = socket.gettime()*1000
math.randomseed(time)
math.random(); math.random(); math.random()



-- load env vars
-- local max_user_index = tonumber(os.getenv("max_user_index")) or 962

-- local post_ids = {}
-- local number_of_post_ids = 0

local max_user_index = 0
local total_followers = 0

local followers = {}

function read_social_graph_analysis_file(path)
    local file = io.open(path, "r");
    local content = file:read "*a" -- *a or *all reads the whole file
    file:close()

    -- Split in lines
    local lines = {}
    for s in content:gmatch("[^\r\n]+") do
      table.insert(lines, s)
    end


    tmp, max_user_index_str = lines[1]:match("(Nodes:) (%d+)")
    tmp, total_followers_str = lines[2]:match("(Total followers:) (%d+)")
    
    max_user_index = tonumber(max_user_index_str)
    total_followers = tonumber(total_followers_str)

    for i = 1, max_user_index do table.insert(followers, 0) end

    for i = 4, #lines do
      -- print(lines[i])
      user_id, number_of_followers = lines[i]:match("(%d+) (%d+)")
      -- print(user_id, number_of_followers)
      followers[tonumber(user_id)] = tonumber(number_of_followers)
    end

    -- print(max_user_index)
    -- print(total_followers)
    -- print(followers[147], followers[204], followers[962], followers[1])
end

-- TODO: Read from somewhere
read_social_graph_analysis_file("experiments/social/socfb/socfb-analysis.txt")


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

local decset = {'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

function decRandom(length)
    if length > 0 then
        return decRandom(length - 1) .. decset[math.random(1, #decset)]
    else
        return ""
    end
end


local function get_random_user_index()
  return math.random(1, max_user_index)
end

local function compose_post()
  local user_index = get_random_user_index()
  -- local username = "username_" .. tostring(user_index)
  local user_id = "User" .. tostring(user_index)
  local text = string.random(256)

  --
  -- For now we just do random text (no mentions/urls/etc)
  --
  -- local num_user_mentions = math.random(0, 5)
  -- local num_urls = math.random(0, 5)
  -- local num_media = math.random(0, 4)
  -- local media_ids = '['
  -- local media_types = '['

  -- for i = 0, num_user_mentions, 1 do
  --   local user_mention_id
  --   while (true) do
  --     user_mention_id = math.random(0, max_user_index - 1)
  --     if user_index ~= user_mention_id then
  --       break
  --     end
  --   end
  --   text = text .. " @username_" .. tostring(user_mention_id)
  -- end

  -- for i = 0, num_urls, 1 do
  --   text = text .. " http://" .. string.random(64)
  -- end

  -- for i = 0, num_media, 1 do
  --   local media_id = decRandom(18)
  --   media_ids = media_ids .. "\"" .. media_id .. "\","
  --   media_types = media_types .. "\"png\","
  -- end

  -- media_ids = media_ids:sub(1, #media_ids - 1) .. "]"
  -- media_types = media_types:sub(1, #media_types - 1) .. "]"

  local method = "POST"
  local headers = {}
  local param = {
      creator_id = user_id,
      text = text
  }
  local body = JSON:encode(param)
  headers["Content-Type"] = "application/json"
  
  -- TODO: Can we automate this?
  local path = "http://localhost:8084/compose_post"
  
  return wrk.format(method, path, headers, body)
end

local function read_user_timeline()
  local user_id = tostring(get_random_user_index())
  local start = tostring(math.random(0, 100))
  local stop = tostring(start + 10)

  local args = "user_id=" .. user_id .. "&start=" .. start .. "&stop=" .. stop
  local method = "GET"
  local headers = {}
  headers["Content-Type"] = "application/x-www-form-urlencoded"
  local path = "http://localhost:8080/wrk2-api/user-timeline/read?" .. args
  return wrk.format(method, path, headers, nil)
end

local function read_home_timeline()
    local user_id = tostring(get_random_user_index())
    local start = tostring(math.random(0, 100))
    local stop = tostring(start + 10)

    local args = "user_id=" .. user_id .. "&start=" .. start .. "&stop=" .. stop
    local method = "GET"
    local headers = {}
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    local path = "http://localhost:8080/wrk2-api/home-timeline/read?" .. args
    return wrk.format(method, path, headers, nil)
  end

-- TODO: Automate paths and ports.
-- TODO: Fix the body to be JSON
-- TODO: Make a populate function that reads is given a social graph
--       (which is retrieved from some real dataset).
-- TODO: Preprocess the graph and add how followed each user is. We can then
--       use this file to determine how many user_timeline requests each user
--       will receive.
-- TODO: Read user timeline should be a skewed and not uniform distribution
--       Most followed should be seen more.
-- TODO: Determine what should the read_home_timeline distribution for users be
-- TODO: Determine what should the compose_post distribution be among users
-- TODO: We want to get a different latency distribution for each type of request

-- TODO: Consider pregenerating requests in init and only looking them up for better wrk performance.

request = function()
  cur_time = math.floor(socket.gettime())
  local read_home_timeline_ratio = 0.60
  local read_user_timeline_ratio = 0.30
  local compose_post_ratio       = 0.10

  -- Just for debugging
  local read_home_timeline_ratio = 0.0
  local read_user_timeline_ratio = 0.0
  local compose_post_ratio       = 1.0

  local coin = math.random()
  if coin < read_home_timeline_ratio then
    return read_home_timeline()
  elseif coin < read_home_timeline_ratio + read_user_timeline_ratio then
    return read_user_timeline()
  else
    return compose_post()
  end
end


--
-- Only uncomment the following when debugging
--


-- max_requests = 0
-- counter = 1

-- function setup(thread)
--     thread:set("id", counter)
    
--     counter = counter + 1
-- end

-- response = function (status, headers, body)
--   io.write("------------------------------\n")
--   io.write("Response ".. counter .." with status: ".. status .." on thread ".. id .."\n")
--   io.write("------------------------------\n")

--   io.write("[response] Headers:\n")

--   -- Loop through passed arguments
--   for key, value in pairs(headers) do
--     io.write("[response]  - " .. key  .. ": " .. value .. "\n")
--   end

--   io.write("[response] Body:\n")
--   io.write(body .. "\n")

--   -- Stop after max_requests if max_requests is a positive number
--   if (max_requests > 0) and (counter > max_requests) then
--     wrk.thread:stop()
--   end
  
--   counter = counter + 1
-- end