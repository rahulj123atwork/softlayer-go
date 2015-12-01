package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"text/template"

	services "github.com/TheWeatherCompany/softlayer-go/services"
	softlayer "github.com/TheWeatherCompany/softlayer-go/softlayer"
)

const (
	SOFTLAYER_API_URL  = "api.softlayer.com/rest/v3"
	TEMPLATE_ROOT_PATH = "templates"
	SL_GO_NON_VERBOSE  = "SL_GO_NON_VERBOSE"
)

type softLayerClient struct {
	username string
	apiKey   string

	templatePath string

	httpClient *http.Client

	softLayerServices map[string]softlayer.Service

	nonVerbose bool
}

func NewSoftLayerClient(username, apiKey string) *softLayerClient {
	pwd, _ := os.Getwd()
	slc := &softLayerClient{
		username: username,
		apiKey:   apiKey,

		templatePath: filepath.Join(pwd, TEMPLATE_ROOT_PATH),

		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},

		nonVerbose: checkNonVerbose(),

		softLayerServices: map[string]softlayer.Service{},
	}

	slc.initSoftLayerServices()

	return slc
}

//softlayer.Client interface methods

func (slc *softLayerClient) GetService(serviceName string) (softlayer.Service, error) {
	slService, ok := slc.softLayerServices[serviceName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("softlayer-go does not support service '%s'", serviceName))
	}

	return slService, nil
}

func (slc *softLayerClient) GetSoftLayer_Account_Service() (softlayer.SoftLayer_Account_Service, error) {
	slService, err := slc.GetService("SoftLayer_Account")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Account_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Virtual_Guest_Service() (softlayer.SoftLayer_Virtual_Guest_Service, error) {
	slService, err := slc.GetService("SoftLayer_Virtual_Guest")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Virtual_Guest_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Virtual_Disk_Image_Service() (softlayer.SoftLayer_Virtual_Disk_Image_Service, error) {
	slService, err := slc.GetService("SoftLayer_Virtual_Disk_Image")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Virtual_Disk_Image_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Security_Ssh_Key_Service() (softlayer.SoftLayer_Security_Ssh_Key_Service, error) {
	slService, err := slc.GetService("SoftLayer_Security_Ssh_Key")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Security_Ssh_Key_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Product_Package_Service() (softlayer.SoftLayer_Product_Package_Service, error) {
	slService, err := slc.GetService("SoftLayer_Product_Package")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Product_Package_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service() (softlayer.SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service, error) {
	slService, err := slc.GetService("SoftLayer_Virtual_Guest_Block_Device_Template_Group")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Virtual_Guest_Block_Device_Template_Group_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Network_Storage_Service() (softlayer.SoftLayer_Network_Storage_Service, error) {
	slService, err := slc.GetService("SoftLayer_Network_Storage")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Network_Storage_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Network_Storage_Allowed_Host_Service() (softlayer.SoftLayer_Network_Storage_Allowed_Host_Service, error) {
	slService, err := slc.GetService("SoftLayer_Network_Storage_Allowed_Host")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Network_Storage_Allowed_Host_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Product_Order_Service() (softlayer.SoftLayer_Product_Order_Service, error) {
	slService, err := slc.GetService("SoftLayer_Product_Order")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Product_Order_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Billing_Item_Cancellation_Request_Service() (softlayer.SoftLayer_Billing_Item_Cancellation_Request_Service, error) {
	slService, err := slc.GetService("SoftLayer_Billing_Item_Cancellation_Request")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Billing_Item_Cancellation_Request_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Hardware_Service() (softlayer.SoftLayer_Hardware_Service, error) {
	slService, err := slc.GetService("SoftLayer_Hardware")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Hardware_Service), nil
}

func (slc *softLayerClient) GetSoftLayer_Dns_Domain_Record_Service() (softlayer.SoftLayer_Dns_Domain_Record_Service, error) {
	slService, err := slc.GetService("SoftLayer_Dns_Domain_ResourceRecord")
	if err != nil {
		return nil, err
	}

	return slService.(softlayer.SoftLayer_Dns_Domain_Record_Service), nil
}

//Public methods

func (slc *softLayerClient) DoRawHttpRequestWithObjectMask(path string, masks []string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s@%s/%s", slc.username, slc.apiKey, SOFTLAYER_API_URL, path)

	url += "?objectMask="
	for i := 0; i < len(masks); i++ {
		url += masks[i]
		if i != len(masks)-1 {
			url += ";"
		}
	}

	return slc.makeHttpRequest(url, requestType, requestBody)
}

func (slc *softLayerClient) DoRawHttpRequestWithObjectFilter(path string, filters string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s@%s/%s", slc.username, slc.apiKey, SOFTLAYER_API_URL, path)
	url += "?objectFilter=" + filters

	return slc.makeHttpRequest(url, requestType, requestBody)
}

func (slc *softLayerClient) DoRawHttpRequestWithObjectFilterAndObjectMask(path string, masks []string, filters string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s@%s/%s", slc.username, slc.apiKey, SOFTLAYER_API_URL, path)

	url += "?objectFilter=" + filters

	url += "&objectMask=filteredMask["
	for i := 0; i < len(masks); i++ {
		url += masks[i]
		if i != len(masks)-1 {
			url += ";"
		}
	}
	url += "]"

	return slc.makeHttpRequest(url, requestType, requestBody)
}

func (slc *softLayerClient) DoRawHttpRequest(path string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	url := fmt.Sprintf("https://%s:%s@%s/%s", slc.username, slc.apiKey, SOFTLAYER_API_URL, path)
	return slc.makeHttpRequest(url, requestType, requestBody)
}

func (slc *softLayerClient) GenerateRequestBody(templateData interface{}) (*bytes.Buffer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	bodyTemplate := template.Must(template.ParseFiles(filepath.Join(cwd, slc.templatePath)))
	body := new(bytes.Buffer)
	bodyTemplate.Execute(body, templateData)

	return body, nil
}

func (slc *softLayerClient) HasErrors(body map[string]interface{}) error {
	if errString, ok := body["error"]; !ok {
		return nil
	} else {
		return errors.New(errString.(string))
	}
}

func (slc *softLayerClient) CheckForHttpResponseErrors(data []byte) error {
	var decodedResponse map[string]interface{}
	err := json.Unmarshal(data, &decodedResponse)
	if err != nil {
		return err
	}

	if err := slc.HasErrors(decodedResponse); err != nil {
		return err
	}

	return nil
}

//Private methods

func (slc *softLayerClient) initSoftLayerServices() {
	slc.softLayerServices["SoftLayer_Account"] = services.NewSoftLayer_Account_Service(slc)
	slc.softLayerServices["SoftLayer_Virtual_Guest"] = services.NewSoftLayer_Virtual_Guest_Service(slc)
	slc.softLayerServices["SoftLayer_Virtual_Disk_Image"] = services.NewSoftLayer_Virtual_Disk_Image_Service(slc)
	slc.softLayerServices["SoftLayer_Security_Ssh_Key"] = services.NewSoftLayer_Security_Ssh_Key_Service(slc)
	slc.softLayerServices["SoftLayer_Product_Package"] = services.NewSoftLayer_Product_Package_Service(slc)
	slc.softLayerServices["SoftLayer_Network_Storage"] = services.NewSoftLayer_Network_Storage_Service(slc)
	slc.softLayerServices["SoftLayer_Network_Storage_Allowed_Host"] = services.NewSoftLayer_Network_Storage_Allowed_Host_Service(slc)
	slc.softLayerServices["SoftLayer_Product_Order"] = services.NewSoftLayer_Product_Order_Service(slc)
	slc.softLayerServices["SoftLayer_Billing_Item_Cancellation_Request"] = services.NewSoftLayer_Billing_Item_Cancellation_Request_Service(slc)
	slc.softLayerServices["SoftLayer_Virtual_Guest_Block_Device_Template_Group"] = services.NewSoftLayer_Virtual_Guest_Block_Device_Template_Group_Service(slc)
	slc.softLayerServices["SoftLayer_Hardware"] = services.NewSoftLayer_Hardware_Service(slc)
	slc.softLayerServices["SoftLayer_Dns_Domain_ResourceRecord"] = services.NewSoftLayer_Dns_Domain_Record_Service(slc)
}

func (slc *softLayerClient) makeHttpRequest(url string, requestType string, requestBody *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest(requestType, url, requestBody)
	if err != nil {
		return nil, err
	}

	bs, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}

	if !slc.nonVerbose {
		fmt.Fprintf(os.Stderr, "\n---\n[softlayer-go] Request:\n%s\n", string(bs))
	}

	resp, err := slc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err = httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	if !slc.nonVerbose {
		fmt.Fprintf(os.Stderr, "[softlayer-go] Response:\n%s\n", string(bs))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

//Private helper methods

func checkNonVerbose() bool {
	slGoNonVerbose := os.Getenv(SL_GO_NON_VERBOSE)
	switch slGoNonVerbose {
	case "yes":
		return true
	case "YES":
		return true
	case "true":
		return true
	case "TRUE":
		return true
	}

	return false
}
