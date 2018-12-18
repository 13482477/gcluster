package ehpc

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

// ImageInfo is a nested struct in ehpc response
type ImageInfo struct {
	Uid               string    `json:"Uid" xml:"Uid"`
	Repository        string    `json:"Repository" xml:"Repository"`
	SkuCode           string    `json:"SkuCode" xml:"SkuCode"`
	ImageId           string    `json:"ImageId" xml:"ImageId"`
	ImageOwnerAlias   string    `json:"ImageOwnerAlias" xml:"ImageOwnerAlias"`
	System            string    `json:"System" xml:"System"`
	PostInstallScript string    `json:"PostInstallScript" xml:"PostInstallScript"`
	ProductCode       string    `json:"ProductCode" xml:"ProductCode"`
	Tag               string    `json:"Tag" xml:"Tag"`
	PricingCycle      string    `json:"PricingCycle" xml:"PricingCycle"`
	ImageName         string    `json:"ImageName" xml:"ImageName"`
	Status            string    `json:"Status" xml:"Status"`
	Description       string    `json:"Description" xml:"Description"`
	Type              string    `json:"Type" xml:"Type"`
	UpdateDateTime    string    `json:"UpdateDateTime" xml:"UpdateDateTime"`
	BaseOsTag         BaseOsTag `json:"BaseOsTag" xml:"BaseOsTag"`
}
