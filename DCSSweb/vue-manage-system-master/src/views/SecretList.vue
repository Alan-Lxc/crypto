<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-cascades"></i> 秘密列表
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="container">
      <el-header height="40px" >
        <router-link to="/newsecret">
          <el-button @click="" type="mini" bgcolor="bule" >新建秘密</el-button>
        </router-link>
        <el-button @click="fresh" type="mini" >Refresh</el-button>
      </el-header>
      <el-table
          :data="tableData"
          style="width: 100%"
          :row-class-name="tableRowClassName"
          @row-click="handleClick">
        <el-table-column
            prop="Secretname"
            label="秘密名称"
            width="180">
        </el-table-column>
        <el-table-column
            prop="ID"
            label="秘密ID"
            width="500">
        </el-table-column>
        <el-table-column
            prop="Degree"
            label="门限阈值">
        </el-table-column>
        <el-table-column
            prop="Counter"
            label="委员会成员数">
        </el-table-column>

      </el-table>
    </div>
  </div>

</template>

<style>
/*.el-table .warning-row {*/
/*  background: oldlace;*/
/*}*/

/*.el-table .success-row {*/
/*  background: #f0f9eb;*/
/*}*/
</style>

<script>
import axios from "axios";
import {useRouter} from "vue-router";
export default {
  data() {
    return {
      tableData: [],
      total:0,
    }
  },
  created() {
    this.getsecretlist()
  },
  methods: {
    tableRowClassName({row, rowIndex}) {
      if (rowIndex === 1) {
        return 'warning-row';
      } else if (rowIndex === 3) {
        return 'success-row';
      }
      return '';
    },
    getsecretlist(){
      let arr = this;
      const url = "http://localhost:8080/api/secret/getsecretlist";
      axios({
        methods: 'get',
        url:url,
        params: {
          "userid": 1,
        },

      }).then(
          function (res) {
            // var arr = this;
            arr.tableData = res.data.data.secretlist;
          }
      ).catch(err =>{
        console.log(err);
      })
    },
    handleClick(row){
      this.$router.push("/secretinfo/?id="+row.ID)
    },
    fresh(){
      this.getsecretlist()
    }
  },
}

</script>