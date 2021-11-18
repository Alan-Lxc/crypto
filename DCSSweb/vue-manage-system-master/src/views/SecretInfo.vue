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
    </div>
    <div class="container">
      <div class="form-box">
        <el-form ref="secretRef" label-width="160px">

          <el-form-item label="门限阈值" >
            {{secretinfo.degree}}
          </el-form-item>
          <el-form-item label="委员会成员数" prop="numberOfN">
            {{secretinfo.counter}}
          </el-form-item>
          <el-form-item label="秘密描述" prop="secret">
            {{secretinfo.description}}
          </el-form-item>
        </el-form>
      </div>
    </div>

    <el-table
        :data=null
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


    </el-table>

  </div>

</template>

<script>
import { ref, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { fetchData } from "../api/index";
import axios from "axios";
export default {

  name: "SecretInfo",
  data() {
    return {
      secretinfo: {
        degree: 0,
        counter: 0,
        user_id: 1,
        description: "",
      },
      nodelist: {
        node
      },
    }
  },
  created() {
    let arr = this;
    let id = this.$route.params.id;
    console.log(id);
    axios.get("http://localhost:8080/api/secret/getsecret",{
      params: {
        id: id
      }
    }).then(
        function (res) {
          console.log(res.data.data.secret);
          arr.secretinfo=res.data.data.secret;
        }
    ).catch(err =>{

    });
  },
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