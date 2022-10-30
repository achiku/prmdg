package taskyapi

import jsval "github.com/lestrrat-go/go-jsval"

var TaskCreateValidator *jsval.JSVal
var TaskInstancesValidator *jsval.JSVal
var TaskSelfValidator *jsval.JSVal
var UserSelfValidator *jsval.JSVal
var M *jsval.ConstraintMap
var R0 jsval.Constraint
var R1 jsval.Constraint

func init() {
	M = &jsval.ConstraintMap{}
	R0 = jsval.Array().
		Items(
			jsval.String(),
		).
		AdditionalItems(
			jsval.EmptyConstraint,
		)
	R1 = jsval.String()
	M.SetReference("#/definitions/task/definitions/tags", R0)
	M.SetReference("#/definitions/task/definitions/title", R1)
	TaskCreateValidator = jsval.New().
		SetName("TaskCreateValidator").
		SetConstraintMap(M).
		SetRoot(
			jsval.Object().
				Required("title").
				AdditionalProperties(
					jsval.EmptyConstraint,
				).
				AddProp(
					"tags",
					jsval.Reference(M).RefersTo("#/definitions/task/definitions/tags"),
				).
				AddProp(
					"title",
					jsval.Reference(M).RefersTo("#/definitions/task/definitions/title"),
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
