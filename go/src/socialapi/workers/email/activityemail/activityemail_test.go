package activityemail

import (
	// "github.com/kr/pretty"
	"koding/db/mongodb/modelhelper"
	"socialapi/config"
	"socialapi/models"
	"socialapi/rest"
	"testing"
	"time"

	"github.com/koding/runner"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSaveDailyDigestNotification(t *testing.T) {
	r := runner.New("email notification")
	if err := r.Init(); err != nil {
		panic(err)
	}
	defer r.Close()

	// initialize mongo
	appConfig := config.MustRead(r.Conf.Path)
	modelhelper.Initialize(appConfig.Mongo)

	// initialize redis
	redisConn := r.Bongo.MustGetRedisConn()

	Convey("User replies to another user who has daily digests", t, func() {
		acc1, err := rest.CreateAccountWithDailyDigest()
		So(err, ShouldBeNil)

		acc2, err := models.CreateAccountInBothDbs()
		So(err, ShouldBeNil)

		groupName := models.RandomGroupName()

		// fetch admin's session
		ses, err := models.FetchOrCreateSession(acc1.Nick, groupName)
		So(err, ShouldBeNil)
		So(ses, ShouldNotBeNil)

		channel, err := rest.CreateChannelByGroupNameAndType(
			acc1.Id,
			groupName,
			models.Channel_TYPE_TOPIC,
			ses.ClientId,
		)
		So(err, ShouldBeNil)

		cp := rest.NewCreatePost(channel.Id, acc1.Id, acc2.Id, ses.ClientId)
		_, err = cp.CreateReplies(2)
		So(err, ShouldBeNil)

		// sleep to wait for notification worker to its thing
		time.Sleep(200 * time.Millisecond)

		key := prepareSetterCacheKey(acc1.Id)
		activityIds, err := redisConn.GetSetMembers(key)
		So(err, ShouldBeNil)

		Convey("it should save user activity in redis", func() {
			So(len(activityIds), ShouldEqual, 2)
		})
	})
}
