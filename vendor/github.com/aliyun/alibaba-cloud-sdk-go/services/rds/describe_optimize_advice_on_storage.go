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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeOptimizeAdviceOnStorage invokes the rds.DescribeOptimizeAdviceOnStorage API synchronously
// api document: https://help.aliyun.com/api/rds/describeoptimizeadviceonstorage.html
func (client *Client) DescribeOptimizeAdviceOnStorage(request *DescribeOptimizeAdviceOnStorageRequest) (response *DescribeOptimizeAdviceOnStorageResponse, err error) {
	response = CreateDescribeOptimizeAdviceOnStorageResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeOptimizeAdviceOnStorageWithChan invokes the rds.DescribeOptimizeAdviceOnStorage API asynchronously
// api document: https://help.aliyun.com/api/rds/describeoptimizeadviceonstorage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeOptimizeAdviceOnStorageWithChan(request *DescribeOptimizeAdviceOnStorageRequest) (<-chan *DescribeOptimizeAdviceOnStorageResponse, <-chan error) {
	responseChan := make(chan *DescribeOptimizeAdviceOnStorageResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeOptimizeAdviceOnStorage(request)
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

// DescribeOptimizeAdviceOnStorageWithCallback invokes the rds.DescribeOptimizeAdviceOnStorage API asynchronously
// api document: https://help.aliyun.com/api/rds/describeoptimizeadviceonstorage.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeOptimizeAdviceOnStorageWithCallback(request *DescribeOptimizeAdviceOnStorageRequest, callback func(response *DescribeOptimizeAdviceOnStorageResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeOptimizeAdviceOnStorageResponse
		var err error
		defer close(result)
		response, err = client.DescribeOptimizeAdviceOnStorage(request)
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

// DescribeOptimizeAdviceOnStorageRequest is the request struct for api DescribeOptimizeAdviceOnStorage
type DescribeOptimizeAdviceOnStorageRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	PageSize             requests.Integer `position:"Query" name:"PageSize"`
	DBInstanceId         string           `position:"Query" name:"DBInstanceId"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	PageNumber           requests.Integer `position:"Query" name:"PageNumber"`
}

// DescribeOptimizeAdviceOnStorageResponse is the response struct for api DescribeOptimizeAdviceOnStorage
type DescribeOptimizeAdviceOnStorageResponse struct {
	*responses.BaseResponse
	RequestId         string                                 `json:"RequestId" xml:"RequestId"`
	DBInstanceId      string                                 `json:"DBInstanceId" xml:"DBInstanceId"`
	TotalRecordsCount int                                    `json:"TotalRecordsCount" xml:"TotalRecordsCount"`
	PageNumber        int                                    `json:"PageNumber" xml:"PageNumber"`
	PageRecordCount   int                                    `json:"PageRecordCount" xml:"PageRecordCount"`
	Items             ItemsInDescribeOptimizeAdviceOnStorage `json:"Items" xml:"Items"`
}

// CreateDescribeOptimizeAdviceOnStorageRequest creates a request to invoke DescribeOptimizeAdviceOnStorage API
func CreateDescribeOptimizeAdviceOnStorageRequest() (request *DescribeOptimizeAdviceOnStorageRequest) {
	request = &DescribeOptimizeAdviceOnStorageRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "DescribeOptimizeAdviceOnStorage", "rds", "openAPI")
	return
}

// CreateDescribeOptimizeAdviceOnStorageResponse creates a response to parse from DescribeOptimizeAdviceOnStorage response
func CreateDescribeOptimizeAdviceOnStorageResponse() (response *DescribeOptimizeAdviceOnStorageResponse) {
	response = &DescribeOptimizeAdviceOnStorageResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}