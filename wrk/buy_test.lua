-- 定义请求参数（JSON 格式）
local request_body = [[
{
    "concert_id": 8,
    "user_id": 650609386997157888,
    "seat_idx": {
        "section": "A区"
    }
}
]]

-- 定义请求函数
function request()
    local path = "/api/v1/buy"
    local method = "POST"
    local token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2NTA2MDkzODY5OTcxNTc4ODgsInVzZXJfbmFtZSI6ImhhaGEiLCJpc3MiOiJjb25jZXJ0IiwiZXhwIjoxNzQ5NTQ2MTMxfQ.JjihzYT-WntrNOmu3QZ7LdoyvEeI_zH2SuJZWBXdpF0"

    local headers = {
        ["Content-Type"] = "application/json",
        ["Content-Length"] = tostring(#request_body),
        ["Authorization"] = "Bearer " .. token
    }

    return wrk.format(method, path, headers, request_body)
end
