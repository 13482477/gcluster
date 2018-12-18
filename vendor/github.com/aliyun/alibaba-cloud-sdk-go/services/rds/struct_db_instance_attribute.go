package rds

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

// DBInstanceAttribute is a nested struct in rds response
type DBInstanceAttribute struct {
	Engine                            string                                             `json:"Engine" xml:"Engine"`
	TempUpgradeTimeStart              string                                             `json:"TempUpgradeTimeStart" xml:"TempUpgradeTimeStart"`
	TempUpgradeRecoveryTime           string                                             `json:"TempUpgradeRecoveryTime" xml:"TempUpgradeRecoveryTime"`
	TempUpgradeRecoveryMaxIOPS        string                                             `json:"TempUpgradeRecoveryMaxIOPS" xml:"TempUpgradeRecoveryMaxIOPS"`
	DBInstanceDiskUsed                string                                             `json:"DBInstanceDiskUsed" xml:"DBInstanceDiskUsed"`
	AdvancedFeatures                  string                                             `json:"AdvancedFeatures" xml:"AdvancedFeatures"`
	DBInstanceClass                   string                                             `json:"DBInstanceClass" xml:"DBInstanceClass"`
	DBInstanceNetType                 string                                             `json:"DBInstanceNetType" xml:"DBInstanceNetType"`
	VpcCloudInstanceId                string                                             `json:"VpcCloudInstanceId" xml:"VpcCloudInstanceId"`
	DBMaxQuantity                     int                                                `json:"DBMaxQuantity" xml:"DBMaxQuantity"`
	DBInstanceCPU                     string                                             `json:"DBInstanceCPU" xml:"DBInstanceCPU"`
	MaxConnections                    int                                                `json:"MaxConnections" xml:"MaxConnections"`
	IncrementSourceDBInstanceId       string                                             `json:"IncrementSourceDBInstanceId" xml:"IncrementSourceDBInstanceId"`
	InstanceNetworkType               string                                             `json:"InstanceNetworkType" xml:"InstanceNetworkType"`
	DBInstanceType                    string                                             `json:"DBInstanceType" xml:"DBInstanceType"`
	TempUpgradeRecoveryClass          string                                             `json:"TempUpgradeRecoveryClass" xml:"TempUpgradeRecoveryClass"`
	DBInstanceId                      string                                             `json:"DBInstanceId" xml:"DBInstanceId"`
	DBInstanceMemory                  int                                                `json:"DBInstanceMemory" xml:"DBInstanceMemory"`
	VpcId                             string                                             `json:"VpcId" xml:"VpcId"`
	DBInstanceStorageType             string                                             `json:"DBInstanceStorageType" xml:"DBInstanceStorageType"`
	SecurityIPList                    string                                             `json:"SecurityIPList" xml:"SecurityIPList"`
	LatestKernelVersion               string                                             `json:"LatestKernelVersion" xml:"LatestKernelVersion"`
	SupportUpgradeAccountType         string                                             `json:"SupportUpgradeAccountType" xml:"SupportUpgradeAccountType"`
	MaxIOPS                           int                                                `json:"MaxIOPS" xml:"MaxIOPS"`
	Tags                              string                                             `json:"Tags" xml:"Tags"`
	EngineVersion                     string                                             `json:"EngineVersion" xml:"EngineVersion"`
	MaintainTime                      string                                             `json:"MaintainTime" xml:"MaintainTime"`
	PayType                           string                                             `json:"PayType" xml:"PayType"`
	DBInstanceStorage                 int                                                `json:"DBInstanceStorage" xml:"DBInstanceStorage"`
	SupportCreateSuperAccount         string                                             `json:"SupportCreateSuperAccount" xml:"SupportCreateSuperAccount"`
	TempDBInstanceId                  string                                             `json:"TempDBInstanceId" xml:"TempDBInstanceId"`
	CurrentKernelVersion              string                                             `json:"CurrentKernelVersion" xml:"CurrentKernelVersion"`
	ZoneId                            string                                             `json:"ZoneId" xml:"ZoneId"`
	ConnectionMode                    string                                             `json:"ConnectionMode" xml:"ConnectionMode"`
	IPType                            string                                             `json:"IPType" xml:"IPType"`
	ReadonlyInstanceSQLDelayedTime    string                                             `json:"ReadonlyInstanceSQLDelayedTime" xml:"ReadonlyInstanceSQLDelayedTime"`
	LockMode                          string                                             `json:"LockMode" xml:"LockMode"`
	CanTempUpgrade                    bool                                               `json:"CanTempUpgrade" xml:"CanTempUpgrade"`
	LockReason                        string                                             `json:"LockReason" xml:"LockReason"`
	Category                          string                                             `json:"Category" xml:"Category"`
	GuardDBInstanceId                 string                                             `json:"GuardDBInstanceId" xml:"GuardDBInstanceId"`
	InsId                             int                                                `json:"InsId" xml:"InsId"`
	DBInstanceDescription             string                                             `json:"DBInstanceDescription" xml:"DBInstanceDescription"`
	AccountType                       string                                             `json:"AccountType" xml:"AccountType"`
	GuardDBInstanceName               string                                             `json:"GuardDBInstanceName" xml:"GuardDBInstanceName"`
	RegionId                          string                                             `json:"RegionId" xml:"RegionId"`
	ResourceGroupId                   string                                             `json:"ResourceGroupId" xml:"ResourceGroupId"`
	TempUpgradeTimeEnd                string                                             `json:"TempUpgradeTimeEnd" xml:"TempUpgradeTimeEnd"`
	ExpireTime                        string                                             `json:"ExpireTime" xml:"ExpireTime"`
	TempUpgradeRecoveryMemory         int                                                `json:"TempUpgradeRecoveryMemory" xml:"TempUpgradeRecoveryMemory"`
	AccountMaxQuantity                int                                                `json:"AccountMaxQuantity" xml:"AccountMaxQuantity"`
	TempUpgradeRecoveryMaxConnections string                                             `json:"TempUpgradeRecoveryMaxConnections" xml:"TempUpgradeRecoveryMaxConnections"`
	Port                              string                                             `json:"Port" xml:"Port"`
	VSwitchId                         string                                             `json:"VSwitchId" xml:"VSwitchId"`
	CreationTime                      string                                             `json:"CreationTime" xml:"CreationTime"`
	MasterInstanceId                  string                                             `json:"MasterInstanceId" xml:"MasterInstanceId"`
	SecurityIPMode                    string                                             `json:"SecurityIPMode" xml:"SecurityIPMode"`
	DBInstanceClassType               string                                             `json:"DBInstanceClassType" xml:"DBInstanceClassType"`
	ReadDelayTime                     string                                             `json:"ReadDelayTime" xml:"ReadDelayTime"`
	DBInstanceStatus                  string                                             `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	ReplicateId                       string                                             `json:"ReplicateId" xml:"ReplicateId"`
	ConnectionString                  string                                             `json:"ConnectionString" xml:"ConnectionString"`
	TempUpgradeRecoveryCpu            int                                                `json:"TempUpgradeRecoveryCpu" xml:"TempUpgradeRecoveryCpu"`
	AvailabilityValue                 string                                             `json:"AvailabilityValue" xml:"AvailabilityValue"`
	ReadOnlyDBInstanceIds             ReadOnlyDBInstanceIdsInDescribeDBInstanceAttribute `json:"ReadOnlyDBInstanceIds" xml:"ReadOnlyDBInstanceIds"`
}
