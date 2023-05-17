<template>
  <div>
    <el-row :gutter="15" class="system_state">
      <el-col :span="12">
        <el-card v-if="state.os" class="card_item">
          <template #header>
            <div>Runtime</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">os:</el-col>
              <el-col :span="12" v-text="state.os.goos" />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">cpu nums:</el-col>
              <el-col :span="12" v-text="state.os.numCpu" />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">compiler:</el-col>
              <el-col :span="12" v-text="state.os.compiler" />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">go version:</el-col>
              <el-col :span="12" v-text="state.os.goVersion" />
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">goroutine nums:</el-col>
              <el-col :span="12" v-text="state.os.numGoroutine" />
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card v-if="state.disk" class="card_item">
          <template #header>
            <div>Disk</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col :span="12" v-text="state.disk.totalMb" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col :span="12" v-text="state.disk.usedMb" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col :span="12" v-text="state.disk.totalGb" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col :span="12" v-text="state.disk.usedGb" />
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                  type="dashboard"
                  :percentage="state.disk.usedPercent"
                  :color="colors"
                />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="15" class="system_state">
      <el-col :span="12">
        <el-card
          v-if="state.cpu"
          class="card_item"
          :body-style="{ height: '180px', 'overflow-y': 'scroll' }"
        >
          <template #header>
            <div>CPU</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">physical number of cores:</el-col>
              <el-col :span="12" v-text="state.cpu.cores" />
            </el-row>
            <el-row v-for="(item, index) in state.cpu.cpus" :key="index" :gutter="10">
              <el-col :span="12">core {{ index }}:</el-col>
              <el-col
                :span="12"
              >
                <el-progress
                  type="line"
                  :percentage="+item.toFixed(0)"
                  :color="colors"
                />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card v-if="state.ram" class="card_item">
          <template #header>
            <div>Ram</div>
          </template>
          <div>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-row :gutter="10">
                  <el-col :span="12">total (MB)</el-col>
                  <el-col :span="12" v-text="state.ram.totalMb" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (MB)</el-col>
                  <el-col :span="12" v-text="state.ram.usedMb" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">total (GB)</el-col>
                  <el-col :span="12" v-text="state.ram.totalMb / 1024" />
                </el-row>
                <el-row :gutter="10">
                  <el-col :span="12">used (GB)</el-col>
                  <el-col
                    :span="12"
                    v-text="(state.ram.usedMb / 1024).toFixed(2)"
                  />
                </el-row>
              </el-col>
              <el-col :span="12">
                <el-progress
                  type="dashboard"
                  :percentage="state.ram.usedPercent"
                  :color="colors"
                />
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="15" class="system_state">
      <el-col :span="12">
        <el-card v-if="state.gpu" class="card_item">
          <template #header>
            <div>GPU</div>
          </template>
          <div>
            <el-col v-for="(gpu, index) in state.gpu.gpu_infos" :key="gpu.index" :gutter="10">
              <el-row :gutter="10" tag="b">GPU {{ gpu.name }}:</el-row>
              <el-row :gutter="10">
                <el-col>
                  <p class="card-text">Memory Used: {{ gpu.memory_used }} / {{ gpu.memory_total }} MB</p>
                </el-col>
                <el-col>
                  <el-progress
                    type="line"
                    :percentage="(gpu.memory_used / gpu.memory_total * 100).toFixed(0)"
                    :color="colors"
                  />
                </el-col>

              </el-row>
              <el-row :gutter="10">
                <el-col>
                  <el-col>
                    <p class="card-text">GPU Utilization: </p>
                  </el-col>
                  <el-progress
                    type="line"
                    :percentage="gpu.utilization_gpu"
                    :color="colors"
                  />
                </el-col>
              </el-row>
              <el-divider />
            </el-col>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card v-if="state.gpu" class="card_item" :body-style="{ height: '180px', 'overflow-y': 'scroll' }">
          <template #header>
            <div>Process</div>
          </template>
          <div>
            <el-col v-for="(process, index) in state.gpu.processes" :key="process.pid" :gutter="10">
              <el-row :gutter="10" tag="b">
                <el-col :span="12">
                  <p class="card-text">Process {{ index + 1 }}:</p>
                </el-col>
                <el-col :span="12">
                  <p class="card-text">PID: {{ process.pid }}</p>
                </el-col>
              </el-row>
              <el-row :gutter="10">Used GPU Memory: {{ process.used_gpu_memory }} MB</el-row>
              <el-row :gutter="10">User: {{ process.user }}</el-row>
              <el-row :gutter="10">Command: {{ process.command }}</el-row>
              <el-divider />
            </el-col>
          </div>
        </el-card>
      </el-col>
    </el-row>

  </div>
</template>

<script setup>
import { getSystemState } from '@/api/system'
import { onUnmounted, ref } from 'vue'

const timer = ref(null)
const state = ref({})
const colors = ref([
  { color: '#5cb87a', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#f56c6c', percentage: 80 }
])

const reload = async() => {
  const { data } = await getSystemState()
  state.value = data.server
}

reload()
timer.value = setInterval(() => {
  reload()
}, 1000 * 10)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
})

</script>

<script>
export default {
  name: 'State',
}
</script>

<style>
.system_state {
    padding: 10px;
}

.card_item {
    height: 280px;
}
</style>
