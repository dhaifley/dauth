package cmd

import (
	"net"
	"os"

	"github.com/spf13/viper"

	"github.com/dhaifley/dauth/server"
	"github.com/dhaifley/dlib"
	"github.com/dhaifley/dlib/ptypes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the application server",
	Long:  "The serve command starts the application server.",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.Server{Log: logrus.New()}
		s.Log.(*logrus.Logger).Out = os.Stdout
		s.Log.(*logrus.Logger).Formatter = new(logrus.JSONFormatter)
		err := s.ConnectSQL(nil)
		if err != nil {
			s.Log.Fatal(err.Error())
		}

		defer s.Close()
		lis, err := net.Listen("tcp", ":3612")
		if err != nil {
			s.Log.Fatal(err)
		}

		var opts []grpc.ServerOption
		creds, err := dlib.GetGRPCServerCredentials(viper.GetString("cert"), viper.GetString("key"))
		if err != nil {
			s.Log.Fatal(err)
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
		grpcServer := grpc.NewServer(opts...)
		ptypes.RegisterAuthServer(grpcServer, &s)
		s.Log.Fatal(grpcServer.Serve(lis))
	},
}
