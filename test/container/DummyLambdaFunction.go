package test

// let _ = require('lodash');

// import { Descriptor } from 'pip-services3-commons-node';
// import { FilterParams } from 'pip-services3-commons-node';
// import { PagingParams} from 'pip-services3-commons-node';
// import { IReferences } from 'pip-services3-commons-node';
// import { ObjectSchema } from 'pip-services3-commons-node';
// import { Schema} from 'pip-services3-commons-node';
// import { MapSchema } from 'pip-services3-commons-node';
// import { TypeCode } from 'pip-services3-commons-node';
// import { FilterParamsSchema } from 'pip-services3-commons-node';
// import { PagingParamsSchema } from 'pip-services3-commons-node';

// import { LambdaFunction } from '../../src/container/LambdaFunction';
// import { IDummyController } from '../IDummyController';
// import { DummyFactory } from '../DummyFactory';
// import { DummySchema } from '../DummySchema';

// export class DummyLambdaFunction extends LambdaFunction {
//     private _controller: IDummyController;

//     public constructor() {
//         super("dummy", "Dummy lambda function");
//         this._dependencyResolver.put('controller', new Descriptor('pip-services-dummies', 'controller', 'default', '*', '*'));
//         this._factories.add(new DummyFactory());
//     }

//     public setReferences(references: IReferences): void {
//         super.setReferences(references);
//         this._controller = this._dependencyResolver.getOneRequired<IDummyController>('controller');
//     }

//     private getPageByFilter(params: any, callback: (err: any, result?: any) => void): void {
//         this._controller.getPageByFilter(
//             params.correlation_id,
//             new FilterParams(params.filter),
//             new PagingParams(params.paging),
//             callback
//         );
//     }

//     private getOneById(params: any, callback: (err: any, result?: any) => void): void {
//         this._controller.getOneById(
//             params.correlation_id,
//             params.dummy_id,
//             callback
//         );
//     }

//     private create(params: any, callback: (err: any, result?: any) => void): void {
//         this._controller.create(
//             params.correlation_id,
//             params.dummy,
//             callback
//         );
//     }

//     private update(params: any, callback: (err: any, result?: any) => void): void {
//         this._controller.update(
//             params.correlation_id,
//             params.dummy,
//             callback
//         );
//     }

//     private deleteById(params: any, callback: (err: any, result?: any) => void): void {
//         this._controller.deleteById(
//             params.correlation_id,
//             params.dummy_id,
//             callback
//         );
//     }

//     protected register() {
//         this.registerAction(
//             'get_dummies',
//             new ObjectSchema(true)
//                 .withOptionalProperty("filter", new FilterParamsSchema())
//                 .withOptionalProperty("paging", new PagingParamsSchema())
//             , this.getPageByFilter);

//         this.registerAction(
//             'get_dummy_by_id',
//             new ObjectSchema(true)
//                 .withOptionalProperty("dummy_id", TypeCode.String)
//             , this.getOneById);

//         this.registerAction(
//             'create_dummy',
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy", new DummySchema())
//             , this.create);

//         this.registerAction(
//             'update_dummy',
//             new ObjectSchema(true)
//                 .withRequiredProperty("dummy", new DummySchema())
//             , this.update);

//         this.registerAction(
//             'delete_dummy',
//             new ObjectSchema(true)
//                 .withOptionalProperty("dummy_id", TypeCode.String)
//             , this.deleteById);
//     }
// }

// export const handler = new DummyLambdaFunction().getHandler();
