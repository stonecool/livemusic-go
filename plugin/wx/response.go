package wx

type PublicAccountResp struct {
	AppMsgCnt  int                    `json:"app_msg_cnt"`
	AppMsgList []PublicAccountRespMsg `json:"app_msg_list"`
}

type PublicAccountRespMsg struct {
	Aid                   string                        `json:"aid"`
	AlbumId               string                        `json:"album_id"`
	AppMsgAlbumInfos      []interface{}                 `json:"appmsg_album_infos"`
	AppMsgId              int                           `json:"appmsgid"`
	Checking              int                           `json:"checking"`
	CopyrightType         int                           `json:"copyright_type"`
	Cover                 string                        `json:"cover"`
	CreateTime            int                           `json:"create_time"`
	Digest                string                        `json:"digest"`
	HasRedPacketCover     int                           `json:"has_red_packet_cover"`
	IsPaySubscribe        int                           `json:"is_pay_subscribe"`
	ItemShowType          int                           `json:"item_show_type"`
	ItemIdx               int                           `json:"itemidx"`
	Link                  string                        `json:"link"`
	MediaDuration         string                        `json:"media_duration"`
	MediaApiPublishStatus int                           `json:"mediaapi_publish_status"`
	PayAlbumInfo          PublicAccountRespMsgAlbumInfo `json:"pay_album_info"`
	TagId                 []interface{}                 `json:"tagid"`
	Title                 string                        `json:"title"`
	UpdateTime            int                           `json:"update_time"`
}

type PublicAccountRespMsgAlbumInfo struct {
	AppMsgAlbumInfos []interface{} `json:"appmsg_album_infos"`
}
