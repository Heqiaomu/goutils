-- type = remote

local json = require("json")
local log = require("log")
local ssh = require("ssh")
local virtualbox_sdk = require("drivers.virtualboxv1_1_9.sdk")

local function read(f)
    local file = io.open(f, "r")
    local s = file:read("*a")
    file:close()
    return s
end

local M = {
    Name = "virtualbox",
    Version = "v1_1_9",

    userName = "",
    password = "",
    name = "",
    form_sdk_zone = "1号机",
}

function M.getInputJson(req)

    local m = json.decode(read("drivers/virtualboxv1_1_9/uiform/" .. req.actionID .. "/input.json"))

    if req.actionID == "create-hosts" then
        for _, y in pairs(m.subInputs) do
            if y.id == "virtualboxv1_1_9/host/Input" then
                for i, v in pairs(y.fields) do
                    if v.id == "virtualboxv1_1_9/host/form_sdk_zone" then
                        y.fields[i].validate.options = {}
                        for k, host in pairs(virtualbox_sdk.getVirtualboxHost()) do
                            local unusedCount = virtualbox_sdk.qryResource(host).data.unusedCount
                            print("剩余资源数量:" .. unusedCount)
                            local a = { id = k, text = k .. " (剩余资源数量:" .. unusedCount .. ")" }
                            table.insert(y.fields[i].validate.options, a)
                        end
                    end
                end
            end
        end
    end

    local response = {}
    response.code = 200
    response.message = "getInputJson ok"
    response.data = m

    return response
end

-- M.init() called by golang
function M.form()
    local form = {
        form_values = {
            {
                name = "userName",
                description = "用户名",
                required = true,
            },
            {
                name = "password",
                description = "密码",
                required = true,
            },
            {
                name = "name",
                description = "虚拟机名",
                required = true,
            },
            {
                name = "form_sdk_zone",
                description = "virtualbox 区域",
                required = true,
            },
        }
    }
    -- return form definition
    return form
end

function M.create(req)
    print(req)

    local response = {}
    response.code = 200
    response.message = "create ok"

    local FullDataFormSdkZone = {
        Key = "form_sdk_zone",
        KeyText = "区域",
        Value = M.form_sdk_zone,
    }
    --local FullDataUserName = {
    --    Key = "userName",
    --    KeyText = "用户名",
    --    Value = M.userName,
    --}
    --local FullDataPassword = {
    --    Key = "password",
    --    KeyText = "密码",
    --    Value = M.password,
    --}

    -- 驱动中的全量数据
    local FullData = {
        Items = {
            FullDataFormSdkZone,
            --FullDataUserName,
            --FullDataPassword
        }
    }

    response.data = {
        State = "3",
        Message = M.name .. " virtualbox is creating in zone of " .. M.form_sdk_zone,
        FullData = json.encode(FullData),
        Name = M.name
    }
    return response
end

function M.init(req)
    local response = {}

    local form = {
        req = req,
        zone = M.form_sdk_zone,
        name = M.name,
        userName = M.userName,
        password = M.password
    }

    virtualbox_sdk.init(form)

    local resData, err = virtualbox_sdk.create(form)
    if err ~= nil then
        response.code = 500
        response.message = err
        return response
    end

    response.code = 200
    response.data = resData
    response.message = "init ok"
    return response
end

function M.stop(req)
    log.debug("lua stop")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.stop(form)

end

function M.remove(req)
    log.debug("lua remove")

    local form = {
        req = req,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.remove(form)

end

function M.inspect(req)
    log.debug("lua inspect")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.inspect(form)
end

function M.list(req)
    log.debug("lua list")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.list(form)
end

function M.restart(req)
    log.debug("lua restart")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.restart(form)

end

function M.start(req)
    log.debug("lua start")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    return virtualbox_sdk.start(form)

end

function M.shell(req)
    log.debug("lua ssh")
    local resp = {}

    log.debug("----------------------------")
    log.debug(req)

    local username, password, serverip, port, shell = "docker", "Abcd1234", req.ip, req.port, req.shell

    log.debug(username, password, serverip, port, true, shell)

    local sshres, ssherr = ssh.connectAndCommand(username, password, serverip, port, true, shell, 30)

    log.debug(sshres)
    log.debug(ssherr)
    log.debug("----------------------------")

    resp.code = 200
    if ssherr then
        resp.code = 500
    end
    resp.data = sshres
    resp.message = ssherr

    return resp
end

function M.update(req)
    log.debug("lua update")

    local form = {
        req = req,
        zone = M.form_sdk_zone,
    }

    virtualbox_sdk.init(form)

    local resData, err = virtualbox_sdk.update(form)

    local response = {}
    if err ~= nil then
        response.code = 500
        response.message = err
        return response
    end

    response.code = 200
    response.data = resData
    response.message = "update ok"
    return response
end

function M.check(req)
    log.debug("lua check")

    local form = {
        req = req,
    }

    log.info(form)

    virtualbox_sdk.init(form)

    local response = {}
    if req.Credential.userName ~= "root" or req.Credential.password ~= "123456" then
        response.code = 500
        response.data = {}
        response.message = "凭证错误"
        return response
    end

    local username, password, serverip, port, ports = "docker", "Abcd1234", req.ip, req.port

    log.debug(username, password, serverip, port, true, shell)

    local shell = "netstat -an | grep ':22 '"
    local sshres, ssherr = ssh.connectAndCommand(username, password, serverip, port, true, shell, 30)
    log.debug(sshres)

    response.code = 200
    response.data = {}
    response.message = "check ok"

    return response
end

return M
