package live

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

// AddCustomLiveStreamTranscode invokes the live.AddCustomLiveStreamTranscode API synchronously
// api document: https://help.aliyun.com/api/live/addcustomlivestreamtranscode.html
func (client *Client) AddCustomLiveStreamTranscode(request *AddCustomLiveStreamTranscodeRequest) (response *AddCustomLiveStreamTranscodeResponse, err error) {
	response = CreateAddCustomLiveStreamTranscodeResponse()
	err = client.DoAction(request, response)
	return
}

// AddCustomLiveStreamTranscodeWithChan invokes the live.AddCustomLiveStreamTranscode API asynchronously
// api document: https://help.aliyun.com/api/live/addcustomlivestreamtranscode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddCustomLiveStreamTranscodeWithChan(request *AddCustomLiveStreamTranscodeRequest) (<-chan *AddCustomLiveStreamTranscodeResponse, <-chan error) {
	responseChan := make(chan *AddCustomLiveStreamTranscodeResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.AddCustomLiveStreamTranscode(request)
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

// AddCustomLiveStreamTranscodeWithCallback invokes the live.AddCustomLiveStreamTranscode API asynchronously
// api document: https://help.aliyun.com/api/live/addcustomlivestreamtranscode.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) AddCustomLiveStreamTranscodeWithCallback(request *AddCustomLiveStreamTranscodeRequest, callback func(response *AddCustomLiveStreamTranscodeResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *AddCustomLiveStreamTranscodeResponse
		var err error
		defer close(result)
		response, err = client.AddCustomLiveStreamTranscode(request)
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

// AddCustomLiveStreamTranscodeRequest is the request struct for api AddCustomLiveStreamTranscode
type AddCustomLiveStreamTranscodeRequest struct {
	*requests.RpcRequest
	App          string           `position:"Query" name:"App"`
	Template     string           `position:"Query" name:"Template"`
	Profile      requests.Integer `position:"Query" name:"Profile"`
	FPS          requests.Integer `position:"Query" name:"FPS"`
	Gop          string           `position:"Query" name:"Gop"`
	OwnerId      requests.Integer `position:"Query" name:"OwnerId"`
	TemplateType string           `position:"Query" name:"TemplateType"`
	AudioBitrate requests.Integer `position:"Query" name:"AudioBitrate"`
	Domain       string           `position:"Query" name:"Domain"`
	Width        requests.Integer `position:"Query" name:"Width"`
	VideoBitrate requests.Integer `position:"Query" name:"VideoBitrate"`
	Height       requests.Integer `position:"Query" name:"Height"`
}

// AddCustomLiveStreamTranscodeResponse is the response struct for api AddCustomLiveStreamTranscode
type AddCustomLiveStreamTranscodeResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateAddCustomLiveStreamTranscodeRequest creates a request to invoke AddCustomLiveStreamTranscode API
func CreateAddCustomLiveStreamTranscodeRequest() (request *AddCustomLiveStreamTranscodeRequest) {
	request = &AddCustomLiveStreamTranscodeRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("live", "2016-11-01", "AddCustomLiveStreamTranscode", "live", "openAPI")
	return
}

// CreateAddCustomLiveStreamTranscodeResponse creates a response to parse from AddCustomLiveStreamTranscode response
func CreateAddCustomLiveStreamTranscodeResponse() (response *AddCustomLiveStreamTranscodeResponse) {
	response = &AddCustomLiveStreamTranscodeResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
