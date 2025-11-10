package aliyun

import sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}
