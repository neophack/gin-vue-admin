<template>
  <div v-loading.fullscreen.lock="fullscreenLoading">
    <div class="gva-table-box">
      <warning-bar title="点击图片可以查看大图。按上传时间顺序后台识别，负载高时请耐心等待" />
      <div class="gva-btn-list">
        <upload-common
          v-model:imageCommon="imageCommon"
          class="upload-btn"
          @on-success="getTableData"
        />
        <upload-image
          v-model:imageUrl="imageUrl"
          :file-size="512"
          :max-w-h="1080"
          class="upload-btn"
          @on-success="getTableData"
        />

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
        <el-table-column align="left" label="上传图像" width="100">
          <template #default="scope">
            <CustomPic
              pic-type="file"
              :pic-src="scope.row.url"
              
              @click="handlePictureCardPreview(scope.row.url)"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="识别结果" width="100">
          <template #default="scope">
            <CustomPic
              pic-type="file"
              :pic-src="scope.row.url_detection"
              v-if="scope.row.url_detection!=''"
              @click="handlePictureCardPreview(scope.row.url_detection)"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="日期" prop="UpdatedAt" width="180">
          <template #default="scope">
            <div>{{ formatDate(scope.row.UpdatedAt) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="文件名/备注" prop="name" min-width="280">
          <template #default="scope">
            <div class="name" >
              {{ scope.row.name }}
            </div>
          </template>
        </el-table-column>
        <!-- <el-table-column align="left" label="链接" prop="url" min-width="300" /> -->
        <el-table-column align="left" label="类型" prop="tag" width="100">
          <template #default="scope">
            <el-tag
              :type="scope.row.tag === 'jpg' ? 'primary' : 'success'"
              disable-transitions
              >{{ scope.row.tag }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="160">
          <template #default="scope">
            <!-- <el-button icon="download" type="primary" link @click="downloadFile(scope.row)">下载</el-button> -->
            <el-button
              icon="delete"
              type="primary"
              link
              @click="deleteFileFunc(scope.row)"
              >删除</el-button
            >
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
  </div>
  <el-dialog v-model="dialogVisible">
    <img w-full :src="dialogImageUrl" alt="Preview Image" style="width: 100%" />
  </el-dialog>
</template>

<script setup>
import { getFileList, deleteFile, editFileName, getRouterName } from "./api/fileUploadAndDownload";
import { downloadImage } from "@/utils/downloadImg";
import CustomPic from "./detectionComponents/customPic/index.vue";
import UploadImage from "./detectionComponents/image.vue";
import UploadCommon from "./detectionComponents/common.vue";
import { formatDate } from "@/utils/format";
import WarningBar from "@/components/warningBar/warningBar.vue";
import { useUserStore } from '@/pinia/modules/user'

import { ElMessage, ElMessageBox } from "element-plus";
import { onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'
const path = ref(import.meta.env.VITE_BASE_API);

const imageUrl = ref("");
const imageCommon = ref("");

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const search = ref({});
const tableData = ref([]);
const userStore = useUserStore()
const route = useRoute()
// 分页
const handleSizeChange = (val) => {
  pageSize.value = val;
  getTableData();
};

const handleCurrentChange = (val) => {
  page.value = val;
  getTableData();
};

// 查询
const getTableData = async () => {
  const table = await getFileList({
    page: page.value,
    pageSize: pageSize.value,
    user:userStore.userInfo.uuid,
    app:route.name,
    ...search.value,
  });
  if (table.code === 0) {
    tableData.value = table.data.list;
    total.value = table.data.total;
    page.value = table.data.page;
    pageSize.value = table.data.pageSize;
  }
};
getTableData();

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

const deleteFileFunc = async (row) => {
  ElMessageBox.confirm("此操作将永久删除文件, 是否继续?", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async () => {
      const res = await deleteFile(row);
      if (res.code === 0) {
        ElMessage({
          type: "success",
          message: "删除成功!",
        });
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--;
        }
        getTableData();
      }
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "已取消删除",
      });
    });
};

const downloadFile = (row) => {
  if (row.url.indexOf("http://") > -1 || row.url.indexOf("https://") > -1) {
    downloadImage(row.url, row.name);
  } else {
    debugger;
    downloadImage(path.value + "/" + row.url, row.name);
  }
};

/**
 * 编辑文件名或者备注
 * @param row
 * @returns {Promise<void>}
 */
const editFileNameFunc = async (row) => {
  ElMessageBox.prompt("请输入文件名或者备注", "编辑", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    inputPattern: /\S/,
    inputErrorMessage: "不能为空",
    inputValue: row.name,
  })
    .then(async ({ value }) => {
      row.name = value;
      // console.log(row)
      const res = await editFileName(row);
      if (res.code === 0) {
        ElMessage({
          type: "success",
          message: "编辑成功!",
        });
        getTableData();
      }
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "取消修改",
      });
    });
};

const dialogImageUrl = ref("");
const dialogVisible = ref(false);
const handlePictureCardPreview = (url) => {
  if (url !== "" && url.slice(0, 4) === "http") {
    dialogImageUrl.value = url;
  }else{
    dialogImageUrl.value = path.value + "/" + url;
  }
  dialogVisible.value = true;
};

const timer = ref(null)

const reload = async() => {
  getTableData();
}

timer.value = setInterval(() => {
  reload()
}, 1000 * 5)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
})

</script>

<script>
export default {
  name: "Upload",
};
</script>
<style scoped>
.name {
  cursor: pointer;
}

.upload-btn + .upload-btn {
  margin-left: 12px;
}
</style>
