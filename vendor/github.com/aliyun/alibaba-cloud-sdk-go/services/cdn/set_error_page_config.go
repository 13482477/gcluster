package cdn

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

// SetErrorPageConfig invokes the cdn.SetErrorPageConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/seterrorpageconfig.html
func (client *Client) SetErrorPageConfig(request *SetErrorPageConfigRequest) (response *SetErrorPageConfigResponse, err error) {
	response = CreateSetErrorPageConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetErrorPageConfigWithChan invokes the cdn.SetErrorPageConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/seterrorpageconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetErrorPageConfigWithChan(request *SetErrorPageConfigRequest) (<-chan *SetErrorPageConfigResponse, <-chan error) {
	responseChan := make(chan *SetErrorPageConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetErrorPageConfig(request)
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

// SetErrorPageConfigWithCallback invokes the cdn.SetErrorPageConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/seterrorpageconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetErrorPageConfigWithCallback(request *SetErrorPageConfigRequest, callback func(response *SetErrorPageConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetErrorPageConfigResponse
		var err error
		defer close(result)
		response, err = client.SetErrorPageConfig(request)
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

// SetErrorPageConfigRequest is the request struct for api SetErrorPageConfig
type SetErrorPageConfigRequest struct {
	*requests.RpcRequest
	PageType      string           `position:"Query" name:"PageType"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	CustomPageUrl string           `position:"Query" name:"CustomPageUrl"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// SetErrorPageConfigResponse is the response struct for api SetErrorPageConfig
type SetErrorPageConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetErrorPageConfigRequest creates a request to invoke SetErrorPageConfig API
func CreateSetErrorPageConfigRequest() (request *SetErrorPageConfigRequest) {
	request = &SetErrorPageConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2014-11-11", "SetErrorPageConfig", "", "")
	return
}

// CreateSetErrorPageConfigResponse creates a response to parse from SetErrorPageConfig response
func CreateSetErrorPageConfigResponse() (response *SetErrorPageConfigResponse) {
	response = &SetErrorPageConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}