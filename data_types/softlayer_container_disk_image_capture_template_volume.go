package data_types

type SoftLayer_Container_Disk_Image_Capture_Template_Volume struct {
	Name       string                                                           `json:"name"`
	Partitions SoftLayer_Container_Disk_Image_Capture_Template_Volume_Partition `json:"partitions"`
}
