package test

// let assert = require('chai').assert;

// import { LogLevel } from 'pip-services3-components-node';
// import { CachedLogger } from 'pip-services3-components-node';

// export class LoggerFixture {
//     private _logger: CachedLogger;

//     public constructor(logger: CachedLogger) {
//         this._logger = logger;
//     }

//     public testLogLevel() {
//         assert.isTrue(this._logger.getLevel() >= LogLevel.None);
//         assert.isTrue(this._logger.getLevel() <= LogLevel.Trace);
//     }

//     public testSimpleLogging(done) {
//         this._logger.setLevel(LogLevel.Trace);

//         this._logger.fatal(null, null, "Fatal error message");
//         this._logger.error(null, null, "Error message");
//         this._logger.warn(null, "Warning message");
//         this._logger.info(null, "Information message");
//         this._logger.debug(null, "Debug message");
//         this._logger.trace(null, "Trace message");

//         this._logger.dump();
//         setTimeout(done, 1000);
//     }

//     public testErrorLogging(done) {
//         try {
//             // Raise an exception
//             throw new Error();
//         } catch (ex) {
//             this._logger.fatal("123", ex, "Fatal error");
//             this._logger.error("123", ex, "Recoverable error");

//             assert.isNotNull(ex);
//         }

//         this._logger.dump();
//         setTimeout(done, 1000);
//     }

// }
