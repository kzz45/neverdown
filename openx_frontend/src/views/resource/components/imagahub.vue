<template>
  <div>
    <div class="hub-item">
      <div class="hub-value">Project</div>
      <div class="hub-value">Domain</div>
      <div class="hub-value">Repository</div>
      <div class="hub-value">Tag</div>
      <el-select
        v-model="ProjectValue"
        filterable
        placeholder="Select"
        @change="changePro"
      >
        <el-option
          v-for="Project in ProjectList"
          :key="Project.metadata.name"
          :label="Project.metadata.name"
          :value="Project.metadata.name"
        />
      </el-select>
      <el-select v-model="DomainValue" filterable placeholder="Select">
        <el-option
          v-for="Domain in DomainList"
          :key="Domain"
          :label="Domain"
          :value="Domain"
        />
      </el-select>
      <el-select
        v-model="RepositoryValue"
        filterable
        placeholder="Select"
        @change="changeRep"
      >
        <el-option
          v-for="Project in RepositoryComputed"
          :key="Project.metadata.name"
          :label="Project.spec.repositoryMeta.repositoryName"
          :value="Project.spec.repositoryMeta.repositoryName"
        />
      </el-select>
      <el-select v-model="TagValue" filterable placeholder="Select">
        <el-option
          v-for="Project in TagComputed"
          :key="Project.metadata.name"
          :label="Project.spec.tag"
          :value="Project.spec.tag"
        />
      </el-select>
    </div>
    <div class="hub-title" style="margin-top: 6px">
      <el-input style="width: 100%" v-model="imageInput" @change="inputImage" />
    </div>
    <el-descriptions :column="3" border>
      <el-descriptions-item label="author" label-align="right" align="center">
        {{ showTagInfo.author }}
      </el-descriptions-item>
      <el-descriptions-item
        label="lastModifiedTime"
        label-align="right"
        align="center"
      >
        {{
          showTagInfo.lastModifiedTime
            ? formatTime(showTagInfo.lastModifiedTime)
            : ""
        }}
      </el-descriptions-item>
      <el-descriptions-item label="sha256" label-align="right" align="center">
        {{ showTagInfo.sha256.slice(0, 7) }}
      </el-descriptions-item>
    </el-descriptions>
  </div>
</template>

<style lang="scss" scoped>
.hub-item {
  display: grid;
  grid-template-columns: 25% 35% 25% 15%;
  gap: 3px;
  padding-right: 20px;
  font-size: 0.7rem;
}
.hub-title {
  font-weight: 600;
}
</style>

<script setup lang="ts">
import { cloneDeep } from "lodash";
import { computed, ref, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { initSocketData, sendSocketMessage, binaryToStr } from "@/api/socket";
import { useStore } from "@/store";
import { formatTime } from "../util";
const store = useStore();
const router = useRouter();
const route = useRoute();

const props = defineProps<{
  image: any;
}>();

let ProjectValue = ref("");
let DomainValue = ref("");
let RepositoryValue = ref("");
let TagValue = ref("");

let imageInput = ref("");

watch(
  () => props.image,
  (newImage) => {
    initShowImage();
  }
);
initShowImage();

function initShowImage() {
  const imageArr = props.image.split("/");
  ProjectValue.value = imageArr[1] || "";
  DomainValue.value = imageArr[0] || "";
  if (imageArr[2]) {
    const tagArr = imageArr[2].split(":");
    RepositoryValue.value = tagArr[0] || "";
    TagValue.value = tagArr[1] || "";
  } else {
    RepositoryValue.value = "";
    TagValue.value = "";
  }
  imageInput.value = imageStr();
}
let showTagInfo = computed(() => {
  let tagName = TagValue.value;
  let tagIndex = TagComputed.value.findIndex((tag: any) => {
    return tag.spec.tag === tagName;
  });
  if (tagIndex >= 0) {
    return TagComputed.value[tagIndex]?.spec?.dockerImage;
  } else {
    return {
      author: "null",
      sha256: "",
    };
  }
});

let getList = function (gvk: string) {
  const nsGvk = route.path.split("/");
  const senddata = initSocketData("Request", nsGvk[1], gvk, "list");
  sendSocketMessage(senddata, store);
};

getList("jingx-v1-Project");
getList("jingx-v1-Repository");
getList("jingx-v1-Tag");

import { returnResource, deleteSocketData, currentPageVue } from "./../util";

let loading = ref(false);
let loadOver = function () {
  loading.value = false;
};
let ProjectList = ref<any>([]);
let RepositoryList = ref([]);
let TagList = ref([]);

let RepositoryComputed: any = computed(() => {
  return RepositoryList.value.filter((rep: any) => {
    return rep.spec.repositoryMeta.projectName === ProjectValue.value;
  });
});
let DomainList = computed(() => {
  const project = ProjectList.value.findIndex((pro: any) => {
    return pro.metadata.name === ProjectValue.value;
  });
  if (project >= 0) {
    return ProjectList.value[project].spec.domains;
  } else {
    return [];
  }
});
let TagComputed: any = computed(() => {
  const proName = ProjectValue.value;
  const repName = RepositoryValue.value;
  let resultList: any = TagList.value.filter((tag: any) => {
    return (
      tag.spec.repositoryMeta.projectName === proName &&
      tag.spec.repositoryMeta.repositoryName === repName
    );
  });
  resultList.sort((ltag: any, rtag: any) => {
    if (ltag.spec.tag === "latest") {
      return -1;
    }
    if (rtag.spec.tag === "latest") {
      return 1;
    }
    let ltagStr = ltag.spec.tag.replace("v", "").split(".");
    let RtagStr = rtag.spec.tag.replace("v", "").split(".");
    if (Number(RtagStr[0]) === Number(ltagStr[0])) {
      if (ltagStr[1] && RtagStr[1] && ltagStr[1] != RtagStr[1]) {
        return Number(RtagStr[1]) - Number(ltagStr[1]);
      } else {
        if (ltagStr[2] && RtagStr[2] && ltagStr[2] != RtagStr[2]) {
          return Number(RtagStr[2]) - Number(ltagStr[2]);
        } else {
          let ltTemp = ltagStr[3] || "0",
            rtTemp = RtagStr[3] || "0";
          if (ltTemp != rtTemp) {
            return Number(rtTemp) - Number(ltTemp);
          } else {
            return 0;
          }
        }
      }
    } else {
      return Number(RtagStr[0]) - Number(ltagStr[0]);
    }
  });

  return resultList;
});

// watch(() => RepositoryComputed.value, ()=>{
//   if(RepositoryComputed.value.length > 0){
//     RepositoryValue.value = RepositoryComputed.value[0].spec.repositoryMeta.repositoryName
//   }else{
//     RepositoryValue.value = ''
//   }
// })
// watch(() => TagComputed.value, ()=>{
//   //console.log('TagComputed.value', TagComputed.value)
//   if(TagComputed.value.length > 0){
//     TagValue.value = TagComputed.value[0].spec.tag
//   }else{
//     TagValue.value = ''
//   }
// })
function changePro() {
  if (RepositoryComputed.value.length > 0) {
    let resIndex = RepositoryComputed.value.findIndex((reposi: any) => {
      return (
        reposi.spec.repositoryMeta.repositoryName === RepositoryValue.value
      );
    });
    if (resIndex >= 0) {
      RepositoryValue.value =
        RepositoryComputed.value[resIndex].spec.repositoryMeta.repositoryName;
    } else {
      RepositoryValue.value =
        RepositoryComputed.value[0].spec.repositoryMeta.repositoryName;
    }
  } else {
    RepositoryValue.value = "";
  }
}
function changeRep() {
  if (TagComputed.value.length > 0) {
    TagValue.value = TagComputed.value[0].spec.tag;
  } else {
    TagValue.value = "";
  }
}

function initMessage(msg: any, type: string) {
  const nsGvk = route.path.split("/");
  const gvkArr = type.split("-");
  let gvkObj = {
    group: gvkArr[0],
    version: gvkArr[1],
    kind: gvkArr[2],
  };
  try {
    let resultList = returnResource(msg, nsGvk[1], gvkObj, loadOver);
    if (resultList) {
      resultList.sort((itemL: any, itemR: any) => {
        const itemLTime = itemL.metadata.creationTimestamp.seconds;
        const itemRTime = itemR.metadata.creationTimestamp.seconds;
        return itemRTime - itemLTime;
      });
      if (type === "jingx-v1-Project") {
        ProjectList.value = resultList;
      }
      if (type === "jingx-v1-Repository") {
        RepositoryList.value = resultList;
      }
      if (type === "jingx-v1-Tag") {
        TagList.value = resultList;
      }
    }
  } catch (e) {
    console.log("error");
  }
}

watch(
  () => store.state.socket.socket.message,
  (msg) => {
    const requestList = [
      "jingx-v1-Project",
      "jingx-v1-Repository",
      "jingx-v1-Tag",
    ];
    for (let requestType of requestList) {
      initMessage(msg, requestType);
    }
  }
);

const emit = defineEmits(["image-change"]);
function inputImage(newImage: any) {
  emit("image-change", newImage);
}

function imageStr() {
  if (
    DomainValue.value &&
    ProjectValue.value &&
    RepositoryValue.value &&
    TagValue.value
  ) {
    return `${DomainValue.value}/${ProjectValue.value}/${RepositoryValue.value}:${TagValue.value}`;
  } else {
    return String(props.image);
  }
}
watch(
  () => imageStr(),
  (newImage) => {
    emit("image-change", newImage);
  }
);
</script>
