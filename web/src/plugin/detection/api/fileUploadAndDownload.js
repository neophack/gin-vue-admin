import service from '@/utils/request'
// @Tags detection
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "分页获取文件户列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /detection/getFileList [post]
export const getFileList = (data) => {
  return service({
    url: '/detection/getFileList',
    method: 'post',
    data
  })
}

// @Tags detection
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body dbModel.detection true "传入文件里面id即可"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /detection/deleteFile [post]
export const deleteFile = (data) => {
  return service({
    url: '/detection/deleteFile',
    method: 'post',
    data
  })
}

/**
 * 编辑文件名或者备注
 * @param data
 * @returns {*}
 */
export const editFileName = (data) => {
  return service({
    url: '/detection/editFileName',
    method: 'post',
    data
  })
}

export const getRouterName = (data) => {
  return service({
    url: '/detection/routerName',
    method: 'post',
    data
  })
}

export const getBatchList = (data) => {
  return service({
    url: '/detection/getBatchList',
    method: 'post',
    data
  })
}

export const deleteBatch = (data) => {
  return service({
    url: '/detection/deleteBatch',
    method: 'post',
    data
  })
}

export const editBatch = (data) => {
  return service({
    url: '/detection/editBatch',
    method: 'post',
    data
  })
}

export const newBatch = (data) => {
  return service({
    url: '/detection/newBatch',
    method: 'post',
    data
  })
}

export const changeStatus = (data) => {
  return service({
    url: '/detection/changeStatus',
    method: 'post',
    data
  })
}

export const downloadFilesZip = (data) => {
  return service({
    url: '/detection/downloadFilesZip',
    method: 'post',
    data
  })
}