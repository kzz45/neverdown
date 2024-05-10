package app

import (
	"fmt"
	"net/http"

	"github.com/kzz45/neverdown/pkg/openx/controller/etcd"
	"github.com/kzz45/neverdown/pkg/openx/controller/loadbalancer"
	"github.com/kzz45/neverdown/pkg/openx/controller/mysql"
	"github.com/kzz45/neverdown/pkg/openx/controller/openx"
	"github.com/kzz45/neverdown/pkg/openx/controller/redis"
	"github.com/kzz45/neverdown/pkg/openx/controller/volcloadbalancer"
)

func startLoadBalancerController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "mysqls"}] {
	//	return nil, false, nil
	//}
	dc, err := loadbalancer.NewLoadBalancerController(
		ctx.OpenXInformerFactory.Openx().V1().AliyunLoadBalancers(),
		ctx.OpenXInformerFactory.Openx().V1().AliyunAccessControls(),
		ctx.ClientBuilder.OpenxClientOrDie("loadbalancer-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating loadbalancer controller: %v", err)
	}
	ctx.loadBalancher = dc
	return nil, true, nil
}

func startVolcLoadBalancerController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "mysqls"}] {
	//	return nil, false, nil
	//}
	dc, err := volcloadbalancer.NewVolcLoadBalancerController(
		ctx.OpenXInformerFactory.Openx().V1().VolcLoadBalancers(),
		ctx.OpenXInformerFactory.Openx().V1().VolcAccessControls(),
		ctx.ClientBuilder.OpenxClientOrDie("volcloadbalancer-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating loadbalancer controller: %v", err)
	}
	ctx.loadBalancher = dc
	return nil, true, nil
}

func startMysqlController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "mysqls"}] {
	//	return nil, false, nil
	//}
	dc, err := mysql.NewMysqlController(
		ctx.loadBalancher,
		ctx.OpenXInformerFactory.Openx().V1().Mysqls(),
		ctx.InformerFactory.Core().V1().Services(),
		ctx.InformerFactory.Apps().V1().StatefulSets(),
		ctx.ClientBuilder.ClientOrDie("mysql-controller-base"),
		ctx.ClientBuilder.OpenxClientOrDie("mysql-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating Mysql controller: %v", err)
	}
	go dc.Run(1, ctx.Stop)
	return nil, true, nil
}

func starRedisController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "mysqls"}] {
	//	return nil, false, nil
	//}
	dc, err := redis.NewRedisController(
		ctx.loadBalancher,
		ctx.OpenXInformerFactory.Openx().V1().Redises(),
		ctx.InformerFactory.Core().V1().Services(),
		ctx.InformerFactory.Apps().V1().StatefulSets(),
		ctx.ClientBuilder.ClientOrDie("redis-controller-base"),
		ctx.ClientBuilder.OpenxClientOrDie("redis-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating Redis controller: %v", err)
	}
	go dc.Run(1, ctx.Stop)
	return nil, true, nil
}

func starOpenxController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "mysqls"}] {
	//	return nil, false, nil
	//}
	tc, err := openx.NewOpenxController(
		ctx.jingx,
		ctx.loadBalancher,
		ctx.OpenXInformerFactory.Openx().V1().Openxes(),
		ctx.InformerFactory.Core().V1().Services(),
		ctx.InformerFactory.Apps().V1().Deployments(),
		ctx.InformerFactory.Autoscaling().V2().HorizontalPodAutoscalers(),
		ctx.ClientBuilder.ClientOrDie("openx-controller-base"),
		ctx.ClientBuilder.OpenxClientOrDie("openx-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating openx controller: %v", err)
	}
	go tc.Run(1, ctx.Stop)
	return nil, true, nil
}

func startNewEtcdController(ctx *ControllerContext) (http.Handler, bool, error) {
	//if !ctx.AvailableResources[schema.GroupVersionResource{Group: "openx", Version: "v1", Resource: "etcds"}] {
	//	return nil, false, nil
	//}
	dc, err := etcd.NewEtcdController(
		ctx.OpenXInformerFactory.Openx().V1().Etcds(),
		ctx.InformerFactory.Core().V1().Services(),
		ctx.InformerFactory.Apps().V1().StatefulSets(),
		ctx.ClientBuilder.ClientOrDie("etcd-controller-base"),
		ctx.ClientBuilder.OpenxClientOrDie("etcd-controller-extension"),
	)
	if err != nil {
		return nil, true, fmt.Errorf("error creating Mysql controller: %v", err)
	}
	go dc.Run(1, ctx.Stop)
	return nil, true, nil
}
