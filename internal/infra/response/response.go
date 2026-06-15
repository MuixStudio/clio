/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package response

import (
	"net/http"

	"github.com/muixstudio/clio/internal/infra/errors"

	"github.com/gin-gonic/gin"
)

// SuccessOK writes a 200 OK response with a plain {"message":"ok"} body.
func SuccessOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func SuccessWithData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    data,
	})
}

func SuccessNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Fail(c *gin.Context, err error) {
	e := errors.ToError(err)
	c.JSON(int(e.GetCode()), gin.H{
		"reason":   e.GetReason(),
		"message":  e.GetMessage(),
		"metadata": e.GetMetadata(),
	})
}
