import moment from "moment";

export function sortByName(list) {
  let nameList = [];
  for (let item of list) {
    nameList.push(item.name);
  }
  nameList.sort();
  let resultList = [];
  for (let name of nameList) {
    const nameIndex = list.findIndex((val) => {
      return val.name === name;
    });
    resultList.push(list[nameIndex]);
  }
  return resultList;
}

export function timeFormat(ts) {
  return moment(ts * 1000).format("YYYY-MM-DD HH:mm:ss");
}
