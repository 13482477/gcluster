package cms

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// CreateMyGroups invokes the cms.CreateMyGroups API synchronously
// api document: https://help.aliyun.com/api/cms/createmygroups.html
func (client *Client) CreateMyGroups(request *CreateMyGroupsRequest) (response *CreateMyGroupsResponse, err error) {
	response = CreateCreateMyGroupsResponse()
	err = client.DoAction(request, response)
	return
}

// CreateMyGroupsWithChan invokes the cms.CreateMyGroups API asynchronously
// api document: https://help.aliyun.com/api/cms/createmygroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateMyGroupsWithChan(request *CreateMyGroupsRequest) (<-chan *CreateMyGroupsResponse, <-chan error) {
	responseChan := make(chan *CreateMyGroupsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateMyGroups(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// CreateMyGroupsWithCallback invokes the cms.CreateMyGroups API asynchronously
// api document: https://help.aliyun.com/api/cms/createmygroups.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateMyGroupsWithCallback(request *CreateMyGroupsRequest, callback func(response *CreateMyGroupsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateMyGroupsResponse
		var err error
		defer close(result)
		response, err = client.CreateMyGroups(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// CreateMyGroupsRequest is the request struct for api CreateMyGroups
type CreateMyGroupsRequest struct {
	*requests.RpcRequest
	ContactGroups string           `position:"Query" name:"ContactGroups"`
	Options       string           `position:"Query" name:"Options"`
	Type          string           `position:"Query" name:"Type"`
	ServiceId     requests.Integer `position:"Query" name:"ServiceId"`
	GroupName     string           `position:"Query" name:"GroupName"`
	BindUrl       string           `position:"Query" name:"BindUrl"`
}

// CreateMyGroupsResponse is the response struct for api CreateMyGroups
type CreateMyGroupsResponse struct {
	*responses.BaseResponse
	RequestId    string `json:"RequestId" xml:"RequestId"`
	Success      bool   `json:"Success" xml:"Success"`
	ErrorCode    int    `json:"ErrorCode" xml:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage" xml:"ErrorMessage"`
	GroupId      int    `json:"GroupId" xml:"GroupId"`
}

// CreateCreateMyGroupsRequest creates a request to invoke CreateMyGroups API
func CreateCreateMyGroupsRequest() (request *CreateMyGroupsRequest) {
	request = &CreateMyGroupsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cms", "2018-03-08", "CreateMyGroups", "cms", "openAPI")
	return
}

// CreateCreateMyGroupsResponse creates a response to parse from CreateMyGroups response
func CreateCreateMyGroupsResponse() (response *CreateMyGroupsResponse) {
	response = &CreateMyGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
