package common

const DefaultPass = "liujun"

// jwt相关
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
const DefaultBackImage = " https://staraway.love/02c05856-e60c-44a0-9e85-6c2ba7ae9e0942d98a82be40d551ce6e123abbd0718292f80a71_raw.jpg"
const DefaultAvatar = "https://staraway.love/031ec513-8976-45b5-80fd-8c725bc7ada7u%3D2169083367%2C64951360%26fm%3D253%26fmt%3Dauto%26app%3D138%26f%3DJPEG.webp"
const AccessKey = "qYSQWz-EEJSJDAsg3Q7QqVy-GmpGLclb5lREXqO1"
const SecretKey = "NkIglvJgXyoO6s1sh7rX2O1KBwltzLjuF4sPj3QY"
const BucketName = "interestoriented"
const QiliuyunDomain = "https://staraway.love/"

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

/*
redis缓存前缀
*/

// 计数服务缓存
const CntCacheUserLikePrefix = "CntCache:user_like_cnt:"
const CntCacheVideoLikedPrefix = "CntCache:video_liked_cnt:"
const CntCacheVideoCommentedPrefix = "CntCache:video_commented_cnt:"

// 点赞记录缓存
const LikeCacheVideoLikedPrefix = "LikeCache:video_liked:"

// 用户获赞数
const CntCacheUserLikedPrefix = "CntCache:User_Liked:"

// 用户作品总数
const CntCacheUserWorkPrefix = "CntCache:User_Work:"

// 用videoId取userID
const VideoCache2User = "VideoCache:UserId:"

// 用户关注数量
const FollowNum = "followNum_"

// 用户粉丝数量
const FollowerNum = "followerNum_"

/*
日志文件名
*/

const USER_SECURITY = "userSecurity"
const UPLOAD_SECURITY = "uploadSecurity"
