import proto from '../../../proto/proto'
import router from "@/router"
import axios from 'axios'
import moment from 'moment'

const protoRequest = proto.github.com.kzz45.neverdown.pkg.openx.aggregator.proto
const protoApi: any = proto.k8s.io.api
const protoOpenx: any = proto.github.com.kzz45.neverdown.pkg.apis.openx
const protoJingx: any = proto.github.com.kzz45.neverdown.pkg.apis.jingx
const protoRbac: any = proto.k8s.io.api.rbac

export async function openTerm(data: any, container: any, type: string, logSeconds: number) {
  const cig = await axios.get('config/config.json')
  let env = cig.data
  localStorage.setItem('termToken', String(localStorage.getItem('token')))
  if (type === 'log') {
    const routeData = router.resolve({
      path: '/log',
      query: { namespace: data.namespace, podname: data.name, containername: container.name, type, logSeconds }
    })
    window.open(routeData.href, '_blank')
  } else {
    const routeData = router.resolve({
      path: '/term',
      query: { namespace: data.namespace, podname: data.name, containername: container.name, type }
    })
    window.open(routeData.href, '_blank')
  }
}
import { saveAs } from 'file-saver'
export async function download(data: any, container: any, type: boolean, currentSeconds: number = 0) {
  let sinceTime = 0
  if (type) {
    sinceTime = 0
    container.downloadPre = true
  } else {
    sinceTime = currentSeconds
    container.downloadCur = true
  }
  const cig = await axios.get('config/config.json')
  let env = cig.data
  const url = `https://${String(env.VITE_BASE_URL)}/log/download/namespace/${data.namespace}/pod/${data.name}/container/${container.name}/previous/${type}/sinceSeconds/${sinceTime}/sinceTime/nil`

  console.log('url', url)
  axios({ url: url, method: 'get', timeout: 30000, headers: { 'Token': localStorage.getItem('token') } }).then(res => {

    const str = new Blob([res.data], { type: 'text/plain;charset=utf-8' })
    const ispre = type ? '_previous' : ''
    type ? (container.downloadPre = false) : (container.downloadCur = false)
    saveAs(str, data.namespace + '_' + data.name + '_' + container.name + ispre + '.log')
    if (type) {
      container.previousDownloading = false
    } else {
      container.currentDownloading = false
    }
  }).catch(err => {
    console.log('err ::', err)
    type ? (container.previousDownloading = false) : (container.currentDownloading = false)
  })
}

export const InfoInGvk: any = {
  'apps-v1-Deployment': ['Containers', 'DeploymentStatus', 'createTime'],
  'apps-v1-DaemonSet': ['DaemonSetStatus', 'createTime'],
  'apps-v1-StatefulSet': ['Containers', 'StatefulSetStatus', 'conditions', 'createTime'],
  'openx.neverdown.io-v1-Affinity': ['nodeSelectorTerms', 'createTime'],
  'openx.neverdown.io-v1-Mysql': ['master/slave', 'init', 'StatefulSetStatus', 'createTime'],
  'openx.neverdown.io-v1-Redis': ['master/slave', 'init', 'StatefulSetStatus', 'createTime'],
  'openx.neverdown.io-v1-Openx': ['applications', 'init', 'DeploymentStatus', 'HPAStatus', 'createTime'],
  'openx.neverdown.io-v1-Etcd': ['StatefulSetStatus', 'createTime'],
  'core-v1-ConfigMap': ['config', 'createTime'],
  'core-v1-Pod': ['phase', 'podIP', 'nodeName', 'cpu/memory', 'volume.timezone', 'container', 'conditions', 'createTime'],
  'core-v1-Event': ['message', 'type', 'createTime'],
  'core-v1-Namespace': ['phase', 'conditions', 'createTime'],
  'core-v1-PersistentVolume': ['phase', 'createTime'],
  'core-v1-Node': ['taints', 'cpu/memory', 'conditions', 'createTime'],
  'core-v1-Secret': ['createTime'],
  'core-v1-Service': ['clusterIP', 'type', 'ports', 'conditions', 'createTime'],
  'core-v1-ServiceAccount': ['createTime']
}

export function formatTime(timeStamp: number) {
  return moment(timeStamp * 1000).format('YYYY-MM-DD HH:mm:ss')
}

export function getInfoInGvk(type: string, detail: any, gvk: string) {
  if (type === 'createTime') {
    return formatTime(detail.metadata.creationTimestamp.seconds)
  }

  if (type === 'conditions') {
    return {
      type: 'conditions'
    }
  }

  if (gvk === 'apps-v1-DaemonSet') {
    if (type === 'DaemonSetStatus') {
      const daemonSetStatus = []
      daemonSetStatus.push(detail.status)
      return {
        type: 'DaemonSetStatus',
        daemonSetStatus
      }
    }
  }

  if (gvk === 'apps-v1-Deployment' || gvk === 'apps-v1-StatefulSet') {
    if (type === 'DeploymentStatus') {
      const deploymentStatus = []
      deploymentStatus.push(detail.status)
      return {
        type: 'DeploymentStatus',
        deploymentStatus
      }
    }
    if (type === 'StatefulSetStatus') {
      const statefulSetStatus = []
      statefulSetStatus.push(detail.status)
      return {
        type: 'StatefulSetStatus',
        statefulSetStatus
      }
    }
    if (type === 'Containers') {
      if (detail.spec.template) {
        return detail.spec.template.spec.containers.length
      } else {
        return 0
      }

    }
  }

  if (gvk === 'core-v1-Node') {
    if (type === 'taints') {
      return {
        type: 'taints',
        appInfo: detail.spec.taints
      }
    }
    if (type === 'cpu/memory') {
      return {
        type: 'cpu/memory',
        name: detail.metadata.name
      }
    }
  }

  if (gvk === 'openx.neverdown.io-v1-Affinity') {
    return detail.spec.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.length
  }

  if (gvk === 'openx.neverdown.io-v1-Mysql' || gvk === 'openx.neverdown.io-v1-Redis') {
    if (type === 'master/slave') {
      return {
        type: 'applications',
        appInfo: [{ appName: 'master', replicas: detail.spec.master.replicas },
        { appName: 'slave', replicas: detail.spec.slave.replicas }]
      }
    }
    if (type === 'init') {
      const appInfo = [{ appName: 'master', replicas: detail.spec.master.replicas },
      { appName: 'slave', replicas: detail.spec.slave.replicas }]
      return {
        type: 'init',
        appInfo
      }
    }
    if (type === 'StatefulSetStatus') {
      const statefulSetStatus = []
      statefulSetStatus.push(detail.status.master, detail.status.slave)
      return {
        type: 'StatefulSetStatus',
        statefulSetStatus
      }
    }
  }

  if (gvk === 'openx.neverdown.io-v1-Openx') {
    if (type === 'applications') {
      const appInfo = []
      let appNameList = []
      if (detail.spec.applications) {
        for (let app of detail.spec.applications) {
          appNameList.push(app.appName)
        }
        appNameList.sort()
        for (let name of appNameList) {
          const appIndex = detail.spec.applications.findIndex((ap: any) => {
            return ap.appName === name
          })
          appInfo.push(detail.spec.applications[appIndex])
        }
      }

      return {
        type: 'applications',
        appInfo
      }
    }
    if (type === 'init') {
      const appInfo = []
      let appNameList = []
      if (detail.spec.applications) {
        for (let app of detail.spec.applications) {
          appNameList.push(app.appName)
        }
        appNameList.sort()
        for (let name of appNameList) {
          const appIndex = detail.spec.applications.findIndex((ap: any) => {
            return ap.appName === name
          })
          appInfo.push(detail.spec.applications[appIndex])
        }
      }
      return {
        type: 'init',
        appInfo
      }
    }
    if (type === 'DeploymentStatus') {
      const deploymentStatus = []
      for (let appName in detail.status.items) {
        deploymentStatus.push(detail.status.items[appName].deploymentStatus)
      }
      return {
        type: 'DeploymentStatus',
        deploymentStatus
      }
    }
    if (type === 'HPAStatus') {
      const deploymentStatus = []
      for (let appName in detail.status.items) {
        const appIndex = detail.spec.applications.findIndex((app: any) => {
          let metadataName = detail.metadata.name
          return `${metadataName}-${app.appName}` === appName
        })
        if (appIndex >= 0 && detail.spec.applications[appIndex].horizontalPodAutoscalerSpec) {
          deploymentStatus.push(detail.status.items[appName].horizontalPodAutoscalerStatus)
        }
      }
      return {
        type: 'HPAStatus',
        deploymentStatus
      }
    }
  }
  if (gvk === 'openx.neverdown.io-v1-Etcd') {
    if (type === 'StatefulSetStatus') {
      const statefulSetStatus = []
      statefulSetStatus.push(detail.status)
      return {
        type: 'StatefulSetStatus',
        statefulSetStatus
      }
    }
  }

  if (gvk === 'core-v1-Event') {
    if (type === 'message') {
      return detail.message
    }
    if (type === 'type') {
      return {
        type: 'type',
        appInfo: detail.type
      }
    }
  }

  if (gvk === 'core-v1-ConfigMap') {
    if (type === 'config') {
      const configKey = []
      for (let ckey in detail.data) {
        configKey.push(ckey)
      }
      return configKey.join(', ')
    }
  }

  if (gvk === 'core-v1-Namespace' || gvk === 'core-v1-PersistentVolume') {
    if (type === 'phase') {
      return detail.status.phase
    }
  }

  if (gvk === 'core-v1-Pod') {
    if (type === 'phase') {
      return detail.status.phase
    }
    if (type === 'podIP') {
      return detail.status.podIP
    }
    if (type === 'nodeName') {
      return detail.spec.nodeName
    }
    if (type === 'cpu/memory') {
      return {
        type: 'cpu/memory',
        name: detail.metadata.name
      }
    }
    if (type === 'volume.timezone') {
      let lt = ''
      if (detail.spec.containers) {
        if (detail.spec.containers[0]?.volumeMounts) {
          let vmIndex = detail.spec.containers[0]?.volumeMounts.findIndex((mount: any) => {
            return mount.mountPath === '/etc/localtime'
          })
          if (vmIndex >= 0) {
            lt = detail.spec.containers[0].volumeMounts[vmIndex].subPath
          }
        }
      }
      return lt
    }
    if (type === 'container') {
      return {
        type: 'container',
        containerStatuses: detail.status.containerStatuses
      }
    }
  }

  if (gvk === 'core-v1-Service') {
    if (type === 'clusterIP') {
      return detail.spec.clusterIP
    }
    if (type === 'type') {
      return detail.spec.type
    }
    if (type === 'ports') {
      const portsStr = []
      for (let port of detail.spec.ports) {
        portsStr.push(`${port.port}:${port.nodePort}/${port.protocol}`)
      }
      return portsStr.join(', ')
    }
  }
}

