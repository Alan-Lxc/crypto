<template>
  <div>
    <div class="crumbs">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>
          <i class="el-icon-lx-people"></i>秘密详情
        </el-breadcrumb-item>
        <el-breadcrumb-item>
          节点信息
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="container">
      <div class="block">

        <el-header>
          <div class="radio">
            排序：
            <el-radio-group v-model="reverse">
              <el-radio :label="reverse">倒序</el-radio>
              <el-radio :label="!reverse">正序</el-radio>
            </el-radio-group>
          </div>
        </el-header>
        <el-timeline :reverse="reverse">
          <el-timeline-item
              v-for="(log, index) in logs"
              :key="index"
              :timestamp="log.timestamp">
            {{log.content}}
          </el-timeline-item>
        </el-timeline>
      </div>
    </div>


  </div>
</template>

<script>
import {reactive, ref} from "vue";
import axios from "axios";

export default {
  name: "NodeInfo",
  data(){
    return{
      logs : [{
        content: 'Node 1 new done,ip = 127.0.0.1:10001',
        timestamp: '2021/10/15 01:58:42'
      },{
        content: ' [Node 1] Now Get start Phase1',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send point message to [Node 2]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send point message to [Node 3]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send point message to [Node 4]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send point message to [Node 5]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' Phase 1 :[Node 1] receive point message from [Node 2]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' Phase 1 :[Node 1] receive point message from [Node 3]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' Phase 1 :[Node 1] receive point message from [Node 4]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' Phase 1 :[Node 1] receive point message from [Node 5]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] read bulletinboard in phase 1',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] receive zero message from [Node 2] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send message to [Node 2] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send message to [Node 3] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send message to [Node 4] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] send message to [Node 5] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] receive zero message from [Node 3] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] receive zero message from [Node 4] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] receive zero message from [Node 5] in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' 1 has finish _0ShareSum',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node %!d(MISSING)] start verification in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] read bulletinboard in phase 2',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] exponentSum: [919527479152314492520017575367191592683263256082405632185091027972819030997943823490183648180157342586948451807488935955194244950067545680105211023924584, 4488708331611103817365843553572127828114720293967289052761970848310698238665156340739995637022331696132887341554588441642277309627493593404039601398517850]',
        timestamp: '2021/10/15 01:58:43'
      },{
        content: ' [Node 1] receive point message from [Node 4] in phase3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' node 1 send point message to node 2 in phase 3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' node 1 send point message to node 3 in phase 3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' node 1 send point message to node 4 in phase 3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' node 1 send point message to node 5 in phase 3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' [Node 1] receive point message from [Node 5] in phase3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' [Node 1] receive point message from [Node 2] in phase3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' [Node 1] receive point message from [Node 3] in phase3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' [Node 1] has finish sharePhase3',
        timestamp: '2021/10/15 01:58:44'
      },{
        content: ' [Node 1] write bulletinboard in phase 3',
        timestamp: '2021/10/15 01:58:44'
      },],
      total: 0,
      reverse:true,
    }
  },
  created() {
    this.getloglist()
  },
  methods: {
    getloglist(){
      let arr = this;
      let secretid = arr.$route.query.secretid;
      let unitid = arr.$route.query.unitid;
      const url = "http://localhost:8080/api/unit/getunitlog";
      axios({
        methods: 'get',
        url:url,
        params: {
          "secretid": secretid,
          "unitid": unitid,
        },
      }).then(
          function (res) {
            arr.logs = res.data.data.logs;
          }
      ).catch(err =>{
        console.log(err);
      })
    },
  }
}
</script>

<style scoped>

</style>