package test
// let assert = require('chai').assert;
// let async = require('async');

// import { ConfigParams } from 'pip-services3-commons-node';

// import { Dummy } from '../Dummy';
// import { DummyController } from '../DummyController';
// import { DummyCommandableLambdaFunction } from './DummyCommandableLambdaFunction';

// suite('DummyCommandableLambdaFunction', () => {
//     let DUMMY1: Dummy = { id: null, key: "Key 1", content: "Content 1" };
//     let DUMMY2: Dummy = { id: null, key: "Key 2", content: "Content 2" };

//     let lambda: DummyCommandableLambdaFunction;

//     suiteSetup((done) => {
//         let config = ConfigParams.fromTuples(
//             'logger.descriptor', 'pip-services:logger:console:default:1.0',
//             'controller.descriptor', 'pip-services-dummies:controller:default:default:1.0'
//         );

//         lambda = new DummyCommandableLambdaFunction();
//         lambda.configure(config);
//         lambda.open(null, done);
//     });

//     suiteTeardown((done) => {
//         lambda.close(null, done);
//     });

//     test('CRUD Operations', (done) => {
//         var dummy1, dummy2;

//         async.series([
//             // Create one dummy
//             (callback) => {
//                 lambda.act(
//                     {
//                         cmd: 'create_dummy',
//                         dummy: DUMMY1
//                     },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.content, DUMMY1.content);
//                         assert.equal(dummy.key, DUMMY1.key);

//                         dummy1 = dummy;

//                         callback();
//                     }
//                 );
//             },
//             // Create another dummy
//             (callback) => {
//                 lambda.act(
//                     {
//                         cmd: 'create_dummy',
//                         dummy: DUMMY2
//                     },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.content, DUMMY2.content);
//                         assert.equal(dummy.key, DUMMY2.key);

//                         dummy2 = dummy;

//                         callback();
//                     }
//                 );
//             },
//             // Get all dummies
//             (callback) => {
//                 lambda.act(
//                     {
//                         cmd: 'get_dummies'
//                     },
//                     (err, dummies) => {
//                         assert.isNull(err);

//                         assert.isObject(dummies);
//                         assert.lengthOf(dummies.data, 2);

//                         callback();
//                     }
//                 );
//             },
//             // Update the dummy
//             (callback) => {
//                 dummy1.content = 'Updated Content 1'
//                 lambda.act(
//                     {
//                         cmd: 'update_dummy',
//                         dummy: dummy1
//                     },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isObject(dummy);
//                         assert.equal(dummy.id, dummy1.id);
//                         assert.equal(dummy.content, dummy1.content);
//                         assert.equal(dummy.key, dummy1.key);

//                         callback();
//                     }
//                 );
//             },
//             // Delete dummy
//             (callback) => {
//                 lambda.act(
//                     {
//                         cmd: 'delete_dummy',
//                         dummy_id: dummy1.id
//                     },
//                     (err) => {
//                         assert.isNull(err);

//                         callback();
//                     }
//                 );
//             },
//             // Try to get delete dummy
//             (callback) => {
//                 lambda.act(
//                     {
//                         cmd: 'get_dummy_by_id',
//                         dummy_id: dummy1.id
//                     },
//                     (err, dummy) => {
//                         assert.isNull(err);

//                         assert.isNull(dummy || null);

//                         callback();
//                     }
//                 );
//             }
//         ], done);
//     });

// });
