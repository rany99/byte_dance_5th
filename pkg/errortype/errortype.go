package errortype

// 通用错误信息
const (
	DataNotMatchErr    string = "数据格式错误"
	ParseUserIdErr     string = "user_id解析错误"
	ParseVideoIdErr    string = "video_id解析错误"
	ParseActionTypeErr string = "action_type解析错误"
	ParsePasswordErr   string = "password解析错误"
	PointerIsNilErr    string = "传入指针为空"
	UserNoExistErr     string = "用户不存在"
)

// token错误信息
const (
	TokenOutDateErr        string = "token已经过期"
	ParseTokenErr          string = "token解析错误"
	UserNoAuthenticatedErr string = "用户未认证"
	UserNameExistErr       string = "用户名已存在"
)

// 登录错误信息
const (
	UserWrongOrNoExistErr string = "账户输入错误或用户不存在"
	UserNameOverMaxLenErr string = "用户名超出最大长度"
	UserNameEmptyErr      string = "用户名为空"
	PasswordWrongErr      string = "密码输入错误"
	PasswordEmptyErr      string = "密码为空"
	SnowFlakeErr          string = "ID数量达到上限"
)

// 视频信息错误
const (
	FavorListEmptyErr string = "点赞列表为空"
	VideoNoExistErr   string = "不存在video_id的视频"
)

// ActionType错误信息
const (
	PostFavorActionTypeErr   string = "只能进行1:点赞;2:取消操作"
	PostCommentActionTypeErr string = "只能进行1:评论;2:删除评论"
	PostFollowActionTypeErr  string = "只能进行1:关注;2:取消关注"
)

// 投稿错误信息
const (
	ParseTitleErr     string = "视频标题解析错误"
	WrongVideoTypeErr string = "不支持视频格式:"
	VideoSaveErr      string = "视频保存失败"
)

// 封面截取错误信息
const (
	SnapShotErr  string = "截取封面失败"
	ImgDecodeErr string = "封面解码错误"
	SaveSnapErr  string = "封面保存失败"
)

// 关注操作
const (
	ParseToUserIdErr     string = "to_user_id解析错误"
	FollowUserNoExistErr string = "被关注的用户不存在"
	CantFollowSelfErr    string = "不能自己关注自己"
	FollowAgainErr       string = "请勿重复关注"
)

// 评论错误信息
const (
	CommentEmptyErr      string = "输入评论为空"
	VideoHasNoCommentErr string = "还没有人发现这里，赶紧抢首评吧"
	VideoListEmptyErr    string = "传入comments列表为空"
)

// 点赞错误信息
const (
	AlreadyPostFavorErr string = "您点的太快了，休息一下吧"
	FavorCountZeroErr   string = "点赞数目为0"
)

// 消息收发错误
const (
	EmptyMsgErr           string = "输入消息为空"
	ParseMsgFromUserIdErr string = "from_user_id解析错误"
	ParseMsgToUserIdErr   string = "to_user_id解析错误"
	FromUserNoExistErr    string = "消息发出者不存在"
	ToUserNoExistErr      string = "消息接受者不存在"
)
