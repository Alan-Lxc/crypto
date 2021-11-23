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
      <el-row>
<!--        <el-col :span="7"></el-col>-->
<!--        <el-col :span="8">-->
<!--          <div class="grid-content bg-purple-light">-->

<!--            <el-button @click="tonewsecret()" type="">新建秘密</el-button>-->
<!--          </div>-->
<!--        </el-col>-->
        <el-col :span="21"></el-col>
        <el-col :span="3">
          <div class="grid-content bg-purple-light">
            <el-button @click="fresh" type="" >Refresh</el-button>
          </div>
        </el-col>
      </el-row>
    </div>
    <div class="container">
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
      return '';rules
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
            arr.tableData = res.data.data.secretlist;
          }
      ).catch(err =>{
        console.log(err);
      })
    },
    handleClick(row){
      this.$router.push({
        path:"/secretinfo",
        query:{
          id : row.ID,
        }
      })
    },
    fresh(){
      this.getsecretlist()
    },
    tonewsecret(){
      let that = this
      this.$router.push({
        path:"/newsecret",
      })
    },
  },
}

</script>