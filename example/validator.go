package exampleapi

import "github.com/lestrrat/go-jsval"

var TaskCreateValidator *jsval.JSVal
var TaskInstancesValidator *jsval.JSVal
var TaskSelfValidator *jsval.JSVal
var UserSelfValidator *jsval.JSVal
var M *jsval.ConstraintMap
var R0 jsval.Constraint
var R1 jsval.Constraint
var R2 jsval.Constraint

func init() {
	M = &jsval.ConstraintMap{}
	R0 = jsval.String()
	R1 = jsval.Reference(M).RefersTo("#/definitions/user/definitions/id")
	R2 = jsval.String().Format("uuid")
	M.SetReference("#/definitions/task/definitions/name", R0)
	M.SetReference("#/definitions/task/definitions/userId", R1)
	M.SetReference("#/definitions/user/definitions/id", R2)
	TaskCreateValidator = jsval.New().
		SetName("TaskCreateValidator").
		SetConstraintMap(M).
		SetRoot(
			jsval.Object().
				Required("name", "userId").
				AdditionalProperties(
					jsval.EmptyConstraint,
				).
				AddProp(
					"name",
					jsval.Reference(M).RefersTo("#/definitions/task/definitions/name"),
				).
				AddProp(
					"userId",
					jsval.Reference(M).RefersTo("#/definitions/task/definitions/userId"),
				),
		)

	TaskInstancesValidator = jsval.New().
		SetName("TaskInstancesValidator").
		SetConstraintMap(M).
		SetRoot(
			jsval.Object().
				AdditionalProperties(
					jsval.EmptyConstraint,
				).
				AddProp(
					"limit",
					jsval.Integer(),
				).
				AddProp(
					"offset",
					jsval.Integer(),
				),
		)

	TaskSelfValidator = jsval.New().
		SetName("TaskSelfValidator").
		SetConstraintMap(M).
		SetRoot(
			jsval.EmptyConstraint,
		)

	UserSelfValidator = jsval.New().
		SetName("UserSelfValidator").
		SetConstraintMap(M).
		SetRoot(
			jsval.EmptyConstraint,
		)

}
