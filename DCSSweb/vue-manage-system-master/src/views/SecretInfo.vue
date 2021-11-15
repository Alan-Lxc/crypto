<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-text"></i> 秘密详情
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="container">
      <el-header>
        <el-button type="mini">修改门限阈值</el-button>
        <el-button type="mini">修改委员会成员数</el-button>
        <el-button type="mini">重构秘密值</el-button>
      </el-header>
      <el-container>
        <el-descriptions style="width: 100%" title="基本信息" :column="3" border>
          <el-descriptions-item
              label="门限阈值">1

          </el-descriptions-item>
          <el-descriptions-item
            label="委员会成员数">3

          </el-descriptions-item>
        </el-descriptions>

      </el-container>

      <el-table
          :data="tableData"
          style="width: 100%"
          :row-class-name="tableRowClassName">
        <el-table-column
            prop="nodeID"
            label="节点ID"
            width="180">
        </el-table-column>
        <el-table-column
            prop="nodeIP"
            label="节点IP"
            width="300">
        </el-table-column>
        <el-table-column>

        </el-table-column>


        <el-table-column
            fixed="right"
            label="操作"
            width="100">
          <!--          <template slot-scope="scope">-->
          <router-link to="/nodeinfo">
            <el-button @click="handleClick(scope.row)" type="text" size="small">查看</el-button>
          </router-link>

          <el-button type="text" size="small">重构</el-button>
          <!--          </template>-->
        </el-table-column>
      </el-table>

    </div>
  </div>
</template>

<script>
import { ref, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { fetchData } from "../api/index";
export default {
  name: "SecretInfo",
  setup() {
    const tableData = reactive( [
      {
        nodeID: 1,
        nodeIP: "193.168.0.1:10001"
      },
      {
        nodeID: 2,
        nodeIP: "193.168.0.1:10002"
      },
      {
        nodeID: 3,
        nodeIP: "193.168.0.1:10003"
      },

    ]);
    return {
      tableData,
    }
  }
}
</script>

<style scoped>

</style>