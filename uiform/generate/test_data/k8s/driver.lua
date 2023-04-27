local json = require("json")
local log = require("log")
local sdk = require("k8s")

local function rand()
    return math.random(1000, 9999)
end

local function genRandUUID()
    local prefix = string.sub(tostring(os.time()), 7, 11)
    return "k8s-" .. prefix .. "-" .. rand() .. "-" .. rand()
end

local M = {
    Name = "k8s",
    Version = "v1_0_7",
    Type = "k8s",
    name = "k8s",
    KubeConfig = "",

    EnvironmentScope = "",
    ApiURL = "",
    CACertificate = "",
    Token = "",

    Namespace = "",
    AutoCreatePV = false,
    ApplyStorageClass = false,
    ApplyPVTemplate = false,
    StorageClassName = "",
    PVTemplate = "",
    PVNameTemplate = "",
    ApplyPV = ""
}

-----------------------------------
--- form input
-----------------------------------
function M.form()
    local form = {
        form_values = {
            {
                name = "ApplyPV",
                description = "k8s PV创建模式",
                required = true,
            },
            {
                name = "name",
                description = "k8s 主机名称",
                required = true,
            },
            {
                name = "Namespace",
                description = "k8s namesapce",
                required = true,
            },
            {
                name = "KubeConfig",
                description = "k8s 连接配置",
                required = true,
            },
            {
                name = "EnvironmentScope",
                description = "k8s 环境域",
                required = true,
            },
            {
                name = "ApiURL",
                description = "k8s API URL",
                required = true,
            },
            {
                name = "CACertificate",
                description = "k8s CA 证书",
                required = true,
            },
            {
                name = "Token",
                description = "k8s 访问令牌",
                required = true,
            },
            {
                name = "AutoCreatePV",
                description = "k8s 资源自动创建 PV",
                required = true,
            },
            {
                name = "ApplyStorageClass",
                description = "自动创建 PV 是否使用 StorageClass",
                required = true,
            },
            {
                name = "ApplyPVTemplate",
                description = "自动创建 PV 是否使用 PVTemplate",
                required = true,
            },
            {
                name = "StorageClassName",
                description = "自动创建 PV 使用的 StorageClass 名称",
                required = true,
            },
            {
                name = "PVTemplate",
                description = "自动创建 PV 使用 PVTemplate，具体模板",
                required = true,
            },
            {
                name = "PVNameTemplate",
                description = "选择不自动创建PV或PVTemplate创建PV，提供PV名字的模板",
                required = true,
            },
        }
    }
    return form
end

local function read(f)
    local file = io.open(f, "r")
    local s = file:read("*a")
    file:close()
    return s
end

function M.getInputJson(req)
    local m = json.decode(read("drivers/k8sv1_0_7/uiform/" .. req.actionID .. "/input.json"))
    log.debug(m)
    local response = {}
    response.code = 200
    response.message = "getInputJson ok"
    response.data = m
    return response
end

-----------------------------------
--- create change input to machine
--- return res
-----------------------------------
function M.create()
    local resp = {}
    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end

    M.PVNameTemplate = "{{chainname}}-{{org-name}}-" .. M.PVNameTemplate
    local k8sInfo = {}
    -- 构造返回值
    k8sInfo.Name = M.name

    local info, err = sdk.getNodeIP(M.ApiURL, M.Token, M.CACertificate)
    if err ~= nil then
        resp.code = 500
        resp.data = k8sInfo
        resp.message = err
        return resp
    end

    k8sInfo.HostIP = info.HostIP
    log.debug("cpu info", info.CPU)
    k8sInfo.CPU = tostring(info.CPU) .. "核"
    log.debug("memory info", info.Memory / 1024 / 1024 / 1024)
    k8sInfo.Memory = tostring(math.ceil(info.Memory / 1024 / 1024 / 1024)) .. "G"
    k8sInfo.Disk = "20G"
    k8sInfo.GitVersion = info.GitVersion
    k8sInfo.GoVersion = info.GOVersion
    k8sInfo.Platform = info.Platform
    -- 驱动中的全量数据
    local FullData = {
        Items = {}
    }
    table.insert(FullData.Items, { Key = "Namespace", KeyText = "命名空间", Value = M.Namespace, })
    table.insert(FullData.Items, { Key = "GitVersion", KeyText = "kubernetes 版本", Value = k8sInfo.GitVersion, })
    table.insert(FullData.Items, { Key = "GoVersion",
                                   KeyText = "k8s 构建 Golang 版本",
                                   Value = k8sInfo.GoVersion, })
    table.insert(FullData.Items, { Key = "Platform", KeyText = "系统平台", Value = k8sInfo.Platform, })
    --table.insert(FullData.Items, { Key = "CPU", KeyText = "CPU", Value = k8sInfo.CPU, })
    --table.insert(FullData.Items, { Key = "Memory", KeyText = "内存", Value = k8sInfo.Memory, })
    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            table.insert(FullData.Items, {
                Key = "StorageClassName",
                KeyText = "StorageClass名称",
                Value = M.StorageClassName,
            })
        elseif apv ~= "ApplyPVTemplate" then
            table.insert(FullData.Items, {
                Key = "StorageClassName",
                KeyText = "StorageClass名称",
                Value = "-",
            })
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end

    --table.insert(FullData.Items, { Key = "PVNameTemplate", KeyText = "PV名字模板", Value = M.PVNameTemplate, })

    k8sInfo.FullData = json.encode(FullData)

    --k8sInfo.Data = json.encode(M) --all server info
    k8sInfo.State = "3"
    k8sInfo.Message = M.name .. " kube is creating with namespace " .. M.Namespace

    log.debug("server install done.")

    resp.code = 200
    resp.data = k8sInfo
    resp.message = "ok"
    return resp
end

-----------------------------------
--- shell 使用配置文件生成 k8s client
--- return res
-----------------------------------
function M.shell(req)
    local resp = {}
    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end
    M.PVNameTemplate = "{{chainname}}-{{org-name}}-" .. M.PVNameTemplate
    local msg, err = sdk.checkEnv(M.ApiURL, M.Token, M.CACertificate)
    if err ~= nil then
        resp.code = 500
        resp.data = msg
        resp.message = err
        return resp
    end

    local agentTmp = req.agentTempYaml
    if agentTmp == "" or agentTmp == nil then
        resp.code = 500
        resp.data = ""
        resp.message = err
        return resp
    end
    local yaml, err2 = sdk.templatingYaml(agentTmp, M)
    if err2 ~= nil then
        resp.code = 500
        resp.data = ""
        resp.message = err2
        return resp
    end
    log.debug("template yaml without agentuuid", yaml)
    local agentUuid, err3 = sdk.getInfo(yaml, M.ApiURL, M.Token, M.CACertificate)
    if err3 ~= nil then
        resp.code = 500
        resp.data = ""
        resp.message = err3
    end

    log.debug("get agent uuid from k8s", agentUuid)

    if agentUuid == "" or agentUuid == nil then
        -- 创建agent
        agentUuid = req.agents
        log.debug("get agent uuid from req", agentUuid)
        if agentUuid == "" or agentUuid == nil then
            resp.code = 500
            resp.data = ""
            resp.message = "agents param missing, for agent uuid"
            return resp
        end
        --local labels = {}
        M.AgentUUID = agentUuid
        M.EnvType = 2

        local yaml2, err4 = sdk.templatingYaml(agentTmp, M)
        if err4 ~= nil then
            log.err("template yaml error", err)
            resp.code = 500
            resp.data = ""
            resp.message = err4
            return resp
        end
        log.debug("template yaml with agent uuid", yaml2)

        --local yaml, err = sdk.addLabel(yaml, labels)
        --if err ~= nil then
        --    log.debug(err)
        --    resp.code = 500
        --    resp.data = ""
        --    resp.message = err
        --    return resp
        --end
        local _, err5 = sdk.applyByYaml(M.Namespace, yaml2, M.ApiURL, M.Token, M.CACertificate)
        if err5 ~= nil then
            log.err("apply yaml error", err5)
            resp.code = 500
            resp.data = ""
            resp.message = err5
            return resp
        end
    end

    local serverInfo = {}
    -- 构造返回值
    serverInfo.Name = string.format("agent-%s", agentUuid)
    serverInfo.agentUuid = agentUuid
    serverInfo.extra = json.encode(M) --all server info
    resp.code = 200
    resp.data = json.encode(serverInfo)
    resp.message = ""
    return resp
end

-----------------------------------
--- init 将传输的参数填入模板，并调用 k8s-client-go 的创建 agent 的 deployment
-----------------------------------
function M.init(req)
    print(req)
    local resp = {}
    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end
    M.PVNameTemplate = "{{chainname}}-{{org-name}}-" .. M.PVNameTemplate
    local msg, err = sdk.checkEnv(M.ApiURL, M.Token, M.CACertificate)
    if err ~= nil then
        resp.code = 500
        resp.data = msg
        resp.message = err
        return resp
    end

    local k8sInfo = {
        UUID = genRandUUID(),
        State = "2",
        Message = M.name .. " kube is available in namespace " .. M.Namespace
    }

    resp.code = 200
    resp.data = k8sInfo
    resp.message = "ok"
    return resp
end

-----------------------------------
--- check 将传输的参数填入模板，并调用 k8s-client-go 的创建 agent 的 deployment
-----------------------------------
function M.check(req)
    -- check agent running
    log.debug(req)
    local resp = {}
    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end
    M.PVNameTemplate = "{{chainname}}-{{org-name}}-" .. M.PVNameTemplate
    local agentTmp = req.agentTempYaml
    -- todo add M's fields to agent env
    local yaml, err = sdk.templatingYaml(agentTmp, M)
    if err ~= nil then
        resp.code = 500
        resp.data = ""
        resp.message = err
        return resp
    end
    log.debug(yaml)

    -- check
    local res, err2 = sdk.checkRunning(M.ApiURL, M.Token, M.CACertificate, yaml)
    if err2 ~= nil then
        log.debug(err2)
        resp.code = 500
        resp.data = ""
        resp.message = err2
        return resp
    end
    resp.code = 200
    resp.data = res
    resp.message = ""
    return resp
end


-----------------------------------
--- restart 根据传入的参数进行 重启 资源
-----------------------------------
function M.restart(req)
    print("start restart machine")
    print(json.encode(req))
    print(json.encode(M))

    local resp = {}

    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end
    resp.code = 200
    resp.data = ""
    resp.message = ""
    return resp
end

-----------------------------------
--- remove 根据传入的参数进行删除 资源
-----------------------------------
function M.remove(req)
    print("start remove machine")
    print(json.encode(req))
    print(json.encode(M))
    local resp = {}

    if M.AutoCreatePV == "true" then
        local apv = M.ApplyPV
        if apv == "ApplyStorageClass" then
            M.ApplyStorageClass = true
            M.ApplyPVTemplate = false
        elseif apv == "ApplyPVTemplate" then
            M.ApplyPVTemplate = true
            M.ApplyStorageClass = false
        else
            resp.code = 500
            resp.data = {}
            resp.message = "got no ApplyPV"
            return resp
        end
    end

    local m = req.Machine
    local label = {
        machine = string.format("m%d", m.ID)
    }
    local res, err = sdk.deleteByLabels(M.ApiURL, M.Token, M.CACertificate, M.Namespace, "deployment", label)
    log.debug("delete deployment res : ", res)
    if err ~= nil then
        resp.code = 500
        resp.data = {}
        resp.message = string.format("delete deployment of machine : m%d, error : %s", m.ID, err)
        return resp
    end
    log.debug(M.ApiURL, M.Token, M.CACertificate, M.Namespace, "statefulSet", label)
    res, err = sdk.deleteByLabels(M.ApiURL, M.Token, M.CACertificate, M.Namespace, "statefulSet", label)
    log.debug("delete statefulSet res : ", res)
    if err ~= nil then
        resp.code = 500
        resp.data = {}
        resp.message = string.format("delete statefulSet of machine : m%d, error : %s", m.ID, err)
        return resp
    end
    res, err = sdk.deleteByLabels(M.ApiURL, M.Token, M.CACertificate, M.Namespace, "service", label)
    log.debug("delete service res : ", res)
    if err ~= nil then
        resp.code = 500
        resp.data = {}
        resp.message = string.format("delete service of machine : m%d, error : %s", m.ID, err)
        return resp
    end
    res, err = sdk.deleteByLabels(M.ApiURL, M.Token, M.CACertificate, M.Namespace, "pvc", label)
    log.debug("delete pvc res : ", res)
    if err ~= nil then
        resp.code = 500
        resp.data = {}
        resp.message = string.format("delete pvc of machine : m%d, error : %s", m.ID, err)
        return resp
    end

    --
    --if M.AutoCreatePV == "true" then
    --    local apv = M.ApplyPV
    --    if apv == "ApplyStorageClass" then
    --        M.ApplyStorageClass = true
    --        M.ApplyPVTemplate = false
    --    elseif apv == "ApplyPVTemplate" then
    --        M.ApplyPVTemplate = true
    --        M.ApplyStorageClass = false
    --    else
    --        resp.code = 500
    --        resp.data = {}
    --        resp.message = "got no ApplyPV"
    --        return resp
    --    end
    --end
    --M.PVNameTemplate = "{{chainname}}-{{org-name}}-" .. M.PVNameTemplate
    --local namespaceVar = req.Namespace
    --if namespaceVar ~= M.Namespace then
    --    resp.code = 500
    --    resp.data = nil
    --    resp.message = "namespace not right"
    --end
    --local yaml = req.yaml
    --local err = sdk.deleteByYaml(M.Namespace, M.KubeConfig, yaml)
    --if err ~= nil then
    --    resp.code = 500
    --    resp.data = nil
    --    resp.message = "delete resource error"
    --end
    resp.code = 200
    resp.data = ""
    resp.message = "ok"
    return resp
end

function M.update(req)
    print("start update machine")
    print(json.encode(req))
    local resp = {}
    -- todo something
    resp.code = 200
    resp.data = req.Machine
    resp.message = "ok"
    return resp
end

function M.check(req)
    log.info(req)
    local response = {}

    response.code = 200
    response.data = {}
    response.message = "check ok"

    return response
end

return M
