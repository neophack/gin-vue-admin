<template>
  <div v-loading.fullscreen.lock="">
    <div class="gva-table-box">
      <div class="upload">
        <div v-show="$refs.upload && $refs.upload.dropActive" class="drop-active">
          <h3>Drop files to upload</h3>
        </div>
        <file-upload
          ref="upload"
          v-model="files"
          class="upload-btn"
          :post-action="postAction"
          :multiple="true"
          :drop="true"
          :drop-directory="true"
          :headers="headers"
          :thread="10"
          :add-index="true"
          @input-filter="inputFilter"
        />
      </div>

      <warning-bar title="按上传时间顺序后台识别，负载高时请耐心等待" />

      <div class="gva-btn-list">
        <!--        <upload-common-->
        <!--          v-model:imageCommon="imageCommon"-->
        <!--          class="upload-btn"-->
        <!--          @on-success="getTableData"-->
        <!--        />-->
        <!--        <upload-image-->
        <!--          v-model:imageUrl="imageUrl"-->
        <!--          :file-size="512"-->
        <!--          :max-w-h="1080"-->
        <!--          class="upload-btn"-->
        <!--          @on-success="getTableData"-->
        <!--        />-->
        <el-button type="primary" @click="uploadDialog()">upload</el-button>
        <el-form ref="searchForm" :inline="true" :model="search">
          <el-form-item label="">
            <el-input
              v-model="search.keyword"
              class="keyword"
              placeholder="请输入文件名或备注"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" icon="search" @click="getTableData">查询</el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table :data="tableData">
        <!--        <el-table-column align="left" label="上传图像" width="100">-->
        <!--          <template #default="scope">-->
        <!--            <CustomPic-->
        <!--              pic-type="file"-->
        <!--              :pic-src="scope.row.url"-->

        <!--              @click="handlePictureCardPreview(scope.row.url)"-->
        <!--            />-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <!--        <el-table-column align="left" label="识别结果" width="100">-->
        <!--          <template #default="scope">-->
        <!--            <CustomPic-->
        <!--              v-if="scope.row.url_detection!=''"-->
        <!--              pic-type="file"-->
        <!--              :pic-src="scope.row.url_detection"-->
        <!--              @click="handlePictureCardPreview(scope.row.url_detection)"-->
        <!--            />-->
        <!--          </template>-->
        <!--        </el-table-column>-->
        <el-table-column align="left" label="batches" prop="name" min-width="280">
          <template #default="scope">
            <div class="name">
              {{ scope.row.batchid }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="update日期" prop="UpdatedAt" width="180">
          <template #default="scope">
            <div>{{ formatDate(scope.row.UpdatedAt) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="count" prop="count" min-width="80">
          <template #default="scope">
            <div>
              {{ scope.row.files_count }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="progress" prop="progress" min-width="180">
          <template #default="scope">
            <div>
              <el-progress :text-inside="true" :stroke-width="22" :percentage="Number(scope.row.progress)" />
            </div>
          </template>
        </el-table-column>

        <!-- <el-table-column align="left" label="链接" prop="url" min-width="300" /> -->
        <el-table-column align="left" label="size" prop="tag" width="100">
          <template #default="scope">
            <div>
              {{ scope.row.files_size }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="status" prop="status" min-width="100">
          <template #default="scope">
            <div>
              {{ scope.row.status }}
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="460">
          <template #default="scope">
            <el-button icon="view" type="primary" @click="viewFiles(scope.row,false)">view</el-button>
            <el-button icon="view" type="primary" @click="viewFiles(scope.row,true)">res</el-button>
            <el-button v-if="scope.row.status=='ready'" type="primary" @click="onChangeStatus(scope.row)">pause</el-button>
            <el-button v-else type="primary" @click="onChangeStatus(scope.row)">run-></el-button>
            <!--            <el-button icon="edit" type="primary" @click="downloadFile(scope.row)">edit</el-button>-->
            <el-button icon="download" type="primary" @click="downloadFile(scope.row)">下载</el-button>
            <el-button
              v-if="isAdmin"
              icon="delete"
              type="danger"
              @click="onDeleteBatch(scope.row)"
            >删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :style="{ float: 'right', padding: '20px' }"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-dialog v-model="dialogVisible">
      <img w-full :src="dialogImageUrl" alt="Preview Image" style="width: 100%">
    </el-dialog>
    <el-dialog v-model="uploadVisible" title="upload">

      <!--      <div class="gva-btn-list">-->
      <el-space :size="8" spacer=" ">
        <el-button @click="onAddFiles">Select files</el-button>
        <el-button @click="onAddFolder">Select folder</el-button>
        <el-button
          v-if="!$refs.upload || !$refs.upload.active"
          type="primary"
          class="btn btn-success"
          @click.prevent="onUploadFiles"
        >
          <i class="fa fa-arrow-up" aria-hidden="true" />
          Start Upload
        </el-button>
        <el-button
          v-else
          type="primary"
          class="btn btn-danger"
          disabled
          @click.prevent="$refs.upload.active = false"
        >
          <i class="fa fa-stop" aria-hidden="true" />
          Stop Upload
        </el-button>
        <el-button type="danger" @click="reset()">Reset</el-button>

      </el-space>
      <el-space :size="8" spacer=" ">
        <el-input v-model="batchid" placeholder="batchid" />

        {{ files.filter(f => f.success).length }}/{{ files.length }}
        <el-progress
          :width="28"
          type="circle"
          :text-inside="true"
          :stroke-width="5"
          :percentage="progress0 = files.filter(f=>f.success).length*100/files.length"
        />
        <span v-show="files.length>0 && $refs.upload && $refs.upload.uploaded && onUploadFinish()">上传完毕</span>
      </el-space>
      <!--      </div>-->

      <span>
        <vxe-table
          ref="xTable"
          keep-source
          :tooltip-config="{}"
          :row-config="{isCurrent: true, isHover: true}"
          :data="files"
          height="500"
          :show-header="true"
          show-overflow="false"
          :mouse-config="{selected: true}"
          :keyboard-config="{isArrow: true, isEnter: true, isChecked: true}"
        >
          <vxe-column type="seq" title="Seq" width="60" />
          <vxe-column field="name" title="Name" show-overflow="ellipsis">
            <template #default="{ row }">
              <span>{{ row.name }}</span>
            </template>
          </vxe-column>
          <vxe-column field="size" title="Size" width="100">
            <template #default="{ row }">
              <span>{{ formatSize(row.size) }}</span>
            </template>
          </vxe-column>
          <vxe-column field="progress" title="Progress" width="200">
            <template #default="{ row }">
              <el-progress
                v-if="row.success"
                :text-inside="true"
                :stroke-width="20"
                :percentage="100"
                status="success"
              />
              <el-progress v-else-if="row.active" :text-inside="true" :stroke-width="22" :percentage="Number(row.progress)" />
              <el-progress
                v-else-if="row.error"
                :text-inside="true"
                :stroke-width="22"
                :percentage="50"
                status="exception"
              />

            </template>
          </vxe-column>
          <vxe-column field="status" title="Status" width="100">
            <template #default="{ row }">
              <span v-if="row.error">{{ row.error }}</span>
              <span v-else-if="row.success">success</span>
              <span v-else-if="row.active">active</span>
            </template>
          </vxe-column>

        </vxe-table>
      </span>
    </el-dialog>
    <el-dialog v-model="viewFileVisible" title="view" width="96%">
      <el-container>
        <el-aside width="20%">
          <vxe-table
            ref="xTable"
            keep-source
            :tooltip-config="{}"
            :data="filetableData"
            height="500"
            :show-header="true"
            :mouse-config="{selected: true}"
            :keyboard-config="{isArrow: true, isEnter: true, isChecked: true}"
            @cell-click="selFileView"
            @cell-selected="selFileView"
          >
            <vxe-column type="seq" title="Seq" width="60" />
            <vxe-column field="name" title="Name" show-overflow="ellipsis">
              <template #default="{ row }">
                <span>{{ row.name }}</span>
              </template>
            </vxe-column>
          </vxe-table>
        </el-aside>
        <el-main>
          <warning-bar :title="current_file" />

        </el-main>
      </el-container>

    </el-dialog>

  </div>

</template>

<script setup>
import { getFileList, deleteBatch, editFileName, newBatch, getBatchList, changeStatus } from './api/fileUploadAndDownload'
import { downloadImage } from '@/utils/downloadImg'
import { formatDate } from '@/utils/format'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { useUserStore } from '@/pinia/modules/user'

import { ElMessage, ElMessageBox } from 'element-plus'
import { getCurrentInstance, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import FileUpload from 'vue-upload-component'
import { api as viewerApi } from 'v-viewer'

const path = ref(import.meta.env.VITE_BASE_API)

const imageUrl = ref('')
const imageCommon = ref('')

const files = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const search = ref({})
const tableData = ref([])
const userStore = useUserStore()
const route = useRoute()
const batchid = ref('0001')
const filetableData = ref([])
const filetotal = ref(0)

const dialogImageUrl = ref('')
const dialogVisible = ref(false)
const uploadVisible = ref(false)
const viewFileVisible = ref(false)
const current_file = ref('')
const current_fileurl = ref('')
const progress0 = ref(0)
const isAdmin = userStore.userInfo.authorityId=='888'

const postAction = path.value + '/detection/upload'
const headers = {
  'x-token': userStore.token,
  'user': userStore.userInfo.uuid,
  'app': route.name,
  'batchid': batchid,
  'progress': progress0
}

const extensions = 'jpeg|jpe|jpg|gif|png|webp'

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}
const currentInstance = getCurrentInstance()

const onAddFolder = async() => {
  if (!currentInstance.refs.upload.features.directory) {
    this.alert('Your browser does not support')
    return
  }
  const input = document.createElement('input')
  input.style = 'background: rgba(255, 255, 255, 0);overflow: hidden;position: fixed;width: 1px;height: 1px;z-index: -1;opacity: 0;'
  input.type = 'file'
  input.setAttribute('allowdirs', true)
  input.setAttribute('directory', true)
  input.setAttribute('webkitdirectory', true)
  input.multiple = true
  document.querySelector('body').appendChild(input)
  input.click()
  input.onchange = (e) => {
    currentInstance.refs.upload.addInputFile(input).then(function() {
      document.querySelector('body').removeChild(input)
    })
  }
}
const onAddFiles = async() => {
  if (!currentInstance.refs.upload.features.directory) {
    this.alert('Your browser does not support')
    return
  }
  const input = document.createElement('input')
  input.style = 'background: rgba(255, 255, 255, 0);overflow: hidden;position: fixed;width: 1px;height: 1px;z-index: -1;opacity: 0;'
  input.type = 'file'
  input.setAttribute('allowdirs', false)
  // input.setAttribute('directory', false)
  // input.setAttribute('webkitdirectory', false)
  input.multiple = true
  document.querySelector('body').appendChild(input)
  input.click()
  input.onchange = (e) => {
    currentInstance.refs.upload.addInputFile(input).then(function() {
      document.querySelector('body').removeChild(input)
    })
  }
}
const onUploadFiles = async() => {
  if (currentInstance.refs.upload.active == false) {
    const table = await getBatchList({
      page: 0,
      pageSize: 1,
      user: userStore.userInfo.uuid,
      app: route.name,
      keyword: batchid.value,
    })
    if (table.code === 0) {
      table.data.list
      for (var i = 0; i < table.data.total; i++) {
        if (batchid.value == table.data.list[i].batchid) {
          ElMessage({
            type: 'error',
            message: 'same batch name!',
          })
          return
        }
      }
    }

    var all_size = 0
    var all_count = 0
    for (var i = 0; i < files.value.length; i++) {
      all_size += files.value[i].size
      all_count += 1
    }
    const res = await newBatch({
      'own': userStore.userInfo.uuid,
      'app': route.name,
      'batchid': batchid.value,
      'files_count': all_count,
      'files_size': formatSize(all_size)
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: 'newbatch成功!',
      })
      currentInstance.refs.upload.active = true
    }
  }
}

const selFileView = (column) => {
  console.log(`单元格点击${column.row.name}`)
  current_file.value = column.row.name
  var url = column.row.url
  if (url !== '' && url.slice(0, 4) === 'http') {
    current_fileurl.value = url
  } else {
    current_fileurl.value = path.value + '/' + url
  }
}

// 查询
const getTableData = async() => {
  const table = await getBatchList({
    page: page.value,
    pageSize: pageSize.value,
    user: userStore.userInfo.uuid,
    app: route.name,
    ...search.value,
  })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

const onUploadFinish = async() => {
  const res = await changeStatus({
    'own': userStore.userInfo.uuid,
    'app': route.name,
    'batchid': batchid.value,
    'status': 'ready',
  })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: 'UploadFinish成功!',
    })
  }
}

const onChangeStatus = async(row) => {
  var newstatus = row.status == 'ready' ? 'pause' : 'ready'
  const res = await changeStatus({
    'own': userStore.userInfo.uuid,
    'app': route.name,
    'batchid': row.batchid,
    'status': newstatus,
  })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: newstatus + '成功!',
    })
    getTableData()
  }
}

const viewFiles = async(row, showres) => {
  const table = await getFileList({
    // page: page.value,
    // pageSize: pageSize.value,
    user: userStore.userInfo.uuid,
    app: route.name,
    batchid: row.batchid,
    ...search.value,
  })
  if (table.code === 0) {
    filetableData.value = table.data.list
    filetotal.value = table.data.total
    // page.value = table.data.page
    // pageSize.value = table.data.pageSize
  }
  // viewFileVisible.value = true
  var sourceImageURLs = []
  var sourceImageNames = []
  for (var i = 0; i < filetotal.value; i++) {
    var url = filetableData.value[i].url
    var name = filetableData.value[i].name
    if (showres) {
      url = filetableData.value[i].url_detection
    }
    if (url !== '') {
      if (url.slice(0, 4) === 'http') {
        sourceImageURLs.push({
          thumbnail: url,
          src: url,
          title: name
        })
      } else {
        sourceImageURLs.push({
          thumbnail: path.value + '/' + url,
          src: path.value + '/' + url,
          title: name
        })
      }
    }
  }
  const viewer = viewerApi({
    options: {
      toolbar: true,
      url: 'src',
      initialViewIndex: 1,
    },
    images: sourceImageURLs,
  })
}

const inputFilter = (newFile, oldFile, prevent) => {
  if (newFile && !oldFile) {
    // Add file

    // Filter non-image file
    // Will not be added to files
    // if (!/\.(jpeg|jpe|jpg|gif|png|webp)$/i.test(newFile.name)) {
    //   return prevent()
    // }
    const regex = new RegExp(`\\.(${extensions})$`, 'i')
    if (!regex.test(newFile.name)) {
      return prevent()
    }

    // Create the 'blob' field for thumbnail preview
    newFile.blob = ''
    const URL = window.URL || window.webkitURL
    if (URL && URL.createObjectURL) {
      newFile.blob = URL.createObjectURL(newFile.file)
    }
  }

  if (newFile && oldFile) {
    // Update file

    // Increase the version number
    if (!newFile.version) {
      newFile.version = 0
    }
    newFile.version++
  }

  if (!newFile && oldFile) {
    // Remove file

    // Refused to remove the file
    // return prevent()
  }
}

// const getRouterData = async () => {
//   const table = await getRouterName({
//     page: page.value,
//     pageSize: pageSize.value,
//     ...search.value,
//   });
//   if (table.code === 0) {
//     ElMessage("Warning",table.data);
//   }
// };
// getRouterData();

const onDeleteBatch = async(row) => {
  ElMessageBox.confirm('此操作将永久删除文件, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async() => {
      const res = await deleteBatch(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!',
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除',
      })
    })
}

const downloadFile = (row) => {
  if (row.url.indexOf('http://') > -1 || row.url.indexOf('https://') > -1) {
    downloadImage(row.url, row.name)
  } else {
    debugger
    downloadImage(path.value + '/' + row.url, row.name)
  }
}

/**
 * 编辑文件名或者备注
 * @param row
 * @returns {Promise<void>}
 */
const editFileNameFunc = async(row) => {
  ElMessageBox.prompt('请输入文件名或者备注', '编辑', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /\S/,
    inputErrorMessage: '不能为空',
    inputValue: row.name,
  })
    .then(async({ value }) => {
      row.name = value
      // console.log(row)
      const res = await editFileName(row)
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '编辑成功!',
        })
        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消修改',
      })
    })
}

const handlePictureCardPreview = (url) => {
  if (url !== '' && url.slice(0, 4) === 'http') {
    dialogImageUrl.value = url
  } else {
    dialogImageUrl.value = path.value + '/' + url
  }
  dialogVisible.value = true
}

const timer = ref(null)

const reload = async() => {
  getTableData()
}

const reset = async() => {
  files.value = []
}

timer.value = setInterval(() => {
  reload()
}, 1000 * 5)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
})
const uploadDialog = async() => {
  uploadVisible.value = true
  if (files.value.length == 0) {
    batchid.value = getCurrentTime()
  }
}

function formatSize(size) {
  if (size > 1024 * 1024 * 1024 * 1024) {
    return (size / 1024 / 1024 / 1024 / 1024).toFixed(2) + ' TB'
  } else if (size > 1024 * 1024 * 1024) {
    return (size / 1024 / 1024 / 1024).toFixed(2) + ' GB'
  } else if (size > 1024 * 1024) {
    return (size / 1024 / 1024).toFixed(2) + ' MB'
  } else if (size > 1024) {
    return (size / 1024).toFixed(2) + ' KB'
  }
  return size.toString() + ' B'
}

function getCurrentTime() {
  const date = new Date()
  const year = date.getFullYear()
  const month = padZero(date.getMonth() + 1)
  const day = padZero(date.getDate())
  const hour = padZero(date.getHours())
  const minute = padZero(date.getMinutes())
  const second = padZero(date.getSeconds())
  return `${year}${month}${day}-${hour}${minute}${second}`
}

function padZero(num) {
  return num < 10 ? `0${num}` : num
}

</script>

<script>
export default {
  name: 'Upload',
}
</script>
<style scoped>
.name {
    cursor: pointer;
}

.upload-btn + .upload-btn {
    margin-left: 12px;
}

.gva-table-box label.btn {
    margin-bottom: 0;
    margin-right: 1rem;
}

.long-word {
    overflow: hidden; /*溢出的部分隐藏*/
    white-space: nowrap; /*文本不换行*/
    text-overflow: ellipsis; /*ellipsis:文本溢出显示省略号（...）；clip：不显示省略标记（...），而是简单的裁切*/
}
</style>
