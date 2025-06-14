-- 定义请求参数（JSON 格式）
local request_body = [[
{
    "concert_id": 13,
    "user_id": 652676240481718272,
    "seat_idx": {
        "section": "B区"
    }
}
]]

-- 定义请求函数
function request()
    local path = "/api/v1/buy"
    local method = "POST"
    local token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2NTI2NzYyNDA0ODE3MTgyNzIsInVzZXJfbmFtZSI6ImhhaGEzIiwiaXNzIjoiY29uY2VydCIsImV4cCI6MTc1MDc1NzQ5OH0.fqOPqkfgYev5z6nlxKT90XHOlhRoJbJ0wgM3W_IQnSs"

    local headers = {
        ["Content-Type"] = "application/json",
        ["Content-Length"] = tostring(#request_body),
        ["Authorization"] = "Bearer " .. token
    }

    return wrk.format(method, path, headers, request_body)
end
