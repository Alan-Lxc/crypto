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
      <el-row>
        <!--        <el-col :span="6"><div class="grid-content bg-purple"><el-button type="">修改门限阈值</el-button></div></el-col>-->
        <el-col :span="8">
          <div class="grid-content bg-purple-light">

              <el-button @click="tochangesecret()" type="">修改委员会成员数</el-button>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content bg-purple">
            <el-button @click="handoffsecret()" type="">交接秘密</el-button>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="grid-content bg-purple-light">
            <el-button @click="reconstructsecret"  type="">重构秘密值</el-button>
          </div>
        </el-col>


      </el-row>
    </div>
    <div class="container">
      <div class="form-box">
        <el-form ref="secretRef" label-width="160px"  :data="secretinfo">

          <el-form-item label="门限阈值" >
            {{secretinfo.degree}}
          </el-form-item>
          <el-form-item label="委员会成员数" >
            {{secretinfo.counter}}
          </el-form-item>
          <el-form-item label="秘密创建时间" >
            {{secretinfo.create_time}}
          </el-form-item>
          <el-form-item label="上一次变更时间" >
            {{secretinfo.last_update_time}}
          </el-form-item>
          <el-form-item label="秘密描述" >
            {{secretinfo.description}}
          </el-form-item>
        </el-form>
      </div>
    </div>

    <el-table
        :data="nodelist"
        style="width: 100%"
        :row-class-name="tableRowClassName"
        @row-click="handleClick">
      <el-table-column
          prop="UnitId"
          label="节点ID"
          width="180">
      </el-table-column>
      <el-table-column
          prop="UnitIp"
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
      secretinfo: {      },
      nodelist: [],
      secretid: this.$route.query.id,
    }
  },
  created() {
    this.getsecretinfoAndUnitList()
  },
  methods :{
    getsecretinfoAndUnitList(){
      let that = this
      let secretid = that.$route.query.id;
      console.log(secretid);
      axios.get("http://localhost:8080/api/secret/getsecret",{
        params: {
          "secretid": secretid,
        },
      }).then(
          function (res) {
            that.secretinfo=res.data.data.secret;
          }
      ).catch(err =>{

      });
      axios.get("http://localhost:8080/api/unit/getunitlist",{
        params:{
          "secretid": secretid,
        }
      }).then(function (res){
        console.log(res.data.data.unitlist);
        that.nodelist = res.data.data.unitlist;
      });
    },
    handleClick(row){
      this.$router.push({
        path:"/unitinfo",
        query:{userid:row[unit_id],secretid:row[unit_ip]}
      })
    },
    tochangesecret(){
      let that = this
      this.$router.push({
        path:"/changesecret",
        query:{
          id:that.secretid,
          oldcounter:that.secretinfo.counter,
          degree:that.secretinfo.degree,
        }
      })
    },
    updatesecret(){
      axios.get("http://localhost:8080/api/secret/updatesecret",{
        params: {
          "id": this.secretid,
          "counter":90,
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
          }
      ).catch(err =>{

      });

    },
    handoffsecret(){
      let that = this;
      axios.get("http://localhost:8080/api/secret/handoffsecret",{
        params: {
          "secretid": that.secretid,
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
            alert("秘密值"+res.data.data.secret)
          }
      ).catch(err =>{});

    },
    reconstructsecret(){
      axios.get("http://localhost:8080/api/secret/reconstructsecret",{
        params: {
          "secretid": this.secretid
        }
      }).then(
          function (res) {
            console.log(res.data.data.secret);
            alert("秘密值"+res.data.data.secret)
          }
      ).catch(err =>{});
    }
  },

}
</script>

<style scoped>

</style>