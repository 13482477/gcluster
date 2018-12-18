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

// SendRoomUserNotification invokes the live.SendRoomUserNotification API synchronously
// api document: https://help.aliyun.com/api/live/sendroomusernotification.html
func (client *Client) SendRoomUserNotification(request *SendRoomUserNotificationRequest) (response *SendRoomUserNotificationResponse, err error) {
	response = CreateSendRoomUserNotificationResponse()
	err = client.DoAction(request, response)
	return
}

// SendRoomUserNotificationWithChan invokes the live.SendRoomUserNotification API asynchronously
// api document: https://help.aliyun.com/api/live/sendroomusernotification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SendRoomUserNotificationWithChan(request *SendRoomUserNotificationRequest) (<-chan *SendRoomUserNotificationResponse, <-chan error) {
	responseChan := make(chan *SendRoomUserNotificationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SendRoomUserNotification(request)
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

// SendRoomUserNotificationWithCallback invokes the live.SendRoomUserNotification API asynchronously
// api document: https://help.aliyun.com/api/live/sendroomusernotification.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SendRoomUserNotificationWithCallback(request *SendRoomUserNotificationRequest, callback func(response *SendRoomUserNotificationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SendRoomUserNotificationResponse
		var err error
		defer close(result)
		response, err = client.SendRoomUserNotification(request)
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

// SendRoomUserNotificationRequest is the request struct for api SendRoomUserNotification
type SendRoomUserNotificationRequest struct {
	*requests.RpcRequest
	Data     string           `position:"Query" name:"Data"`
	ToAppUid string           `position:"Query" name:"ToAppUid"`
	AppUid   string           `position:"Query" name:"AppUid"`
	OwnerId  requests.Integer `position:"Query" name:"OwnerId"`
	Priority requests.Integer `position:"Query" name:"Priority"`
	RoomId   string           `position:"Query" name:"RoomId"`
	AppId    string           `position:"Query" name:"AppId"`
}

// SendRoomUserNotificationResponse is the response struct for api SendRoomUserNotification
type SendRoomUserNotificationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	MessageId string `json:"MessageId" xml:"MessageId"`
}

// CreateSendRoomUserNotificationRequest creates a request to invoke SendRoomUserNotification API
func CreateSendRoomUserNotificationRequest() (request *SendRoomUserNotificationRequest) {
	request = &SendRoomUserNotificationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("live", "2016-11-01", "SendRoomUserNotification", "live", "openAPI")
	return
}

// CreateSendRoomUserNotificationResponse creates a response to parse from SendRoomUserNotification response
func CreateSendRoomUserNotificationResponse() (response *SendRoomUserNotificationResponse) {
	response = &SendRoomUserNotificationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
