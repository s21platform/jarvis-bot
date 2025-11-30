package model

import "github.com/samber/lo"

const (
	COMMUNITY_PROJECT = "COM"
	MATERIALS_PROJECT = "MTR"
	JARVIS_PROJECT    = "JRV"
	FRONTEND_PROJECT  = "FNT"
	TECH_STREAM       = "TECH"
	EVO_PROJECT       = "EVO"
	USER_PROJECT      = "USR"
	STORAGE_PROJECT   = "STG"
	ADVERT_PROJECT    = "ADV"
	AUTH_PROJECT      = "AUTH"
	FEED_PROJECT      = "FED"
	SEARCH_PROJECT    = "SRC"
	SOCIETY_PROJECT   = "STY"
	OPTIONHUB_PROJECT = "OPT"
	FLAGMAN_PROJECT   = "FLG"
)

var (
	LOGGER_LABEL  = lo.ToPtr("logger")
	COMMON_LABEL  = lo.ToPtr("common")
	METRICS_LABEL = lo.ToPtr("metrics")
	GATEWAY_LABEL = lo.ToPtr("gateway")
	KAFKA_LABEL   = lo.ToPtr("kafka")
)
