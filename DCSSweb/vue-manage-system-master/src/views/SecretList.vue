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
          :row-class-name="tableRowClassName">
        <el-table-column
            prop="secretname"
            label="秘密名称"
            width="180">
        </el-table-column>
        <el-table-column
            prop="secretid"
            label="秘密ID"
            width="500">
        </el-table-column>
        <el-table-column
            prop="degree"
            label="门限阈值">
        </el-table-column>
        <el-table-column
            prop="counter"
            label="委员会成员数">
        </el-table-column>
        <el-table-column
            fixed="right"
            label="操作"
            width="100">
<!--          <template slot-scope="scope">-->
            <router-link to="/secretinfo">
              <el-button @click="handleClick(scope.row)" type="text" size="small">查看</el-button>
            </router-link>

            <el-button  type="text" size="small">编辑</el-button>
<!--          </template>-->
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
       const url = "http://localhost:8080/api/getsecretlist"
      axios({
        methods: 'get',
        url:url
      }).then(
          function (res) {
            // var arr = this;
            console.log(res);
            console.log(res.data.total);
            console.log(res.data.datalist);
            arr.tableData = res.data.datalist;
          }
      ).catch(err =>{

      })
    },
    handleClick(){

    },
    fresh(){
      this.getsecretlist()
    }
  },
}

</script>