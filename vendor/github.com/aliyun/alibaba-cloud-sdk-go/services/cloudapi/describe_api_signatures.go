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

// DescribeApiSignatures invokes the cloudapi.DescribeApiSignatures API synchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapisignatures.html
func (client *Client) DescribeApiSignatures(request *DescribeApiSignaturesRequest) (response *DescribeApiSignaturesResponse, err error) {
	response = CreateDescribeApiSignaturesResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeApiSignaturesWithChan invokes the cloudapi.DescribeApiSignatures API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapisignatures.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiSignaturesWithChan(request *DescribeApiSignaturesRequest) (<-chan *DescribeApiSignaturesResponse, <-chan error) {
	responseChan := make(chan *DescribeApiSignaturesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeApiSignatures(request)
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

// DescribeApiSignaturesWithCallback invokes the cloudapi.DescribeApiSignatures API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/describeapisignatures.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeApiSignaturesWithCallback(request *DescribeApiSignaturesRequest, callback func(response *DescribeApiSignaturesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeApiSignaturesResponse
		var err error
		defer close(result)
		response, err = client.DescribeApiSignatures(request)
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

// DescribeApiSignaturesRequest is the request struct for api DescribeApiSignatures
type DescribeApiSignaturesRequest struct {
	*requests.RpcRequest
	StageName     string           `position:"Query" name:"StageName"`
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	GroupId       string           `position:"Query" name:"GroupId"`
	PageSize      requests.Integer `position:"Query" name:"PageSize"`
	PageNumber    requests.Integer `position:"Query" name:"PageNumber"`
	ApiIds        string           `position:"Query" name:"ApiIds"`
}

// DescribeApiSignaturesResponse is the response struct for api DescribeApiSignatures
type DescribeApiSignaturesResponse struct {
	*responses.BaseResponse
	RequestId     string        `json:"RequestId" xml:"RequestId"`
	TotalCount    int           `json:"TotalCount" xml:"TotalCount"`
	PageSize      int           `json:"PageSize" xml:"PageSize"`
	PageNumber    int           `json:"PageNumber" xml:"PageNumber"`
	ApiSignatures ApiSignatures `json:"ApiSignatures" xml:"ApiSignatures"`
}

// CreateDescribeApiSignaturesRequest creates a request to invoke DescribeApiSignatures API
func CreateDescribeApiSignaturesRequest() (request *DescribeApiSignaturesRequest) {
	request = &DescribeApiSignaturesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "DescribeApiSignatures", "apigateway", "openAPI")
	return
}

// CreateDescribeApiSignaturesResponse creates a response to parse from DescribeApiSignatures response
func CreateDescribeApiSignaturesResponse() (response *DescribeApiSignaturesResponse) {
	response = &DescribeApiSignaturesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
