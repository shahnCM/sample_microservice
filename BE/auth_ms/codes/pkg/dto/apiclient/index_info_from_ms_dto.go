package apiclient

type IndexInfoFromMsDto struct {
	Status     any                         `json:"status"`
	Data       *DataForIndexInfo           `json:"data"`
	DataBn     *DataForIndexInfo           `json:"data_bn"`
	Additional *AdditionalDataForIndexInfo `json:"additional"`
}

type DataForIndexInfo struct {
	Name       string   `json:"name"`
	DeepLink   string   `json:"deep_link"`
	Slug       string   `json:"slug"`
	Image2x    string   `json:"image2x"`
	Image3x    string   `json:"image3x"`
	Tnc        string   `json:"tnc"`
	Tags       []string `json:"tags"`
	Status     any      `json:"status"`
	Attributes any      `json:"attributes"`
	SortOrder  any      `json:"sort_order"`
}

type AdditionalDataForIndexInfo struct {
	IosVersionMin      string `json:"ios_version_min"`
	IosVersionMax      string `json:"ios_version_max"`
	AndroidVersionMin  string `json:"android_version_min"`
	AndroidsVersionMax string `json:"android_version_max"`
}
