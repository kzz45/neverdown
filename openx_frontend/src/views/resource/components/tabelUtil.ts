import { cloneDeep } from 'lodash'

export function rowEdit(row: any, isSet: boolean, allData: any, saveObj: any) {
  for (let item of allData) {
    if (!isSet && item.isSet) {
      return
    }
  }
  saveObj.value = cloneDeep(row)
  row.isSet = !Boolean(row.isSet)
}

export function rowDelete(allData: any, delIndex: any) {
  allData.splice(delIndex, 1)
}

export function rowCancelEdit(row: any, saveObj: any) {
  const copyData: any = saveObj.value
  for (let dataIndex in copyData) {
    row[dataIndex] = copyData[dataIndex]
  }
  row.isSet = false
}

export function initLimits(containerOne: any) {
  if (!containerOne.resources.limits.cpu) {
    containerOne.resources.limits.cpu = {
      string: '0'
    }
  }
  if (!containerOne.resources.limits.memory) {
    containerOne.resources.limits.memory = {
      string: '0'
    }
  }
  if (!containerOne.resources.requests.cpu) {
    containerOne.resources.requests.cpu = {
      string: '0'
    }
  }
  if (!containerOne.resources.requests.memory) {
    containerOne.resources.requests.memory = {
      string: '0'
    }
  }
}

import proto from '../../../../proto/proto'
const protoApi: any = proto.k8s.io.api.core.v1
export function initPod(pod: any) {
  if (!pod.spec.securityContext) {
    pod.spec.securityContext = protoApi.PodSecurityContext.create()
  }
  for (let securIndex in pod.spec.securityContext) {
    if (pod.spec.securityContext[securIndex] === null) {
      if (securIndex === 'seLinuxOptions') {
        pod.spec.securityContext[securIndex] = protoApi.SELinuxOptions.create()
      }
      if (securIndex === 'windowsOptions') {
        pod.spec.securityContext[securIndex] = protoApi.WindowsSecurityContextOptions.create()
      }
      if (securIndex === 'seccompProfile') {
        pod.spec.securityContext[securIndex] = protoApi.SeccompProfile.create()
        pod.spec.securityContext[securIndex].type = 'Localhost'
      }
    }
  }
}
export function initCon(con: any) {
  if (!con.securityContext) {
    let newSecur = protoApi.SecurityContext.create()
    con.securityContext = newSecur
  }
  for (let securIndex in con.securityContext) {
    if (con.securityContext[securIndex] === null) {
      if (securIndex === 'capabilities') {
        con.securityContext[securIndex] = protoApi.Capabilities.create()
      }
      if (securIndex === 'seLinuxOptions') {
        con.securityContext[securIndex] = protoApi.SELinuxOptions.create()
      }
      if (securIndex === 'windowsOptions') {
        con.securityContext[securIndex] = protoApi.WindowsSecurityContextOptions.create()
      }
      if (securIndex === 'seccompProfile') {
        con.securityContext[securIndex] = protoApi.SeccompProfile.create()
        con.securityContext[securIndex].type = 'Localhost'
      }
    }
  }
}

export function clickSilder(docId: string) {
  let el: any = document.querySelector(`#${docId}`);

  console.log('clickSilder', docId, el.offsetTop)
  // chrome
  document.body.scrollTop = el.offsetTop;
  // firefox
  document.documentElement.scrollTop = el.offsetTop
}

export function goToNav(elementId: string) {

  document.getElementById('id-' + elementId)?.focus()
  document.getElementById('id-' + elementId)?.scrollIntoView({
    behavior: 'smooth'
  })
}
