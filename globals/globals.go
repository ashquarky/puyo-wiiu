package globals

import (
	"database/sql"
	pb "github.com/PretendoNetwork/grpc/go/account"
	"github.com/PretendoNetwork/nex-go/v2"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Postgres *sql.DB
var MatchmakingManager *common_globals.MatchmakingManager
var DatastoreManager *common_globals.DataStoreManager

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password

var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint

var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint

var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb.AccountClient
var GRPCAccountCommonMetadata metadata.MD

var MinIOClient *minio.Client
var S3Presigner *common_globals.MinIOPresigner
var S3Bucket string
var S3KeyBase string

var TokenAESKey []byte
var LocalAuthMode bool
