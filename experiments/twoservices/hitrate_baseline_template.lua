local JSON = require("JSON")

wrk.method = "POST"

param = {
    k = 1,
}
wrk.body   = JSON:encode(param)
wrk.headers["Content-Type"] = "application/json"
