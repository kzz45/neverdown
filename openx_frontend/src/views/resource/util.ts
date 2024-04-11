import { ElNotification } from 'element-plus'
import { binaryToStr } from "@/api/socket"
import proto from '../../../proto/proto'
import router from "@/router"
import moment from 'moment'
import { throttle } from 'lodash'
import { cloneDeep } from 'lodash'

const protoRequest = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto
const protoApi: any = proto.k8s.io.api
const protoOpenx: any = proto.github.com.kzz45.neverdown.pkg.apis.openx
const protoJingx: any = proto.github.com.kzz45.neverdown.pkg.apis.jingx
const protoRbac: any = proto.k8s.io.api.rbac
const protoMetrics: any = proto.k8s.io.metrics.pkg.apis.metrics

interface groupVersionKind {
  group: string,
  version: string,
  kind: string
}

export function returnResource(
  msg: any, nameSpace: string, gvk: groupVersionKind,
  loadCallback: any, refreshCallback?: any, getList?: any) {
  const result: any = protoRequest.Response.decode(msg)
  if (result.verb === 'ping') return
  let resultGvk = result.groupVersionKind
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else if (gvk.group.startsWith('metrics')) {
    protoFetch = protoMetrics[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  if (!resultGvk.group) {
    resultGvk.group = 'core'
  }
  const gvkStr = `${resultGvk.group}-${resultGvk.version}-${resultGvk.kind}`
  const requestGvk = `${gvk.group}-${gvk.version}-${gvk.kind}`
  if (gvkStr === requestGvk && result.namespace === nameSpace) {
    loadCallback()
    if (result.code === 401) {
      ElNotification({
        title: '权限不足',
        message: `gvk: ${gvkStr}, verb: ${result.verb}`,
        type: 'error', duration: 0
      })
      router.replace({ name: 'Index' })
      return
    }
    switch (result.verb) {
      case 'create':
        if (result.code === 0) {
          ElNotification({ title: '新增成功', message: 'success', type: 'success', duration: 2000 })
        } else {
          ElNotification({ title: '新增失败', message: binaryToStr(result.raw), type: 'error', duration: 3000 })
        }
        if (typeof getList === 'function') {
          getList()
        }
        break
      case 'update':
        if (result.code === 0) {
          ElNotification({ title: '修改成功', message: 'success', type: 'success', duration: 2000 })
        } else {
          ElNotification({ title: '修改失败', message: binaryToStr(result.raw), type: 'error', duration: 3000 })
        }
        if (typeof getList === 'function') {
          getList()
        }
        break
      case 'delete':
        if (result.code === 0) {
          ElNotification({ title: '删除成功', message: 'success', type: 'success', duration: 2000 })
        } else {
          ElNotification({ title: '删除失败', message: binaryToStr(result.raw), type: 'error', duration: 3000 })
        }
        if (typeof getList === 'function') {
          getList()
        }
        break
      case 'list':
        return protoFetch[`${gvk.kind}List`].decode(result.raw).items
      case 'watch':
        const watchEvent = protoRequest.WatchEvent.decode(result.raw)
        if (typeof refreshCallback === 'function') {
          const decodeRaw = protoFetch[`${gvk.kind}`].decode(watchEvent.raw)
          refreshCallback(watchEvent.type, decodeRaw)
        }
        break
      default: break
    }
  }
}
export function returnResourceList(msg: any, nameSpace: string, gvk: groupVersionKind) {
  const result: any = protoRequest.Response.decode(msg)
  if (result.verb === 'ping') return
  let resultGvk = result.groupVersionKind
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else if (gvk.group.startsWith('metrics')) {
    protoFetch = protoMetrics[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  if (!resultGvk.group) {
    resultGvk.group = 'core'
  }
  const gvkStr = `${resultGvk.group}-${resultGvk.version}-${resultGvk.kind}`
  const requestGvk = `${gvk.group}-${gvk.version}-${gvk.kind}`
  if (gvkStr === requestGvk && result.namespace === nameSpace) {
    switch (result.verb) {
      case 'list':
        return protoFetch[`${gvk.kind}List`].decode(result.raw).items
        break
      default: break
    }
  }
}

export function deleteSocketData(gvk: any, item: any) {
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  const createData = {
    metadata: item.metadata
  }
  const message = protoFetch[gvk.kind].create(createData)
  const param = protoFetch[gvk.kind].encode(message).finish()
  return param
}
export function updateSocketData(gvk: any, item: any) {
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  const message = cloneDeep(protoFetch[gvk.kind].create(item))
  if (message.spec?.applications) {
    checkOpenx(message.spec.applications)
  }
  if (validate(gvk, message)) {
    const param = protoFetch[gvk.kind].encode(message).finish()
    return param
  }
  return false
}
function checkOpenx(apps: any) {
  for (let app of apps) {
    if (app.horizontalPodAutoscalerSpec) {
      for (let met of app.horizontalPodAutoscalerSpec.metrics) {
        const average = met.resource.target?.averageValue?.string
        if (!average) {
          delete met.resource.target.averageValue
        }
      }
    }
    let nodeSelectorTerms = app.pod.spec?.affinity?.nodeAffinity?.requiredDuringSchedulingIgnoredDuringExecution?.nodeSelectorTerms
    let hasAff = (nodeSelectorTerms && nodeSelectorTerms.length > 0)
    if (!hasAff) {
      app.pod.spec.affinity = {}
    }
    for (let con of app.pod.spec.containers) {
      if (con.resources?.limits?.cpu?.string == '0' && con.resources?.limits?.memory?.string == '0') {
        delete con.resources.limits
      }
      if (con.resources?.requests?.cpu?.string == '0' && con.resources?.requests?.memory?.string == '0') {
        delete con.resources.requests
      }
    }
  }
}


export const encodeify = (gvk: any, item: any) => {
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  const message = protoFetch[gvk.kind].create(item)
  const fetchEncode = protoFetch[gvk.kind].encode(message).finish()
  return fetchEncode
}
export const decodeify = (gvk: any, item: any) => {
  let protoFetch: any = {}
  if (gvk.group.startsWith('openx')) {
    protoFetch = protoOpenx[gvk.version]
  } else if (gvk.group.startsWith('rbac')) {
    protoFetch = protoRbac[gvk.version]
  } else if (gvk.group === 'jingx') {
    protoFetch = protoJingx[gvk.version]
  } else {
    protoFetch = protoApi[gvk.group][gvk.version]
  }
  return protoFetch[gvk.kind].decode(item)
}

import valiFnc from './util_vali'
import { ElMessage } from 'element-plus'
function validate(gvk: any, message: any) {
  let resultList = []
  const metaCheck = valiFnc.metadata_check(message.metadata)
  const specCheck = valiFnc.spec_check(message.spec)

  resultList.push(...metaCheck, ...specCheck)
  if (resultList.length > 0) {
    ElMessage({
      message: resultList[0].tips,
      type: 'error',
    })
    return false
  } else {
    return true
  }
}

export function formatTime(timeStamp: number) {
  return moment(timeStamp * 1000).format('YYYY-MM-DD HH:mm:ss')
}

export function initLabels(ans: any) {
  let tags = []
  for (let key in ans) {
    tags.push({
      label: key, value: ans[key]
    })
  }
  return tags
}

export function showLabel(anno: any) {
  return `${anno.label} : ${anno.value}`
}
export function fetchLabel(anno: any) {
  if (anno.label.length < 20) {
    const fetchKey = anno.label
    const fetchValue = anno.value.length < 20 ? anno.value : anno.value.slice(0, 15) + '...'
    return `${fetchKey} : ${fetchValue}`
  } else {
    const fetchKey = anno.label.slice(0, 18) + '...'
    const fetchValue = anno.value.length < 20 ? anno.value : anno.value.slice(0, 15) + '...'
    return `${fetchKey} : ${fetchValue}`
  }
}
export function initAnnotations(ans: any) {
  let tags = []
  for (let key in ans) {
    tags.push({
      label: key, value: ans[key]
    })
  }
  return tags
}
const podTemp = {
  metadata: {
    labels: {}
  },
  spec: {
    volumes: [],
    containers: [],
    nodeSelector: {},
    serviceAccountName: '',
    imagePullSecrets: [],
    affinity: [],
    tolerations: []
  }
}
const serviceTemp = {
  metadata: {
    labels: {}
  },
  spec: {
    clusterIP: '',
    type: '',
    ports: []
  }
}

const metadata: any = {
  name: '',
  namespace: '',
  annotations: {},
  labels: {},
  creationTimestamp: {
    seconds: 0
  }
}
const Deployment: any = {
  metadata,
  spec: {
    replicas: 0,
    selector: {
      matchLabels: {}
    },
    template: podTemp
  }
}
const Affinity: any = {
  metadata,
  spec: {
    affinity: {
      nodeAffinity: {
        preferredDuringSchedulingIgnoredDuringExecution: [],
        requiredDuringSchedulingIgnoredDuringExecution: {
          nodeSelectorTerms: []
        }
      },
      podAffinity: {
        preferredDuringSchedulingIgnoredDuringExecution: [],
        requiredDuringSchedulingIgnoredDuringExecution: []
      },
      podAntiAffinity: {
        preferredDuringSchedulingIgnoredDuringExecution: [],
        requiredDuringSchedulingIgnoredDuringExecution: []
      }
    }
  }
}
const AliyunAccessControl: any = {
  metadata,
  spec: {
    instance: {
      key: '', value: ''
    },
    status: {
      key: '', value: ''
    },
    type: {
      key: '', value: ''
    }
  }
}
const AliyunLoadBalancer: any = {
  metadata,
  spec: {
    instance: {
      key: '', value: ''
    },
    overrideListeners: {
      key: '', value: ''
    }
  }
}
const NodeSelector: any = {
  metadata,
  spec: {
    nodeSelector: {}
  }
}
const Toleration: any = {
  metadata,
  spec: {
    toleration: {
      effect: "",
      key: "",
      operator: "",
      value: "",
    }
  }
}
const ConfigMap: any = {
  metadata,
  data: {},
  binaryData: {}
}
const Endpoints: any = {
  metadata,
  subsets: []
}
const Mysql: any = {
  metadata,
  spec: {
    master: {
      cloudNetworkConfig: {
        aliyunSLB: {}
      },
      persistentStorage: {
        storageVolumePath: ''
      },
      pod: {
        metadata: {
          labels: {}
        },
        spec: {
          volumes: [],
          containers: [],
          nodeSelector: {},
          serviceAccountName: '',
          imagePullSecrets: [],
          affinity: [],
          tolerations: []
        }
      },
      service: {
        metadata: {
          labels: {}
        },
        spec: {
          clusterIP: '',
          type: '',
          ports: []
        }
      }
    },
    slave: {
      cloudNetworkConfig: {
        aliyunSLB: {}
      },
      persistentStorage: {
        storageVolumePath: ''
      },
      pod: {
        metadata: {
          labels: {}
        },
        spec: {
          volumes: [],
          containers: [],
          nodeSelector: {},
          serviceAccountName: '',
          imagePullSecrets: [],
          affinity: [],
          tolerations: []
        }
      },
      service: {
        metadata: {
          labels: {}
        },
        spec: {
          clusterIP: '',
          type: '',
          ports: []
        }
      }
    }
  }
}
const NameSpace = {
  metadata,
  spec: {
    finalizers: []
  }
}
const NodeTemp = {
  metadata,
  spec: {
    externalID: '',
    podCIDR: '',
    podCIDRs: [],
    providerID: '',
    taints: [],
    unschedulable: false
  }
}
const PersistentVolume = {
  metadata,
  spec: {
    accessModes: [],
    capacity: {
      storage: {
        string: ''
      }
    },
    claimRef: {
      name: ""
    },
    mountOptions: [],
    persistentVolumeReclaimPolicy: '',
    persistentVolumeSource: {},
    storageClassName: "",
    volumeMode: ""
  }
}
const PersistentVolumeClaim = {
  apiVersion: "",
  fieldPath: "",
  kind: "",
  name: "",
  namespace: "",
  resourceVersion: "",
  uid: ""
}
const Pod = {
  metadata,
  spec: {
    volumes: [],
    containers: [],
    nodeSelector: {},
    serviceAccountName: '',
    imagePullSecrets: [],
    affinity: [],
    tolerations: []
  }
}
const Secret = {
  metadata,
  data: {},
  type: '',
  stringData: {}
}
const Service = {
  metadata,
  spec: {
    clusterIP: '',
    type: '',
    selector: {},
    ports: []
  }
}
const serviceAccount = {
  metadata,
  imagePullSecrets: [],
  secrets: []
}
const openx: any = {
  metadata,
  spec: {
    applications: []
  }
}
const Etcd = {
  metadata,
  spec: {
    replicas: 0,
    persistentStorage: {
      storageVolumePath: ''
    },
    pod: podTemp
  }
}
const ClusterRoleBinding = {
  metadata,
  roleRef: {
    apiGroup: '', kind: '', name: ''
  },
  subjects: []
}
const ClusterRole = {
  metadata,
  rules: []
}

const objectInitInfo: any = {
  'apps-v1-Deployment': Deployment,
  'apps-v1-StatefulSet': Deployment,
  'openx.neverdown.io-v1-Affinity': Affinity,
  'openx.neverdown.io-v1-AliyunAccessControl': AliyunAccessControl,
  'openx.neverdown.io-v1-AliyunLoadBalancer': AliyunLoadBalancer,
  'openx.neverdown.io-v1-NodeSelector': NodeSelector,
  'openx.neverdown.io-v1-Toleration': Toleration,
  'openx.neverdown.io-v1-Mysql': Mysql,
  'openx.neverdown.io-v1-Redis': Mysql,
  'rbac.authorization.k8s.io-v1-ClusterRoleBinding': ClusterRoleBinding,
  'rbac.authorization.k8s.io-v1-ClusterRole': ClusterRole,
  'core-v1-ConfigMap': ConfigMap,
  'core-v1-Namespace': NameSpace,
  'core-v1-Endpoints': Endpoints,
  'core-v1-Node': NodeTemp,
  'core-v1-PersistentVolume': PersistentVolume,
  'core-v1-PersistentVolumeClaim': PersistentVolumeClaim,
  'core-v1-Pod': Pod,
  'core-v1-Secret': Secret,
  'core-v1-Service': Service,
  'core-v1-ServiceAccount': serviceAccount,
  'openx.neverdown.io-v1-Openx': openx,
  'openx.neverdown.io-v1-Etcd': Etcd
}

export function initObject(gvk: string) {
  return objectInitInfo[gvk]
}

import appsv1Deployment from './appsv1Deployment.vue'
import appsv1StatefulSet from './appsv1StatefulSet.vue'
import corev1ConfigMap from './corev1ConfigMap.vue'
import corev1Endpoints from './corev1Endpoints.vue'
import corev1Event from './corev1Event.vue'
import corev1Namespace from './corev1Namespace.vue'
import corev1Node from './corev1Node.vue'
import corev1PersistentVolume from './corev1PersistentVolume.vue'
import corev1PersistentVolumeClaim from './corev1PersistentVolumeClaim.vue'
import corev1Pod from './corev1Pod.vue'
import corev1Secret from './corev1Secret.vue'
import corev1Service from './corev1Service.vue'
import corev1ServiceAccount from './corev1ServiceAccount.vue'

import rbacv1ClusterRoleBinding from './rbac.authorization.k8s.iov1ClusterRoleBinding.vue'
import rbacv1ClusterRole from './rbac.authorization.k8s.iov1ClusterRole.vue'

import openxv1Affinity from './openx.neverdown.iov1Affinity.vue'
import openxv1AliyunAccessControl from './openx.neverdown.iov1AliyunAccessControl.vue'
import openxv1AliyunLoadBalancer from './openx.neverdown.iov1AliyunLoadBalancer.vue'
import openxv1Mysql from './openx.neverdown.iov1Mysql.vue'
import openxv1NodeSelector from './openx.neverdown.iov1NodeSelector.vue'
import openxv1Redis from './openx.neverdown.iov1Redis.vue'
import openxv1Openx from './openx.neverdown.iov1Openx.vue'
import openxv1Etcd from './openx.neverdown.iov1Etcd.vue'
import openxv1Toleration from './openx.neverdown.iov1Toleration.vue'

const modules: any = {
  appsv1Deployment, appsv1StatefulSet,
  corev1ConfigMap, corev1Endpoints, corev1Event, corev1Namespace, corev1Node, corev1PersistentVolume,
  corev1PersistentVolumeClaim, corev1Pod, corev1Secret, corev1Service, corev1ServiceAccount,
  openxv1Affinity, openxv1AliyunAccessControl, openxv1AliyunLoadBalancer,
  openxv1Mysql, openxv1NodeSelector, openxv1Redis, openxv1Openx, openxv1Toleration,
  openxv1Etcd, rbacv1ClusterRoleBinding, rbacv1ClusterRole
}
export function currentPageVue(path: any) {
  const nsGvk = path.split('/')
  if (nsGvk.length < 3) {
    return
  }
  const gvkArr = nsGvk[2].split('-')
  let group = gvkArr[0]
  if (group.startsWith('openx.neverdown')) {
    group = 'openx'
  }
  if (group.startsWith('rbac.authorization')) {
    group = 'rbac'
  }
  let groupVersionKind = `${group}${gvkArr[1]}${gvkArr[2]}`
  return modules[groupVersionKind]
}

export function routerToNamespace(ns: string, toKind?: any) {
  let rule = JSON.parse(String(localStorage.getItem('clusterRole')))
  const gvkRule = rule[ns]
  let gvkList = []
  for (let gv in gvkRule) {
    const kinds = gvkRule[gv].gvk
    for (let gvkind of kinds) {
      let gvSplit = gv.split('/')
      if (gvSplit[0] !== 'jingx' && gvSplit[0] !== 'native') {
        gvkList.push({
          gv: `${gvSplit[0] || 'core'}-${gvSplit[1]}`,
          kind: gvkind.kind,
          verbs: gvkind.verbs
        })
      }
    }
  }
  if (gvkList.length <= 0) {
    ElNotification({ title: '权限不足', message: '请联系管理员', type: 'error', duration: 3000 })
    return
  }
  gvkList.sort((gvka: any, gvkb: any) => {
    const nameA = `${gvka.gv}-${gvka.kind}`.toUpperCase()
    const nameB = `${gvkb.gv}-${gvkb.kind}`.toUpperCase()
    if (nameA < nameB) {
      return -1;
    }
    if (nameA > nameB) {
      return 1;
    }
    return 0;
  })
  let to: any = gvkList[0]
  const podIndex = gvkList.findIndex((gvkind: any) => {
    return gvkind.kind === 'Pod'
  })
  if (podIndex > 0) {
    to = gvkList[podIndex]
  }
  if (toKind) {
    const findTo = gvkList.findIndex(gvk => {
      return `${gvk.gv}${gvk.kind}` === `${toKind.gv}${toKind.kind}`
    })
    if (findTo >= 0) {
      to = toKind
    }
  }
  localStorage.setItem('gvkList', JSON.stringify(gvkList))
  router.addRoute({ name: 'dashboard', path: '/dashboard', component: () => import("./../Dashboard.vue") })
  for (let one of gvkList) {
    router.addRoute('dashboard',
      {
        name: `${ns}${one.gv}${one.kind}`,
        path: `/${ns}/${one.gv}-${one.kind}`,
        component: () => import(`./list.vue`)
      })
  }
  router.push({ name: `${ns}${to.gv}${to.kind}` })
}


export class TimerUtil {
  static async sleepSeconds(seconds: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, seconds * 1000));
  }
}
