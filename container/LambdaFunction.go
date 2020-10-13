package container

// /*@module container */
// /*@hidden */
// let _ = require('lodash');
// /*@hidden */
// let process = require('process');

// import { ConfigParams } from 'pip-services3-commons-node';
// import { IReferences } from 'pip-services3-commons-node';
// import { DependencyResolver } from 'pip-services3-commons-node';
// import { Schema } from 'pip-services3-commons-node';
// import { UnknownException } from 'pip-services3-commons-node';
// import { BadRequestException } from 'pip-services3-commons-node';
// import { Container } from 'pip-services3-container-node';
// import { Timing } from 'pip-services3-components-node';
// import { ConsoleLogger } from 'pip-services3-components-node';
// import { CompositeCounters } from 'pip-services3-components-node';

// /*
//  Abstract AWS Lambda function, that acts as a container to instantiate and run components
//  and expose them via external entry point.
//  *
//  When handling calls "cmd" parameter determines which what action shall be called, while
//  other parameters are passed to the action itself.
//  *
//  Container configuration for this Lambda function is stored in "./config/config.yml" file.
//  But this path can be overriden by CONFIG_PATH environment variable.
//  *
//  ### Configuration parameters ###
//  *
//  - dependencies:
//      - controller:                  override for Controller dependency
//  - connections:
//      - discovery_key:               (optional) a key to retrieve the connection from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]]
//      - region:                      (optional) AWS region
//  - credentials:
//      - store_key:                   (optional) a key to retrieve the credentials from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/auth.icredentialstore.html ICredentialStore]]
//      - access_id:                   AWS access/client id
//      - access_key:                  AWS access/client id
//  *
//  ### References ###
//  *
//  - \*:logger:\*:\*:1.0            (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages
//  - \*:counters:\*:\*:1.0          (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters]] components to pass collected measurements
//  - \*:discovery:\*:\*:1.0         (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]] services to resolve connection
//  - \*:credential-store:\*:\*:1.0  (optional) Credential stores to resolve credentials
//  *
//  See [[LambdaClient]]
//  *
//  ### Example ###
//  *
//      class MyLambdaFunction extends LambdaFunction {
//          private _controller: IMyController;
//          ...
//          public constructor() {
//              base("mygroup", "MyGroup lambda function");
//              this._dependencyResolver.put(
//                  "controller",
//                  new Descriptor("mygroup","controller","*","*","1.0")
//              );
//          }
//  *
//          public setReferences(references: IReferences): void {
//              base.setReferences(references);
//              this._controller = this._dependencyResolver.getRequired<IMyController>("controller");
//          }
//  *
//          public register(): void {
//              registerAction("get_mydata", null, (params, callback) => {
//                  let correlationId = params.correlation_id;
//                  let id = params.id;
//                  this._controller.getMyData(correlationId, id, callback);
//              });
//              ...
//          }
//      }
//  *
//      let lambda = new MyLambdaFunction();
//  *
//      service.run((err) => {
//          console.log("MyLambdaFunction is started");
//      });
//  */
// export abstract class LambdaFunction extends Container {
//     /*
//      The performanc counters.
//      */
//     protected _counters = new CompositeCounters();
//     /*
//      The dependency resolver.
//      */
//     protected _dependencyResolver = new DependencyResolver();
//     /*
//      The map of registred validation schemas.
//      */
//     protected _schemas: { [id: string]: Schema } = {};
//     /*
//      The map of registered actions.
//      */
//     protected _actions: { [id: string]: any } = {};
//     /*
//      The default path to config file.
//      */
//     protected _configPath: string = './config/config.yml';

//     /*
//      Creates a new instance of this lambda function.
//      *
//      - name          (optional) a container name (accessible via ContextInfo)
//      - description   (optional) a container description (accessible via ContextInfo)
//      */
//     public constructor(name?: string, description?: string) {
//         super(name, description);

//         this._logger = new ConsoleLogger();
//     }

//     private getConfigPath(): string {
//         return process.env.CONFIG_PATH || this._configPath;
//     }

//     private getParameters(): ConfigParams {
//         let parameters = ConfigParams.fromValue(process.env);
//         return parameters;
//     }

//     private captureErrors(correlationId: string): void {
//         // Log uncaught exceptions
//         process.on('uncaughtException', (ex) => {
//             this._logger.fatal(correlationId, ex, "Process is terminated");
//             process.exit(1);
//         });
//     }

//     private captureExit(correlationId: string): void {
//         this._logger.info(correlationId, "Press Control-C to stop the microservice...");

//         // Activate graceful exit
//         process.on('SIGINT', () => {
//             process.exit();
//         });

//         // Gracefully shutdown
//         process.on('exit', () => {
//             this.close(correlationId);
//             this._logger.info(correlationId, "Goodbye!");
//         });
//     }

// 	/*
// 	 Sets references to dependent components.
// 	 *
// 	 - references 	references to locate the component dependencies.
// 	 */
//     public setReferences(references: IReferences): void {
//         super.setReferences(references);
//         this._counters.setReferences(references);
//         this._dependencyResolver.setReferences(references);

//         this.register();
//     }

//     /*
//      Adds instrumentation to log calls and measure call time.
//      It returns a Timing object that is used to end the time measurement.
//      *
//      - correlationId     (optional) transaction id to trace execution through call chain.
//      - name              a method name.
//      Returns Timing object to end the time measurement.
//      */
//     protected instrument(correlationId: string, name: string): Timing {
//         this._logger.trace(correlationId, "Executing %s method", name);
//         return this._counters.beginTiming(name + ".exec_time");
//     }

//     /*
//      Runs this lambda function, loads container configuration,
//      instantiate components and manage their lifecycle,
//      makes this function ready to access action calls.
//      *
//      - callback callback function that receives error or null for success.
//      */
//     public run(callback?: (err: any) => void): void {
//         let correlationId = this._info.name;

//         let path = this.getConfigPath();
//         let parameters = this.getParameters();
//         this.readConfigFromFile(correlationId, path, parameters);

//         this.captureErrors(correlationId);
//         this.captureExit(correlationId);
//     	this.open(correlationId, callback);
//     }

//     /*
//      Registers all actions in this lambda function.
//      *
//      This method is called by the service and must be overriden
//      in child classes.
//      */
//     protected abstract register(): void;

//     /*
//      Registers an action in this lambda function.
//      *
//      - cmd           a action/command name.
//      - schema        a validation schema to validate received parameters.
//      - action        an action function that is called when action is invoked.
//      */
//     protected registerAction(cmd: string, schema: Schema,
//         action: (params: any, callback: (err: any, result: any) => void) => void): void {
//         if (cmd == '')
//             throw new UnknownException(null, 'NO_COMMAND', 'Missing command');

//         if (action == null)
//             throw new UnknownException(null, 'NO_ACTION', 'Missing action');

//         if (!_.isFunction(action))
//             throw new UnknownException(null, 'ACTION_NOT_FUNCTION', 'Action is not a function');

//         // Hack!!! Wrapping action to preserve prototyping context
//         let actionCurl = (params, callback) => {
//             // Perform validation
//             if (schema != null) {
//                 let correlationId = params.correlaton_id;
//                 let err = schema.validateAndReturnException(correlationId, params, false);
//                 if (err != null) {
//                     callback(err, null);
//                     return;
//                 }
//             }

//             // Todo: perform verification?
//             action.call(this, params, callback);
//         };

//         this._actions[cmd] = actionCurl;
//     }

//     private execute(event: any, context: any) {
//         let cmd: string = event.cmd;
//         let correlationId = event.correlation_id;

//         if (cmd == null) {
//             let err = new BadRequestException(
//                 correlationId,
//                 'NO_COMMAND',
//                 'Cmd parameter is missing'
//             );

//             context.done(err, null);
//             return;
//         }

//         let action: any = this._actions[cmd];
//         if (action == null) {
//             let err = new BadRequestException(
//                 correlationId,
//                 'NO_ACTION',
//                 'Action ' + cmd + ' was not found'
//             )
//             .withDetails('command', cmd);

//             context.done(err, null);
//             return;
//         }

//         action(event, context.done);
//     }

//     private handler(event: any, context: any) {
//         // If already started then execute
//         if (this.isOpen()) {
//             this.execute(event, context);
//         }
//         // Start before execute
//         else {
//             this.run((err) => {
//                 if (err) context.done(err, null);
//                 else this.execute(event, context);
//             });
//         }
//     }

//     /*
//      Gets entry point into this lambda function.
//      *
//      - event     an incoming event object with invocation parameters.
//      - context   a context object with local references.
//      */
//     public getHandler(): (event: any, context: any) => void {
//         let self = this;

//         // Return plugin function
//         return function (event, context) {
//             // Calling run with changed context
//             return self.handler.call(self, event, context);
//         }
//     }

//     /*
//      Calls registered action in this lambda function.
//      "cmd" parameter in the action parameters determin
//      what action shall be called.
//      *
//      This method shall only be used in testing.
//      *
//      - params action parameters.
//      - callback callback function that receives action result or error.
//      */
//     public act(params: any, callback: (err: any, result: any) => void): void {
//         let context = {
//             done: callback
//         };
//         this.getHandler()(params, context);
//     }

// }
