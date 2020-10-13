package test

// import { FilterParams } from 'pip-services3-commons-node';
// import { PagingParams } from 'pip-services3-commons-node';
// import { DataPage } from 'pip-services3-commons-node';

// import { LambdaClient } from '../../src/clients/LambdaClient';
// import { IDummyClient } from '../IDummyClient';
// import { Dummy } from '../Dummy';

// export class DummyLambdaClient extends LambdaClient implements IDummyClient {

//     public constructor() {
//         super();
//     }

//     public getDummies(correlationId: string, filter: FilterParams, paging: PagingParams,
//         callback: (err: any, result: DataPage<Dummy>) => void): void {
//         this.call(
//             'get_dummies',
//             correlationId,
//             {
//                 filter: filter,
//                 paging: paging
//             },
//             (err, result) => {
//                 callback(err, result);
//             }
//         );
//     }

//     public getDummyById(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call(
//             'get_dummy_by_id',
//             correlationId,
//             {
//                 dummy_id: dummyId
//             },
//             (err, result) => {
//                 callback(err, result);
//             }
//         );
//     }

//     public createDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call(
//             'create_dummy',
//             correlationId,
//             {
//                 dummy: dummy
//             },
//             (err, result) => {
//                 callback(err, result);
//             }
//         );
//     }

//     public updateDummy(correlationId: string, dummy: any,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call(
//             'update_dummy',
//             correlationId,
//             {
//                 dummy: dummy
//             },
//             (err, result) => {
//                 callback(err, result);
//             }
//         );
//     }

//     public deleteDummy(correlationId: string, dummyId: string,
//         callback: (err: any, result: Dummy) => void): void {
//         this.call(
//             'delete_dummy',
//             correlationId,
//             {
//                 dummy_id: dummyId
//             },
//             (err, result) => {
//                 callback(err, result);
//             }
//         );
//     }

// }
