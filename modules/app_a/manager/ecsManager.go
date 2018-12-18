package manager

import (
	"context"
	"fmt"
	"time"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/utils"
	ali "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"mcloud/thirdas.v2"
	"mcloud/ecs.v2/model"
	"strings"
	"github.com/kataras/iris/core/errors"
	tencent "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	log "github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"mcloud/ecs.v2/http"
	"mcloud/public.v2/rpc"
	ecsRpc "mcloud/ecs.v2/rpc"
	publicModel "mcloud/public.v2/model"
	model2 "mcloud/mcloud.v1/model"
)

type EcsManager struct {
	Db      *gorm.DB
	Clients []*EcsMClient
}

type EcsMClient struct {
	AliClient     *ali.Client
	TencentClient *tencent.Client
	Platform      publicModel.Platform
	Region        string
	AccountName   string
	AccountId     int64
}

type EcsMInstance struct {
	aliInstance     ali.Instance
	tencentInstance tencent.Instance
	AccountName     string
	AccountId       int64
}

type EcsConfig struct {
	Ak     string
	Sk     string
	Region string
	Logger log.Logger
	DbUri  string
}

const (
	PrePaid  = "PrePaid"
	PostPaid = "PostPaid"
)

func (manager *EcsManager) StartMcloudManager() error {

	var clients []*EcsMClient
	as := make([]thirdas.ThirdAS, 0)
	manager.Db.Find(&as)
	for _, ins := range as {
		if regions, ok := publicModel.RegionMap[publicModel.Platform(ins.Platform)]; ok {
			for _, region := range regions {
				emc := new(EcsMClient)
				emc.Region = region
				emc.AccountId = int64(ins.ID)
				emc.AccountName = ins.AccountName
				emc.Platform = publicModel.Platform(ins.Platform)
				if ins.Platform == int(publicModel.PlatformAliCloud) {
					if client, err := ali.NewClientWithAccessKey(region, ins.AccessKey, ins.SecretKey); err != nil {
						log.WithError(err)
						continue
					} else {
						emc.AliClient = client
					}
				} else if ins.Platform == int(publicModel.PlatformTencentCloud) {
					if client, err := tencent.NewClientWithSecretId(ins.AccessKey, ins.SecretKey, region); err != nil {
						log.WithError(err)
						continue
					} else {
						emc.TencentClient = client
					}
				}
				clients = append(clients, emc)
			}
		}
	}

	manager.Clients = clients
	return nil
}

func (manager *EcsManager) BindGroup(ctx context.Context, input http.EcsBindGroupInput) error {
	ecs := model.Ecs{}
	if err := manager.Db.Where("id = ?", input.EcsId).First(&ecs).Error; err != nil {
		return err
	}
	ecs.GroupId = int64(input.GroupId)
	ecs.UpdatedTime = time.Now()
	if err := manager.Db.Save(&ecs).Error; err != nil {
		return err
	}
	return nil
}

func (manager *EcsManager) DescribeInstanceAttribute(ctx context.Context, input http.EcsDescribeInstanceAttrInput) (*model.Ecs, error) {
	instance := &model.Ecs{}
	if err := manager.Db.Where(&model.Ecs{
		Platform:   publicModel.Platform(input.Platform),
		InstanceId: input.InstanceId,
	}).First(instance).Error; err != nil {
		return nil, err
	}

	ak, sk := manager.GetAkSk(input.AccountId)
	if publicModel.Platform(input.Platform) == publicModel.PlatformAliCloud {
		if err := describeInstanceAttributeFromAli(instance, ak, sk, input.Region, input.InstanceId); err != nil {
			return nil, err
		}
	} else if publicModel.Platform(input.Platform) == publicModel.PlatformTencentCloud {
		if err := describeInstanceAttributeFromTencent(instance, ak, sk, input.Region, input.InstanceId); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("invalide platform, platform=%d", input.Platform))
	}

	if err := manager.Db.Save(instance).Error; err != nil {
		return nil, err
	}

	return instance, nil
}

func describeInstanceAttributeFromAli(instance *model.Ecs, ak, sk, region, instanceId string) (err error) {
	client, err := ali.NewClientWithAccessKey(region, ak, sk)
	if err != nil {
		return
	}

	req := ali.CreateDescribeInstanceAttributeRequest()
	req.InstanceId = instanceId
	if resp, err := client.DescribeInstanceAttribute(req); err != nil {
		return err
	} else {
		return resolveAliInstanceAttribute(instance, resp)
	}
}

func describeInstanceAttributeFromTencent(instance *model.Ecs, ak, sk, region, instanceId string) (err error) {
	client, err := tencent.NewClientWithSecretId(ak, sk, region)
	if err != nil {
		return
	}

	req := tencent.NewDescribeInstancesRequest()
	req.Limit = common.Int64Ptr(1)
	req.Filters = []*tencent.Filter{
		{
			Name:   common.StringPtr("instance-id"),
			Values: common.StringPtrs([]string{instanceId}),
		},
	}
	if resp, err := client.DescribeInstances(req); err != nil {
		return err
	} else if *resp.Response.TotalCount == 0 {
		return errors.New(fmt.Sprintf("cat not found ecs info from tencent, instanceId=%s", instanceId))
	} else {
		return resolveTencentInstanceAttribute(instance, resp.Response.InstanceSet[0], region)
	}
}

func (manager *EcsManager) ModifyInstanceType(ctx context.Context, input *http.ModifyInstanceTypeInput) error {
	ak, sk := manager.GetAkSk(input.AccountId)
	client, _ := ali.NewClientWithAccessKey(input.Region, ak, sk)
	req := ali.CreateModifyInstanceSpecRequest()
	req.InstanceId = input.InstanceId
	req.InstanceType = input.InstanceType
	_, err := client.ModifyInstanceSpec(req)
	if err != nil {
		return err
	}
	return nil
}

func (manager *EcsManager) StartInstance(ctx context.Context, input *http.EcsStartInput) (string, error) {
	input.Uuid = utils.GetUUIDV4()
	ak, sk := manager.GetAkSk(input.AccountId)
	startInstanceRequest := ali.CreateStartInstanceRequest()
	startInstanceRequest.InstanceId = input.InstanceId

	client, _ := ali.NewClientWithAccessKey(input.Region, ak, sk)
	resp, err := client.StartInstance(startInstanceRequest)
	if err != nil {
		log.WithError(err)
		return "", err
	} else {
		log.WithField(resp.RequestId, "success").Info()
		return resp.RequestId, nil
	}
}

func (manager *EcsManager) StopInstance(ctx context.Context, input *http.EcsStopInput) (string, error) {
	input.Uuid = utils.GetUUIDV4()
	ak, sk := manager.GetAkSk(input.AccountId)
	stopInstanceRequest := ali.CreateStopInstanceRequest()
	stopInstanceRequest.InstanceId = input.InstanceId

	client, _ := ali.NewClientWithAccessKey(input.Region, ak, sk)
	resp, err := client.StopInstance(stopInstanceRequest)
	if err != nil {
		log.WithError(err)
		return "", err
	} else {
		log.WithField(resp.RequestId, "success").Info()
		return resp.RequestId, nil
	}
}

func (manager *EcsManager) RebootInstance(ctx context.Context, input *http.EcsRebootInput) (string, error) {
	input.Uuid = utils.GetUUIDV4()
	ak, sk := manager.GetAkSk(input.AccountId)
	rebootInstanceRequest := ali.CreateRebootInstanceRequest()
	rebootInstanceRequest.InstanceId = input.InstanceId

	client, _ := ali.NewClientWithAccessKey(input.Region, ak, sk)
	resp, err := client.RebootInstance(rebootInstanceRequest)
	if err != nil {
		log.WithError(err)
		return "", err
	} else {
		log.WithField(resp.RequestId, "success").Info()
		return resp.RequestId, nil
	}
}

func (manager *EcsManager) GetDashboard(ctx context.Context, input *http.DashboardInput) (*http.Dashboard, error) {
	log.WithField("in", "getDashboard").Info()
	cpuMemList, err := manager.getCpuMemList(input)
	if err != nil {
		return &http.Dashboard{}, err
	}
	instanceTypeList, err := manager.getInstanceTypeList(input)
	if err != nil {
		return &http.Dashboard{}, err
	}
	return &http.Dashboard{
		CpuMemList:       cpuMemList,
		InstanceTypeList: instanceTypeList,
	}, nil

}

func (manager *EcsManager) getInstanceTypeList(input *http.DashboardInput) ([]*http.EcsCount, error) {
	ret := make([]*http.EcsCount, 0)
	sentence := manager.Db.Table("ecs_model").Select("instance_type as item, count(*) as count")
	if input.AccountId != 0 {
		sentence = sentence.Where("account_id = ?", input.AccountId)
	}
	if input.GroupId != 0 {
		sentence = sentence.Where("group_id = ?", input.GroupId)
	}

	if err := sentence.Group("instance_type").Scan(&ret).Error; err != nil {
		return []*http.EcsCount{}, err
	}
	return ret, nil
}

func (manager *EcsManager) getCpuMemList(input *http.DashboardInput) ([]*http.EcsCount, error) {
	ret := make([]*http.EcsCount, 0)
	sqlPrefix := "select concat(concat(cpu,'core '), concat(ROUND(mem/1024),'g'))as item, count(*) as count from ecs_model where 1=1 "
	sqlSuffix := " group by item"
	if input.AccountId != 0 {
		sqlPrefix = fmt.Sprintf("%s and account_id = %d", sqlPrefix, input.AccountId)
	}
	if input.GroupId != 0 {
		sqlPrefix = fmt.Sprintf("%s and group_id = %d", sqlPrefix, input.GroupId)
	}
	sql := sqlPrefix + sqlSuffix

	if err := manager.Db.Raw(sql).Scan(&ret).Error; err != nil {
		return []*http.EcsCount{}, err
	}
	return ret, nil
}

func (manager *EcsManager) Echo(_ context.Context, s string) (string, error) {
	return s, nil
}

func (manager *EcsManager) Create(ctx context.Context, input *http.EcsCreateInput) (string, error) {
	input.Uuid = utils.GetUUIDV4()
	ak, sk := manager.GetAkSk(input.AccountId)
	createInstanceRequest := ali.CreateCreateInstanceRequest()
	createInstanceRequest.ClientToken = input.Uuid
	createInstanceRequest.VSwitchId = input.VSwitchId
	createInstanceRequest.ImageId = input.ImageId
	createInstanceRequest.InstanceName = input.InstanceName
	createInstanceRequest.SecurityGroupId = input.SecurityGroupId
	createInstanceRequest.InstanceType = input.IType
	createInstanceRequest.ZoneId = input.ZoneId
	createInstanceRequest.SystemDiskSize = requests.NewInteger(input.SystemDiskSize)
	if input.InstanceChargeType == PrePaid {
		createInstanceRequest.InstanceChargeType = input.InstanceChargeType
		createInstanceRequest.Period = requests.NewInteger(1)
		createInstanceRequest.AutoRenew = requests.NewBoolean(true)
		createInstanceRequest.AutoRenewPeriod = requests.NewInteger(1)
	}
	client, _ := ali.NewClientWithAccessKey(input.Region, ak, sk)
	response, err := client.CreateInstance(createInstanceRequest)
	if err != nil {
		log.WithError(err)
		return "", err
	} else {
		newInstance := model.Ecs{
			InstanceId:   response.InstanceId,
			InstanceName: input.InstanceName,
			AccountId:    int64(input.AccountId),
			RegionId:     input.Region,
			BaseModel: publicModel.BaseModel{
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
			},
		}
		if err := manager.Db.Create(&newInstance).Error; err != nil {
			log.WithError(err)
			return "", err
		} else {
			time.Sleep(time.Second * 60)
			manager.StartInstance(ctx, &http.EcsStartInput{
				AccountId:  input.AccountId,
				Region:     input.Region,
				InstanceId: response.InstanceId,
			})
		}
	}
	return response.RequestId, nil
}

func (manager *EcsManager) describeInstancesFromRemoteServer() (ecss []*model.Ecs) {
	for _, cli := range manager.Clients {
		log.WithFields(log.Fields{"平台 ====>>>": cli.Platform, "区域 =====>>>>>": cli.Region, "现在是 ====>>>": cli.AccountName}).Info()
		if cli.Platform == publicModel.PlatformAliCloud {
			ecss = append(ecss, manager.describeInstancesFromAliCloud(cli)...)
		} else if cli.Platform == publicModel.PlatformTencentCloud {
			ecss = append(ecss, manager.describeInstancesFromTencentCloud(cli)...)
		} else {
			log.WithField("platform", cli.Platform).Error("invalid platform")
		}
	}
	return
}

func (manager *EcsManager) describeInstancesFromAliCloud(cli *EcsMClient) (ecss []*model.Ecs) {
	req := ali.CreateDescribeInstancesRequest()
	req.RegionId = cli.Region
	req.PageNumber = requests.NewInteger(1)
	req.PageSize = requests.NewInteger(100)
	resp, err := cli.AliClient.DescribeInstances(req)
	if err != nil {
		log.WithField("account_name", cli.AccountName).WithError(err)
		return
	}
	for _, ins := range resp.Instances.Instance {
		log.WithFields(log.Fields{"account_name": cli.AccountName, "instance_name": ins.InstanceName}).Info()
		ecs := &model.Ecs{}
		resolveAliInstance(ecs, &ins, cli.AccountName, cli.AccountId)
		ecss = append(ecss, ecs)
	}

	if resp.PageNumber > 1 {
		pageTotal := (resp.TotalCount / resp.PageSize) + 1

		for i := 2; i <= pageTotal; i++ {
			req := ali.CreateDescribeInstancesRequest()
			req.RegionId = cli.Region
			req.PageNumber = requests.NewInteger(i)
			res, err := cli.AliClient.DescribeInstances(req)
			if err != nil {
				log.WithError(err)
			}
			for _, ins := range res.Instances.Instance {
				log.WithFields(log.Fields{"account_name": cli.AccountName, "instance_name": ins.InstanceName}).Info()
				ecs := &model.Ecs{}
				resolveAliInstance(ecs, &ins, cli.AccountName, cli.AccountId)
				ecss = append(ecss, ecs)
			}
		}
	}
	return
}

func (manager *EcsManager) describeInstancesFromTencentCloud(cli *EcsMClient) (ecss []*model.Ecs) {
	client := cli.TencentClient
	req := tencent.NewDescribeInstancesRequest()
	req.Limit = common.Int64Ptr(100)

	resp, err := client.DescribeInstances(req)
	if err != nil {
		log.WithField("account_name", cli.AccountName).WithError(err)
		return
	}

	for _, ins := range resp.Response.InstanceSet {
		log.WithFields(log.Fields{"account_name": cli.AccountName, "instance_name": *ins.InstanceName}).Info()
		ecs := &model.Ecs{}
		resolveTencentInstance(ecs, ins, cli.AccountName, cli.AccountId, cli.Region)
		ecss = append(ecss, ecs)
	}

	if totalPageNumber := int(*resp.Response.TotalCount/100 + 1); totalPageNumber > 1 {
		for pageNumber := 2; pageNumber <= totalPageNumber; pageNumber ++ {
			req := tencent.NewDescribeInstancesRequest()
			req.Offset = common.Int64Ptr(int64((pageNumber - 1) * 100))

			resp, err := client.DescribeInstances(req)
			if err != nil {
				log.WithField("account_name", cli.AccountName).Error()
			}

			for _, ins := range resp.Response.InstanceSet {
				log.WithFields(log.Fields{"account_name": cli.AccountName, "instance_name": ins.InstanceName}).Error()
				ecs := &model.Ecs{}
				resolveTencentInstance(ecs, ins, cli.AccountName, cli.AccountId, cli.Region)
				ecss = append(ecss, ecs)
			}
		}
	}
	return
}

func (manager *EcsManager) DescribeInstances() error {
	instances := manager.describeInstancesFromRemoteServer()
	insertOption := `
		on duplicate key update
		instance_name=values(instance_name),
		public_ip=values(public_ip),
		private_ip=values(private_ip),
		os_type=values(os_type),
		cpu=values(cpu),
		mem=values(mem),
		instance_type=values(instance_type),
		zone_id=values(zone_id),
		account_name=values(account_name),
		account_id=values(account_id),
		region_id=values(region_id),
		updated_time=values(updated_time)
`
	for _, ins := range instances {
		if err := manager.Db.Set("gorm:insert_option", insertOption).Create(ins).Error; err != nil {
			log.WithError(err)
		}
	}
	return nil
}

//func (e EcsManager) CostEcs(instanceList []string) (float64, error) {
//	bills, err := NewMonthlyBillMgr(&MonthlyBillConfig{
//		e.Logger,
//		e.DbUri,
//	})
//	if err != nil {
//		e.Logger.Log("err", err.Error())
//	}
//	cost, err := bills.GetCostByInstanceList(instanceList)
//	if err != nil {
//		return 0, err
//	}
//	return cost.Cost, nil
//}

func (manager *EcsManager) ListInstances(ctx context.Context, input *http.ListInstancesInput) ([]*http.EcsModelWithCost, error) {
	IsAdmin := ctx.Value("Is-Admin")
	UserGroupList := ctx.Value("User-Group-List")
	GroupListString := ctx.Value("Group-List")
	UserId := ctx.Value("User-Id")
	fmt.Println("IsAdmin:", IsAdmin, "UserGroupList:", UserGroupList, "UserId:", UserId, "GroupList:", GroupListString)
	sentance := manager.Db.Table("ecs_model")
	if IsAdmin != nil && !IsAdmin.(bool) && GroupListString != nil && GroupListString.(string) != "" {
		sentance = sentance.Where("group_id in (?)", strings.Split(GroupListString.(string), ","))
	}

	ecsModel := model.Ecs{
		AccountId: int64(input.AccountId),
		RegionId:  input.RegionId,
		GroupId:   int64(input.GroupId),
	}

	ret := make([]*http.EcsModelWithCost, 0)
	ecsList := make([]*http.EcsModelWithCost, 0)

	if err := sentance.Where(&ecsModel).Find(&ecsList).Error; err != nil {
		return nil, err
	}

	instanceList := &[]string{}
	for _, ecsItem := range ecsList {
		*instanceList = append(*instanceList, ecsItem.InstanceId)
	}
	// get current month cost
	billingCycleCurrentMonth := time.Now().Format("2006-01")
	//fmt.Println("instanceList", instanceList, "billingCycleCurrentMonth", billingCycleCurrentMonth)

	rpcManager := rpc.GetRpcManager()

	currentMonthCost, err := rpcManager.Call(ctx, ecsRpc.ServiceNameBill, ecsRpc.PathBillList, ecsRpc.GetCostByListRequest{*instanceList, billingCycleCurrentMonth})
	if err != nil {
		return nil, err
	}
	//fmt.Println("currentMonthCost", currentMonthCost)
	currentMonthCostResp := currentMonthCost.(*ecsRpc.GetCostByListResponse)

	// get previous month cost
	billingCyclePreviousMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")
	//fmt.Println("instanceList", instanceList, "billingCycleCurrentMonth", billingCycleCurrentMonth)
	previousMonthCost, err := rpcManager.Call(ctx, "bill", "/ibills/list", ecsRpc.GetCostByListRequest{*instanceList, billingCyclePreviousMonth})
	if err != nil {
		return ret, err
	}
	//fmt.Println("previousMonthCost", previousMonthCost)
	previousMonthCostResp := previousMonthCost.(*ecsRpc.GetCostByListResponse)
	//fmt.Println("length of ecsList:", len(ecsList))
	//fmt.Println("length of currentMonthCostResp.Data:", len(currentMonthCostResp.Data))
	//fmt.Println("length of previousMonthCostResp.Data:", len(previousMonthCostResp.Data))
	// update ecsModel data
	for _, ecsItem := range ecsList {
		for _, currentMonthCostItem := range currentMonthCostResp.Data {
			if ecsItem.InstanceId == currentMonthCostItem.InstanceId {
				ecsItem.CurrentMonthCost = currentMonthCostItem.Cost
				//fmt.Println("ecsInstanceId:", ecsItem.InstanceId, "costInstanceId:", currentMonthCostItem.InstanceId, "currentMonthCost:", currentMonthCostItem.Cost)
			}
		}
		for _, previousMonthCostItem := range previousMonthCostResp.Data {
			if ecsItem.InstanceId == previousMonthCostItem.InstanceId {
				ecsItem.PreviousMonthCost = previousMonthCostItem.Cost
				//fmt.Println("ecsInstanceId:", ecsItem.InstanceId, "costInstanceId:", previousMonthCostItem.InstanceId, "previousMonthCost:", previousMonthCostItem.Cost)
			}
		}
		ret = append(ret, ecsItem)
	}

	return ret, nil
}

func (manager *EcsManager) EcsCount(ctx context.Context, req *http.EcsCountRequest) (int, error) {
	var count int
	ecsList := make([]model.Ecs, 0)
	sql := manager.Db.Table("ecs_model")
	if req.AccountId != 0 {
		sql = sql.Where("account_id = ?", req.AccountId)
	}
	if req.RegionId != "" {
		sql = sql.Where("region_id = ?", req.RegionId)
	}
	if req.GroupId != 0 {
		sql = sql.Where("group_id = ?", req.GroupId)
	}
	if err := sql.Find(&ecsList).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (manager *EcsManager) EcsCost(ctx context.Context, req *http.EcsCostRequest) (float64, error) {
	rpcManager := rpc.GetRpcManager()
	costResp, err := rpcManager.Call(ctx, ecsRpc.ServiceNameBill, ecsRpc.PathBillCost, req)
	if err != nil {
		return 0, err
	}
	return costResp.(ecsRpc.EcsCostResponse).Data.Cost, nil
}

func (manager *EcsManager) EcsRegister(ctx context.Context, req *http.EcsRegisterRequest) (bool, error) {
	if req.CheckUrl == "" {
		return false, errors.New("check url not found")
	}
	PrivateIp := strings.Split(req.CheckUrl, "/")[len(strings.Split(req.CheckUrl, "/"))-1]
	if PrivateIp == "" {
		return false, errors.New("not match private ip found")
	}
	ecsModel := model.Ecs{}
	if err := manager.Db.Where(model.Ecs{PrivateIp: PrivateIp}).First(&ecsModel).Error; err != nil {
		return false, err
	}
	ecsAgent := model.EcsAgent{
		EcsId:        ecsModel.ID,
		PrivateIp:    PrivateIp,
		CheckUrl:     req.CheckUrl,
		CheckUrlSign: req.CheckUrlSign,
	}
	if err := manager.Db.Where(&ecsAgent).Attrs(model.EcsAgent{BaseModel: publicModel.BaseModel{CreatedTime: time.Now()}}).Assign(model.EcsAgent{BaseModel: publicModel.BaseModel{UpdatedTime: time.Now()}}).FirstOrCreate(&ecsAgent).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (manager *EcsManager) EcsGetAgentStatus(ctx context.Context, req *http.EcsGetAgentStatusRequest) (status int, err error) {
	if req.EcsId == 0 {
		return -1, errors.New("ecs_id is required")
	}
	ecsAgent := model.EcsAgent{}
	if err := manager.Db.Where(&model.EcsAgent{EcsId: int64(req.EcsId),}).First(&ecsAgent).Error; err != nil {
		return -1, nil
	} else {
		return ecsAgent.AgentStatus, nil
	}
}

func (manager *EcsManager) GetAkSk(accountId int) (string, string) {
	var account model2.ThirdAS
	if err := manager.Db.Table("third_as").Where("id = ?", accountId).First(&account).Error; err != nil {
		manager.Db.Rollback()
		return "", ""
	}
	return account.AccessKey, account.SecretKey
}

func NewEcsMgr(db *gorm.DB) (*EcsManager, error) {
	return &EcsManager{
		Db: db,
	}, nil
}
