package model

const (
	RegionAliBeijing           = "cn-beijing"
	RegionAliShanghai          = "cn-shanghai"
	RegionAliHongkong          = "cn-hongkong"
	RegionAliHuhehaote         = "cn-huhehaote"
	RegionAliZhangjiakou       = "cn-zhangjiakou"
	RegionTencentBangkok       = "ap-bangkok"
	RegionTencentBeijing       = "ap-beijing"
	RegionTencentChengdu       = "ap-chengdu"
	RegionTencentChongqing     = "ap-chongqing"
	RegionTencentGuangzhou     = "ap-guangzhou"
	RegionTencentGuangzhouOpen = "ap-guangzhou-open"
	RegionTencentHongkong      = "ap-hongkong"
	RegionTencentMumbai        = "ap-mumbai"
	RegionTencentSeoul         = "ap-seoul"
	RegionTencentShanghai      = "ap-shanghai"
	RegionTencentShanghaiFsi   = "ap-shanghai-fsi"
	RegionTencentShenzhenFsi   = "ap-shenzhen-fsi"
	RegionTencentSingapore     = "ap-singapore"
	RegionTencentFrankfurt     = "eu-frankfurt"
	RegionTencentMoscow        = "eu-moscow"
	RegionTencentAshburn       = "na-ashburn"
	RegionTencentSiliconvalley = "na-siliconvalley"
	RegionTencentToronto       = "na-toronto"
)

var (
	RegionMap = map[Platform][]string{
		PlatformAliCloud: {
			RegionAliBeijing,
			RegionAliShanghai,
			RegionAliHongkong,
			RegionAliHuhehaote,
			RegionAliZhangjiakou,
		},
		PlatformTencentCloud: {
			RegionTencentBangkok,
			RegionTencentBeijing,
			RegionTencentChengdu,
			RegionTencentChongqing,
			RegionTencentGuangzhou,
			RegionTencentGuangzhouOpen,
			RegionTencentHongkong,
			RegionTencentMumbai,
			RegionTencentSeoul,
			RegionTencentShanghai,
			RegionTencentShanghaiFsi,
			RegionTencentShenzhenFsi,
			RegionTencentSingapore,
			RegionTencentFrankfurt,
			RegionTencentMoscow,
			RegionTencentAshburn,
			RegionTencentSiliconvalley,
			RegionTencentToronto,
		},
	}
)
