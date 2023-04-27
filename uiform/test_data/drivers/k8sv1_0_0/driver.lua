local json = require("json")
local log = require("log")
local sdk = require("k8s")

local M = {
    Name = "k8s",
    Version = "1.0.0",
    Type = "k8s",

    Memery = "",
    CPU = 0,
    KubeConfig = "",
    Namespace = "",
    AutoCreatePV = false,
    ApplyStorageClass = false,
    ApplyPVTemplate = false,
    StorageClassName = "",
    PVTemplate = "",
    PVNameTemplate = "",
}


-----------------------------------
--- form input
-----------------------------------

function M.form()
    local form = {
        form_values = {
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
                description = "当选择不自动创建 PV 或者选择了 PVTemplate 创建 PV，需要提供 PV 名字的模板",
                required = true,
            },
        }
    }
    return form
end

-----------------------------------
--- create change input to machine
--- return res
-----------------------------------
function M.create()
    local k8sInfo = {}
    -- 构造返回值
    k8sInfo.name = M.Name
    k8sInfo.memery = M.Memery
    k8sInfo.cpu = M.CPU
    k8sInfo.extra = json.encode(M) --all server info

    log.debug("server install done.")

    local resp = {}
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
    local msg, err = sdk.checkEnv(M.KubeConfig)
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
    local yaml, err = sdk.templatingYaml(agentTmp, M)
    if err ~= nil then
        resp.code = 500
        resp.data = ""
        resp.message = err
        return resp
    end
    log.debug(yaml)
    local agentUuid, err = sdk.getInfo(yaml, M.KubeConfig)
    if err ~= nil then
        resp.code = 500
        resp.data = ""
        resp.message = err
    end

    if agentUuid == "" or agentUuid == nil then
        -- 创建agent
        agentUuid = req.agents
        log.debug(agentUuid)
        if agentUuid == "" or agentUuid == nil then
            resp.code = 500
            resp.data = ""
            resp.message = "agents param missing, for agent uuid"
            return resp
        end
        local labels = {}
        labels.agentUuid = agentUuid
        local yaml, err = sdk.addLabel(yaml, labels)
        if err ~= nil then
            log.debug(err)
            resp.code = 500
            resp.data = ""
            resp.message = err
            return resp
        end
        local msg, err = sdk.applyByYaml(M.Namespace, M.KubeConfig, yaml)
        if err ~= nil then
            log.debug(err)
            resp.code = 500
            resp.data = ""
            resp.message = err
            return resp
        end
    end

    local serverInfo = {}
    -- 构造返回值
    serverInfo.Name = "fabric-agent-" + agentUuid
    serverInfo.agentUuid = agentUuid
    serverInfo.extra = json.encode(M) --all server info
    resp.code = 200
    resp.data = serverInfo
    resp.message = "ok"
    return resp
end

-----------------------------------
--- init 将传输的参数填入模板，并调用 k8s-client-go 的创建 agent 的 deployment
-----------------------------------
function M.init(req)
   local resp = {}
   local msg, err = sdk.checkEnv(M.KubeConfig)
   if err ~= nil then
       resp.code = 500
       resp.data = msg
       resp.message = err
       return resp
   end
   resp.code = 200
   resp.data = ""
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
    local res, err = sdk.checkRunning(M.kubeConfig, yaml)
    if err ~= nil then
        log.debug(err)
        resp.code = 500
        resp.data = ""
        resp.message = err
        return resp
    end
    resp.code = 200
    resp.data = res
    resp.message = ""
    return resp
end



-----------------------------------
--- remove 根据传入的参数进行删除 资源
-----------------------------------

function M.remove(req)
    local resp = {}
    local namespaceVar = req.Namespace
    if namespaceVar ~= M.Namespace then
        resp.code = 500
        resp.data = nil
        resp.message = "namespace not right"
    end
    local yaml = req.yaml
    local err = sdk.deleteByYaml(M.Namespace, M.KubeConfig, yaml)
    if err ~= nil then
        resp.code = 500
        resp.data = nil
        resp.message = "delete resource error"
    end
    resp.code = 200
    resp.data = ""
    resp.message = "ok"
    return resp
end

return M