package test

// import { ConfigParams } from 'pip-services3-commons-node';
// import { References } from 'pip-services3-commons-node';
// import { ContextInfo } from 'pip-services3-components-node';
// import { Descriptor } from 'pip-services3-commons-node';

// import { CloudWatchCounters } from '../../src/count/CloudWatchCounters';
// import { CountersFixture } from './CountersFixture';

// suite('CloudWatchCounters', ()=> {
//     let _counters: CloudWatchCounters;
//     let _fixture: CountersFixture;

//     let AWS_REGION = process.env["AWS_REGION"] || "";
//     let AWS_ACCESS_ID = process.env["AWS_ACCESS_ID"] || "";
//     let AWS_ACCESS_KEY = process.env["AWS_ACCESS_KEY"] || "";

//     if (!AWS_REGION || !AWS_ACCESS_ID || !AWS_ACCESS_KEY)
//         return;

//     setup((done) => {

//         _counters = new CloudWatchCounters();
//         _fixture = new CountersFixture(_counters);

//         let config = ConfigParams.fromTuples(
//             "interval", "5000",
//             "connection.region", AWS_REGION,
//             "credential.access_id", AWS_ACCESS_ID,
//             "credential.access_key", AWS_ACCESS_KEY
//         );
//         _counters.configure(config);

//         var contextInfo = new ContextInfo();
//         contextInfo.name = "Test";
//         contextInfo.description = "This is a test container";

//         var references = References.fromTuples(
//             new Descriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
//             new Descriptor("pip-services", "counters", "cloudwatch", "default", "1.0"), _counters
//         );
//         _counters.setReferences(references);

//         _counters.open(null, (err) => {
//              done(err);
//         });
//     });

//     teardown((done) => {
//         _counters.close(null, done);
//     });

//     test('Simple Counters', (done) => {
//         _fixture.testSimpleCounters(done);
//     });

//     test('Measure Elapsed Time', (done) => {
//         _fixture.testMeasureElapsedTime(done);
//     });

// });
