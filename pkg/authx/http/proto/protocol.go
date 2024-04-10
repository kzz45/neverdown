package proto

type ServiceRoute string

const (
	Login                    ServiceRoute = "/login"
	AppsList                 ServiceRoute = "/apps/list"
	AppsCreate               ServiceRoute = "/apps/create"
	AppsDelete               ServiceRoute = "/apps/delete"
	AppsUpdate               ServiceRoute = "/apps/update"
	RbacServiceAccountList   ServiceRoute = "/rbacaccount/list"
	RbacServiceAccountCreate ServiceRoute = "/rbacaccount/create"
	RbacServiceAccountDelete ServiceRoute = "/rbacaccount/delete"
	RbacServiceAccountUpdate ServiceRoute = "/rbacaccount/update"
	AppServiceAccountList    ServiceRoute = "/appaccount/list"
	AppServiceAccountCreate  ServiceRoute = "/appaccount/create"
	AppServiceAccountDelete  ServiceRoute = "/appaccount/delete"
	AppServiceAccountUpdate  ServiceRoute = "/appaccount/update"
	ClusterRoleList          ServiceRoute = "/clusterrole/list"
	ClusterRoleCreate        ServiceRoute = "/clusterrole/create"
	ClusterRoleDelete        ServiceRoute = "/clusterrole/delete"
	ClusterRoleUpdate        ServiceRoute = "/clusterrole/update"
	GroupVersionKindRuleList ServiceRoute = "/kind/list"
)

type Request struct {
	ServiceRoute ServiceRoute `protobuf:"bytes,1,opt,name=serviceRoute,casttype=ServiceRoute"`
	Data         []byte       `protobuf:"bytes,2,opt,name=data"`
}

type Response struct {
	ServiceRoute ServiceRoute `protobuf:"bytes,1,opt,name=serviceRoute,casttype=ServiceRoute"`
	Code         int32        `protobuf:"varint,2,opt,name=code"`
	Message      string       `protobuf:"bytes,3,opt,name=message"`
	Data         []byte       `protobuf:"bytes,4,opt,name=data"`
}

type Context struct {
	Token    string `protobuf:"bytes,1,opt,name=token"`
	IsAdmin  bool   `protobuf:"varint,2,opt,name=isAdmin"`
	ExpireAt int32  `protobuf:"varint,3,opt,name=expireAt"`
}
