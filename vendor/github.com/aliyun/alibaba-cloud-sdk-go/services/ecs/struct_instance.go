package ecs

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Instance is a nested struct in ecs response
type Instance struct {
	ImageId                 string                               `json:"ImageId" xml:"ImageId"`
	InstanceType            string                               `json:"InstanceType" xml:"InstanceType"`
	AutoReleaseTime         string                               `json:"AutoReleaseTime" xml:"AutoReleaseTime"`
	DeviceAvailable         bool                                 `json:"DeviceAvailable" xml:"DeviceAvailable"`
	InstanceNetworkType     string                               `json:"InstanceNetworkType" xml:"InstanceNetworkType"`
	LocalStorageAmount      int                                  `json:"LocalStorageAmount" xml:"LocalStorageAmount"`
	InstanceChargeType      string                               `json:"InstanceChargeType" xml:"InstanceChargeType"`
	ClusterId               string                               `json:"ClusterId" xml:"ClusterId"`
	InstanceName            string                               `json:"InstanceName" xml:"InstanceName"`
	CreditSpecification     string                               `json:"CreditSpecification" xml:"CreditSpecification"`
	GPUAmount               int                                  `json:"GPUAmount" xml:"GPUAmount"`
	StartTime               string                               `json:"StartTime" xml:"StartTime"`
	ZoneId                  string                               `json:"ZoneId" xml:"ZoneId"`
	InternetChargeType      string                               `json:"InternetChargeType" xml:"InternetChargeType"`
	InternetMaxBandwidthIn  int                                  `json:"InternetMaxBandwidthIn" xml:"InternetMaxBandwidthIn"`
	HostName                string                               `json:"HostName" xml:"HostName"`
	Cpu                     int                                  `json:"Cpu" xml:"Cpu"`
	Status                  string                               `json:"Status" xml:"Status"`
	SpotPriceLimit          float64                              `json:"SpotPriceLimit" xml:"SpotPriceLimit"`
	OSName                  string                               `json:"OSName" xml:"OSName"`
	SerialNumber            string                               `json:"SerialNumber" xml:"SerialNumber"`
	RegionId                string                               `json:"RegionId" xml:"RegionId"`
	InternetMaxBandwidthOut int                                  `json:"InternetMaxBandwidthOut" xml:"InternetMaxBandwidthOut"`
	IoOptimized             bool                                 `json:"IoOptimized" xml:"IoOptimized"`
	ResourceGroupId         string                               `json:"ResourceGroupId" xml:"ResourceGroupId"`
	InstanceTypeFamily      string                               `json:"InstanceTypeFamily" xml:"InstanceTypeFamily"`
	InstanceId              string                               `json:"InstanceId" xml:"InstanceId"`
	GPUSpec                 string                               `json:"GPUSpec" xml:"GPUSpec"`
	Description             string                               `json:"Description" xml:"Description"`
	Recyclable              bool                                 `json:"Recyclable" xml:"Recyclable"`
	SaleCycle               string                               `json:"SaleCycle" xml:"SaleCycle"`
	ExpiredTime             string                               `json:"ExpiredTime" xml:"ExpiredTime"`
	OSType                  string                               `json:"OSType" xml:"OSType"`
	Memory                  int                                  `json:"Memory" xml:"Memory"`
	CreationTime            string                               `json:"CreationTime" xml:"CreationTime"`
	KeyPairName             string                               `json:"KeyPairName" xml:"KeyPairName"`
	HpcClusterId            string                               `json:"HpcClusterId" xml:"HpcClusterId"`
	LocalStorageCapacity    int                                  `json:"LocalStorageCapacity" xml:"LocalStorageCapacity"`
	VlanId                  string                               `json:"VlanId" xml:"VlanId"`
	StoppedMode             string                               `json:"StoppedMode" xml:"StoppedMode"`
	SpotStrategy            string                               `json:"SpotStrategy" xml:"SpotStrategy"`
	SecurityGroupIds        SecurityGroupIdsInDescribeInstances  `json:"SecurityGroupIds" xml:"SecurityGroupIds"`
	InnerIpAddress          InnerIpAddressInDescribeInstances    `json:"InnerIpAddress" xml:"InnerIpAddress"`
	PublicIpAddress         PublicIpAddressInDescribeInstances   `json:"PublicIpAddress" xml:"PublicIpAddress"`
	RdmaIpAddress           RdmaIpAddress                        `json:"RdmaIpAddress" xml:"RdmaIpAddress"`
	EipAddress              EipAddress                           `json:"EipAddress" xml:"EipAddress"`
	DedicatedHostAttribute  DedicatedHostAttribute               `json:"DedicatedHostAttribute" xml:"DedicatedHostAttribute"`
	VpcAttributes           VpcAttributes                        `json:"VpcAttributes" xml:"VpcAttributes"`
	NetworkInterfaces       NetworkInterfacesInDescribeInstances `json:"NetworkInterfaces" xml:"NetworkInterfaces"`
	OperationLocks          OperationLocksInDescribeInstances    `json:"OperationLocks" xml:"OperationLocks"`
	Tags                    TagsInDescribeInstances              `json:"Tags" xml:"Tags"`
}
