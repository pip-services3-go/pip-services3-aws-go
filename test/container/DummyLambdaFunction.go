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
	c.LambdaFunction = awscont.InheriteLambdaFunction(c, "dummy", "Dummy lambda function")

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
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetPageByFilter(
		correlationId,
		cdata.NewFilterParamsFromValue(params["filter"]),
		cdata.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaFunction) getOneById(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetOneById(
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaFunction) create(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)

	c.Logger().Debug(correlationId, "Create method called Dummy %v", entity)

	res, err := c.controller.Create(
		correlationId,
		entity,
	)

	c.Logger().Debug(correlationId, "Create method called Result: %v Err: %v", res, err)

	return res, err
}

func (c *DummyLambdaFunction) update(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity = awstest.Dummy{}
	json.Unmarshal(val, &entity)
	return c.controller.Update(
		correlationId,
		entity,
	)
}

func (c *DummyLambdaFunction) deleteById(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)

	c.Logger().Debug(correlationId, "DeleteById method called Id %v", params["dummy_id"].(string))

	res, err := c.controller.DeleteById(
		correlationId,
		params["dummy_id"].(string),
	)
	c.Logger().Debug(correlationId, "DeleteById method called Result: %v Err: %v", res, err)

	return res, err
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
