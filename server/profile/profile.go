package profile

import (
	// "fmt"
	configs "github.com/aiwantaozi/microservice_profile/configs"
	proto_profile "github.com/aiwantaozi/microservice_profile/proto/profile"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
)

// Implements of ActivityServiceServer
type profileServer struct{}

// const (
//  SEX_MALE       = 0
//  SEX_FEMALE     = 1
//  DEFAULT_WEIGHT = 140
//  DEFAULT_HEIGHT = 64
//  DEFAULT_AGE    = 32
//  DEFAULT_SEX    = SEX_MALE
// )

func NewProfileServer() proto_profile.ProfileServiceServer {
	return new(profileServer)
}

func (s *profileServer) ProfileInstancePost(stream proto_profile.ProfileService_ProfileInstancePostServer) error {
	profileCollection := profileCol()
	defer profileCollection.Database.Session.Close()

	msg, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	md, _ := metadata.FromContext(stream.Context())
	msg.Uid = md["uid"][0]

	_, err = profileCollection.Upsert(bson.M{"uid": msg.Uid}, msg)
	if err != nil {
		return err
	}
	return stream.Send(msg)
}

func (s *profileServer) ProfileInstanceGet(_ *empty.Empty, stream proto_profile.ProfileService_ProfileInstanceGetServer) error {
	profileCollection := profileCol()
	defer profileCollection.Database.Session.Close()

	md, _ := metadata.FromContext(stream.Context())
	uid := md["uid"][0]

	var profile proto_profile.ProfileInstance
	err := profileCollection.Find(bson.M{"uid": uid}).One(&profile)
	if err != nil {
		return err
	}

	return stream.Send(&profile)
}

func profileCol() *mgo.Collection {
	return configs.GetDatabase().C("profiles")
}

func validate(msg proto_profile.ProfileInstance) error {
	return nil
}
