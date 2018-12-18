package cloudapi

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

// DescribeApiTrafficData invokes the cloudapi.DescribeApiTrafficData API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapitrafficdata.html
func (client *Client) DescribeApiTrafficData(request *DescribeApiTrafficDataRequest) (response *DescribeApiTrafficDataResponse, err error) {
	response = CreateDescribeApiTrafficDataResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeApiTrafficDataWithChan invokes the cloudapi.DescribeApiTrafficData API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapitrafficdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiTrafficDataWithChan(request *DescribeApiTrafficDataRequest) (<-chan *DescribeApiTrafficDataResponse, <-chan error) {
	responseChan := make(chan *DescribeApiTrafficDataResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeApiTrafficData(request)
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

// DescribeApiTrafficDataWithCallback invokes the cloudapi.DescribeApiTrafficData API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapitrafficdata.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiTrafficDataWithCallback(request *DescribeApiTrafficDataRequest, callback func(response *DescribeApiTrafficDataResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeApiTrafficDataResponse
		var err error
		defer close(result)
		response, err = client.DescribeApiTrafficData(request)
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

// DescribeApiTrafficDataRequest is the request struct for api DescribeApiTrafficData
type DescribeApiTrafficDataRequest struct {
	*requests.RpcRequest
	SecurityToken string `position:"Query" name:"SecurityToken"`
	GroupId       string `position:"Query" name:"GroupId"`
	EndTime       string `position:"Query" name:"EndTime"`
	StartTime     string `position:"Query" name:"StartTime"`
	ApiId         string `position:"Query" name:"ApiId"`
}

// DescribeApiTrafficDataResponse is the response struct for api DescribeApiTrafficData
type DescribeApiTrafficDataResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	CallUploads   CallUploads   `json:"CallUploads" xml:"CallUploads"`
	CallDownloads CallDownloads `json:"CallDownloads" xml:"CallDownloads"`
}

// CreateDescribeApiTrafficDataRequest creates a request to invoke DescribeApiTrafficData API
func CreateDescribeApiTrafficDataRequest() (request *DescribeApiTrafficDataRequest) {
	request = &DescribeApiTrafficDataRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeApiTrafficData", "apigateway", "openAPI")
	return
}

// CreateDescribeApiTrafficDataResponse creates a response to parse from DescribeApiTrafficData response
func CreateDescribeApiTrafficDataResponse() (response *DescribeApiTrafficDataResponse) {
	response = &DescribeApiTrafficDataResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
