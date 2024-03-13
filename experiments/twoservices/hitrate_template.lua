local JSON = require("JSON")

wrk.method = "POST"

param = {
    k = 1,
    hit_rate = ${HITRATE},
}
wrk.body   = JSON:encode(param)
wrk.headers["Content-Type"] = "application/json"
