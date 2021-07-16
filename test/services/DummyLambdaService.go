package test_services

import (
	"encoding/json"

	awsserv "github.com/pip-services3-go/pip-services3-aws-go/services"
	awstest "github.com/pip-services3-go/pip-services3-aws-go/test"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type DummyLambdaService struct {
	*awsserv.LambdaService
	controller awstest.IDummyController
}

func NewDummyLambdaService() *DummyLambdaService {
	c := &DummyLambdaService{}
	c.LambdaService = awsserv.InheritLambdaService(c, "dummy")

	c.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *DummyLambdaService) SetReferences(references cref.IReferences) {
	c.LambdaService.SetReferences(references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(awstest.IDummyController)
	}
}

func (c *DummyLambdaService) getPageByFilter(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetPageByFilter(
		correlationId,
		cdata.NewFilterParamsFromValue(params["filter"]),
		cdata.NewPagingParamsFromValue(params["paging"]),
	)
}

func (c *DummyLambdaService) getOneById(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.GetOneById(
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaService) create(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Create(
		correlationId,
		entity,
	)
}

func (c *DummyLambdaService) update(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	val, _ := json.Marshal(params["dummy"])
	var entity awstest.Dummy
	json.Unmarshal(val, &entity)
	return c.controller.Update(
		correlationId,
		entity,
	)
}

func (c *DummyLambdaService) deleteById(params map[string]interface{}) (interface{}, error) {
	correlationId, _ := params["correlation_id"].(string)
	return c.controller.DeleteById(
		correlationId,
		params["dummy_id"].(string),
	)
}

func (c *DummyLambdaService) Register() {

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
