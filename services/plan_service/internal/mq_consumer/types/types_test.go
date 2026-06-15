package types

import (
	"testing"

	mqTypes "github.com/sakamoto-max/rabbit_mq/types"
	"github.com/sakamoto-max/wt_2_proto/shared/enum"
	"github.com/stretchr/testify/assert"
)

func Test_ToData(t *testing.T) {

	mqData := mqTypes.Data{
		DbId:          "123",
		SentBy:        "sent",
		TaskName:      "task",
		Payload:       map[string]any{"user": "123"},
		TargetService: "target",
	}

	data := ToData(mqData)

	assert.Equal(t, data.DbId, mqData.DbId)
	assert.Equal(t, data.SentBy, mqData.SentBy)
	assert.Equal(t, data.TaskName, mqData.TaskName)
	assert.Equal(t, data.Payload, mqData.Payload)
	assert.Equal(t, data.TargetService, mqData.TargetService)
}

func Test_GetUserId(t *testing.T) {
	data := Data{
		DbId:          "123",
		SentBy:        "sent",
		TaskName:      "task",
		Payload:       map[string]any{enum.QueueFields_USER_ID.String(): "123"},
		TargetService: "target",
	}

	userId, err := data.GetUserId()

	assert.NoError(t, err)

	assert.Equal(t, userId, "123")
}

func Test_GetPlanName(t *testing.T) {

	data := Data{
		DbId:          "123",
		SentBy:        "sent",
		TaskName:      "task",
		Payload:       map[string]any{enum.QueueFields_PLAN_NAME.String(): "123"},
		TargetService: "target",
	}

	planName, err := data.GetPlanName()
	assert.NoError(t, err)

	assert.Equal(t, planName, "123")
}

func Test_GetNewExercises(t *testing.T) {

	exers := []any{"exer_1", "exer_2", "exer_3"}

	payload := map[string]any{
		enum.QueueFields_EXERCISE_NAMES.String(): exers,
	}

	data := Data{
		DbId:          "123",
		SentBy:        "sent",
		TaskName:      "task",
		Payload:       payload,
		TargetService: "target",
	}

	allExers, err := data.GetNewExercises()
	assert.NoError(t, err)

	assert.Equal(t, allExers[0], "exer_1")
	assert.Equal(t, allExers[1], "exer_2")
	assert.Equal(t, allExers[2], "exer_3")

}
