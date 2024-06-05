package openx

import (
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CloudLoadBalancerObject struct {
	Key   string `protobuf:"bytes,1,opt,name=key"`
	Value string `protobuf:"bytes,2,opt,name=value"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type LoadBalancer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              LoadBalancerSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// 负载均衡
type LoadBalancerSpec struct {
	Instance          CloudLoadBalancerObject `protobuf:"bytes,1,opt,name=instance"`
	OverrideListeners CloudLoadBalancerObject `protobuf:"bytes,2,opt,name=overrideListeners"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []LoadBalancer `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type AccessControl struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              AccessControlSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

// 访问控制
type AccessControlSpec struct {
	Instance CloudLoadBalancerObject `protobuf:"bytes,1,opt,name=instance"`
	Status   CloudLoadBalancerObject `protobuf:"bytes,2,opt,name=status"`
	Type     CloudLoadBalancerObject `protobuf:"bytes,3,opt,name=type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AccessControlList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []AccessControl `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type CloudLoadBalancerStatus string

const (
	CloudLoadBalancerOn  CloudLoadBalancerStatus = "on"
	CloudLoadBalancerOff CloudLoadBalancerStatus = "off"
)

type CloudSLB struct {
	LoadBalancerId    string                  `protobuf:"bytes,1,opt,name=loadBalancerId"`
	OverrideListeners bool                    `protobuf:"varint,2,opt,name=overrideListeners"`
	Status            CloudLoadBalancerStatus `protobuf:"bytes,3,opt,name=status,casttype=CloudLoadBalancerStatus"`
	AccessControlId   string                  `protobuf:"bytes,4,opt,name=accessControlId"`
}

// CloudNetworkConfig was the config for the cloud service provider, e.g. Aliyun
type CloudNetworkConfig struct {
	CloudSLB *CloudSLB `protobuf:"bytes,1,opt,name=cloudSLB"`
}

type ClusterRole string

const (
	ClusterRoleMaster ClusterRole = "master"
	ClusterRoleSlave  ClusterRole = "slave"
)

type MysqlConfig struct {
	Role               ClusterRole        `protobuf:"bytes,1,opt,name=role,casttype=ClusterRole"`
	MysqlServerConfig  MysqlServerConfig  `protobuf:"bytes,2,opt,name=mysqlServerConfig"`
	Pod                *corev1.Pod        `protobuf:"bytes,3,opt,name=pod"`
	Service            *corev1.Service    `protobuf:"bytes,4,opt,name=service"`
	PersistentStorage  PersistentStorage  `protobuf:"bytes,5,opt,name=persistentStorage"`
	CloudNetworkConfig CloudNetworkConfig `protobuf:"bytes,6,opt,name=cloudNetworkConfig"`
	Replicas           *int32             `protobuf:"varint,7,opt,name=replicas"`
}

// PersistentStorage will provide a persistent volume path, such as mount a Nas.
type PersistentStorage struct {
	StorageVolumePath string `protobuf:"bytes,1,opt,name=storageVolumePath"`
}

type MysqlServerConfig struct {
	ServerId    *int32 `json:"server_id" protobuf:"varint,1,opt,name=server_id,json=serverId"`
	Host        string `json:"host" protobuf:"bytes,2,opt,name=host"`
	User        string `json:"user" protobuf:"bytes,3,opt,name=user"`
	Password    string `json:"password" protobuf:"bytes,4,opt,name=password"`
	LogFile     string `json:"log_file" protobuf:"bytes,5,opt,name=log_file,json=logFile"`
	LogPosition string `json:"log_position" protobuf:"bytes,6,opt,name=log_position,json=logPosition"`
}

type MysqlSpec struct {
	Master MysqlConfig `protobuf:"bytes,1,opt,name=master"`
	Slave  MysqlConfig `protobuf:"bytes,2,opt,name=slave"`
}

type ClusterStatus struct {
	Master appsv1.StatefulSetStatus `json:"master,omitempty" protobuf:"bytes,1,opt,name=master"`
	Slave  appsv1.StatefulSetStatus `json:"slave,omitempty" protobuf:"bytes,2,opt,name=slave"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

type Mysql struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              MysqlSpec     `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            ClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MysqlList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Mysql `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              RedisSpec     `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            ClusterStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type RedisSpec struct {
	Master RedisConfig `protobuf:"bytes,1,opt,name=master"`
	Slave  RedisConfig `protobuf:"bytes,2,opt,name=slave"`
}

type RedisConfig struct {
	Role               ClusterRole        `protobuf:"bytes,1,opt,name=role,casttype=ClusterRole"`
	Pod                *corev1.Pod        `protobuf:"bytes,2,opt,name=pod"`
	Service            *corev1.Service    `protobuf:"bytes,3,opt,name=service"`
	PersistentStorage  PersistentStorage  `protobuf:"bytes,4,opt,name=persistentStorage"`
	CloudNetworkConfig CloudNetworkConfig `protobuf:"bytes,5,opt,name=cloudNetworkConfig"`
	Replicas           *int32             `protobuf:"varint,6,opt,name=replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Redis `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Openx describes a Openx resource
type Openx struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              OpenxSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            OpenxStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

type OpenxStatus struct {
	Items map[string]AppStatus `protobuf:"bytes,1,rep,name=items"`
}

type AppStatus struct {
	DeploymentStatus              appsv1.DeploymentStatus                     `protobuf:"bytes,1,opt,name=deploymentStatus"`
	HorizontalPodAutoscalerStatus autoscalingv2.HorizontalPodAutoscalerStatus `protobuf:"bytes,2,opt,name=horizontalPodAutoscalerStatus"`
}

// OpenxSpec is the spec for a Openx resource
type OpenxSpec struct {
	Applications []App `protobuf:"bytes,1,rep,name=applications"`
}

type WatchPolicy string

const (
	WatchPolicyManual         WatchPolicy = "manual"           // 手动升级
	WatchPolicyInPlaceUpgrade WatchPolicy = "in-place-upgrade" // 原地升级
	WatchPolicyRollingUpgrade WatchPolicy = "rolling-upgrade"  // 滚动升级
)

type App struct {
	AppName                     string                                     `protobuf:"bytes,1,opt,name=appName"`
	Pod                         *corev1.Pod                                `protobuf:"bytes,2,opt,name=pod"`
	Service                     *corev1.Service                            `protobuf:"bytes,3,opt,name=service"`
	PersistentStorage           PersistentStorage                          `protobuf:"bytes,4,opt,name=persistentStorage"`
	CloudNetworkConfig          CloudNetworkConfig                         `protobuf:"bytes,5,opt,name=cloudNetworkConfig"`
	Replicas                    *int32                                     `protobuf:"varint,6,opt,name=replicas"`
	WatchPolicy                 WatchPolicy                                `protobuf:"bytes,7,opt,name=watchPolicy,casttype=WatchPolicy"`
	HorizontalPodAutoscalerSpec *autoscalingv2.HorizontalPodAutoscalerSpec `protobuf:"bytes,8,opt,name=horizontalPodAutoscalerSpec"`
	ExtensionService            *corev1.Service                            `protobuf:"bytes,9,opt,name=extensionService"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// OpenxList is a list of Openx resources

type OpenxList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Openx `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type EtcdSpec struct {
	Replicas          *int32            `protobuf:"varint,1,opt,name=replicas"`
	PersistentStorage PersistentStorage `protobuf:"bytes,2,opt,name=persistentStorage"`
	Pod               *corev1.Pod       `protobuf:"bytes,3,opt,name=pod"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

type Etcd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              EtcdSpec                 `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status            appsv1.StatefulSetStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EtcdList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Etcd `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type Affinity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              AffinitySpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type AffinitySpec struct {
	Affinity *corev1.Affinity `protobuf:"bytes,1,opt,name=affinity"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type AffinityList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Affinity `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type Toleration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              TolerationSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type TolerationSpec struct {
	Toleration corev1.Toleration `protobuf:"bytes,1,opt,name=toleration"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TolerationList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Toleration `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:nostatus

type NodeSelector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec              NodeSelectorSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type NodeSelectorSpec struct {
	NodeSelector map[string]string `protobuf:"bytes,1,rep,name=nodeSelector"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeSelectorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []NodeSelector `json:"items" protobuf:"bytes,2,rep,name=items"`
}
