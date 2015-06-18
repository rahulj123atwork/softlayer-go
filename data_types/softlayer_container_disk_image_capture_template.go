package data_types

type SoftLayer_Container_Disk_Image_Capture_Template struct {
	Description string                                                 `json:"description"`
	Name        string                                                 `json:"name"`
	Summary     string                                                 `json:"summary"`
	Volumes     SoftLayer_Container_Disk_Image_Capture_Template_Volume `json:"volumes"`
}
