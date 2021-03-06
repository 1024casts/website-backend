package course

import (
	"strconv"

	"github.com/1024casts/1024casts/pkg/app"
	"github.com/1024casts/1024casts/pkg/errno"
	"github.com/1024casts/1024casts/util"

	"github.com/1024casts/1024casts/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// @Summary Update a course info by the course identifier
// @Description Update a course by ID
// @Tags course
// @Accept  json
// @Produce  json
// @Param id path uint64 true "The course's database id index num"
// @Param user body model.CourseModel true "The course info"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /courses/{id} [put]
func Update(c *gin.Context) {
	log.Info("Update function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	// Get the course id from the url parameter.
	courseId, _ := strconv.Atoi(c.Param("id"))

	// Binding the course data.
	var course CreateRequest
	if err := c.Bind(&course); err != nil {
		app.Response(c, errno.ErrBind, nil)
		return
	}

	srv := service.NewCourseService()
	//user, err := model.GetUserById(userId)
	_, err := srv.GetCourseById(courseId)
	if err != nil {
		app.Response(c, errno.ErrCourseNotFound, nil)
		log.Warn("course info", lager.Data{"id": courseId})
		return
	}

	// Validate the data.
	//if err := u.Validate(); err != nil {
	//	app.Response(c, errno.ErrValidation, nil)
	//	return
	//}

	// Save changed fields.
	courseMap := make(map[string]interface{}, 0)
	courseMap["name"] = course.Name
	courseMap["type"] = course.Type
	courseMap["keywords"] = course.Keywords
	courseMap["description"] = course.Description
	courseMap["slug"] = course.Slug
	courseMap["content"] = course.Content
	courseMap["cover_key"] = course.CoverKey
	courseMap["update_status"] = course.UpdateStatus

	if err := srv.UpdateCourse(courseMap, courseId); err != nil {
		app.Response(c, errno.InternalServerError, nil)
		return
	}

	app.Response(c, nil, nil)
}
