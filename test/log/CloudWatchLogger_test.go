package test

// import { ConfigParams } from 'pip-services3-commons-node';
// import { References } from 'pip-services3-commons-node';
// import { ContextInfo } from 'pip-services3-components-node';
// import { Descriptor } from 'pip-services3-commons-node';

// import { CloudWatchLogger } from '../../src/log/CloudWatchLogger';
// import { LoggerFixture } from './LoggerFixture';

// suite('CloudWatchLogger', ()=> {
//     let _logger: CloudWatchLogger;
//     let _fixture: LoggerFixture;

//     let AWS_REGION = process.env["AWS_REGION"] || "";
//     let AWS_ACCESS_ID = process.env["AWS_ACCESS_ID"] || "";
//     let AWS_ACCESS_KEY = process.env["AWS_ACCESS_KEY"] || "";

//     if (!AWS_REGION || !AWS_ACCESS_ID || !AWS_ACCESS_KEY)
//         return;

//     setup((done) => {

//         _logger = new CloudWatchLogger();
//         _fixture = new LoggerFixture(_logger);

//         let config = ConfigParams.fromTuples(
//             "group", "TestGroup",
//             "connection.region", AWS_REGION,
//             "credential.access_id", AWS_ACCESS_ID,
//             "credential.access_key", AWS_ACCESS_KEY
//         );
//         _logger.configure(config);

//         var contextInfo = new ContextInfo();
//         contextInfo.name = "TestStream";

//         var references = References.fromTuples(
//             new Descriptor("pip-services", "context-info", "default", "default", "1.0"), contextInfo,
//             new Descriptor("pip-services", "counters", "cloudwatch", "default", "1.0"), _logger
//         );
//         _logger.setReferences(references);

//         _logger.open(null, (err) => {
//              done(err);
//         });
//     });

//     teardown((done) => {
//         _logger.close(null, done);
//     });

//     test('Log Level', () => {
//         _fixture.testLogLevel();
//     });

//     test('Simple Logging', (done) => {
//         _fixture.testSimpleLogging(done);
//     });

//     test('Error Logging', (done) => {
//         _fixture.testErrorLogging(done);
//     });

// });
