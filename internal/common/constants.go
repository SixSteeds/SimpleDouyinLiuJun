package common

const DefaultPass = "liujun"

/*
jwt相关
*/

const JwtSecret = "bGl1anVu"
const JwtExpire = 7200
const OrganizationName = " SixSteeds"

/*
MinIO相关
*/

const HTTP = "http://"
const MinIOEndPoint = "8.137.50.160:9000"
const MinIOAccessKey = "taozixun"
const MinIOSecretKey = "taozixun"
const MinIOVideoBucketName = "dousheng-video"
const MinIOCoverBucketName = "dousheng-cover"

/*
七牛云相关
*/

//const DefaultBackImage = " https://staraway.love/02c05856-e60c-44a0-9e85-6c2ba7ae9e0942d98a82be40d551ce6e123abbd0718292f80a71_raw.jpg"
//const DefaultAvatar = "https://staraway.love/031ec513-8976-45b5-80fd-8c725bc7ada7u%3D2169083367%2C64951360%26fm%3D253%26fmt%3Dauto%26app%3D138%26f%3DJPEG.webp"
//const AccessKey = "qYSQWz-EEJSJDAsg3Q7QqVy-GmpGLclb5lREXqO1"
//const SecretKey = "NkIglvJgXyoO6s1sh7rX2O1KBwltzLjuF4sPj3QY"
//const BucketName = "interestoriented"
//const QiliuyunDomain = "https://staraway.love/"

/*
雪花算法机器id
*/

const UserApiMachineId = 1
const UserRpcMachineId = 2
const MediaApiMachineId = 3
const MediaRpcMachineId = 4
const ContentApiMachineId = 5
const ContentRpcMachineId = 6
const ChatApiMachineId = 7
const ChatRpcMachineId = 8
const OtherMachineId = 9

/*
redis缓存前缀
*/

// CntCacheUserLikePrefix 计数服务缓存
const CntCacheUserLikePrefix = "CntCache:user_like_cnt:"
const CntCacheVideoLikedPrefix = "CntCache:video_liked_cnt:"
const CntCacheVideoCommentedPrefix = "CntCache:video_commented_cnt:"

// LikeCacheVideoLikedPrefix 点赞记录缓存
const LikeCacheVideoLikedPrefix = "LikeCache:video_liked:"

// CntCacheUserLikedPrefix 用户获赞数
const CntCacheUserLikedPrefix = "CntCache:User_Liked:"

// CntCacheUserWorkPrefix 用户作品总数
const CntCacheUserWorkPrefix = "CntCache:User_Work:"

// VideoCache2User 用videoId取userID
const VideoCache2User = "VideoCache:UserId:"

// FollowNum 用户关注数量
const FollowNum = "followNum_"

// FollowerNum 用户粉丝数量
const FollowerNum = "followerNum_"

const UploadLockPrefix = "UploadLock:"

/*
日志文件名
*/

const UserSecurity = "userSecurity"
const UploadSecurity = "uploadSecurity"

/*
雪花算法相关
*/

const Twepoch = int64(1483228800000)                  //开始时间截 (2017-01-01)
const WorkeridBits = uint(10)                         //机器id所占的位数
const SequenceBits = uint(12)                         //序列所占的位数
const WorkeridMax = int64(-1 ^ (-1 << WorkeridBits))  //支持的最大机器id数量
const SequenceMask = int64(-1 ^ (-1 << SequenceBits)) //
const WorkeridShift = SequenceBits                    //机器id左移位数
const TimestampShift = SequenceBits + WorkeridBits    //时间戳左移位数
