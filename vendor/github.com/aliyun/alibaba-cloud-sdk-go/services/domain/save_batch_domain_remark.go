package domain

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

// SaveBatchDomainRemark invokes the domain.SaveBatchDomainRemark API synchronously
// api document: https://help.aliyun.com/api/domain/savebatchdomainremark.html
func (client *Client) SaveBatchDomainRemark(request *SaveBatchDomainRemarkRequest) (response *SaveBatchDomainRemarkResponse, err error) {
	response = CreateSaveBatchDomainRemarkResponse()
	err = client.DoAction(request, response)
	return
}

// SaveBatchDomainRemarkWithChan invokes the domain.SaveBatchDomainRemark API asynchronously
// api document: https://help.aliyun.com/api/domain/savebatchdomainremark.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchDomainRemarkWithChan(request *SaveBatchDomainRemarkRequest) (<-chan *SaveBatchDomainRemarkResponse, <-chan error) {
	responseChan := make(chan *SaveBatchDomainRemarkResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SaveBatchDomainRemark(request)
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

// SaveBatchDomainRemarkWithCallback invokes the domain.SaveBatchDomainRemark API asynchronously
// api document: https://help.aliyun.com/api/domain/savebatchdomainremark.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SaveBatchDomainRemarkWithCallback(request *SaveBatchDomainRemarkRequest, callback func(response *SaveBatchDomainRemarkResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SaveBatchDomainRemarkResponse
		var err error
		defer close(result)
		response, err = client.SaveBatchDomainRemark(request)
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

// SaveBatchDomainRemarkRequest is the request struct for api SaveBatchDomainRemark
type SaveBatchDomainRemarkRequest struct {
	*requests.RpcRequest
	InstanceIds string `position:"Query" name:"InstanceIds"`
	Remark      string `position:"Query" name:"Remark"`
	Lang        string `position:"Query" name:"Lang"`
}

// SaveBatchDomainRemarkResponse is the response struct for api SaveBatchDomainRemark
type SaveBatchDomainRemarkResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSaveBatchDomainRemarkRequest creates a request to invoke SaveBatchDomainRemark API
func CreateSaveBatchDomainRemarkRequest() (request *SaveBatchDomainRemarkRequest) {
	request = &SaveBatchDomainRemarkRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Domain", "2018-01-29", "SaveBatchDomainRemark", "", "")
	return
}

// CreateSaveBatchDomainRemarkResponse creates a response to parse from SaveBatchDomainRemark response
func CreateSaveBatchDomainRemarkResponse() (response *SaveBatchDomainRemarkResponse) {
	response = &SaveBatchDomainRemarkResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
