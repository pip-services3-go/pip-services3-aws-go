package test

// import { CommandSet } from 'pip-services3-commons-node';
// import { ICommand } from 'pip-services3-commons-node';
// import { Command } from 'pip-services3-commons-node';
// import { Parameters } from 'pip-services3-commons-node';
// import { FilterParams } from 'pip-services3-commons-node';
// import { PagingParams } from 'pip-services3-commons-node';
// import { ObjectSchema } from 'pip-services3-commons-node';
// import { Schema} from 'pip-services3-commons-node';
// import { MapSchema } from 'pip-services3-commons-node';
// import { TypeCode } from 'pip-services3-commons-node';
// import { FilterParamsSchema } from 'pip-services3-commons-node';
// import { PagingParamsSchema } from 'pip-services3-commons-node';

// import { Dummy } from './Dummy';
// import { IDummyController } from './IDummyController';
// import { DummySchema } from './DummySchema';

// export class DummyCommandSet extends CommandSet {
//     private _controller: IDummyController;

// 	constructor(controller: IDummyController) {
// 		super();

// 		this._controller = controller;

// 		this.addCommand(this.makeGetPageByFilterCommand());
// 		this.addCommand(this.makeGetOneByIdCommand());
// 		this.addCommand(this.makeCreateCommand());
// 		this.addCommand(this.makeUpdateCommand());
// 		this.addCommand(this.makeDeleteByIdCommand());
// 	}

// 	private makeGetPageByFilterCommand(): ICommand {
// 		return new Command(
// 			"get_dummies",
// 			new ObjectSchema(true)
//                 .withOptionalProperty("filter", new FilterParamsSchema())
//                 .withOptionalProperty("paging", new PagingParamsSchema()),
// 			(correlationId: string, args: Parameters, callback: (err: any, result: any) => void) => {
// 				let filter = FilterParams.fromValue(args.get("filter"));
// 				let paging = PagingParams.fromValue(args.get("paging"));
// 				this._controller.getPageByFilter(correlationId, filter, paging, callback);
// 			}
// 		);
// 	}

// 	private makeGetOneByIdCommand(): ICommand {
// 		return new Command(
// 			"get_dummy_by_id",
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy_id", TypeCode.String),
// 			(correlationId: string, args: Parameters, callback: (err: any, result: any) => void) => {
// 				let id = args.getAsString("dummy_id");
// 				this._controller.getOneById(correlationId, id, callback);
// 			}
// 		);
// 	}

// 	private makeCreateCommand(): ICommand {
// 		return new Command(
// 			"create_dummy",
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy", new DummySchema()),
// 			(correlationId: string, args: Parameters, callback: (err: any, result: any) => void) => {
// 				let entity: Dummy = args.get("dummy");
// 				this._controller.create(correlationId, entity, callback);
// 			}
// 		);
// 	}

// 	private makeUpdateCommand(): ICommand {
// 		return new Command(
// 			"update_dummy",
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy", new DummySchema()),
// 			(correlationId: string, args: Parameters, callback: (err: any, result: any) => void) => {
// 				let entity: Dummy = args.get("dummy");
// 				this._controller.update(correlationId, entity, callback);
// 			}
// 		);
// 	}

// 	private makeDeleteByIdCommand(): ICommand {
// 		return new Command(
// 			"delete_dummy",
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy_id", TypeCode.String),
// 			(correlationId: string, args: Parameters, callback: (err: any, result: any) => void) => {
// 				let id = args.getAsString("dummy_id");
// 				this._controller.deleteById(correlationId, id, callback);
// 			}
// 		);
// 	}

// }
