<template>
  <div>

    <div class="container-fluid" style="background: #212529;color:#fff;padding: 1.5rem">
      <div class="row" v-for="(item,key) in detail" style="min-height: 40px">
        <div class="col-1" style="font-weight: bold">
          {{key}}
        </div>
        <div class="col" style="white-space: pre-wrap;">
          {{item}}
        </div>
      </div>
    </div>

  </div>
</template>
<script setup>

import {reactive,toRefs,onMounted} from "vue";
import { useRoute,useRouter } from 'vueRouter';
import request  from "request";

function getDetail(id,msgType){
  return request.get("/detail",{"params":{"id":id,"msgType":msgType}})
}

let data = reactive({
  detail:{}
});

const uRoute = useRoute();
onMounted(async ()=>{
  let id = uRoute.params.id;
  let msgType = uRoute.params.msgType;
  let res = await getDetail(id,msgType);
  data.detail = res.data;
})
const {detail} = toRefs(data);
</script>