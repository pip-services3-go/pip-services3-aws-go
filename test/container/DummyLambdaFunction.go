package test_container

import (
	"encoding/json"

	awscont "github.com/pip-services3-go/pip-services3-aws-go/container"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type DummyLambdaFunction struct {
	*awscont.LambdaFunction
	controller awstest.IDummyController
}

func NewDummyLambdaFunction() *DummyLambdaFunction {
	c := &DummyLambdaFunction{}
	c.LambdaFunction = awscont.InheriteLambdaFunction("dummy", "Dummy lambda function", c)

	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	c.AddFactory(awstest.NewDummyFactory())
	return c
}

func (c *DummyLambdaFunction) SetReferences(references cref.IReferences) {
	c.LambdaFunction.SetReferences(references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(awstest.IDummyController)
	}
}

func (c *DummyLambdaFunction) getPageByFilter(params map[string]interface{}) (interface{}, error) {
	return c.controller.GetPageByFilter(
		params["correlation_id"].(string),
		cdata.NewFilterParamsFromValue(params["filter"]),
		cdata.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaFunction) getOneById(params map[string]interface{}) (interface{}, error) {
	return c.controller.GetOneById(
		params["correlation_id"].(string),
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaFunction) create(params map[string]interface{}) (interface{}, error) {
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Create(
		params["correlation_id"].(string),
		entity,
	)
}

func (c *DummyLambdaFunction) update(params map[string]interface{}) (interface{}, error) {
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Update(
		params["correlation_id"].(string),
		entity,
	)
}

func (c *DummyLambdaFunction) deleteById(params map[string]interface{}) (interface{}, error) {
	return c.controller.DeleteById(
		params["correlation_id"].(string),
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaFunction) Register() {
	c.RegisterAction(
		"get_dummies",
		&cvalid.NewObjectSchema().
			WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).
			WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()).Schema,
		c.getPageByFilter)

	c.RegisterAction(
		"get_dummy_by_id",
		&cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.getOneById)

	c.RegisterAction(
		"create_dummy",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.create)

	c.RegisterAction(
		"update_dummy",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("dummy", awstest.NewDummySchema()).Schema,
		c.update)

	c.RegisterAction(
		"delete_dummy",
		&cvalid.NewObjectSchema().
			WithOptionalProperty("dummy_id", cconv.String).Schema,
		c.deleteById)
}

// export const handler = new DummyLambdaFunction().getHandler();
