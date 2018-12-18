package manager

import (
	"mcloud/ecs.v2/model"
	"time"
	ali "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	publicModel "mcloud/public.v2/model"
)

func resolveAliInstance(ecs *model.Ecs, instance *ali.Instance, accountName string, accountId int64) {
	publicIp := ""
	privateIp := ""
	if len(instance.PublicIpAddress.IpAddress) > 0 {
		publicIp = instance.PublicIpAddress.IpAddress[0]
	} else {
		publicIp = instance.EipAddress.IpAddress
	}

	if len(instance.InnerIpAddress.IpAddress) > 0 {
		privateIp = instance.InnerIpAddress.IpAddress[0]
	} else {
		if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
			privateIp = instance.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
	}
	ecs.Platform = publicModel.PlatformAliCloud
	ecs.InstanceId = instance.InstanceId
	ecs.InstanceName = instance.InstanceName
	ecs.PublicIp = publicIp
	ecs.PrivateIp = privateIp
	ecs.OsType = instance.OSType
	ecs.Cpu = instance.Cpu
	ecs.Mem = instance.Memory
	ecs.InstanceType = instance.InstanceType
	ecs.ZoneId = instance.ZoneId
	ecs.AccountName = accountName
	ecs.AccountId = accountId
	ecs.RegionId = instance.RegionId
	ecs.CreatedTime = time.Now()
	ecs.UpdatedTime = time.Now()
}

func resolveTencentInstance(ecs *model.Ecs, instance *tencent.Instance, accountName string, accountId int64, regionId string) {
	publicIp := ""
	privateIp := ""
	if len(instance.PublicIpAddresses) > 0 {
		publicIp = *instance.PublicIpAddresses[0]
	}

	if len(instance.PrivateIpAddresses) > 0 {
		privateIp = *instance.PrivateIpAddresses[0]
	} else {
		if len(instance.VirtualPrivateCloud.PrivateIpAddresses) > 0 {
			privateIp = *instance.VirtualPrivateCloud.PrivateIpAddresses[0]
		}
	}

	ecs.Platform = publicModel.PlatformTencentCloud
	ecs.InstanceId = *instance.InstanceId
	ecs.InstanceName = *instance.InstanceName
	ecs.PublicIp = publicIp
	ecs.PrivateIp = privateIp
	ecs.OsType = *instance.OsName
	ecs.Cpu = int(*instance.CPU)
	ecs.Mem = int(*instance.Memory)
	ecs.InstanceType = *instance.InstanceType
	ecs.ZoneId = *instance.Placement.Zone
	ecs.AccountName = accountName
	ecs.AccountId = accountId
	ecs.RegionId = regionId
	ecs.CreatedTime = time.Now()
	ecs.UpdatedTime = time.Now()
}

func resolveAliInstanceAttribute(ecs *model.Ecs, resp *ali.DescribeInstanceAttributeResponse) (err error) {
	publicIp := ""
	privateIp := ""
	if len(resp.PublicIpAddress.IpAddress) > 0 {
		publicIp = resp.PublicIpAddress.IpAddress[0]
	} else {
		publicIp = resp.EipAddress.IpAddress
	}

	if len(resp.InnerIpAddress.IpAddress) > 0 {
		privateIp = resp.InnerIpAddress.IpAddress[0]
	} else {
		if len(resp.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
			privateIp = resp.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
	}

	ecs.Platform = publicModel.PlatformAliCloud
	ecs.InstanceId = resp.InstanceId
	ecs.InstanceName = resp.InstanceName
	ecs.PublicIp = publicIp
	ecs.PrivateIp = privateIp
	ecs.Cpu = resp.Cpu
	ecs.Mem = resp.Memory
	ecs.InstanceType = resp.InstanceType
	ecs.ZoneId = resp.ZoneId
	ecs.RegionId = resp.RegionId
	ecs.Status = resolveAliStatus(resp.Status)
	ecs.UpdatedTime = time.Now()
	return
}

func resolveAliStatus(status string) model.EcsStatus {
	if status == "Pending" {
		return model.EcsStatusCreated
	} else if status == "Starting" {
		return model.EcsStatusStarting
	} else if status == "Running" {
		return model.EcsStatusRunning
	} else if status == "Stopping" {
		return model.EcsStatusStopping
	} else if status == "Stopped" {
		return model.EcsStatusStopped
	} else {
		return model.EcsStatusNon
	}

}

func resolveTencentInstanceAttribute(ecs *model.Ecs, respInstance *tencent.Instance, regionId string) (err error) {
	publicIp := ""
	privateIp := ""
	if len(respInstance.PublicIpAddresses) > 0 {
		publicIp = *respInstance.PublicIpAddresses[0]
	}
	if len(respInstance.PrivateIpAddresses) > 0 {
		privateIp = *respInstance.PrivateIpAddresses[0]
	}
	ecs.Platform = publicModel.PlatformTencentCloud
	ecs.InstanceId = *respInstance.InstanceId
	ecs.InstanceName = *respInstance.InstanceName
	ecs.PublicIp = publicIp
	ecs.PrivateIp = privateIp
	ecs.Cpu = int(*respInstance.CPU)
	ecs.Mem = int(*respInstance.Memory)
	ecs.InstanceType = *respInstance.InstanceType
	ecs.ZoneId = *respInstance.Placement.Zone
	ecs.RegionId = regionId
	ecs.Status = resolveTencentStatus(*respInstance.InstanceState)
	ecs.UpdatedTime = time.Now()
	return
}

func resolveTencentStatus(status string) model.EcsStatus {
	if status == "PENDING" {
		return model.EcsStatusCreated
	} else if status == "LAUNCH_FAILED" {
		return model.EcsStatusCreateFailed
	} else if status == "STARTING" {
		return model.EcsStatusStarting
	} else if status == "RUNNING" {
		return model.EcsStatusRunning
	} else if status == "STOPPING" {
		return model.EcsStatusStopping
	} else if status == "STOPPED" {
		return model.EcsStatusStopped
	} else if status == "REBOOTING" {
		return model.EcsStatusStarting
	} else if status == "SHUTDOWN" {
		return model.EcsStatusStopped
	} else if status == "TERMINATING" {
		return model.EcsStatusStopped
	} else {
		return model.EcsStatusNon
	}
}
