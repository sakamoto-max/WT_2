package mappings

import exerpb "github.com/sakamoto-max/wt_2_proto/shared/exercise"

type GetExerciseByName struct {
	UserId       string
	ExerciseName string
}

func ToGetExerciesByName(in *exerpb.SendExerciseName) GetExerciseByName {
	return GetExerciseByName{
		UserId:       in.UserId,
		ExerciseName: in.ExerciseName,
	}
}

type CreateExercise struct {
	UserId        string
	ExerciseName  string
	BodyPartName  string
	EquipmentName string
}

func ToCreateExercise(in *exerpb.CreateExerciseReq) CreateExercise {
	return CreateExercise{
		UserId:        in.UserId,
		ExerciseName:  in.ExerciseName,
		BodyPartName:  in.BodyPart,
		EquipmentName: in.Equipment,
	}
}

type DeleteExercise struct {
	UserId       string
	ExerciseName string
}

func ToDeleteExercise(in *exerpb.SendExerciseName) DeleteExercise {
	return DeleteExercise{
		UserId:       in.UserId,
		ExerciseName: in.ExerciseName,
	}
}

type ExerciseExistsReturnId struct {
	UserId       string
	ExerciseName string
}

func ToExerciseExistsReturnId(in *exerpb.SendExerciseName) ExerciseExistsReturnId {
	return ExerciseExistsReturnId{
		UserId: in.UserId,
		ExerciseName: in.ExerciseName,
	}
}